package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWindowStateRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "desktop-state.json")

	if err := (windowState{Width: 900, Height: 600}).save(path); err != nil {
		t.Fatalf("save: %v", err)
	}
	got := loadWindowState(path)
	if got.Width != 900 || got.Height != 600 {
		t.Errorf("loaded %+v, want 900x600", got)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Error("temporary file left behind after save")
	}
}

func TestLoadWindowStateDefaults(t *testing.T) {
	dir := t.TempDir()
	cases := map[string]string{
		"missing":       "",
		"corrupt":       "{not json",
		"too small":     `{"width": 100, "height": 50}`,
		"zero (hidden)": `{"width": 0, "height": 0}`,
	}
	for name, body := range cases {
		path := filepath.Join(dir, name+".json")
		if body != "" {
			if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
				t.Fatalf("write: %v", err)
			}
		}
		got := loadWindowState(path)
		if got.Width != defaultWidth || got.Height != defaultHeight {
			t.Errorf("%s: loaded %+v, want default %dx%d", name, got, defaultWidth, defaultHeight)
		}
	}
}

func TestSaveIgnoresHiddenWindowSize(t *testing.T) {
	path := filepath.Join(t.TempDir(), "desktop-state.json")
	if err := (windowState{Width: 800, Height: 500}).save(path); err != nil {
		t.Fatalf("save: %v", err)
	}

	// A hidden window can report 0x0; that must not replace the real size.
	if err := (windowState{Width: 0, Height: 0}).save(path); err != nil {
		t.Fatalf("save: %v", err)
	}
	got := loadWindowState(path)
	if got.Width != 800 || got.Height != 500 {
		t.Errorf("loaded %+v, want the earlier 800x500", got)
	}
}

func TestDefaultStatePathIsBesideConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("APPDATA", dir)

	got := defaultStatePath()
	if filepath.Base(got) != "desktop-state.json" || filepath.Base(filepath.Dir(got)) != "cairn" {
		t.Errorf("defaultStatePath() = %q, want .../cairn/desktop-state.json", got)
	}
}
