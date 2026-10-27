package main

import (
	"embed"
	"log"
	"os"
	"runtime"
	"slices"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"

	"github.com/ndy40/cairn/internal/config"
	"github.com/ndy40/cairn/internal/store"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cfgManager := config.NewManager()
	if err := cfgManager.Load("", ""); err != nil {
		log.Fatalf("load configuration: %v", err)
	}
	cfg := cfgManager.Get()

	s, err := store.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer func() { _ = s.Close() }()

	var (
		app      *application.App
		launcher *Launcher
	)
	svc := NewBookmarkService(s,
		func(url string) error { return app.Browser.OpenURL(url) },
		func() string {
			text, _ := app.Clipboard.Text()
			return text
		},
	)

	platform := detectPlatform(runtime.GOOS, os.Getenv)
	windows := NewWindowService(platform, defaultSettingsPath())

	app = application.New(application.Options{
		Name:        "Cairn",
		Description: "Quick launcher for Cairn bookmarks",
		Services: []application.Service{
			application.NewService(svc),
			application.NewService(windows),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		// Later launches (a desktop shortcut bound to `cairn-desktop --toggle`,
		// or running it again) are forwarded here instead of starting a new app.
		// Each launch sends its own activation token, which the running
		// instance needs to raise its window on Wayland.
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "io.github.ndy40.cairn.desktop",
			AdditionalData: map[string]string{
				dataActivationToken: os.Getenv("XDG_ACTIVATION_TOKEN"),
				dataStartupID:       os.Getenv("DESKTOP_STARTUP_ID"),
			},
			OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
				if launcher != nil {
					launcher.HandleArgs(data.Args, activationToken(data.AdditionalData))
				}
			},
		},
		OnShutdown: func() {
			if launcher != nil {
				launcher.SaveState()
			}
		},
		Mac: application.MacOptions{
			// Run as a menu-bar app: no Dock icon, stays alive with no windows.
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	syncer := NewSyncer(s,
		func() { app.Event.Emit(changedEvent, "pull") },
		func(msg string) { app.Event.Emit(syncWarningEvent, msg) },
	)
	svc.afterWrite = syncer.Push

	launcher = NewLauncher(app, defaultStatePath(), platform)
	launcher.onShow = syncer.Pull
	launcher.background = windows.backgroundMode
	windows.launcher = launcher
	launcher.SetupTray()
	launcher.RegisterHotkey(cfg.DesktopHotkey)

	if debugEnabled {
		// Shows whether an existing login item (e.g. one written by an older
		// installer) is detected; see research R9.
		app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
			status, err := app.Autostart.Status()
			debugf("autostart status: %+v (err: %v)", status, err)
		})
	}

	args := os.Args[1:]
	switch {
	case slices.Contains(args, flagBackground), slices.Contains(args, flagHidden):
		// Started at login or restarted by the installer: no window, but
		// catch up with remote changes. --hidden is the old name, still
		// used by login items created by earlier installers.
		app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
			syncer.Pull()
		})
	case slices.Contains(args, flagAdd):
		app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
			launcher.BringForward(viewAdd, "")
		})
	default:
		app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
			launcher.BringForward(viewSearch, "")
		})
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
