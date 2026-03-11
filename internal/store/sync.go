package store

import (
	"database/sql"
	"fmt"
	"time"
)

// PendingChange represents a bookmark modification waiting to be synced.
type PendingChange struct {
	ID           int64
	BookmarkUUID string
	Operation    string // "add", "update", "delete"
	Payload      string // JSON snapshot of bookmark
	CreatedAt    time.Time
	RetryCount   int
}

// InsertPendingChange records a pending sync change atomically within the given transaction.
// If tx is nil, a standalone exec is used.
func (s *Store) InsertPendingChange(tx *sql.Tx, bookmarkUUID, operation, payload string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO pending_sync(bookmark_uuid, operation, payload, created_at) VALUES (?, ?, ?, ?)`
	if tx != nil {
		_, err := tx.Exec(query, bookmarkUUID, operation, payload, now)
		return err
	}
	_, err := s.db.Exec(query, bookmarkUUID, operation, payload, now)
	return err
}

// ListPendingChanges returns all pending changes ordered by creation time.
func (s *Store) ListPendingChanges() ([]*PendingChange, error) {
	rows, err := s.db.Query(
		`SELECT id, bookmark_uuid, operation, payload, created_at, retry_count
		 FROM pending_sync ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list pending changes: %w", err)
	}
	defer rows.Close()

	var changes []*PendingChange
	for rows.Next() {
		var c PendingChange
		var createdAt string
		if err := rows.Scan(&c.ID, &c.BookmarkUUID, &c.Operation, &c.Payload, &createdAt, &c.RetryCount); err != nil {
			return nil, err
		}
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			c.CreatedAt = t
		}
		changes = append(changes, &c)
	}
	return changes, rows.Err()
}

// DeletePendingChange removes a pending change by ID.
func (s *Store) DeletePendingChange(id int64) error {
	_, err := s.db.Exec(`DELETE FROM pending_sync WHERE id = ?`, id)
	return err
}

// ClearPendingChanges removes all pending changes.
func (s *Store) ClearPendingChanges() error {
	_, err := s.db.Exec(`DELETE FROM pending_sync`)
	return err
}

// IncrementRetryCount increments the retry counter for a pending change.
func (s *Store) IncrementRetryCount(id int64) error {
	_, err := s.db.Exec(`UPDATE pending_sync SET retry_count = retry_count + 1 WHERE id = ?`, id)
	return err
}

// PendingChangeCount returns the number of pending sync changes.
func (s *Store) PendingChangeCount() (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM pending_sync`).Scan(&count)
	return count, err
}

// ExportAll returns all bookmarks (including archived) for sync export.
func (s *Store) ExportAll() ([]*Bookmark, error) {
	rows, err := s.db.Query(
		`SELECT id, uuid, url, domain, title, description, created_at, updated_at, tags, last_visited_at, is_permanent, is_archived, archived_at
		 FROM bookmarks ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("export all bookmarks: %w", err)
	}
	defer rows.Close()
	return scanBookmarks(rows)
}

// GetByUUID returns the bookmark with the given UUID or ErrNotFound.
func (s *Store) GetByUUID(bookmarkUUID string) (*Bookmark, error) {
	row := s.db.QueryRow(
		`SELECT id, uuid, url, domain, title, description, created_at, updated_at, tags, last_visited_at, is_permanent, is_archived, archived_at
		 FROM bookmarks WHERE uuid = ?`, bookmarkUUID,
	)
	b, err := scanBookmark(row)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return b, err
}

// BeginTx starts a new database transaction.
func (s *Store) BeginTx() (*sql.Tx, error) {
	return s.db.Begin()
}
