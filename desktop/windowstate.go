package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/ndy40/cairn/internal/config"
)

// Default and minimum window sizes.
const (
	defaultWidth  = 720
	defaultHeight = 520
	minWidth      = 520
	minHeight     = 360
)

// windowState is the window state remembered between runs. It lives in its own
// file rather than cairn.json so the app never rewrites the user's config.
type windowState struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	// X and Y are the window position, saved only where apps can place their
	// own windows (macOS and X11; not Wayland).
	X *int `json:"x,omitempty"`
	Y *int `json:"y,omitempty"`
	// Screen identifies the display the window was on; see screenRect.ID.
	Screen string `json:"screen,omitempty"`
}

// screenRect is a display's work area, as needed to restore the window.
type screenRect struct {
	// ID is the display's name, or its index when it has none.
	ID                  string
	X, Y, Width, Height int
	Primary             bool
}

// defaultStatePath returns desktop-state.json next to cairn.json.
func defaultStatePath() string {
	cfgPath := config.DefaultConfigPath()
	if cfgPath == "" {
		return ""
	}
	return filepath.Join(filepath.Dir(cfgPath), "desktop-state.json")
}

// loadWindowState reads the saved state, falling back to the default size when
// the file is missing, unreadable or holds a size below the minimum.
func loadWindowState(path string) windowState {
	state := windowState{Width: defaultWidth, Height: defaultHeight}
	if path == "" {
		return state
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return state
	}
	var saved windowState
	if json.Unmarshal(data, &saved) != nil || !saved.valid() {
		return state
	}
	return saved
}

// restorePlan decides the size and position to open the window at. The size
// is clamped to the work area of the saved screen, or of the primary screen
// when that one is gone. A position is only returned when the platform lets
// apps place windows (not Wayland) and the saved screen is still connected;
// otherwise hasPos is false and the caller centres the window or leaves it to
// the compositor.
func restorePlan(state windowState, wayland bool, screens []screenRect) (w, h, x, y int, hasPos bool) {
	w, h = state.Width, state.Height

	var saved, target *screenRect
	for i := range screens {
		s := &screens[i]
		if state.Screen != "" && s.ID == state.Screen && saved == nil {
			saved = s
		}
		if s.Primary && target == nil {
			target = s
		}
	}
	if saved != nil {
		target = saved
	}
	if target == nil && len(screens) > 0 {
		target = &screens[0]
	}
	if target == nil {
		return w, h, 0, 0, false
	}

	w, h = min(w, target.Width), min(h, target.Height)
	if wayland || saved == nil || state.X == nil || state.Y == nil {
		return w, h, 0, 0, false
	}
	// Keep the whole window inside the work area, preferring its top-left
	// corner when it is larger than the screen.
	x = max(min(*state.X, target.X+target.Width-w), target.X)
	y = max(min(*state.Y, target.Y+target.Height-h), target.Y)
	return w, h, x, y, true
}

// save writes the state atomically so a crash mid-write can't corrupt it.
func (s windowState) save(path string) error {
	if path == "" {
		return errors.New("no state path")
	}
	if !s.valid() {
		return nil
	}
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// valid reports whether the size is usable. Windows report 0x0 (or tiny sizes)
// while hidden on some platforms, which must not overwrite a real size.
func (s windowState) valid() bool {
	return s.Width >= minWidth && s.Height >= minHeight
}
