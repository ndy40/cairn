package main

import "testing"

func TestShortcutAction(t *testing.T) {
	cases := []struct {
		name                         string
		goos                         string
		visible, focused, background bool
		want                         action
	}{
		{"hidden on macOS", "darwin", false, false, false, actionShow},
		{"behind another app on macOS", "darwin", true, false, false, actionShow},
		{"in front on macOS", "darwin", true, true, false, actionHideApp},
		{"hidden on Linux", "linux", false, false, false, actionShow},
		{"minimised or behind on Linux", "linux", true, false, false, actionShow},
		{"in front on Linux", "linux", true, true, false, actionMinimise},
		{"in front on Linux in background mode", "linux", true, true, true, actionHideWindow},
	}
	for _, tc := range cases {
		if got := shortcutAction(tc.goos, tc.visible, tc.focused, tc.background); got != tc.want {
			t.Errorf("%s: shortcutAction = %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestCloseAction(t *testing.T) {
	cases := []struct {
		name                           string
		goos                           string
		runInBackground, trayAvailable bool
		want                           action
	}{
		{"macOS", "darwin", false, false, actionHideWindow},
		{"macOS with background set", "darwin", true, true, actionHideWindow},
		{"Linux with background and tray", "linux", true, true, actionHideWindow},
		{"Linux with background, no tray", "linux", true, false, actionQuit},
		{"Linux with tray, background off", "linux", false, true, actionQuit},
		{"Linux with neither", "linux", false, false, actionQuit},
	}
	for _, tc := range cases {
		if got := closeAction(tc.goos, tc.runInBackground, tc.trayAvailable); got != tc.want {
			t.Errorf("%s: closeAction = %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestActivationToken(t *testing.T) {
	cases := []struct {
		name string
		data map[string]string
		want string
	}{
		{"none", nil, ""},
		{"Wayland", map[string]string{dataActivationToken: "wl-1"}, "wl-1"},
		{"X11", map[string]string{dataStartupID: "x11-1"}, "x11-1"},
		{"both prefer Wayland", map[string]string{dataActivationToken: "wl-1", dataStartupID: "x11-1"}, "wl-1"},
		{"empty Wayland falls back", map[string]string{dataActivationToken: "", dataStartupID: "x11-1"}, "x11-1"},
	}
	for _, tc := range cases {
		if got := activationToken(tc.data); got != tc.want {
			t.Errorf("%s: activationToken = %q, want %q", tc.name, got, tc.want)
		}
	}
}
