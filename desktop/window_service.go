package main

import (
	"errors"
	"fmt"
	"sync"
)

// Setting keys accepted by WindowService.SetSetting.
const (
	settingShowInMenuBar   = "showInMenuBar"
	settingRunInBackground = "runInBackground"
	settingHideAfterOpen   = "hideAfterOpen"
)

// errNoTray is returned when "Run in background" is turned on with no tray
// to return to.
var errNoTray = errors.New("run in background needs a system tray, and none is available")

// WindowService exposes window, platform and settings calls to the frontend.
type WindowService struct {
	// launcher is set once the window exists; nil before app.Run.
	launcher *Launcher

	settingsPath string
	// mu guards platform, which changes when a tray host comes or goes, and
	// settings.
	mu       sync.Mutex
	platform PlatformInfo
	settings DesktopSettings
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
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.platform
}

// setTrayAvailable records whether a tray host is running and returns the
// updated PlatformInfo, and whether it changed.
func (w *WindowService) setTrayAvailable(available bool) (PlatformInfo, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	changed := w.platform.TrayAvailable != available
	w.platform.TrayAvailable = available
	return w.platform, changed
}

// Settings returns the effective settings, with the per-desktop default for
// runInBackground resolved.
func (w *WindowService) Settings() DesktopSettings {
	w.mu.Lock()
	settings, desktop := w.settings, w.platform.Desktop
	w.mu.Unlock()
	background := settings.resolvedRunInBackground(desktop)
	settings.RunInBackground = &background
	return settings
}

// SetSetting changes one setting and saves it. Turning on runInBackground
// fails on Linux while no tray is available.
func (w *WindowService) SetSetting(key string, value bool) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	next := w.settings
	switch key {
	case settingShowInMenuBar:
		next.ShowInMenuBar = value
	case settingRunInBackground:
		if value && !w.platform.TrayAvailable {
			return errNoTray
		}
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

// CloseWindow does what the window's close button does (Cmd/Ctrl+W): hide
// the window, or quit where closing quits.
func (w *WindowService) CloseWindow() {
	if w.launcher != nil {
		w.launcher.Close()
	}
}

// Quit ends the app, first showing "Finishing sync…" and waiting briefly if
// a sync is running.
func (w *WindowService) Quit() {
	if w.launcher != nil {
		w.launcher.Quit()
	}
}

// closeAction decides what closing the window does, from the platform, the
// "Run in background" setting and whether a tray is available now.
func (w *WindowService) closeAction() action {
	w.mu.Lock()
	defer w.mu.Unlock()
	return closeAction(w.platform.OS, w.settings.resolvedRunInBackground(w.platform.Desktop), w.platform.TrayAvailable)
}

// backgroundMode reports whether Linux "Run in background" is in effect: on,
// and with a tray to return to.
func (w *WindowService) backgroundMode() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.platform.OS != "linux" || !w.platform.TrayAvailable {
		return false
	}
	return w.settings.resolvedRunInBackground(w.platform.Desktop)
}
