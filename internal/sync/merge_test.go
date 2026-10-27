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

	result := MergeBookmarks(local, remote, nil)

	// b-local adopts the earlier remote copy's UUID, so its heal targets b-remote.
	want := map[string]time.Time{"a": early, "b-remote": early}
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
	if len(result.ToRekey) != 1 || result.ToRekey[0] != (UUIDFix{From: "b-local", To: "b-remote"}) {
		t.Errorf("ToRekey = %+v, want b-local → b-remote", result.ToRekey)
	}
	if len(result.ToInsert) != 0 || len(result.ToUpdate) != 0 {
		t.Errorf("unexpected inserts/updates: %+v", result)
	}
}

func TestMergeBookmarks_DivergedUUIDs(t *testing.T) {
	early := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	late := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		local     *store.Bookmark
		remote    BookmarkEntry
		wantRekey bool
	}{
		{
			name:      "remote earlier adopts remote uuid",
			local:     &store.Bookmark{UUID: "local", URL: "https://x.example", CreatedAt: late},
			remote:    BookmarkEntry{UUID: "remote", URL: "https://x.example", CreatedAt: early},
			wantRekey: true,
		},
		{
			name:   "local earlier keeps local uuid",
			local:  &store.Bookmark{UUID: "local", URL: "https://x.example", CreatedAt: early},
			remote: BookmarkEntry{UUID: "remote", URL: "https://x.example", CreatedAt: late},
		},
		{
			name:      "tie goes to smaller remote uuid",
			local:     &store.Bookmark{UUID: "bbb", URL: "https://x.example", CreatedAt: early},
			remote:    BookmarkEntry{UUID: "aaa", URL: "https://x.example", CreatedAt: early},
			wantRekey: true,
		},
		{
			name:   "tie keeps smaller local uuid",
			local:  &store.Bookmark{UUID: "aaa", URL: "https://x.example", CreatedAt: early},
			remote: BookmarkEntry{UUID: "bbb", URL: "https://x.example", CreatedAt: early},
		},
		{
			name:   "zero remote date never wins",
			local:  &store.Bookmark{UUID: "local", URL: "https://x.example", CreatedAt: late},
			remote: BookmarkEntry{UUID: "remote", URL: "https://x.example"},
		},
		{
			name:   "empty remote uuid is ignored",
			local:  &store.Bookmark{UUID: "local", URL: "https://x.example", CreatedAt: late},
			remote: BookmarkEntry{URL: "https://x.example", CreatedAt: early},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeBookmarks([]*store.Bookmark{tt.local}, &SyncRecord{Bookmarks: []BookmarkEntry{tt.remote}}, nil)
			if got := len(result.ToRekey) == 1; got != tt.wantRekey {
				t.Fatalf("ToRekey = %+v, want rekey %v", result.ToRekey, tt.wantRekey)
			}
			if tt.wantRekey && result.ToRekey[0] != (UUIDFix{From: tt.local.UUID, To: tt.remote.UUID}) {
				t.Errorf("ToRekey = %+v", result.ToRekey[0])
			}
		})
	}
}

func TestMergeBookmarks_SkipsLocallyDeleted(t *testing.T) {
	remote := &SyncRecord{
		Bookmarks: []BookmarkEntry{
			{UUID: "pending-uuid", URL: "https://a.example"},
			{UUID: "other-uuid", URL: "https://b.example"}, // same URL as a pending delete
			{UUID: "tombstoned", URL: "https://c.example"},
			{UUID: "fresh", URL: "https://d.example"},
		},
		Tombstones: []TombstoneEntry{{UUID: "tombstoned"}},
	}
	pending := []*store.PendingChange{
		{BookmarkUUID: "pending-uuid", Operation: "delete", Payload: `{"url":"https://a.example"}`},
		{BookmarkUUID: "diverged-uuid", Operation: "delete", Payload: `{"url":"https://b.example"}`},
		{BookmarkUUID: "fresh", Operation: "update", Payload: "{}"},
	}

	result := MergeBookmarks(nil, remote, pending)

	if len(result.ToInsert) != 1 || result.ToInsert[0].UUID != "fresh" {
		t.Errorf("ToInsert = %+v, want only fresh", result.ToInsert)
	}
}

func TestMergeBookmarks_UpdateUsesLocalUUID(t *testing.T) {
	early := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	late := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Local copy is older (keeps its UUID) but the remote edit is newer.
	local := []*store.Bookmark{{UUID: "local", URL: "https://x.example", Title: "Old", CreatedAt: early, UpdatedAt: early}}
	remote := &SyncRecord{Bookmarks: []BookmarkEntry{
		{UUID: "remote", URL: "https://x.example", Title: "New", CreatedAt: late, UpdatedAt: late},
	}}

	result := MergeBookmarks(local, remote, nil)

	if len(result.ToUpdate) != 1 || result.ToUpdate[0].UUID != "local" || result.ToUpdate[0].Title != "New" {
		t.Errorf("ToUpdate = %+v, want New under local UUID", result.ToUpdate)
	}
	if len(result.ToRekey) != 0 {
		t.Errorf("ToRekey = %+v, want none", result.ToRekey)
	}
}
