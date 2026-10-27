//go:build !darwin

package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// openAtLogin reports whether an XDG autostart entry starts this binary.
// Wails matches entries by executable, so one written by the installer counts.
func openAtLogin(app *application.App) (bool, error) {
	status, err := app.Autostart.Status()
	if err != nil {
		return false, err
	}
	return status.Enabled, nil
}

// setOpenAtLogin writes or removes ~/.config/autostart/cairn-desktop.desktop.
func setOpenAtLogin(app *application.App, enabled bool) error {
	if !enabled {
		return app.Autostart.Disable()
	}
	return app.Autostart.EnableWithOptions(application.AutostartOptions{
		Identifier: autostartID,
		Arguments:  []string{flagBackground},
	})
}
