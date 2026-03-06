package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// ErrDuplicate is returned by Insert when the URL already exists.
var ErrDuplicate = errors.New("bookmark already exists")

// ErrNotFound is returned when a bookmark cannot be found by ID.
var ErrNotFound = errors.New("bookmark not found")

// Bookmark represents a saved web page.
type Bookmark struct {
	ID            int64
	URL           string
	Domain        string
	Title         string
	Description   string
	CreatedAt     time.Time
	Tags          []string
	LastVisitedAt *time.Time
	IsPermanent   bool
	IsArchived    bool
	ArchivedAt    *time.Time
}

// Insert saves a new bookmark. Returns ErrDuplicate if the URL already exists.
func (s *Store) Insert(rawURL, title, description string, tags []string) (*Bookmark, error) {
	domain := extractDomain(rawURL)
	now := time.Now().UTC()
	createdAt := now.Format(time.RFC3339)

	normTags := NormaliseTags(tags)
	tagsJSON, err := json.Marshal(normTags)
	if err != nil {
		return nil, fmt.Errorf("encode tags: %w", err)
	}

	res, err := s.db.Exec(
		`INSERT INTO bookmarks(url, domain, title, description, created_at, tags) VALUES (?, ?, ?, ?, ?, ?)`,
		rawURL, domain, title, description, createdAt, string(tagsJSON),
	)
	if err != nil {
		if isDuplicateErr(err) {
			return nil, ErrDuplicate
		}
		return nil, fmt.Errorf("insert bookmark: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Bookmark{
		ID:          id,
		URL:         rawURL,
		Domain:      domain,
		Title:       title,
		Description: description,
		CreatedAt:   now,
		Tags:        normTags,
	}, nil
}

// List returns all active (non-archived) bookmarks ordered by created_at descending.
func (s *Store) List() ([]*Bookmark, error) {
	rows, err := s.db.Query(
		`SELECT id, url, domain, title, description, created_at, tags, last_visited_at, is_permanent, is_archived, archived_at
		 FROM bookmarks WHERE is_archived = 0 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list bookmarks: %w", err)
	}
	defer rows.Close()
	return scanBookmarks(rows)
}

// DeleteByID removes a bookmark by its ID. Returns ErrNotFound if not found.
func (s *Store) DeleteByID(id int64) error {
	res, err := s.db.Exec(`DELETE FROM bookmarks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete bookmark: %w", err)
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

// ExistsByURL returns true if a bookmark with the given URL already exists.
func (s *Store) ExistsByURL(rawURL string) (bool, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(1) FROM bookmarks WHERE url = ?`, rawURL).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetByID returns the bookmark with the given ID or ErrNotFound.
func (s *Store) GetByID(id int64) (*Bookmark, error) {
	row := s.db.QueryRow(
		`SELECT id, url, domain, title, description, created_at, tags, last_visited_at, is_permanent, is_archived, archived_at
		 FROM bookmarks WHERE id = ?`, id,
	)
	b, err := scanBookmark(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return b, err
}

// ListByIDs returns bookmarks matching the provided IDs, preserving order.
func (s *Store) ListByIDs(ids []int64) ([]*Bookmark, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := fmt.Sprintf(
		`SELECT id, url, domain, title, description, created_at, tags, last_visited_at, is_permanent, is_archived, archived_at
		 FROM bookmarks WHERE id IN (%s) AND is_archived = 0`,
		strings.Join(placeholders, ","),
	)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBookmarks(rows)
}

// NormaliseTags splits, trims, lowercases, deduplicates, truncates to 32 runes,
// filters empty strings, and returns at most 3 tags.
func NormaliseTags(raw []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, t := range raw {
		t = strings.ToLower(strings.TrimSpace(t))
		if t == "" {
			continue
		}
		runes := []rune(t)
		if len(runes) > 32 {
			t = string(runes[:32])
		}
		if seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
		if len(out) == 3 {
			break
		}
	}
	return out
}

// UpdateTags replaces the tags on the bookmark with the given ID.
// Tags are normalised (lowercase, dedup, truncate, max 3) before saving.
func (s *Store) UpdateTags(id int64, tags []string) error {
	normTags := NormaliseTags(tags)
	tagsJSON, err := json.Marshal(normTags)
	if err != nil {
		return fmt.Errorf("encode tags: %w", err)
	}
	_, err = s.db.Exec(`UPDATE bookmarks SET tags = ? WHERE id = ?`, string(tagsJSON), id)
	if err != nil {
		return fmt.Errorf("update tags: %w", err)
	}
	return nil
}

// NormaliseTagsFromString splits a comma-separated tag string and normalises it.
func NormaliseTagsFromString(input string) []string {
	parts := strings.Split(input, ",")
	return NormaliseTags(parts)
}

// scanBookmarks reads all rows into a Bookmark slice.
func scanBookmarks(rows *sql.Rows) ([]*Bookmark, error) {
	var bookmarks []*Bookmark
	for rows.Next() {
		b, err := scanRow(rows.Scan)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, rows.Err()
}

// scanBookmark reads a single row into a Bookmark.
func scanBookmark(row *sql.Row) (*Bookmark, error) {
	return scanRow(row.Scan)
}

// scanRow is the shared scanning logic for both single and multi-row queries.
func scanRow(scan func(...any) error) (*Bookmark, error) {
	var b Bookmark
	var createdAt string
	var tagsJSON string
	var lastVisitedAt sql.NullString
	var archivedAt sql.NullString
	var isPermanent, isArchived int

	if err := scan(
		&b.ID, &b.URL, &b.Domain, &b.Title, &b.Description, &createdAt,
		&tagsJSON, &lastVisitedAt, &isPermanent, &isArchived, &archivedAt,
	); err != nil {
		return nil, err
	}

	if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
		b.CreatedAt = t
	}
	if err := json.Unmarshal([]byte(tagsJSON), &b.Tags); err != nil {
		b.Tags = []string{}
	}
	if lastVisitedAt.Valid {
		if t, err := time.Parse(time.RFC3339, lastVisitedAt.String); err == nil {
			b.LastVisitedAt = &t
		}
	}
	b.IsPermanent = isPermanent == 1
	b.IsArchived = isArchived == 1
	if archivedAt.Valid {
		if t, err := time.Parse(time.RFC3339, archivedAt.String); err == nil {
			b.ArchivedAt = &t
		}
	}
	return &b, nil
}

// extractDomain extracts a clean domain from a URL (strips www. prefix).
func extractDomain(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	h := u.Hostname()
	return strings.TrimPrefix(h, "www.")
}

// isDuplicateErr returns true if the error is a SQLite UNIQUE constraint error.
func isDuplicateErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}
