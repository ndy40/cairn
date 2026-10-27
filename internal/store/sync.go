package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	defer func() { _ = rows.Close() }()

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
	defer func() { _ = rows.Close() }()
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

// InsertSynced saves a bookmark pulled from the sync snapshot, keeping its UUID,
// timestamps, and pin/archive state. Unlike Insert it does not record a pending
// change, since the bookmark already exists remotely. The FTS index is kept in
// sync by the bookmarks_ai trigger. Returns ErrDuplicate if the URL or UUID
// already exists.
func (s *Store) InsertSynced(b *Bookmark) error {
	bookmarkUUID := b.UUID
	if bookmarkUUID == "" {
		bookmarkUUID = uuid.New().String()
	}
	domain := b.Domain
	if domain == "" {
		domain = extractDomain(b.URL)
	}
	createdAt := b.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	updatedAt := b.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	var archivedAt any
	if b.ArchivedAt != nil {
		archivedAt = b.ArchivedAt.UTC().Format(time.RFC3339)
	}

	tagsJSON, err := json.Marshal(NormaliseTags(b.Tags))
	if err != nil {
		return fmt.Errorf("encode tags: %w", err)
	}

	_, err = s.db.Exec(
		`INSERT INTO bookmarks(uuid, url, domain, title, description, tags, created_at, updated_at, is_permanent, is_archived, archived_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		bookmarkUUID, b.URL, domain, b.Title, b.Description, string(tagsJSON),
		createdAt.UTC().Format(time.RFC3339), updatedAt.UTC().Format(time.RFC3339),
		boolToInt(b.IsPermanent), boolToInt(b.IsArchived), archivedAt,
	)
	if err != nil {
		if isDuplicateErr(err) {
			return ErrDuplicate
		}
		return fmt.Errorf("insert synced bookmark: %w", err)
	}
	return nil
}

// SetCreatedAt overwrites the created_at of the bookmark with the given UUID.
// Used by sync to adopt an earlier creation date from another device; it leaves
// updated_at alone and records no pending change.
func (s *Store) SetCreatedAt(bookmarkUUID string, createdAt time.Time) error {
	res, err := s.db.Exec(
		`UPDATE bookmarks SET created_at = ? WHERE uuid = ?`,
		createdAt.UTC().Format(time.RFC3339), bookmarkUUID,
	)
	if err != nil {
		return fmt.Errorf("set created_at: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// SetUUID renames a bookmark's UUID, along with any pending changes queued
// under the old one. Used by sync to converge copies of the same URL on a single
// UUID; it records no pending change. Returns ErrNotFound if oldUUID does not
// exist and ErrDuplicate if newUUID is already taken.
func (s *Store) SetUUID(oldUUID, newUUID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.Exec(`UPDATE bookmarks SET uuid = ? WHERE uuid = ?`, newUUID, oldUUID)
	if err != nil {
		if isDuplicateErr(err) {
			return ErrDuplicate
		}
		return fmt.Errorf("set uuid: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}

	if _, err := tx.Exec(`UPDATE pending_sync SET bookmark_uuid = ? WHERE bookmark_uuid = ?`, newUUID, oldUUID); err != nil {
		return fmt.Errorf("rename pending changes: %w", err)
	}

	return tx.Commit()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
