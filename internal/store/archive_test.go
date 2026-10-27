package store

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"
)

// openTestStore creates a Store backed by a temporary SQLite database.
func openTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "bookmarks.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

// TestListArchivedRoundtrip guards against the column-count regression where
// ListArchived selected fewer columns than the shared row scanner expected,
// causing every scan to fail and the archive screen to render empty.
func TestListArchivedRoundtrip(t *testing.T) {
	s := openTestStore(t)

	b, err := s.Insert("https://example.com/page", "Example", "desc", []string{"go"})
	if err != nil {
		t.Fatalf("insert: %v", err)
	}

	// Archive the row directly (ArchiveStale only touches 30-day-old rows).
	if _, err := s.db.Exec(
		`UPDATE bookmarks SET is_archived = 1, archived_at = datetime('now') WHERE id = ?`, b.ID,
	); err != nil {
		t.Fatalf("archive row: %v", err)
	}

	archived, err := s.ListArchived()
	if err != nil {
		t.Fatalf("ListArchived returned error: %v", err)
	}
	if len(archived) != 1 {
		t.Fatalf("expected 1 archived bookmark, got %d", len(archived))
	}

	got := archived[0]
	if got.ID != b.ID || got.UUID != b.UUID || got.URL != b.URL {
		t.Errorf("scanned fields mismatch: got id=%d uuid=%q url=%q", got.ID, got.UUID, got.URL)
	}
	if !got.IsArchived {
		t.Error("expected IsArchived to be true")
	}
	if got.ArchivedAt == nil {
		t.Error("expected ArchivedAt to be parsed from datetime('now') format, got nil")
	}

	// Archived rows must not leak into the active list.
	active, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(active) != 0 {
		t.Fatalf("expected archived bookmark excluded from active list, got %d active", len(active))
	}
}

// TestRestoreByID verifies an archived bookmark returns to the active list.
func TestRestoreByID(t *testing.T) {
	s := openTestStore(t)

	b, err := s.Insert("https://example.com/restore", "Restore", "", nil)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if _, err := s.db.Exec(
		`UPDATE bookmarks SET is_archived = 1, archived_at = datetime('now') WHERE id = ?`, b.ID,
	); err != nil {
		t.Fatalf("archive row: %v", err)
	}

	if err := s.RestoreByID(b.ID); err != nil {
		t.Fatalf("RestoreByID: %v", err)
	}

	archived, err := s.ListArchived()
	if err != nil {
		t.Fatalf("ListArchived: %v", err)
	}
	if len(archived) != 0 {
		t.Fatalf("expected no archived bookmarks after restore, got %d", len(archived))
	}

	active, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(active) != 1 {
		t.Fatalf("expected 1 active bookmark after restore, got %d", len(active))
	}
	if active[0].ArchivedAt != nil {
		t.Error("expected ArchivedAt cleared after restore")
	}
}

// firstSeenAt reads the local-only first_seen_at column for a bookmark.
func firstSeenAt(t *testing.T, s *Store, id int64) string {
	t.Helper()
	var v sql.NullString
	if err := s.db.QueryRow(`SELECT first_seen_at FROM bookmarks WHERE id = ?`, id).Scan(&v); err != nil {
		t.Fatalf("read first_seen_at: %v", err)
	}
	return v.String
}

