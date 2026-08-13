package store

import (
	"path/filepath"
	"testing"
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
