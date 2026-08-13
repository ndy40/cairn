package search

import (
	"sort"
	"strings"

	"github.com/ndy40/cairn/internal/store"
	"github.com/sahilm/fuzzy"
)

const (
	weightTitle       = 3
	weightDomain      = 2
	weightDescription = 1
	weightTags        = 2

	// substringBase makes an exact case-insensitive substring match outrank any
	// loose fuzzy (subsequence) match. Fuzzy scores are small (tens); a substring
	// hit is a stronger signal and should always win and always be included.
	substringBase = 10000
)

// Search performs multi-field, case-insensitive search over bookmarks.
// Returns the full slice unchanged when query is empty.
// Results are sorted by composite score descending.
//
// Two signals contribute, in priority order:
//  1. An exact case-insensitive substring match in any field (e.g. "tau" in
//     "Restaurant"). This always counts and always ranks above fuzzy-only hits.
//  2. A fuzzy subsequence match, for typo-tolerant/loose ranking.
//
// The substring pass exists because fuzzy matching alone anchors greedily on the
// first occurrence of each character and can score a valid mid-word substring so
// low that it gets filtered out.
func Search(query string, bookmarks []*store.Bookmark) []*store.Bookmark {
	if query == "" || len(bookmarks) == 0 {
		return bookmarks
	}

	q := strings.ToLower(query)

	fields := []struct {
		get    func(*store.Bookmark) string
		weight int
	}{
		{func(b *store.Bookmark) string { return b.Title }, weightTitle},
		{func(b *store.Bookmark) string { return b.Domain }, weightDomain},
		{func(b *store.Bookmark) string { return b.Description }, weightDescription},
		{func(b *store.Bookmark) string { return strings.Join(b.Tags, " ") }, weightTags},
	}

	scores := make(map[int64]int, len(bookmarks))

	// Signal 1: exact case-insensitive substring match (guaranteed recall).
	for _, b := range bookmarks {
		for _, f := range fields {
			idx := strings.Index(strings.ToLower(f.get(b)), q)
			if idx < 0 {
				continue
			}
			// Earlier matches rank slightly higher within the same field weight.
			s := substringBase*f.weight - idx
			if s > scores[b.ID] {
				scores[b.ID] = s
			}
		}
	}

	// Signal 2: fuzzy subsequence match (ranks below any substring hit).
	for _, f := range fields {
		get := f.get
		score(q, bookmarks, func(b *store.Bookmark) string { return strings.ToLower(get(b)) }, f.weight, scores)
	}

	// Only include bookmarks that matched at least one field.
	var matched []*store.Bookmark
	for _, b := range bookmarks {
		if scores[b.ID] > 0 {
			matched = append(matched, b)
		}
	}

	// Sort by composite score descending.
	sort.Slice(matched, func(i, j int) bool {
		return scores[matched[i].ID] > scores[matched[j].ID]
	})

	return matched
}

// score runs fuzzy.FindFrom against a single field and accumulates weighted scores.
func score(query string, bookmarks []*store.Bookmark, field func(*store.Bookmark) string, weight int, scores map[int64]int) {
	src := &bookmarkSource{bookmarks: bookmarks, field: field}
	results := fuzzy.FindFrom(query, src)
	for _, r := range results {
		b := bookmarks[r.Index]
		weighted := r.Score * weight
		if weighted > scores[b.ID] {
			scores[b.ID] = weighted
		}
	}
}

// bookmarkSource adapts a []*store.Bookmark slice to fuzzy.Source.
type bookmarkSource struct {
	bookmarks []*store.Bookmark
	field     func(*store.Bookmark) string
}

func (s *bookmarkSource) String(i int) string { return s.field(s.bookmarks[i]) }
func (s *bookmarkSource) Len() int            { return len(s.bookmarks) }
