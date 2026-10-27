package main

import (
	"context"
	"log"
	"strings"

	"github.com/godbus/dbus/v5"
)

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

// trayWatcherName is the D-Bus name a tray host (KDE, or GNOME with an
// AppIndicator extension) owns while it can show tray icons.
const trayWatcherName = "org.kde.StatusNotifierWatcher"

// watchTray reports whether a tray host is available on Linux. It calls
// onChange once with the current availability before returning, then from a
// goroutine each time a tray host appears or disappears, until ctx is done.
// Without a session bus there is no tray. Call it on Linux only.
func watchTray(ctx context.Context, onChange func(available bool)) {
	// A private connection, so closing it can't affect the one Wails' tray uses.
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Printf("tray detection: no session bus: %v", err)
		onChange(false)
		return
	}

	// Subscribe before asking, so a change in between is not missed.
	if err := conn.AddMatchSignal(
		dbus.WithMatchObjectPath("/org/freedesktop/DBus"),
		dbus.WithMatchInterface("org.freedesktop.DBus"),
		dbus.WithMatchMember("NameOwnerChanged"),
		dbus.WithMatchArg(0, trayWatcherName),
	); err != nil {
		log.Printf("tray detection: watch %s: %v", trayWatcherName, err)
	}
	signals := make(chan *dbus.Signal, 10)
	conn.Signal(signals)

	var available bool
	if err := conn.BusObject().Call("org.freedesktop.DBus.NameHasOwner", 0, trayWatcherName).Store(&available); err != nil {
		log.Printf("tray detection: %v", err)
	}
	debugf("tray available: %t", available)
	onChange(available)

	go func() {
		<-ctx.Done()
		// Closing the connection also closes signals, ending the loop below.
		_ = conn.Close()
	}()
	go func() {
		for sig := range signals {
			now, ok := trayOwnerChange(sig.Name, sig.Body)
			if !ok || now == available {
				continue
			}
			available = now
			debugf("tray available: %t", available)
			onChange(available)
		}
	}()
}

// trayOwnerChange reads a NameOwnerChanged signal (body: name, old owner, new
// owner) and reports whether the tray host now has an owner. ok is false for
// any other signal.
func trayOwnerChange(name string, body []any) (available, ok bool) {
	if name != "org.freedesktop.DBus.NameOwnerChanged" || len(body) < 3 {
		return false, false
	}
	if owned, _ := body[0].(string); owned != trayWatcherName {
		return false, false
	}
	newOwner, isString := body[2].(string)
	if !isString {
		return false, false
	}
	return newOwner != "", true
}
