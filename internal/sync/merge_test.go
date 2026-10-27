package sync

import (
	"testing"
	"time"

	"github.com/ndy40/cairn/internal/store"
)

func TestMergeBookmarks_EarliestCreatedAtWins(t *testing.T) {
	early := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	late := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	local := []*store.Bookmark{
		{UUID: "a", URL: "https://a.example", CreatedAt: late, UpdatedAt: late},
		{UUID: "b-local", URL: "https://b.example", CreatedAt: late, UpdatedAt: late},
		{UUID: "c", URL: "https://c.example", CreatedAt: early, UpdatedAt: late},
		{UUID: "d", URL: "https://d.example", CreatedAt: late, UpdatedAt: late},
	}
	remote := &SyncRecord{Bookmarks: []BookmarkEntry{
		{UUID: "a", URL: "https://a.example", CreatedAt: early, UpdatedAt: late},        // matched by UUID, earlier
		{UUID: "b-remote", URL: "https://b.example", CreatedAt: early, UpdatedAt: late}, // matched by URL, earlier
		{UUID: "c", URL: "https://c.example", CreatedAt: late, UpdatedAt: late},         // remote later: keep local
		{UUID: "d", URL: "https://d.example", UpdatedAt: late},                          // zero remote date: ignored
	}}

	result := MergeBookmarks(local, remote)

	want := map[string]time.Time{"a": early, "b-local": early}
	if len(result.ToHeal) != len(want) {
		t.Fatalf("ToHeal = %+v, want %d entries", result.ToHeal, len(want))
	}
	for _, fix := range result.ToHeal {
		w, ok := want[fix.UUID]
		if !ok {
			t.Errorf("unexpected heal for %q", fix.UUID)
			continue
		}
		if !fix.CreatedAt.Equal(w) {
			t.Errorf("heal %q CreatedAt = %v, want %v", fix.UUID, fix.CreatedAt, w)
		}
	}
	if len(result.ToInsert) != 0 || len(result.ToUpdate) != 0 {
		t.Errorf("unexpected inserts/updates: %+v", result)
	}
}
