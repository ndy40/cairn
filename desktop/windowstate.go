package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/ndy40/cairn/internal/config"
)

// Default and minimum popup sizes.
const (
	defaultWidth  = 680
	defaultHeight = 440
	minWidth      = 480
	minHeight     = 300
)

// windowState is the popup state remembered between runs. It lives in its own
// file rather than cairn.json so the app never rewrites the user's config.
type windowState struct {
	Width  int `json:"width"`
	Height int `json:"height"`
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
