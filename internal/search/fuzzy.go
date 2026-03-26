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
)

// Search performs multi-field fuzzy search over bookmarks.
// Returns the full slice unchanged when query is empty.
// Results are sorted by composite score descending.
// Matching is case-insensitive.
func Search(query string, bookmarks []*store.Bookmark) []*store.Bookmark {
	if query == "" || len(bookmarks) == 0 {
		return bookmarks
	}

	q := strings.ToLower(query)

	// Score each bookmark across all three fields.
	scores := make(map[int64]int, len(bookmarks))

	score(q, bookmarks, func(b *store.Bookmark) string { return strings.ToLower(b.Title) }, weightTitle, scores)
	score(q, bookmarks, func(b *store.Bookmark) string { return strings.ToLower(b.Domain) }, weightDomain, scores)
	score(q, bookmarks, func(b *store.Bookmark) string { return strings.ToLower(b.Description) }, weightDescription, scores)

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
