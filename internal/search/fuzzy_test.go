package search

import (
	"testing"

	"github.com/ndy40/cairn/internal/store"
)

func ids(bs []*store.Bookmark) map[int64]bool {
	m := make(map[int64]bool, len(bs))
	for _, b := range bs {
		m[b.ID] = true
	}
	return m
}

// TestSearchCaseInsensitiveMidWord guards the search behavior a restored
// bookmark relies on: matching is case-insensitive and finds the query as a
// subsequence anywhere in a field, including mid-word (e.g. "tau" in
// "Restaurant"). This is the recall that FTS5 prefix-matching used to drop.
func TestSearchCaseInsensitiveMidWord(t *testing.T) {
	bookmarks := []*store.Bookmark{
		{ID: 1, Title: "Best Restaurant Guide", Domain: "example.com"},
		{ID: 2, Title: "Kubernetes Notes", Domain: "k8s.io"},
		{ID: 3, Title: "Centaur systems", Domain: "trails.org"},
	}

	cases := []struct {
		name  string
		query string
		want  []int64
	}{
		// "tau" appears mid-word in "Res(tau)rant" and "Cen(tau)r".
		{"mid-word lowercase", "tau", []int64{1, 3}},
		{"mid-word uppercase query", "TAU", []int64{1, 3}},
		{"mixed case title match", "restaurant", []int64{1}},
		{"uppercase query mixed title", "ReStAuRaNt", []int64{1}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ids(Search(tc.query, bookmarks))
			for _, id := range tc.want {
				if !got[id] {
					t.Errorf("query %q: expected bookmark %d in results, got %v", tc.query, id, got)
				}
			}
		})
	}
}
