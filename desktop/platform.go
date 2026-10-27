package main

import "strings"

// Desktop environments the frontend styles for; sent as PlatformInfo.Desktop.
const (
	desktopGNOME = "gnome"
	desktopKDE   = "kde"
	desktopOther = "other"
)

// PlatformInfo tells the frontend which platform variant to draw.
type PlatformInfo struct {
	// OS is "darwin" or "linux".
	OS string `json:"os"`
	// Desktop is "gnome", "kde" or "other" on Linux, and "" on macOS.
	Desktop string `json:"desktop"`
	// Wayland is true in a Wayland session, where apps can't place their windows.
	Wayland bool `json:"wayland"`
	// TrayAvailable reports whether a tray or menu bar icon can be shown.
	TrayAvailable bool `json:"trayAvailable"`
	// AccentColor is the system accent colour as #rrggbb, or "" when unknown.
	AccentColor string `json:"accentColor,omitempty"`
}

// detectPlatform derives PlatformInfo from the OS and session environment.
// getenv is os.Getenv outside tests.
func detectPlatform(goos string, getenv func(string) string) PlatformInfo {
	info := PlatformInfo{OS: goos}
	if goos == "darwin" {
		info.TrayAvailable = true
		return info
	}

	info.Wayland = strings.EqualFold(getenv("XDG_SESSION_TYPE"), "wayland")
	info.Desktop = desktopOther
	// XDG_CURRENT_DESKTOP is a colon-separated list, e.g. "ubuntu:GNOME".
	for _, name := range strings.Split(getenv("XDG_CURRENT_DESKTOP"), ":") {
		switch strings.ToLower(strings.TrimSpace(name)) {
		case "gnome":
			info.Desktop = desktopGNOME
			return info
		case "kde":
			info.Desktop = desktopKDE
			return info
		}
	}
	return info
}
