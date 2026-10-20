package main

import (
	_ "embed"
	"log"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

// defaultHotkey toggles the popup when desktop_hotkey is not configured.
const defaultHotkey = "CmdOrCtrl+Shift+Space"

// saveSizeDelay waits for a resize to settle before saving the window size.
const saveSizeDelay = 500 * time.Millisecond

// Events sent to the frontend.
const (
	// showEvent tells the frontend the popup was shown and which view to open.
	showEvent = "launcher:show"
	// changedEvent tells the frontend bookmarks changed outside it (e.g. a sync pull).
	changedEvent = "bookmarks:changed"
	// syncWarningEvent carries a non-fatal sync message to display.
	syncWarningEvent = "sync:warning"
)

// Views the popup can open in; sent as the showEvent payload.
const (
	viewSearch = "search"
	viewAdd    = "add"
)

// Command-line flags. A second launch forwards its flags to the running instance.
const (
	flagToggle = "--toggle" // show/hide the popup (the default for a second launch)
	flagAdd    = "--add"    // show the popup on the add form
	flagHidden = "--hidden" // first launch only: start in the tray without showing the popup
)

//go:embed build/appicon.png
var trayIcon []byte

func init() {
	application.RegisterEvent[string](showEvent)
	application.RegisterEvent[string](changedEvent)
	application.RegisterEvent[string](syncWarningEvent)
}

// Launcher owns the popup window and the ways of summoning it.
type Launcher struct {
	app    *application.App
	window *application.WebviewWindow

	mu sync.Mutex
	// ready is set once the frontend runtime has loaded; until then the view
	// requested by Show is held in pendingView so the event is not lost.
	ready       bool
	pendingView string

	// onShow, when set, is called each time the popup is shown (used to pull sync).
	onShow func()

	// statePath is where the window size is remembered between runs.
	statePath string
	saveTimer *time.Timer

	// focus decides whether losing focus hides the popup.
	focus *focusGuard
}

// NewLauncher creates the hidden popup window at the size saved in statePath.
// Closing it hides it instead of quitting.
func NewLauncher(app *application.App, statePath string, focus *focusGuard) *Launcher {
	state := loadWindowState(statePath)
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "launcher",
		Title:  "Cairn",
		Width:  state.Width,
		Height: state.Height,
		// Resizable from its edges; the size is remembered between runs.
		MinWidth:        minWidth,
		MinHeight:       minHeight,
		URL:             "/",
		Frameless:       true,
		AlwaysOnTop:     true,
		Hidden:          true,
		InitialPosition: application.WindowCentered,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
	})

	l := &Launcher{app: app, window: window, statePath: statePath, focus: focus}

	// Hide when focus moves to another app. Wails' HideOnFocusLost is not used
	// because on GNOME/Wayland it hides the popup on any click inside it.
	window.OnWindowEvent(events.Common.WindowLostFocus, func(*application.WindowEvent) {
		debugf("focus lost")
		time.AfterFunc(focusCheckDelay, l.hideIfUnfocused)
	})
	window.OnWindowEvent(events.Common.WindowFocus, func(*application.WindowEvent) {
		debugf("focus gained")
	})

	window.OnWindowEvent(events.Common.WindowDidResize, func(*application.WindowEvent) {
		l.scheduleSaveSize()
	})

	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		window.Hide()
		e.Cancel()
	})

	window.OnWindowEvent(events.Common.WindowRuntimeReady, func(*application.WindowEvent) {
		l.mu.Lock()
		l.ready = true
		view := l.pendingView
		l.pendingView = ""
		l.mu.Unlock()
		if view != "" {
			l.app.Event.Emit(showEvent, view)
		}
	})

	return l
}

// Show centres and focuses the popup on the given view.
func (l *Launcher) Show(view string) {
	l.window.Center()
	l.window.Show()
	l.window.Focus()
	if l.onShow != nil {
		l.onShow()
	}

	l.mu.Lock()
	ready := l.ready
	if !ready {
		l.pendingView = view
	}
	l.mu.Unlock()
	if ready {
		l.app.Event.Emit(showEvent, view)
	}
}

// hideIfUnfocused hides the popup if focus really moved to another app.
func (l *Launcher) hideIfUnfocused() {
	visible, focused := l.window.IsVisible(), l.window.IsFocused()
	hide, reason := l.focus.shouldHide(visible, focused)
	debugf("focus check: visible=%t focused=%t hide=%t (%s)", visible, focused, hide, reason)
	if hide {
		l.window.Hide()
	}
}

// scheduleSaveSize saves the window size once resizing has paused, so a drag
// doesn't write the state file on every intermediate size.
func (l *Launcher) scheduleSaveSize() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.saveTimer != nil {
		l.saveTimer.Stop()
	}
	l.saveTimer = time.AfterFunc(saveSizeDelay, l.SaveSize)
}

// SaveSize writes the current window size to the state file. Sizes reported
// while the window is hidden (often 0x0) are ignored.
func (l *Launcher) SaveSize() {
	width, height := l.window.Size()
	if err := (windowState{Width: width, Height: height}).save(l.statePath); err != nil {
		log.Printf("save window size: %v", err)
	}
}

// Toggle hides the popup if it is visible, otherwise shows the search view.
func (l *Launcher) Toggle() {
	if l.window.IsVisible() {
		l.window.Hide()
		return
	}
	l.Show(viewSearch)
}

// HandleArgs applies flags passed to a second launch of the app.
func (l *Launcher) HandleArgs(args []string) {
	if slices.Contains(args, flagAdd) {
		l.Show(viewAdd)
		return
	}
	l.Toggle()
}

// SetupTray adds a tray icon whose click toggles the popup, with a menu for the
// common actions. Some desktops (e.g. stock GNOME) need an AppIndicator
// extension to display tray icons; the hotkey and --toggle work regardless.
func (l *Launcher) SetupTray() {
	tray := l.app.SystemTray.New()
	if runtime.GOOS == "darwin" {
		tray.SetTemplateIcon(icons.SystrayMacTemplate)
	} else {
		tray.SetIcon(trayIcon)
	}
	tray.SetTooltip("Cairn")

	menu := l.app.NewMenu()
	menu.Add("Search bookmarks").OnClick(func(*application.Context) { l.Show(viewSearch) })
	menu.Add("Add bookmark").OnClick(func(*application.Context) { l.Show(viewAdd) })
	menu.AddSeparator()
	menu.Add("Quit Cairn").OnClick(func(*application.Context) { l.app.Quit() })

	tray.SetMenu(menu)
	tray.OnClick(l.Toggle)
}

// RegisterHotkey binds the global shortcut that toggles the popup. hotkey is
// the configured value: empty uses defaultHotkey and "off" disables it.
// Failure is logged rather than fatal so the tray and --toggle still work.
func (l *Launcher) RegisterHotkey(hotkey string) {
	hotkey = strings.TrimSpace(hotkey)
	if strings.EqualFold(hotkey, "off") {
		return
	}
	if hotkey == "" {
		hotkey = defaultHotkey
	}
	if err := l.app.GlobalShortcut.Register(hotkey, l.Toggle); err != nil {
		log.Printf("global hotkey %q unavailable: %v (use `cairn-desktop --toggle` instead)", hotkey, err)
	}
}
