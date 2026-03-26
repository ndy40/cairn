package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// buildBinary compiles the cairn binary into a temp directory and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "cairn")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

// syncEnv returns environment variables that point XDG_CONFIG_HOME to a temp
// directory so the test never touches the real user config.
func syncEnv(t *testing.T, xdgHome string) []string {
	t.Helper()
	env := os.Environ()
	filtered := env[:0:len(env)]
	for _, e := range env {
		if len(e) >= 15 && e[:15] == "XDG_CONFIG_HOME" {
			continue
		}
		filtered = append(filtered, e)
	}
	return append(filtered, "XDG_CONFIG_HOME="+xdgHome)
}

// writeDeclinedSync writes a sync.json with SyncDeclined=true to xdgHome so
// the first-run prompt is suppressed without needing Dropbox credentials.
func writeDeclinedSync(t *testing.T, xdgHome string) {
	t.Helper()
	dir := filepath.Join(xdgHome, "cairn")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(map[string]bool{"sync_declined": true})
	if err := os.WriteFile(filepath.Join(dir, "sync.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
}

// TestListWithoutDropbox verifies that "cairn list" exits 0 when no Dropbox
// credentials are present.
func TestListWithoutDropbox(t *testing.T) {
	bin := buildBinary(t)
	xdg := t.TempDir()
	writeDeclinedSync(t, xdg)
	db := filepath.Join(t.TempDir(), "test.db")

	cmd := exec.Command(bin, "-db", db, "list")
	cmd.Env = syncEnv(t, xdg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cairn list failed: %v\n%s", err, out)
	}
}

// TestSearchWithoutDropbox verifies that "cairn search <query>" exits 0 when
// no Dropbox credentials are present.
func TestSearchWithoutDropbox(t *testing.T) {
	bin := buildBinary(t)
	xdg := t.TempDir()
	writeDeclinedSync(t, xdg)
	db := filepath.Join(t.TempDir(), "test.db")

	cmd := exec.Command(bin, "-db", db, "search", "golang")
	cmd.Env = syncEnv(t, xdg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cairn search failed: %v\n%s", err, out)
	}
}

// TestSyncStatusWithoutDropbox verifies that "cairn sync status" reports "not
// configured" cleanly (exit 0) when no sync credentials exist.
func TestSyncStatusWithoutDropbox(t *testing.T) {
	bin := buildBinary(t)
	xdg := t.TempDir()
	writeDeclinedSync(t, xdg)
	db := filepath.Join(t.TempDir(), "test.db")

	cmd := exec.Command(bin, "-db", db, "sync", "status")
	cmd.Env = syncEnv(t, xdg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cairn sync status failed: %v\n%s", err, out)
	}
}

// TestSyncPushWithoutDropbox verifies that "cairn sync push" exits with a
// non-zero code and does not panic when sync is not configured.
func TestSyncPushWithoutDropbox(t *testing.T) {
	bin := buildBinary(t)
	xdg := t.TempDir()
	writeDeclinedSync(t, xdg)
	db := filepath.Join(t.TempDir(), "test.db")

	cmd := exec.Command(bin, "-db", db, "sync", "push")
	cmd.Env = syncEnv(t, xdg)
	out, err := cmd.CombinedOutput()
	// We expect a non-zero exit because sync is not configured, but the
	// process must not crash with a signal (exitError with non-signal code).
	if err == nil {
		t.Fatalf("expected non-zero exit for unconfigured sync push, output:\n%s", out)
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		code := exitErr.ExitCode()
		if code <= 0 {
			t.Errorf("expected positive exit code, got %d", code)
		}
	}
}

// TestSyncPullWithoutDropbox verifies that "cairn sync pull" exits with a
// non-zero code and does not panic when sync is not configured.
func TestSyncPullWithoutDropbox(t *testing.T) {
	bin := buildBinary(t)
	xdg := t.TempDir()
	writeDeclinedSync(t, xdg)
	db := filepath.Join(t.TempDir(), "test.db")

	cmd := exec.Command(bin, "-db", db, "sync", "pull")
	cmd.Env = syncEnv(t, xdg)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit for unconfigured sync pull, output:\n%s", out)
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		code := exitErr.ExitCode()
		if code <= 0 {
			t.Errorf("expected positive exit code, got %d", code)
		}
	}
}

// TestListWithNoSyncConfigAtAll verifies that commands work even when the
// sync config file doesn't exist at all (fresh install, first run suppressed
// via non-interactive stdin).
func TestListWithNoSyncConfigAtAll(t *testing.T) {
	bin := buildBinary(t)
	xdg := t.TempDir() // empty — no sync.json exists
	db := filepath.Join(t.TempDir(), "test.db")

	cmd := exec.Command(bin, "-db", db, "list")
	cmd.Env = syncEnv(t, xdg)
	// Pipe empty stdin so the first-run prompt reads EOF and declines.
	cmd.Stdin = nil
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cairn list failed with no sync config: %v\n%s", err, out)
	}
}
