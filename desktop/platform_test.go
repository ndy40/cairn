package main

import "testing"

func TestDetectPlatform(t *testing.T) {
	cases := []struct {
		name string
		goos string
		env  map[string]string
		want PlatformInfo
	}{
		{
			name: "macOS",
			goos: "darwin",
			env:  map[string]string{"XDG_CURRENT_DESKTOP": "GNOME"},
			want: PlatformInfo{OS: "darwin", TrayAvailable: true},
		},
		{
			name: "GNOME on Wayland",
			goos: "linux",
			env:  map[string]string{"XDG_CURRENT_DESKTOP": "GNOME", "XDG_SESSION_TYPE": "wayland"},
			want: PlatformInfo{OS: "linux", Desktop: desktopGNOME, Wayland: true},
		},
		{
			name: "Ubuntu GNOME on X11",
			goos: "linux",
			env:  map[string]string{"XDG_CURRENT_DESKTOP": "ubuntu:GNOME", "XDG_SESSION_TYPE": "x11"},
			want: PlatformInfo{OS: "linux", Desktop: desktopGNOME},
		},
		{
			name: "KDE, lower case",
			goos: "linux",
			env:  map[string]string{"XDG_CURRENT_DESKTOP": "kde", "XDG_SESSION_TYPE": "Wayland"},
			want: PlatformInfo{OS: "linux", Desktop: desktopKDE, Wayland: true},
		},
		{
			name: "other desktop",
			goos: "linux",
			env:  map[string]string{"XDG_CURRENT_DESKTOP": "Hyprland", "XDG_SESSION_TYPE": "wayland"},
			want: PlatformInfo{OS: "linux", Desktop: desktopOther, Wayland: true},
		},
		{
			name: "no session variables",
			goos: "linux",
			want: PlatformInfo{OS: "linux", Desktop: desktopOther},
		},
	}
	for _, tc := range cases {
		got := detectPlatform(tc.goos, func(k string) string { return tc.env[k] })
		if got != tc.want {
			t.Errorf("%s: detectPlatform = %+v, want %+v", tc.name, got, tc.want)
		}
	}
}

func TestTrayOwnerChange(t *testing.T) {
	const signal = "org.freedesktop.DBus.NameOwnerChanged"
	cases := []struct {
		name          string
		signal        string
		body          []any
		available, ok bool
	}{
		{"tray host appears", signal, []any{trayWatcherName, "", ":1.42"}, true, true},
		{"tray host goes away", signal, []any{trayWatcherName, ":1.42", ""}, false, true},
		{"other name", signal, []any{"org.example.Other", "", ":1.7"}, false, false},
		{"other signal", "org.freedesktop.Notifications.ActionInvoked", []any{uint32(1), "default"}, false, false},
		{"short body", signal, []any{trayWatcherName}, false, false},
		{"wrong type", signal, []any{trayWatcherName, "", 42}, false, false},
	}
	for _, tc := range cases {
		available, ok := trayOwnerChange(tc.signal, tc.body)
		if available != tc.available || ok != tc.ok {
			t.Errorf("%s: trayOwnerChange = (%t, %t), want (%t, %t)", tc.name, available, ok, tc.available, tc.ok)
		}
	}
}
