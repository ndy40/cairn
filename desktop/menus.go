package main

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

// helpURL is the Desktop App user guide, opened from Help › Cairn Help.
const helpURL = "https://ndy40.github.io/cairn/docs/desktop-app/"

// Payloads of menuActionEvent: macOS File and Edit menu items that act on
// the selected row, or on the frontend's own state.
const (
	menuOpen = "open"
	menuCopy = "copy"
	menuPin  = "pin"
	menuSync = "sync"
)

// trayStatusInterval is how often the tray's "Synced N minutes ago" line is
// refreshed while nothing else changes it.
const trayStatusInterval = 30 * time.Second

// syncBadge is shown next to the macOS menu bar icon while sync has a problem.
const syncBadge = "⚠︎"

// buildAppMenu returns the macOS menu bar: Cairn, File, Edit, Window and
// Help. Items that act on the list emit menuActionEvent; Quit and Close go
// through the launcher so a running sync can finish. Linux has no menu bar:
// the same actions are in the main menu button (frontend).
func buildAppMenu(app *application.App, l *Launcher) *application.Menu {
	menu := app.NewMenu()
	action := func(name string) func(*application.Context) {
		return func(*application.Context) { app.Event.Emit(menuActionEvent, name) }
	}
	show := func(view string) func(*application.Context) {
		return func(*application.Context) { l.BringForward(view, "") }
	}

	appMenu := menu.AddSubmenu("Cairn")
	appMenu.Add("About Cairn").OnClick(show(viewAbout))
	appMenu.AddSeparator()
	appMenu.Add("Settings…").SetAccelerator("CmdOrCtrl+,").OnClick(show(viewSettings))
	appMenu.AddSeparator()
	appMenu.AddRole(application.ServicesMenu)
	appMenu.AddSeparator()
	appMenu.AddRole(application.Hide)
	appMenu.AddRole(application.HideOthers)
	appMenu.AddRole(application.UnHide)
	appMenu.AddSeparator()
	// Not the Quit role, which calls app.Quit directly and would skip
	// "Finishing sync…".
	appMenu.Add("Quit Cairn").SetAccelerator("CmdOrCtrl+Q").OnClick(func(*application.Context) { l.Quit() })

	file := menu.AddSubmenu("File")
	file.Add("New Bookmark…").SetAccelerator("CmdOrCtrl+N").OnClick(show(viewAdd))
	file.AddSeparator()
	file.Add("Open Link").OnClick(action(menuOpen))
	file.Add("Copy Link").SetAccelerator("Shift+CmdOrCtrl+C").OnClick(action(menuCopy))
	file.Add("Pin/Unpin").OnClick(action(menuPin))
	file.AddSeparator()
	file.Add("Sync Now").OnClick(action(menuSync))
	file.AddSeparator()
	file.Add("Close Window").SetAccelerator("CmdOrCtrl+W").OnClick(func(*application.Context) { l.Close() })

	edit := menu.AddSubmenu("Edit")
	edit.AddRole(application.Undo)
	edit.AddRole(application.Redo)
	edit.AddSeparator()
	edit.AddRole(application.Cut)
	edit.AddRole(application.Copy)
	edit.AddRole(application.Paste)
	edit.AddRole(application.PasteAndMatchStyle)
	edit.AddRole(application.Delete)
	edit.AddRole(application.SelectAll)
	edit.AddSeparator()
	edit.Add("Find Bookmarks").SetAccelerator("CmdOrCtrl+F").OnClick(show(viewSearch))

	menu.AddRole(application.WindowMenu)

	help := menu.AddSubmenu("Help")
	help.Add("Cairn Help").OnClick(func(*application.Context) { _ = app.Browser.OpenURL(helpURL) })
	help.Add("Keyboard Shortcuts").SetAccelerator("CmdOrCtrl+/").OnClick(show(viewShortcuts))

	return menu
}

// trayMenu is the tray (Linux) or menu bar (macOS) icon and its menu.
type trayMenu struct {
	tray   *application.SystemTray
	menu   *application.Menu
	status *application.MenuItem

	// mu serialises refreshes from syncs and the timer.
	mu sync.Mutex
	// statusText and badge are what is shown now, so unchanged refreshes
	// don't rebuild the menu.
	statusText string
	badge      string
}