// TestArchiveStaleWaitsForFirstSeen checks that a bookmark whose created_at was
// moved back by sync is not archived until it has been on this device 30 days.
func TestArchiveStaleWaitsForFirstSeen(t *testing.T) {
	s := openTestStore(t)

	old := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	// Healed: local row whose created_at sync moved back.
	healed, err := s.Insert("https://example.com/healed", "Healed", "", nil)
	if err != nil {
		t.Fatalf("insert healed: %v", err)
	}
	if err := s.SetCreatedAt(healed.UUID, old); err != nil {
		t.Fatalf("SetCreatedAt: %v", err)
	}
	if got := firstSeenAt(t, s, healed.ID); got != healed.CreatedAt.Format(time.RFC3339) {
		t.Errorf("first_seen_at after heal = %q, want insert time", got)
	}

	// Pulled: old bookmark inserted from the snapshot today.
	pulled := &Bookmark{UUID: "pulled", URL: "https://example.com/pulled", Title: "Pulled", CreatedAt: old, UpdatedAt: old}
	if err := s.InsertSynced(pulled); err != nil {
		t.Fatalf("InsertSynced: %v", err)
	}

	// Seen long ago: eligible on both counts.
	stale, err := s.Insert("https://example.com/stale", "Stale", "", nil)
	if err != nil {
		t.Fatalf("insert stale: %v", err)
	}
	oldStr := old.Format(time.RFC3339)
	if _, err := s.db.Exec(`UPDATE bookmarks SET created_at = ?, first_seen_at = ? WHERE id = ?`, oldStr, oldStr, stale.ID); err != nil {
		t.Fatalf("age stale row: %v", err)
	}

	n, err := s.ArchiveStale()
	if err != nil {
		t.Fatalf("ArchiveStale: %v", err)
	}
	if n != 1 {
		t.Fatalf("archived %d bookmarks, want 1", n)
	}
	archived, err := s.ListArchived()
	if err != nil {
		t.Fatalf("ListArchived: %v", err)
	}
	if len(archived) != 1 || archived[0].ID != stale.ID {
		t.Errorf("archived = %+v, want only the stale row", archived)
	}
}

// TestMigrateV5BackfillsFirstSeen checks that existing rows get first_seen_at
// from their current created_at.
func TestMigrateV5BackfillsFirstSeen(t *testing.T) {
	s := openTestStore(t)

	b, err := s.Insert("https://example.com/existing", "Existing", "", nil)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	const created = "2024-05-06T07:08:09Z"

	// Roll the database back to v4 with a pre-existing row.
	for _, stmt := range []string{
		`ALTER TABLE bookmarks DROP COLUMN first_seen_at`,
		`DELETE FROM schema_version WHERE version = 5`,
	} {
		if _, err := s.db.Exec(stmt); err != nil {
			t.Fatalf("roll back to v4: %v", err)
		}
	}
	if _, err := s.db.Exec(`UPDATE bookmarks SET created_at = ? WHERE id = ?`, created, b.ID); err != nil {
		t.Fatalf("set created_at: %v", err)
	}

	if err := s.migrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if got := firstSeenAt(t, s, b.ID); got != created {
		t.Errorf("first_seen_at = %q, want %q", got, created)
	}
}

// TestSetPermanentQueuesSync checks that pinning counts as an edit so it syncs.
func TestSetPermanentQueuesSync(t *testing.T) {
	s := openTestStore(t)

	b, err := s.Insert("https://example.com/pin", "Pin", "", nil)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if _, err := s.db.Exec(`UPDATE bookmarks SET updated_at = '2020-01-01T00:00:00Z' WHERE id = ?`, b.ID); err != nil {
		t.Fatalf("age updated_at: %v", err)
	}
	if err := s.ClearPendingChanges(); err != nil {
		t.Fatalf("clear pending: %v", err)
	}

	if err := s.SetPermanent(b.ID, true); err != nil {
		t.Fatalf("SetPermanent: %v", err)
	}

	got, err := s.GetByID(b.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if !got.IsPermanent {
		t.Error("expected IsPermanent to be true")
	}
	if got.UpdatedAt.Year() == 2020 {
		t.Error("updated_at was not bumped")
	}
	changes, err := s.ListPendingChanges()
	if err != nil {
		t.Fatalf("ListPendingChanges: %v", err)
	}
	if len(changes) != 1 || changes[0].Operation != "update" || changes[0].BookmarkUUID != b.UUID {
		t.Errorf("pending changes = %+v, want one update for %s", changes, b.UUID)
	}

	if err := s.SetPermanent(9999, true); err != ErrNotFound {
		t.Errorf("missing id err = %v, want ErrNotFound", err)
	}
}
