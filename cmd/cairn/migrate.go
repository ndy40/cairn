package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ndy40/cairn/internal/config"
)

// migrateDBPath performs a one-time migration of the database file from the
// legacy "bookmark-manager" directory to the new "cairn" directory.
//
// It only runs when the resolved DB path equals the current default (i.e. the
// user has not set a custom path via config file, env var, or --db flag).
// If the new path already exists the function is a no-op.
func migrateDBPath(resolvedDB string) {
	// Skip if the user has configured an explicit path.
	if resolvedDB != config.DefaultDBPath() {
		return
	}

	newPath := resolvedDB
	oldPath := config.LegacyDBPath()

	// Nothing to do if the new path already has a database.
	if fileExists(newPath) {
		return
	}

	// Nothing to migrate if the old path doesn't exist either.
	if !fileExists(oldPath) {
		return
	}

	// Ensure the destination directory exists.
	if err := os.MkdirAll(filepath.Dir(newPath), 0o700); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "cairn: migrate: create directory: %v\n", err)
		return
	}

	// Prefer an atomic rename; fall back to copy+delete for cross-device moves.
	if err := os.Rename(oldPath, newPath); err == nil {
		fmt.Fprintf(os.Stderr, "cairn: migrated database from %s to %s\n", oldPath, newPath)
		removeEmptyDir(filepath.Dir(oldPath))
		return
	}

	if err := copyFile(oldPath, newPath); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "cairn: migrate: copy database: %v\n", err)
		return
	}

	_ = os.Remove(oldPath)
	fmt.Fprintf(os.Stderr, "cairn: migrated database from %s to %s\n", oldPath, newPath)
	removeEmptyDir(filepath.Dir(oldPath))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// copyFile copies src to dst, creating dst if it does not exist.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, in); err != nil {
		_ = os.Remove(dst)
		return err
	}
	return out.Close()
}

// removeEmptyDir removes dir only if it contains no files, so we never
// silently destroy anything the user may have put there.
func removeEmptyDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) > 0 {
		return
	}
	_ = os.Remove(dir)
}
