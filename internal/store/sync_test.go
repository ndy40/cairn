package store

import (
	"encoding/json"
	"testing"
	"time"
)

// TestInsertSyncedPreservesFields checks that a bookmark pulled from the sync
// snapshot keeps its identity and metadata and is not queued for upload again.
func TestInsertSyncedPreservesFields(t *testing.T) {
	s := openTestStore(t)

	created := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	updated := time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)
	archived := time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC)
	in := &Bookmark{
		UUID:        "11111111-2222-3333-4444-555555555555",
		URL:         "https://example.com/synced",
		Domain:      "example.com",
		Title:       "Synced",
		Description: "from the cloud",
		Tags:        []string{"go", "sync"},
		CreatedAt:   created,
		UpdatedAt:   updated,
		IsPermanent: true,
		IsArchived:  true,
		ArchivedAt:  &archived,
	}
	if err := s.InsertSynced(in); err != nil {
		t.Fatalf("InsertSynced: %v", err)
	}

	got, err := s.GetByUUID(in.UUID)
	if err != nil {
		t.Fatalf("GetByUUID: %v", err)
	}
	if got.URL != in.URL || got.Domain != in.Domain || got.Title != in.Title || got.Description != in.Description {
		t.Errorf("text fields mismatch: %+v", got)
	}
	if len(got.Tags) != 2 || got.Tags[0] != "go" || got.Tags[1] != "sync" {
		t.Errorf("tags = %v, want [go sync]", got.Tags)
	}
	if !got.CreatedAt.Equal(created) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, created)
	}
	if !got.UpdatedAt.Equal(updated) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, updated)
	}
	if !got.IsPermanent || !got.IsArchived {
		t.Errorf("IsPermanent=%v IsArchived=%v, want both true", got.IsPermanent, got.IsArchived)
	}
	if got.ArchivedAt == nil || !got.ArchivedAt.Equal(archived) {
		t.Errorf("ArchivedAt = %v, want %v", got.ArchivedAt, archived)
	}

	n, err := s.PendingChangeCount()
	if err != nil {
		t.Fatalf("PendingChangeCount: %v", err)
	}
	if n != 0 {
		t.Errorf("pending changes = %d, want 0", n)
	}

	// The FTS index must pick up the row like a normal insert.
	var hits int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM bookmarks_fts WHERE bookmarks_fts MATCH 'cloud'`).Scan(&hits); err != nil {
		t.Fatalf("fts query: %v", err)
	}
	if hits != 1 {
		t.Errorf("fts hits = %d, want 1", hits)
	}

	if err := s.InsertSynced(in); err != ErrDuplicate {
		t.Errorf("second InsertSynced err = %v, want ErrDuplicate", err)
	}
}

// TestSetCreatedAt checks that healing a creation date leaves updated_at alone
// and does not queue a pending change.
func TestSetCreatedAt(t *testing.T) {
	s := openTestStore(t)

	b, err := s.Insert("https://example.com/heal", "Heal", "", nil)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if err := s.ClearPendingChanges(); err != nil {
		t.Fatalf("clear pending: %v", err)
	}

	earlier := time.Date(2020, 5, 6, 7, 8, 9, 0, time.UTC)
	if err := s.SetCreatedAt(b.UUID, earlier); err != nil {
		t.Fatalf("SetCreatedAt: %v", err)
	}

	got, err := s.GetByUUID(b.UUID)
	if err != nil {
		t.Fatalf("GetByUUID: %v", err)
	}
	if !got.CreatedAt.Equal(earlier) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, earlier)
	}
	if !got.UpdatedAt.Equal(b.UpdatedAt.Truncate(time.Second)) {
		t.Errorf("UpdatedAt changed: got %v, want %v", got.UpdatedAt, b.UpdatedAt)
	}
	if n, _ := s.PendingChangeCount(); n != 0 {
		t.Errorf("pending changes = %d, want 0", n)
	}

	if err := s.SetCreatedAt("missing", earlier); err != ErrNotFound {
		t.Errorf("missing uuid err = %v, want ErrNotFound", err)
	}
}

// TestSetUUID checks that renaming a UUID carries pending changes along and
// queues nothing new.
func TestSetUUID(t *testing.T) {
	s := openTestStore(t)

	b, err := s.Insert("https://example.com/rekey", "Rekey", "", nil)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	other, err := s.Insert("https://example.com/other", "Other", "", nil)
	if err != nil {
		t.Fatalf("insert other: %v", err)
	}

	const newUUID = "99999999-8888-7777-6666-555555555555"
	if err := s.SetUUID(b.UUID, newUUID); err != nil {
		t.Fatalf("SetUUID: %v", err)
	}

	got, err := s.GetByUUID(newUUID)
	if err != nil || got.ID != b.ID {
		t.Fatalf("GetByUUID(new) = %+v, %v; want id %d", got, err, b.ID)
	}
	if _, err := s.GetByUUID(b.UUID); err != ErrNotFound {
		t.Errorf("old UUID still present: err = %v", err)
	}

	changes, err := s.ListPendingChanges()
	if err != nil {
		t.Fatalf("ListPendingChanges: %v", err)
	}
	if len(changes) != 2 {
		t.Fatalf("pending changes = %d, want 2 (the two inserts)", len(changes))
	}
	for _, c := range changes {
		if c.BookmarkUUID == b.UUID {
			t.Errorf("pending change still under old UUID: %+v", c)
		}
	}

	if err := s.SetUUID(newUUID, other.UUID); err != ErrDuplicate {
		t.Errorf("SetUUID to taken UUID err = %v, want ErrDuplicate", err)
	}
	if err := s.SetUUID("missing", "fresh"); err != ErrNotFound {
		t.Errorf("SetUUID missing err = %v, want ErrNotFound", err)
	}
}

// TestDeleteByIDRecordsURL checks that a pending delete remembers the URL, so
// sync can recognise copies of the bookmark stored under other UUIDs.
func TestDeleteByIDRecordsURL(t *testing.T) {
	s := openTestStore(t)

	b, err := s.Insert("https://example.com/gone", "Gone", "", nil)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if _, err := s.DeleteByID(b.ID); err != nil {
		t.Fatalf("DeleteByID: %v", err)
	}

	changes, err := s.ListPendingChanges()
	if err != nil {
		t.Fatalf("ListPendingChanges: %v", err)
	}
	var found bool
	for _, c := range changes {
		if c.Operation != "delete" {
			continue
		}
		found = true
		var p DeletePayload
		if err := json.Unmarshal([]byte(c.Payload), &p); err != nil {
			t.Fatalf("decode payload %q: %v", c.Payload, err)
		}
		if c.BookmarkUUID != b.UUID || p.URL != b.URL {
			t.Errorf("delete change = uuid %q url %q, want %q %q", c.BookmarkUUID, p.URL, b.UUID, b.URL)
		}
	}
	if !found {
		t.Fatal("no pending delete recorded")
	}
}

// TestUpdateFromSync checks that a remote edit is applied with the remote
// updated_at and without queuing a pending change.
func TestUpdateFromSync(t *testing.T) {
	s := openTestStore(t)

	b, err := s.Insert("https://example.com/old", "Old", "old desc", []string{"a"})
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	other, err := s.Insert("https://example.com/taken", "Taken", "", nil)
	if err != nil {
		t.Fatalf("insert other: %v", err)
	}
	if err := s.ClearPendingChanges(); err != nil {
		t.Fatalf("clear pending: %v", err)
	}

	updated := time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC)
	err = s.UpdateFromSync(&Bookmark{
		UUID: b.UUID, URL: "https://www.new.example/page", Title: "New", Description: "fresh words",
		Tags: []string{"B", "c"}, IsPermanent: true, UpdatedAt: updated,
	})
	if err != nil {
		t.Fatalf("UpdateFromSync: %v", err)
	}

	got, err := s.GetByUUID(b.UUID)
	if err != nil {
		t.Fatalf("GetByUUID: %v", err)
	}
	if got.URL != "https://www.new.example/page" || got.Domain != "new.example" || got.Title != "New" || got.Description != "fresh words" {
		t.Errorf("fields not applied: %+v", got)
	}
	if len(got.Tags) != 2 || got.Tags[0] != "b" || got.Tags[1] != "c" {
		t.Errorf("tags = %v, want [b c]", got.Tags)
	}
	if !got.IsPermanent {
		t.Error("pin not applied")
	}
	if !got.UpdatedAt.Equal(updated) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, updated)
	}
	if n, _ := s.PendingChangeCount(); n != 0 {
		t.Errorf("pending changes = %d, want 0", n)
	}

	var hits int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM bookmarks_fts WHERE bookmarks_fts MATCH 'fresh'`).Scan(&hits); err != nil {
		t.Fatalf("fts query: %v", err)
	}
	if hits != 1 {
		t.Errorf("fts hits = %d, want 1", hits)
	}

	if err := s.UpdateFromSync(&Bookmark{UUID: b.UUID, URL: other.URL, UpdatedAt: updated}); err != ErrDuplicateURL {
		t.Errorf("URL clash err = %v, want ErrDuplicateURL", err)
	}
	if err := s.UpdateFromSync(&Bookmark{UUID: "missing", URL: "https://example.com/z", UpdatedAt: updated}); err != ErrNotFound {
		t.Errorf("missing uuid err = %v, want ErrNotFound", err)
	}
}
