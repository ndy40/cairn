package main

import (
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
