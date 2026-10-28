package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// openAtLogin reports whether Cairn's LaunchAgent exists.
func openAtLogin(*application.App) (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	return launchAgentExists(launchAgentPath(home))
}

// setOpenAtLogin writes Cairn's LaunchAgent, or removes it.
func setOpenAtLogin(_ *application.App, enabled bool) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := launchAgentPath(home)
	if !enabled {
		return removeLaunchAgent(path)
	}
	exe, err := resolvedExecutable()
	if err != nil {
		return fmt.Errorf("find executable: %w", err)
	}
	return writeLaunchAgent(path, exe)
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
