package main

import (
	_ "embed"
	"log"
	"os"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

// defaultHotkey toggles the window when desktop_hotkey is not configured.
const defaultHotkey = "CmdOrCtrl+Shift+Space"

// toggleCommand is the fallback users can bind to a system shortcut when the
// global shortcut can't be registered.
const toggleCommand = "cairn-desktop --toggle"

// saveStateDelay waits for a resize or move to settle before saving the window state.
const saveStateDelay = 500 * time.Millisecond

// Events sent to the frontend.
const (
	// showEvent tells the frontend the window came forward and which view to open.
	showEvent = "launcher:show"
	// changedEvent tells the frontend bookmarks changed outside it (e.g. a sync pull).
	changedEvent = "bookmarks:changed"
	// syncWarningEvent carries a non-fatal sync message to display.
	syncWarningEvent = "sync:warning"
)

// Views the window can open in; sent as the showEvent payload.
const (
	viewSearch = "search"
	viewAdd    = "add"
)

// Command-line flags. A second launch forwards its flags to the running instance.
const (
	flagToggle     = "--toggle"     // bring the window forward, or put it away if it is in front
	flagAdd        = "--add"        // bring the window forward on the add form
	flagBackground = "--background" // first launch only: start without showing the window (used at login)
	flagHidden     = "--hidden"     // deprecated alias for --background, kept for existing login items
)

//go:embed build/appicon.png
var trayIcon []byte

func init() {
	application.RegisterEvent[string](showEvent)
	application.RegisterEvent[string](changedEvent)
	application.RegisterEvent[string](syncWarningEvent)
}

// debugEnabled turns on window and lifecycle logging with CAIRN_DESKTOP_DEBUG=1.
var debugEnabled = os.Getenv("CAIRN_DESKTOP_DEBUG") != ""

func debugf(format string, args ...any) {
	if debugEnabled {
		log.Printf("debug: "+format, args...)
	}
}

// ShortcutStatus reports whether the global shortcut is working, so Settings
// can explain a failure and offer the --toggle fallback.
type ShortcutStatus struct {
	// Hotkey is the effective accelerator, or "off" when disabled.
	Hotkey string `json:"hotkey"`
	// Registered is true when the shortcut was registered successfully.
	Registered bool `json:"registered"`
	// Error says why registration failed (taken, refused, unsupported).
	Error string `json:"error"`
	// ToggleCommand is the command to bind to a system shortcut instead.
	ToggleCommand string `json:"toggleCommand"`
}

// Launcher owns the main window and the ways of summoning it.
type Launcher struct {
	app      *application.App
	window   *application.WebviewWindow
	platform PlatformInfo

	mu sync.Mutex
	// ready is set once the frontend runtime has loaded; until then the view
	// requested by BringForward is held in pendingView so the event is not lost.
	ready       bool
	pendingView string
	// placed is set once the saved size and position have been applied, which
	// happens on the first show when the screens are known.
	placed   bool
	shortcut ShortcutStatus

	// onShow, when set, is called each time the window comes forward (used to pull sync).
	onShow func()
	// background, when set, reports whether Linux "Run in background" is in
	// effect (on, with a tray to return to).
	background func() bool

	// statePath is where the window size and position are remembered between runs.
	statePath string
	state     windowState
	saveTimer *time.Timer
}

// NewLauncher creates the hidden main window. It is a normal framed window
// with the platform's own controls; its saved size and position are applied
// when it is first shown. Closing it hides it instead of quitting.
func NewLauncher(app *application.App, statePath string, platform PlatformInfo) *Launcher {
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
		Hidden:          true,
		InitialPosition: application.WindowCentered,
		// The search field sits in a transparent title band next to the
		// traffic lights; only that band shows the translucent material.
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBarHiddenInsetUnified,
			Backdrop: application.MacBackdropTranslucent,
		},
		// No Linux menu: GTK draws its own title bar and Cairn's main menu is
		// a button in the search bar, per the GNOME HIG.
	})

	l := &Launcher{app: app, window: window, platform: platform, statePath: statePath, state: state}

	window.OnWindowEvent(events.Common.WindowDidResize, func(*application.WindowEvent) {
		l.scheduleSaveState()
	})
	window.OnWindowEvent(events.Common.WindowDidMove, func(*application.WindowEvent) {
		l.scheduleSaveState()
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

// BringForward shows, restores and focuses the window on the given view.
// token is a window activation token forwarded by a second launch, needed on
// Wayland to be allowed to raise the window; "" when there is none.
func (l *Launcher) BringForward(view, token string) {
	l.mu.Lock()
	first := !l.placed
	l.placed = true
	state := l.state
	l.mu.Unlock()

	var x, y int
	var hasPos bool
	if first {
		var w, h int
		w, h, x, y, hasPos = restorePlan(state, l.platform.Wayland, l.screens())
		l.window.SetSize(w, h)
	}

	applyActivationToken(l.window, token)
	if runtime.GOOS == "darwin" {
		// Undo an app-wide hide from the shortcut.
		l.app.Show()
	}
	if l.window.IsMinimised() {
		l.window.Restore()
	}
	l.window.Show()
	if first {
		// Position after showing: GTK can only move a window that has a surface.
		if hasPos {
			l.window.SetPosition(x, y)
		} else {
			l.window.Center()
		}
	}
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

// Toggle brings the window forward, or puts it away when it is already in
// front. Used by the global shortcut, the tray icon and --toggle.
func (l *Launcher) Toggle(token string) {
	visible, focused := l.window.IsVisible(), l.window.IsFocused()
	act := shortcutAction(runtime.GOOS, visible, focused, l.backgroundMode())
	debugf("toggle: visible=%t focused=%t action=%d", visible, focused, act)
	if act == actionShow {
		l.BringForward(viewSearch, token)
		return
	}
	l.apply(act)
}

// PutAway gets Cairn out of the way the platform's way (used after opening a
// link when "Hide Cairn after opening a link" is on).
func (l *Launcher) PutAway() {
	l.apply(putAwayAction(runtime.GOOS, l.backgroundMode()))
}

func (l *Launcher) apply(act action) {
	switch act {
	case actionShow:
		l.BringForward(viewSearch, "")
	case actionHideApp:
		l.app.Hide()
	case actionMinimise:
		l.window.Minimise()
	case actionHideWindow:
		l.window.Hide()
	}
}

func (l *Launcher) backgroundMode() bool {
	return l.background != nil && l.background()
}

// screens returns the connected displays' work areas.
func (l *Launcher) screens() []screenRect {
	var rects []screenRect
	for _, s := range l.app.Screen.GetAll() {
		rects = append(rects, screenRect{
			ID:      screenID(s),
			X:       s.WorkArea.X,
			Y:       s.WorkArea.Y,
			Width:   s.WorkArea.Width,
			Height:  s.WorkArea.Height,
			Primary: s.IsPrimary,
		})
	}
	return rects
}

// screenID identifies a display across runs. The name is preferred because
// on Linux the ID is only the monitor's index.
func screenID(s *application.Screen) string {
	if s.Name != "" {
		return s.Name
	}
	return s.ID
}

// scheduleSaveState saves the window state once resizing or moving has
// paused, so a drag doesn't write the state file on every step.
func (l *Launcher) scheduleSaveState() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.saveTimer != nil {
		l.saveTimer.Stop()
	}
	l.saveTimer = time.AfterFunc(saveStateDelay, l.SaveState)
}

// SaveState writes the current window size, and its position and screen where
// the platform allows placing windows, to the state file. Sizes reported while
// the window is hidden (often 0x0) are ignored, and the position is only
// updated while the window is visible.
func (l *Launcher) SaveState() {
	l.mu.Lock()
	placed, state := l.placed, l.state
	l.mu.Unlock()
	if !placed {
		// Never shown (e.g. started in the tray): keep the saved state as is.
		return
	}

	state.Width, state.Height = l.window.Size()
	if !l.platform.Wayland && l.window.IsVisible() {
		x, y := l.window.Position()
		state.X, state.Y = &x, &y
		if s, err := l.window.GetScreen(); err == nil && s != nil {
			state.Screen = screenID(s)
		}
	}
	if !state.valid() {
		return
	}
	if err := state.save(l.statePath); err != nil {
		log.Printf("save window state: %v", err)
		return
	}
	l.mu.Lock()
	l.state = state
	l.mu.Unlock()
}

// HandleArgs applies flags passed to a second launch of the app. token is the
// window activation token that launch forwarded.
func (l *Launcher) HandleArgs(args []string, token string) {
	switch {
	case slices.Contains(args, flagAdd):
		l.BringForward(viewAdd, token)
	case slices.Contains(args, flagToggle):
		l.Toggle(token)
	default:
		l.BringForward(viewSearch, token)
	}
}

// SetupTray adds a tray icon whose click toggles the window, with a menu for
// the common actions. Some desktops (e.g. stock GNOME) need an AppIndicator
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
	menu.Add("Search bookmarks").OnClick(func(*application.Context) { l.BringForward(viewSearch, "") })
	menu.Add("Add bookmark").OnClick(func(*application.Context) { l.BringForward(viewAdd, "") })
	menu.AddSeparator()
	menu.Add("Quit Cairn").OnClick(func(*application.Context) { l.app.Quit() })

	tray.SetMenu(menu)
	tray.OnClick(func() { l.Toggle("") })
}

// RegisterHotkey binds the global shortcut that toggles the window. hotkey is
// the configured value: empty uses defaultHotkey and "off" disables it.
// Failure is recorded for Settings rather than fatal, so the tray and
// --toggle still work.
func (l *Launcher) RegisterHotkey(hotkey string) {
	status := ShortcutStatus{ToggleCommand: toggleCommand}
	hotkey = strings.TrimSpace(hotkey)
	switch {
	case strings.EqualFold(hotkey, "off"):
		status.Hotkey = "off"
	default:
		if hotkey == "" {
			hotkey = defaultHotkey
		}
		status.Hotkey = hotkey
		if err := l.app.GlobalShortcut.Register(hotkey, func() { l.Toggle("") }); err != nil {
			status.Error = err.Error()
			log.Printf("global hotkey %q unavailable: %v (use `%s` instead)", hotkey, err, toggleCommand)
		} else {
			status.Registered = true
		}
	}

	l.mu.Lock()
	l.shortcut = status
	l.mu.Unlock()
}

// ShortcutStatus returns the outcome of RegisterHotkey.
func (l *Launcher) ShortcutStatus() ShortcutStatus {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.shortcut
}
