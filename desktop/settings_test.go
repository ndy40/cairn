package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSettingsRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "desktop-settings.json")
	on := true
	want := DesktopSettings{ShowInMenuBar: false, RunInBackground: &on, HideAfterOpen: true}

	if err := want.save(path); err != nil {
		t.Fatalf("save: %v", err)
	}
	got := loadSettings(path)
	if got.ShowInMenuBar || !got.HideAfterOpen || got.RunInBackground == nil || !*got.RunInBackground {
		t.Errorf("loaded %+v, want %+v", got, want)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Error("temporary file left behind after save")
	}
}

func TestLoadSettingsDefaults(t *testing.T) {
	dir := t.TempDir()
	cases := map[string]string{
		"missing":     "",
		"corrupt":     "{not json",
		"empty":       "{}",
		"unknown key": `{"somethingElse": 3}`,
	}
	for name, body := range cases {
		path := filepath.Join(dir, name+".json")
		if body != "" {
			if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
				t.Fatalf("write: %v", err)
			}
		}
		got := loadSettings(path)
		if !got.ShowInMenuBar || got.RunInBackground != nil || got.HideAfterOpen {
			t.Errorf("%s: loaded %+v, want defaults", name, got)
		}
	}
}

func TestLoadSettingsKeepsDefaultsForMissingKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "desktop-settings.json")
	if err := os.WriteFile(path, []byte(`{"hideAfterOpen": true}`), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	got := loadSettings(path)
	if !got.HideAfterOpen || !got.ShowInMenuBar {
		t.Errorf("loaded %+v, want hideAfterOpen on and showInMenuBar at its default", got)
	}
}

func TestLoadSettingsLeavesInvalidFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "desktop-settings.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	loadSettings(path)
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "{not json" {
		t.Errorf("invalid settings file was changed: %q, %v", data, err)
	}
}

func TestResolvedRunInBackground(t *testing.T) {
	on, off := true, false
	cases := []struct {
		name    string
		value   *bool
		desktop string
		want    bool
	}{
		{"unset on KDE", nil, desktopKDE, true},
		{"unset on GNOME", nil, desktopGNOME, false},
		{"unset elsewhere", nil, desktopOther, false},
		{"off on KDE", &off, desktopKDE, false},
		{"on on GNOME", &on, desktopGNOME, true},
	}
	for _, tc := range cases {
		got := DesktopSettings{RunInBackground: tc.value}.resolvedRunInBackground(tc.desktop)
		if got != tc.want {
			t.Errorf("%s: resolvedRunInBackground = %t, want %t", tc.name, got, tc.want)
		}
	}
}

func TestDefaultSettingsPathIsBesideConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("APPDATA", dir)

	got := defaultSettingsPath()
	if filepath.Base(got) != "desktop-settings.json" || filepath.Base(filepath.Dir(got)) != "cairn" {
		t.Errorf("defaultSettingsPath() = %q, want .../cairn/desktop-settings.json", got)
	}
}
