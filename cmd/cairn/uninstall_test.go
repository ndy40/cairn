package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func touch(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestCollectDataPaths(t *testing.T) {
	dir := t.TempDir()
	db := filepath.Join(dir, "data", "bookmarks.db")
	legacy := filepath.Join(dir, "legacy", "bookmarks.db")
	cfg := filepath.Join(dir, "config", "cairn.json")
	syncFile := filepath.Join(dir, "config", "sync.json")
	lock := filepath.Join(dir, "config", "cairn-sync-pull.lock")

	touch(t, db)
	touch(t, db+"-wal")
	touch(t, cfg)
	touch(t, lock)
	// legacy, db-shm and sync.json deliberately absent.

	got := collectDataPaths(db, legacy, cfg, syncFile)
	want := []string{db, db + "-wal", cfg, lock}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("collectDataPaths:\n got  %v\n want %v", got, want)
	}

	// The same file passed twice is listed once.
	if got := collectDataPaths(db, db, cfg, cfg); !reflect.DeepEqual(got, want) {
		t.Errorf("expected duplicates removed:\n got  %v\n want %v", got, want)
	}
}

func TestExecuteUninstall(t *testing.T) {
	setup := func(t *testing.T) (*uninstallPlan, string) {
		dir := t.TempDir()
		bin := filepath.Join(dir, "bin", "cairn")
		ext := filepath.Join(dir, "vicinae", "extensions", "cairn")
		cairnDir := filepath.Join(dir, "share", "cairn")
		customDir := filepath.Join(dir, "Documents")
		db := filepath.Join(customDir, "bookmarks.db")
		cfg := filepath.Join(cairnDir, "cairn.json")
		app := filepath.Join(dir, "Applications", "cairn-desktop.app")
		desktopFiles := []string{
			filepath.Join(cairnDir, "desktop-version.txt"),
			filepath.Join(dir, "LaunchAgents", "io.github.ndy40.cairn.desktop.plist"),
		}

		touch(t, bin)
		touch(t, filepath.Join(ext, "package.json"))
		touch(t, filepath.Join(app, "Contents", "MacOS", "cairn-desktop"))
		touch(t, db)
		touch(t, cfg)
		for _, p := range desktopFiles {
			touch(t, p)
		}

		return &uninstallPlan{
			binary:       bin,
			extension:    ext,
			desktop:      app,
			desktopFiles: desktopFiles,
			data:         []string{db, cfg},
			ownDirs:      []string{cairnDir},
		}, dir
	}

	t.Run("keeps data by default", func(t *testing.T) {
		plan, _ := setup(t)
		if err := executeUninstall(plan, false); err != nil {
			t.Fatal(err)
		}
		if fileExists(plan.binary) {
			t.Error("binary not removed")
		}
		if fileExists(plan.extension) {
			t.Error("extension not removed")
		}
		for _, p := range append([]string{plan.desktop}, plan.desktopFiles...) {
			if fileExists(p) {
				t.Errorf("desktop app file %s not removed", p)
			}
		}
		for _, p := range plan.data {
			if !fileExists(p) {
				t.Errorf("data file %s was removed", p)
			}
		}
	})

	t.Run("purge removes data and empty cairn dirs only", func(t *testing.T) {
		plan, dir := setup(t)
		if err := executeUninstall(plan, true); err != nil {
			t.Fatal(err)
		}
		for _, p := range plan.data {
			if fileExists(p) {
				t.Errorf("data file %s not removed", p)
			}
		}
		if fileExists(plan.ownDirs[0]) {
			t.Error("empty cairn directory not removed")
		}
		// A custom --db directory must survive even when left empty.
		if !fileExists(filepath.Join(dir, "Documents")) {
			t.Error("custom database directory was removed")
		}
	})

	t.Run("purge keeps cairn dir holding other files", func(t *testing.T) {
		plan, _ := setup(t)
		other := filepath.Join(plan.ownDirs[0], "notes.txt")
		touch(t, other)
		if err := executeUninstall(plan, true); err != nil {
			t.Fatal(err)
		}
		if !fileExists(other) {
			t.Error("unrelated file in cairn directory was removed")
		}
	})
}

