package store

import "fmt"

// FTSSearch queries the bookmarks_fts full-text index and returns matching
// bookmark IDs. Returns all IDs from the bookmarks table when the term is
// shorter than 3 characters (no FTS filtering applied).
func (s *Store) FTSSearch(term string) ([]int64, error) {
	if len([]rune(term)) < 3 {
		return s.allIDs()
	}

	// Escape special FTS5 characters.
	escaped := ftsEscape(term)

	rows, err := s.db.Query(
		`SELECT rowid FROM bookmarks_fts WHERE bookmarks_fts MATCH ? ORDER BY rank`,
		escaped,
	)
	if err != nil {
		// On FTS syntax error, fall back to returning all IDs.
		return s.allIDs()
	}
	defer func() { _ = rows.Close() }()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// allIDs returns all bookmark IDs ordered by created_at DESC.
func (s *Store) allIDs() ([]int64, error) {
	rows, err := s.db.Query(`SELECT id FROM bookmarks ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("allIDs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// ftsEscape wraps the query in double quotes for FTS5 phrase matching,
// escaping any embedded double quotes.
func ftsEscape(term string) string {
	escaped := ""
	for _, r := range term {
		if r == '"' {
			escaped += `""`
		} else {
			escaped += string(r)
		}
	}
	return `"` + escaped + `"`
}
