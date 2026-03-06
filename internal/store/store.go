package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "modernc.org/sqlite"
)

// Store holds the database connection.
type Store struct {
	db *sql.DB
}

// DefaultPath returns the OS-appropriate default database path.
func DefaultPath() (string, error) {
	var base string
	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, "Library", "Application Support", "cairn")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = home
		}
		base = filepath.Join(appData, "cairn")
	default: // linux and others
		xdg := os.Getenv("XDG_DATA_HOME")
		if xdg == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			xdg = filepath.Join(home, ".local", "share")
		}
		base = filepath.Join(xdg, "cairn")
	}
	return filepath.Join(base, "bookmarks.db"), nil
}

// Open opens (or creates) the SQLite database at the given path, runs
// migrations, and returns a Store ready for use.
func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, fmt.Errorf("create data directory: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Enable WAL mode for concurrent reads.
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable WAL mode: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

// Close closes the underlying database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// DB returns the underlying *sql.DB for use by subpackage functions.
func (s *Store) DB() *sql.DB {
	return s.db
}

// migrate creates the schema_version table and runs any pending migrations.
func (s *Store) migrate() error {
	// Ensure the version table exists.
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS schema_version (
		version    INTEGER PRIMARY KEY,
		applied_at TEXT NOT NULL
	)`)
	if err != nil {
		return err
	}

	var current int
	row := s.db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_version`)
	if err := row.Scan(&current); err != nil {
		return err
	}

	for _, m := range migrations {
		if m.version > current {
			if err := m.run(s.db); err != nil {
				return fmt.Errorf("migration v%d: %w", m.version, err)
			}
			_, err := s.db.Exec(
				`INSERT INTO schema_version(version, applied_at) VALUES (?, datetime('now'))`,
				m.version,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type migration struct {
	version int
	run     func(*sql.DB) error
}

var migrations = []migration{
	{version: 1, run: migrateV1},
	{version: 2, run: migrateV2},
}

func migrateV2(db *sql.DB) error {
	stmts := []string{
		`ALTER TABLE bookmarks ADD COLUMN tags            TEXT    NOT NULL DEFAULT '[]'`,
		`ALTER TABLE bookmarks ADD COLUMN last_visited_at TEXT`,
		`ALTER TABLE bookmarks ADD COLUMN is_permanent    INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE bookmarks ADD COLUMN is_archived     INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE bookmarks ADD COLUMN archived_at     TEXT`,
		`CREATE INDEX IF NOT EXISTS idx_bookmarks_is_archived  ON bookmarks(is_archived)`,
		`CREATE INDEX IF NOT EXISTS idx_bookmarks_archived_at  ON bookmarks(archived_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_bookmarks_is_permanent ON bookmarks(is_permanent)`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

func migrateV1(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS bookmarks (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			url         TEXT    NOT NULL UNIQUE,
			domain      TEXT    NOT NULL,
			title       TEXT    NOT NULL,
			description TEXT    NOT NULL DEFAULT '',
			created_at  TEXT    NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_bookmarks_domain ON bookmarks(domain)`,
		`CREATE INDEX IF NOT EXISTS idx_bookmarks_created_at ON bookmarks(created_at DESC)`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS bookmarks_fts USING fts5(
			title,
			description,
			domain,
			content='bookmarks',
			content_rowid='id'
		)`,
		`CREATE TRIGGER IF NOT EXISTS bookmarks_ai AFTER INSERT ON bookmarks BEGIN
			INSERT INTO bookmarks_fts(rowid, title, description, domain)
			VALUES (new.id, new.title, new.description, new.domain);
		END`,
		`CREATE TRIGGER IF NOT EXISTS bookmarks_ad AFTER DELETE ON bookmarks BEGIN
			INSERT INTO bookmarks_fts(bookmarks_fts, rowid, title, description, domain)
			VALUES ('delete', old.id, old.title, old.description, old.domain);
		END`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}
