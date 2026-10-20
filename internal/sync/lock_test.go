package sync

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// isolateLockDir points DefaultConfigPath at a temp directory and creates it.
func isolateLockDir(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	if err := os.MkdirAll(filepath.Dir(DefaultConfigPath()), 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
}

func TestAcquireLockIsExclusive(t *testing.T) {
	isolateLockDir(t)

	path, f, ok := AcquireLock("push")
	if !ok || f == nil {
		t.Fatal("first AcquireLock failed")
	}
	defer func() { _ = f.Close() }()

	if _, f2, ok := AcquireLock("push"); ok || f2 != nil {
		t.Error("second AcquireLock succeeded while the lock is held by a live process")
	}
	if _, f3, ok := AcquireLock("pull"); !ok {
		t.Error("pull lock should be independent of the push lock")
	} else {
		_ = f3.Close()
	}

	_ = os.Remove(path)
	if _, f4, ok := AcquireLock("push"); !ok {
		t.Error("AcquireLock failed after the lock was released")
	} else {
		_ = f4.Close()
	}
}

func TestAcquireLockReclaimsStaleLocks(t *testing.T) {
	isolateLockDir(t)
	path := filepath.Join(filepath.Dir(DefaultConfigPath()), "cairn-sync-push.lock")

	// Lock left behind by a process that no longer exists.
	if err := os.WriteFile(path, []byte("999999999\n"), 0o600); err != nil {
		t.Fatalf("write lock: %v", err)
	}
	if _, f, ok := AcquireLock("push"); !ok {
		t.Error("lock held by a dead PID was not reclaimed")
	} else {
		_ = f.Close()
	}

	// Lock held by a live process (PID 1) but older than the TTL.
	if err := os.WriteFile(path, []byte("1\n"), 0o600); err != nil {
		t.Fatalf("write lock: %v", err)
	}
	old := time.Now().Add(-lockTTL - time.Minute)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
	if _, f, ok := AcquireLock("push"); !ok {
		t.Error("lock older than the TTL was not reclaimed")
	} else {
		_ = f.Close()
	}
}