func TestConfirm(t *testing.T) {
	cases := map[string]bool{
		"y\n":   true,
		"YES\n": true,
		" y \n": true,
		"n\n":   false,
		"\n":    false,
		"":      false, // EOF
		"y":     true,  // no trailing newline
	}
	for input, want := range cases {
		if got := confirm(strings.NewReader(input), ""); got != want {
			t.Errorf("confirm(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestRmCommand(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "it's.db")
	touch(t, file)
	if got, want := rmCommand([]string{file}), "rm "+shellQuote(file); got != want {
		t.Errorf("rmCommand(file) = %s, want %s", got, want)
	}
	if got, want := rmCommand([]string{file, dir}), "rm -r "+shellQuote(file)+" "+shellQuote(dir); got != want {
		t.Errorf("rmCommand(file, dir) = %s, want %s", got, want)
	}
}

func TestShellQuote(t *testing.T) {
	got := shellQuote("/Users/me/Library/Application Support/it's/cairn.json")
	want := `'/Users/me/Library/Application Support/it'\''s/cairn.json'`
	if got != want {
		t.Errorf("shellQuote = %s, want %s", got, want)
	}
}

// uninstallEnv points HOME and the XDG directories at root so the built binary
// never sees, and can never delete, the real user's data. PATH is limited to
// root/bin so a real cairn-desktop on the machine is never found.
func uninstallEnv(root string) []string {
	var env []string
	for _, e := range os.Environ() {
		switch strings.SplitN(e, "=", 2)[0] {
		case "HOME", "PATH", "XDG_DATA_HOME", "XDG_CONFIG_HOME", "CAIRN_DB_PATH", "SUDO_USER":
			continue
		}
		env = append(env, e)
	}
	return append(env,
		"HOME="+root,
		"PATH="+filepath.Join(root, "bin"),
		"XDG_DATA_HOME="+filepath.Join(root, "share"),
		"XDG_CONFIG_HOME="+filepath.Join(root, "config"),
	)
}

// uninstallFixture is a built cairn binary with a database, sync config,
// Vicinae extension and desktop app laid out as the installer would.
type uninstallFixture struct {
	root, bin, db, syncFile, ext, desktop string
	desktopFiles                          []string
}

func newUninstallFixture(t *testing.T, syncJSON string) *uninstallFixture {
	t.Helper()
	f := &uninstallFixture{bin: buildBinary(t), root: t.TempDir()}
	share, config := filepath.Join(f.root, "share"), filepath.Join(f.root, "config")
	f.db = filepath.Join(share, "cairn", "bookmarks.db")
	f.syncFile = filepath.Join(config, "cairn", "sync.json")
	f.ext = filepath.Join(share, "vicinae", "extensions", "cairn")
	f.desktop = filepath.Join(filepath.Dir(f.bin), "cairn-desktop")
	f.desktopFiles = []string{
		filepath.Join(share, "cairn", "desktop-version.txt"),
		filepath.Join(share, "cairn", "cairn-desktop.png"),
		filepath.Join(share, "applications", "cairn-desktop.desktop"),
		filepath.Join(config, "autostart", "cairn-desktop.desktop"),
	}
	for _, p := range append([]string{f.db, f.syncFile, f.desktop, filepath.Join(f.ext, "version.txt")}, f.desktopFiles...) {
		touch(t, p)
	}
	if err := os.WriteFile(f.syncFile, []byte(syncJSON), 0o600); err != nil {
		t.Fatal(err)
	}
	return f
}

func (f *uninstallFixture) run(t *testing.T, stdin string, args ...string) string {
	t.Helper()
	cmd := exec.Command(f.bin, append([]string{"uninstall"}, args...)...)
	cmd.Dir = f.root
	cmd.Env = uninstallEnv(f.root)
	cmd.Stdin = strings.NewReader(stdin)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cairn uninstall %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func TestUninstallCommand(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("relies on XDG paths")
	}

	const linkedSync = `{"backend":"dropbox","dropbox":{"access_token":"t","app_key":"k"}}`
	const declinedSync = `{"sync_declined":true}`

	t.Run("default keeps data", func(t *testing.T) {
		f := newUninstallFixture(t, linkedSync)
		out := f.run(t, "", "-y")
		for _, p := range append([]string{f.bin, f.ext, f.desktop}, f.desktopFiles...) {
			if fileExists(p) {
				t.Errorf("%s not removed", p)
			}
		}
		if !fileExists(f.db) || !fileExists(f.syncFile) {
			t.Error("data removed without --purge")
		}
		for _, s := range []string{f.db, "contains your Dropbox access token", "rm '" + f.db + "'", "Dropbox has not been deleted"} {
			if !strings.Contains(out, s) {
				t.Errorf("output missing %q:\n%s", s, out)
			}
		}
	})

	t.Run("purge removes data", func(t *testing.T) {
		f := newUninstallFixture(t, linkedSync)
		out := f.run(t, "", "--purge", "--yes")
		if fileExists(f.bin) || fileExists(f.db) || fileExists(f.syncFile) {
			t.Errorf("files left after --purge:\n%s", out)
		}
		// The desktop icon and version file share this directory with the database.
		if fileExists(filepath.Dir(f.db)) {
			t.Error("empty data directory not removed")
		}
	})

	t.Run("declined sync has no Dropbox notes", func(t *testing.T) {
		f := newUninstallFixture(t, declinedSync)
		out := f.run(t, "", "-y")
		if !strings.Contains(out, f.syncFile) {
			t.Errorf("sync.json not listed:\n%s", out)
		}
		for _, s := range []string{"Dropbox access token", "Dropbox has not been deleted"} {
			if strings.Contains(out, s) {
				t.Errorf("output mentions %q although sync was declined:\n%s", s, out)
			}
		}
	})

	t.Run("desktop app elsewhere on PATH is left alone", func(t *testing.T) {
		f := newUninstallFixture(t, linkedSync)
		other := filepath.Join(f.root, "bin", "cairn-desktop")
		touch(t, other)
		if err := os.Chmod(other, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(f.desktop); err != nil {
			t.Fatal(err)
		}
		out := f.run(t, "", "-y")
		if !fileExists(other) {
			t.Error("desktop app on PATH was removed")
		}
		for _, p := range f.desktopFiles {
			if !fileExists(p) {
				t.Errorf("%s belongs to the other desktop app but was removed", p)
			}
		}
		if !strings.Contains(out, "is left in place") {
			t.Errorf("output does not mention the kept desktop app:\n%s", out)
		}
	})

	t.Run("declining the prompt changes nothing", func(t *testing.T) {
		f := newUninstallFixture(t, linkedSync)
		out := f.run(t, "n\n")
		for _, p := range append([]string{f.bin, f.ext, f.desktop, f.db}, f.desktopFiles...) {
			if !fileExists(p) {
				t.Errorf("%s removed after answering no", p)
			}
		}
		if !strings.Contains(out, "Uninstall cancelled.") {
			t.Errorf("unexpected output:\n%s", out)
		}
	})
}
