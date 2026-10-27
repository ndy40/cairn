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

func TestLoadWindowStateWithoutPosition(t *testing.T) {
	// Files written before positions were saved hold only the size.
	path := filepath.Join(t.TempDir(), "desktop-state.json")
	if err := os.WriteFile(path, []byte(`{"width": 800, "height": 500}`), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	got := loadWindowState(path)
	if got.Width != 800 || got.Height != 500 || got.X != nil || got.Y != nil || got.Screen != "" {
		t.Errorf("loaded %+v, want 800x500 with no position", got)
	}
}

func TestWindowStatePositionRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "desktop-state.json")
	x, y := 40, -10
	if err := (windowState{Width: 800, Height: 500, X: &x, Y: &y, Screen: "DELL U2720Q"}).save(path); err != nil {
		t.Fatalf("save: %v", err)
	}
	got := loadWindowState(path)
	if got.X == nil || *got.X != 40 || got.Y == nil || *got.Y != -10 || got.Screen != "DELL U2720Q" {
		t.Errorf("loaded %+v, want position 40,-10 on DELL U2720Q", got)
	}
}

func TestRestorePlan(t *testing.T) {
	ptr := func(v int) *int { return &v }
	laptop := screenRect{ID: "Built-in", Width: 1440, Height: 875, Primary: true}
	external := screenRect{ID: "DELL", X: 1440, Width: 2560, Height: 1415}
	both := []screenRect{laptop, external}

	cases := []struct {
		name       string
		state      windowState
		wayland    bool
		screens    []screenRect
		w, h, x, y int
		hasPos     bool
	}{
		{
			name:    "saved screen present",
			state:   windowState{Width: 800, Height: 500, X: ptr(1600), Y: ptr(200), Screen: "DELL"},
			screens: both,
			w:       800, h: 500, x: 1600, y: 200, hasPos: true,
		},
		{
			name:    "saved screen missing",
			state:   windowState{Width: 800, Height: 500, X: ptr(1600), Y: ptr(200), Screen: "DELL"},
			screens: []screenRect{laptop},
			w:       800, h: 500,
		},
		{
			name:    "Wayland keeps size only",
			state:   windowState{Width: 800, Height: 500, X: ptr(1600), Y: ptr(200), Screen: "DELL"},
			wayland: true,
			screens: both,
			w:       800, h: 500,
		},
		{
			name:    "no saved position",
			state:   windowState{Width: 800, Height: 500},
			screens: both,
			w:       800, h: 500,
		},
		{
			name:    "size clamped to the primary screen",
			state:   windowState{Width: 3000, Height: 2000},
			screens: both,
			w:       1440, h: 875,
		},
		{
			name:    "size clamped to the saved screen",
			state:   windowState{Width: 3000, Height: 2000, X: ptr(1440), Y: ptr(0), Screen: "DELL"},
			screens: both,
			w:       2560, h: 1415, x: 1440, y: 0, hasPos: true,
		},
		{
			name:    "position pulled back inside the work area",
			state:   windowState{Width: 800, Height: 500, X: ptr(1300), Y: ptr(700), Screen: "Built-in"},
			screens: both,
			w:       800, h: 500, x: 640, y: 375, hasPos: true,
		},
		{
			name:    "position left of the screen",
			state:   windowState{Width: 800, Height: 500, X: ptr(1000), Y: ptr(-50), Screen: "DELL"},
			screens: both,
			w:       800, h: 500, x: 1440, y: 0, hasPos: true,
		},
		{
			name:  "no screen information",
			state: windowState{Width: 3000, Height: 2000, X: ptr(10), Y: ptr(10), Screen: "DELL"},
			w:     3000, h: 2000,
		},
	}
	for _, tc := range cases {
		w, h, x, y, hasPos := restorePlan(tc.state, tc.wayland, tc.screens)
		if w != tc.w || h != tc.h || hasPos != tc.hasPos || (hasPos && (x != tc.x || y != tc.y)) {
			t.Errorf("%s: restorePlan = %dx%d at %d,%d (hasPos %t), want %dx%d at %d,%d (hasPos %t)",
				tc.name, w, h, x, y, hasPos, tc.w, tc.h, tc.x, tc.y, tc.hasPos)
		}
	}
}
