package main

import (
	"fmt"
	"sync"
)

// Setting keys accepted by WindowService.SetSetting.
const (
	settingShowInMenuBar   = "showInMenuBar"
	settingRunInBackground = "runInBackground"
	settingHideAfterOpen   = "hideAfterOpen"
)

// WindowService exposes window, platform and settings calls to the frontend.
type WindowService struct {
	// launcher is set once the window exists; nil before app.Run.
	launcher *Launcher
	platform PlatformInfo

	settingsPath string
	mu           sync.Mutex
	settings     DesktopSettings
}

// NewWindowService loads the saved settings from settingsPath.
func NewWindowService(platform PlatformInfo, settingsPath string) *WindowService {
	return &WindowService{
		platform:     platform,
		settingsPath: settingsPath,
		settings:     loadSettings(settingsPath),
	}
}

// Platform tells the frontend which platform variant to draw.
func (w *WindowService) Platform() PlatformInfo {
	return w.platform
}

// Settings returns the effective settings, with the per-desktop default for
// runInBackground resolved.
func (w *WindowService) Settings() DesktopSettings {
	w.mu.Lock()
	settings := w.settings
	w.mu.Unlock()
	background := settings.resolvedRunInBackground(w.platform.Desktop)
	settings.RunInBackground = &background
	return settings
}

// SetSetting changes one setting and saves it.
func (w *WindowService) SetSetting(key string, value bool) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	next := w.settings
	switch key {
	case settingShowInMenuBar:
		next.ShowInMenuBar = value
	case settingRunInBackground:
		next.RunInBackground = &value
	case settingHideAfterOpen:
		next.HideAfterOpen = value
	default:
		return fmt.Errorf("unknown setting %q", key)
	}
	if err := next.save(w.settingsPath); err != nil {
		return fmt.Errorf("save settings: %w", err)
	}
	w.settings = next
	return nil
}

// ShortcutStatus reports whether the global shortcut was registered.
func (w *WindowService) ShortcutStatus() ShortcutStatus {
	if w.launcher == nil {
		return ShortcutStatus{ToggleCommand: toggleCommand}
	}
	return w.launcher.ShortcutStatus()
}

// AfterOpen is called by the frontend after a link opens. It puts Cairn away
// only when "Hide Cairn after opening a link" is on.
func (w *WindowService) AfterOpen() {
	w.mu.Lock()
	hide := w.settings.HideAfterOpen
	w.mu.Unlock()
	if hide && w.launcher != nil {
		w.launcher.PutAway()
	}
}

// backgroundMode reports whether Linux "Run in background" is in effect: on,
// and with a tray to return to.
func (w *WindowService) backgroundMode() bool {
	if w.platform.OS != "linux" || !w.platform.TrayAvailable {
		return false
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.settings.resolvedRunInBackground(w.platform.Desktop)
}
