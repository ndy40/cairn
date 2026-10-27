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