// buildTray creates the tray icon with Search Bookmarks…, Add Bookmark…, a
// disabled sync status line, Settings… and Quit Cairn. Clicking the icon
// toggles the window. Call it without holding l.mu: once the app runs,
// creating a tray waits for the main thread.
func buildTray(app *application.App, l *Launcher) *trayMenu {
	tray := app.SystemTray.New()
	if runtime.GOOS == "darwin" {
		tray.SetTemplateIcon(icons.SystrayMacTemplate)
	} else {
		tray.SetIcon(trayIcon)
	}
	tray.SetTooltip("Cairn")

	menu := app.NewMenu()
	menu.Add("Search Bookmarks…").OnClick(func(*application.Context) { l.BringForward(viewSearch, "") })
	menu.Add("Add Bookmark…").OnClick(func(*application.Context) { l.BringForward(viewAdd, "") })
	menu.AddSeparator()
	t := &trayMenu{tray: tray, menu: menu}
	t.statusText = syncStatusText(l.syncState(), time.Now())
	t.status = menu.Add(t.statusText).SetEnabled(false)
	menu.AddSeparator()
	menu.Add("Settings…").OnClick(func(*application.Context) { l.BringForward(viewSettings, "") })
	menu.Add("Quit Cairn").OnClick(func(*application.Context) { l.Quit() })

	tray.SetMenu(menu)
	tray.OnClick(func() { l.Toggle("") })
	return t
}

// refresh updates the sync status line and, on macOS, the warning badge.
func (t *trayMenu) refresh(state syncState) {
	t.mu.Lock()
	defer t.mu.Unlock()
	text := syncStatusText(state, time.Now())
	if text != t.statusText {
		t.statusText = text
		// Menu items are native objects; change them on the main thread.
		application.InvokeSync(func() { t.status.SetLabel(text) })
		// Linux builds the D-Bus menu from the labels; set it again to publish.
		t.tray.SetMenu(t.menu)
	}

	badge := ""
	if state.Warning != "" && runtime.GOOS == "darwin" {
		badge = syncBadge
	}
	if badge != t.badge {
		t.badge = badge
		t.tray.SetLabel(badge)
	}
}

// SetTrayAvailable records whether a tray can show Cairn's icon: always on
// macOS, and on Linux while a tray host is running (stock GNOME has none
// without an AppIndicator extension; the hotkey and --toggle work
// regardless).
func (l *Launcher) SetTrayAvailable(available bool) {
	l.mu.Lock()
	l.trayAvailable = available
	l.mu.Unlock()
	l.updateTray()
}

// updateTray adds the tray icon once it is wanted, and on macOS shows or hides
// it per "Show in menu bar". On Linux the icon is never removed (Wails can't
// hide it there): if the tray host goes away nothing shows it, and Wails
// registers it again when a host returns.
func (l *Launcher) updateTray() {
	l.mu.Lock()
	wanted := l.trayAvailable && (runtime.GOOS != "darwin" || l.showInMenuBar == nil || l.showInMenuBar())
	tray := l.tray
	l.mu.Unlock()

	switch {
	case tray == nil && wanted:
		tray = buildTray(l.app, l)
		l.mu.Lock()
		l.tray = tray
		l.mu.Unlock()
	case tray != nil && runtime.GOOS == "darwin":
		if wanted {
			tray.tray.Show()
		} else {
			tray.tray.Hide()
		}
	}
}

// SyncChanged refreshes the tray's sync status after a sync.
func (l *Launcher) SyncChanged() {
	l.mu.Lock()
	tray := l.tray
	l.mu.Unlock()
	if tray != nil {
		tray.refresh(l.syncState())
	}
}

// keepTrayStatusFresh refreshes "Synced N minutes ago" until ctx is done.
func (l *Launcher) keepTrayStatusFresh(ctx context.Context) {
	ticker := time.NewTicker(trayStatusInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				l.SyncChanged()
			}
		}
	}()
}

func (l *Launcher) syncState() syncState {
	if l.syncer == nil {
		return syncState{}
	}
	return l.syncer.State()
}
