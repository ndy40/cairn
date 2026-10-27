package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// launchAgentLabel names Cairn's macOS LaunchAgent. It is the label the
// installer used, so a login item it created reads as on after an upgrade.
const launchAgentLabel = "io.github.ndy40.cairn.desktop"

// autostartID names the Linux XDG autostart entry (cairn-desktop.desktop),
// the file name the installer used.
const autostartID = "cairn-desktop"

// launchAgentPlist returns the LaunchAgent that starts executable in the
// background at login.
func launchAgentPlist(executable string) []byte {
	var b bytes.Buffer
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>`)
	_ = xml.EscapeText(&b, []byte(launchAgentLabel))
	b.WriteString(`</string>
	<key>ProgramArguments</key>
	<array>
		<string>`)
	_ = xml.EscapeText(&b, []byte(executable))
	b.WriteString(`</string>
		<string>` + flagBackground + `</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
</dict>
</plist>
`)
	return b.Bytes()
}

// resolvedExecutable is the running binary's path with symlinks resolved, so
// a login item keeps working if a symlink to it changes.
func resolvedExecutable() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		return resolved, nil
	}
	return exe, nil
}

// On macOS Cairn writes its own LaunchAgent rather than using app.Autostart,
// which picks SMAppService for a bundled app on macOS 13+ and can't pass
// --background (research R9). The file handling lives here, not in the darwin
// file, so it is built and tested on every platform.

// launchAgentPath is ~/Library/LaunchAgents/<label>.plist under home.
func launchAgentPath(home string) string {
	return filepath.Join(home, "Library", "LaunchAgents", launchAgentLabel+".plist")
}

// launchAgentExists reports whether the LaunchAgent at path exists.
func launchAgentExists(path string) (bool, error) {
	_, err := os.Stat(path)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	default:
		return false, err
	}
}

// writeLaunchAgent atomically writes the LaunchAgent for executable to path.
func writeLaunchAgent(path, executable string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create login item: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, launchAgentPlist(executable), 0o644); err != nil {
		return fmt.Errorf("write login item: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("write login item: %w", err)
	}
	return nil
}

// removeLaunchAgent deletes the LaunchAgent at path; a missing one is fine.
func removeLaunchAgent(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove login item: %w", err)
	}
	return nil
}
