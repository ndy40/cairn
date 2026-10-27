package main

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestWindowServiceSetSetting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "desktop-settings.json")
	w := NewWindowService(PlatformInfo{OS: "linux", Desktop: desktopGNOME}, path)

	if err := w.SetSetting(settingHideAfterOpen, true); err != nil {
		t.Fatalf("SetSetting: %v", err)
	}
	if err := w.SetSetting("nonsense", true); err == nil {
		t.Error("SetSetting accepted an unknown key")
	}
	if !loadSettings(path).HideAfterOpen {
		t.Error("hideAfterOpen was not saved")
	}
	if got := NewWindowService(w.platform, path).Settings(); !got.HideAfterOpen {
		t.Errorf("reloaded settings %+v, want hideAfterOpen on", got)
	}
}

func TestWindowServiceSettingsResolvesBackground(t *testing.T) {
	dir := t.TempDir()
	kde := NewWindowService(PlatformInfo{OS: "linux", Desktop: desktopKDE}, filepath.Join(dir, "kde.json"))
	if got := kde.Settings().RunInBackground; got == nil || !*got {
		t.Errorf("KDE runInBackground = %v, want on by default", got)
	}
	gnome := NewWindowService(PlatformInfo{OS: "linux", Desktop: desktopGNOME}, filepath.Join(dir, "gnome.json"))
	if got := gnome.Settings().RunInBackground; got == nil || *got {
		t.Errorf("GNOME runInBackground = %v, want off by default", got)
	}
}

func TestWindowServiceBackgroundModeNeedsTray(t *testing.T) {
	path := filepath.Join(t.TempDir(), "desktop-settings.json")
	noTray := NewWindowService(PlatformInfo{OS: "linux", Desktop: desktopKDE}, path)
	if noTray.backgroundMode() {
		t.Error("backgroundMode on without a tray")
	}
	withTray := NewWindowService(PlatformInfo{OS: "linux", Desktop: desktopKDE, TrayAvailable: true}, path)
	if !withTray.backgroundMode() {
		t.Error("backgroundMode off on KDE with a tray")
	}
	mac := NewWindowService(PlatformInfo{OS: "darwin", TrayAvailable: true}, path)
	if mac.backgroundMode() {
		t.Error("backgroundMode on macOS")
	}
}

func TestWindowServiceRunInBackgroundNeedsTray(t *testing.T) {
	path := filepath.Join(t.TempDir(), "desktop-settings.json")
	w := NewWindowService(PlatformInfo{OS: "linux", Desktop: desktopGNOME}, path)

	if err := w.SetSetting(settingRunInBackground, true); !errors.Is(err, errNoTray) {
		t.Errorf("SetSetting(runInBackground, true) without a tray = %v, want errNoTray", err)
	}
	if err := w.SetSetting(settingRunInBackground, false); err != nil {
		t.Errorf("SetSetting(runInBackground, false) without a tray: %v", err)
	}

	if _, changed := w.setTrayAvailable(true); !changed {
		t.Error("setTrayAvailable(true) reported no change")
	}
	if _, changed := w.setTrayAvailable(true); changed {
		t.Error("setTrayAvailable(true) twice reported a change")
	}
	if err := w.SetSetting(settingRunInBackground, true); err != nil {
		t.Fatalf("SetSetting(runInBackground, true) with a tray: %v", err)
	}
	if got := w.closeAction(); got != actionHideWindow {
		t.Errorf("closeAction with background and tray = %d, want hide", got)
	}

	// The tray host going away makes closing quit again.
	if info, _ := w.setTrayAvailable(false); info.TrayAvailable {
		t.Error("Platform still reports a tray")
	}
	if got := w.closeAction(); got != actionQuit {
		t.Errorf("closeAction after the tray went away = %d, want quit", got)
	}
}
