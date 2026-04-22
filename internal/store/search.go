package store

import (
	"fmt"
	"strings"
)

// FTSSearch queries the bookmarks_fts full-text index and returns matching
// bookmark IDs. Returns all IDs from the bookmarks table when the term is
// shorter than 3 characters (no FTS filtering applied).
func (s *Store) FTSSearch(term string) ([]int64, error) {
	if len([]rune(term)) < 3 {
		return s.allIDs()
	}

	escaped := ftsEscape(term)
	if escaped == "" {
		return s.allIDs()
	}

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

// ftsEscape converts a query into an FTS5 prefix-match expression.
// Each whitespace-delimited token becomes token* so partial terms match.
// FTS5 special characters are stripped to avoid syntax errors.
func ftsEscape(term string) string {
	tokens := strings.Fields(term)
	parts := make([]string, 0, len(tokens))
	for _, t := range tokens {
		clean := ftsStripSpecial(t)
		if clean != "" {
			parts = append(parts, clean+"*")
		}
	}
	return strings.Join(parts, " ")
}

// ftsStripSpecial removes characters that carry special meaning in FTS5 query syntax.
func ftsStripSpecial(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '"', '(', ')', '^', '*', ':', '{', '}', '[', ']', '~':
			// skip
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
