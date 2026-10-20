package main

import (
	"embed"
	"log"
	"os"
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

	focus := newFocusGuard()

	app = application.New(application.Options{
		Name:        "Cairn",
		Description: "Quick launcher for Cairn bookmarks",
		Services: []application.Service{
			application.NewService(svc),
			application.NewService(&WindowService{guard: focus}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		// Later launches (a desktop shortcut bound to `cairn-desktop --toggle`,
		// or running it again) are forwarded here instead of starting a new app.
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "io.github.ndy40.cairn.desktop",
			OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
				if launcher != nil {
					launcher.HandleArgs(data.Args)
				}
			},
		},
		OnShutdown: func() {
			if launcher != nil {
				launcher.SaveSize()
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

	launcher = NewLauncher(app, defaultStatePath(), focus)
	launcher.onShow = syncer.Pull
	launcher.SetupTray()
	launcher.RegisterHotkey(cfg.DesktopHotkey)

	args := os.Args[1:]
	switch {
	case slices.Contains(args, flagHidden):
		// Autostart: wait in the tray, but catch up with remote changes.
		app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
			syncer.Pull()
		})
	case slices.Contains(args, flagAdd):
		app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
			launcher.Show(viewAdd)
		})
	default:
		app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
			launcher.Show(viewSearch)
		})
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
