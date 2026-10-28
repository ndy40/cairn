package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/ndy40/cairn/internal/config"
)

// DesktopSettings are the user's choices from the Settings panel. Like the
// window state they live in their own file so the app never rewrites
// cairn.json. Open at login is not stored here: the OS login item is the
// source of truth.
type DesktopSettings struct {
	// ShowInMenuBar shows the menu bar icon (macOS).
	ShowInMenuBar bool `json:"showInMenuBar"`
	// RunInBackground keeps the app in the tray when its window closes (Linux).
	// nil means the per-desktop default; see resolvedRunInBackground.
	RunInBackground *bool `json:"runInBackground"`
	// HideAfterOpen puts Cairn away after a link is opened.
	HideAfterOpen bool `json:"hideAfterOpen"`
}

// defaultSettings returns the settings used when none are saved.
func defaultSettings() DesktopSettings {
	return DesktopSettings{ShowInMenuBar: true}
}

// defaultSettingsPath returns desktop-settings.json next to cairn.json.
func defaultSettingsPath() string {
	cfgPath := config.DefaultConfigPath()
	if cfgPath == "" {
		return ""
	}
	return filepath.Join(filepath.Dir(cfgPath), "desktop-settings.json")
}

// loadSettings reads the saved settings. A missing file gives the defaults; an
// invalid one is logged and also gives the defaults, and is left in place
// until the user changes a setting.
func loadSettings(path string) DesktopSettings {
	settings := defaultSettings()
	if path == "" {
		return settings
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Printf("read desktop settings: %v", err)
		}
		return settings
	}
	// Unmarshal over the defaults so keys missing from the file keep them.
	if err := json.Unmarshal(data, &settings); err != nil {
		log.Printf("desktop settings %s are invalid, using defaults: %v", path, err)
		return defaultSettings()
	}
	return settings
}

// save writes the settings atomically so a crash mid-write can't corrupt them.
func (s DesktopSettings) save(path string) error {
	if path == "" {
		return errors.New("no settings path")
	}
	data, err := json.MarshalIndent(s, "", "  ")
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

// resolvedRunInBackground returns the effective RunInBackground choice. When
// the user hasn't chosen, KDE (which always has a tray) defaults to on.
func (s DesktopSettings) resolvedRunInBackground(desktop string) bool {
	if s.RunInBackground != nil {
		return *s.RunInBackground
	}
	return desktop == desktopKDE
}
