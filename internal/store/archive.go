package store

import "fmt"

// ArchiveStale moves eligible stale bookmarks to the archive.
// A bookmark is eligible if it is not permanent AND not already archived AND
// created_at is 30 or more days ago.
// Returns the number of bookmarks archived.
func (s *Store) ArchiveStale() (int, error) {
	res, err := s.db.Exec(`
		UPDATE bookmarks
		SET    is_archived = 1,
		       archived_at = datetime('now')
		WHERE  is_permanent = 0
		  AND  is_archived  = 0
		  AND  created_at  <= datetime('now', '-30 days')
	`)
	if err != nil {
		return 0, fmt.Errorf("archive stale: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

// ListArchived returns all archived bookmarks ordered by archived_at descending.
func (s *Store) ListArchived() ([]*Bookmark, error) {
	rows, err := s.db.Query(
		`SELECT id, url, domain, title, description, created_at, tags, last_visited_at, is_permanent, is_archived, archived_at
		 FROM bookmarks WHERE is_archived = 1 ORDER BY archived_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list archived: %w", err)
	}
	defer rows.Close()
	return scanBookmarks(rows)
}

// RestoreByID moves an archived bookmark back to the active list.
func (s *Store) RestoreByID(id int64) error {
	_, err := s.db.Exec(
		`UPDATE bookmarks SET is_archived = 0, archived_at = NULL WHERE id = ?`, id,
	)
	if err != nil {
		return fmt.Errorf("restore bookmark: %w", err)
	}
	return nil
}

// SetPermanent sets or clears the permanent flag on a bookmark.
func (s *Store) SetPermanent(id int64, permanent bool) error {
	val := 0
	if permanent {
		val = 1
	}
	_, err := s.db.Exec(
		`UPDATE bookmarks SET is_permanent = ? WHERE id = ?`, val, id,
	)
	if err != nil {
		return fmt.Errorf("set permanent: %w", err)
	}
	return nil
}
