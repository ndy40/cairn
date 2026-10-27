package main

import (
	"context"
	"fmt"
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

// The XDG Settings portal, which publishes the desktop's accent colour.
const (
	portalDest       = "org.freedesktop.portal.Desktop"
	portalPath       = "/org/freedesktop/portal/desktop"
	portalSettings   = "org.freedesktop.portal.Settings"
	appearanceNS     = "org.freedesktop.appearance"
	accentColorKey   = "accent-color"
	settingChangedFn = "SettingChanged"
)

// watchAccentColor reports the desktop's accent colour (#rrggbb, or "" when
// the desktop sets none). It calls onChange once with the current colour
// before returning, then from a goroutine whenever it changes, until ctx is
// done. Call it on Linux only; macOS uses CSS AccentColor.
func watchAccentColor(ctx context.Context, onChange func(color string)) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		onChange("")
		return
	}

	if err := conn.AddMatchSignal(
		dbus.WithMatchObjectPath(portalPath),
		dbus.WithMatchInterface(portalSettings),
		dbus.WithMatchMember(settingChangedFn),
		dbus.WithMatchArg(0, appearanceNS),
		dbus.WithMatchArg(1, accentColorKey),
	); err != nil {
		debugf("accent colour: watch: %v", err)
	}
	signals := make(chan *dbus.Signal, 10)
	conn.Signal(signals)

	color := readAccentColor(conn)
	debugf("accent colour: %q", color)
	onChange(color)

	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()
	go func() {
		for sig := range signals {
			// Body: namespace, key, value.
			if sig.Name != portalSettings+"."+settingChangedFn || len(sig.Body) < 3 {
				continue
			}
			if ns, _ := sig.Body[0].(string); ns != appearanceNS {
				continue
			}
			if key, _ := sig.Body[1].(string); key != accentColorKey {
				continue
			}
			now, _ := accentHex(sig.Body[2])
			if now == color {
				continue
			}
			color = now
			debugf("accent colour: %q", color)
			onChange(color)
		}
	}()
}

// readAccentColor asks the Settings portal for the accent colour, using
// ReadOne (portal version 2) and falling back to the older Read.
func readAccentColor(conn *dbus.Conn) string {
	portal := conn.Object(portalDest, portalPath)
	var value dbus.Variant
	err := portal.Call(portalSettings+".ReadOne", 0, appearanceNS, accentColorKey).Store(&value)
	if err != nil {
		err = portal.Call(portalSettings+".Read", 0, appearanceNS, accentColorKey).Store(&value)
	}
	if err != nil {
		debugf("accent colour: %v", err)
		return ""
	}
	color, _ := accentHex(value)
	return color
}

// accentHex converts the portal's accent colour, an (r, g, b) tuple of
// doubles in 0–1 (possibly wrapped in variants), to #rrggbb. ok is false when
// the value is missing or out of range, which is how desktops say "no accent".
func accentHex(v any) (color string, ok bool) {
	for {
		variant, isVariant := v.(dbus.Variant)
		if !isVariant {
			break
		}
		v = variant.Value()
	}
	rgb, isSlice := v.([]any)
	if !isSlice || len(rgb) != 3 {
		return "", false
	}
	var out [3]int
	for i, c := range rgb {
		f, isFloat := c.(float64)
		if !isFloat || f < 0 || f > 1 {
			return "", false
		}
		out[i] = int(f*255 + 0.5)
	}
	return fmt.Sprintf("#%02x%02x%02x", out[0], out[1], out[2]), true
}
