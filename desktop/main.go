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
				launcher.Shutdown()
			}
		},
		Mac: application.MacOptions{
			// A regular app with a Dock icon and menu bar. Closing the window
			// hides it, so the app stays alive with no windows.
			ActivationPolicy: application.ActivationPolicyRegular,
		},
	})

	syncer := NewSyncer(s,
		func() { app.Event.Emit(changedEvent, "pull") },
		func(msg string) {
			app.Event.Emit(syncWarningEvent, msg)
			launcher.SyncChanged()
		},
		func() {
			app.Event.Emit(syncOKEvent)
			launcher.SyncChanged()
		},
	)
	svc.afterWrite = syncer.Push
	svc.syncNow = syncer.SyncNow
	svc.writeClipboard = func(text string) bool { return app.Clipboard.SetText(text) }
	svc.changed = func(reason string) { app.Event.Emit(changedEvent, reason) }

	launcher = NewLauncher(app, defaultStatePath(), platform)
	launcher.onShow = syncer.Pull
	launcher.background = windows.backgroundMode
	launcher.closing = windows.closeAction
	launcher.syncer = syncer
	launcher.showInMenuBar = windows.showInMenuBar
	windows.app = app
	windows.launcher = launcher
	if platform.OS == "darwin" {
		// Linux gets no menu bar (GNOME HIG); its main menu is a button.
		app.Menu.Set(buildAppMenu(app, launcher))
	}
	if platform.OS == "linux" {
		// Show the tray icon, and allow running in the background, only
		// while a tray host can display it.
		watchTray(app.Context(), func(available bool) {
			info, changed := windows.setTrayAvailable(available)
			launcher.SetTrayAvailable(available)
			if changed {
				app.Event.Emit(platformChangedEvent, info)
			}
		})
		// The selection follows the desktop's accent colour (macOS uses
		// CSS AccentColor instead).
		watchAccentColor(app.Context(), func(color string) {
			if info, changed := windows.setAccentColor(color); changed {
				app.Event.Emit(platformChangedEvent, info)
			}
		})
	} else {
		launcher.SetTrayAvailable(true)
	}
	launcher.RegisterHotkey(cfg.DesktopHotkey)

	// Clicking the Dock icon with no window open brings the window back.
	app.Event.OnApplicationEvent(events.Mac.ApplicationShouldHandleReopen, func(e *application.ApplicationEvent) {
		if !e.Context().HasVisibleWindows() {
			launcher.BringForward(viewSearch, "")
		}
	})

	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		launcher.keepTrayStatusFresh(app.Context())
		if debugEnabled {
			// Shows whether an existing login item (e.g. one written by an
			// older installer) is detected; see research R9.
			enabled, err := openAtLogin(app)
			debugf("open at login: %t (err: %v)", enabled, err)
		}
	})

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
