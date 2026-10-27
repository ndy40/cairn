package main

import (
	"fmt"
	"os"

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
