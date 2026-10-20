package main

import (
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	// focusCheckDelay is how long to wait after a focus loss before deciding
	// whether to hide, so a brief loss-and-regain does not hide the popup.
	focusCheckDelay = 250 * time.Millisecond

	// pointerGrace is how long after a mouse press inside the popup a focus
	// loss is attributed to that press. On GNOME/Wayland pressing inside the
	// window (click, drag, edge resize) starts a compositor grab that GTK
	// reports as losing focus, which would otherwise hide the popup.
	pointerGrace = 600 * time.Millisecond
)

// focusGuard decides whether a focus loss should hide the popup.
type focusGuard struct {
	mu                sync.Mutex
	lastPointerDown   time.Time
	now               func() time.Time
	focusFollowsMouse bool
}

func newFocusGuard() *focusGuard {
	return &focusGuard{now: time.Now, focusFollowsMouse: focusFollowsMouse()}
}

// pointerDown records a mouse press inside the popup.
func (g *focusGuard) pointerDown() {
	g.mu.Lock()
	g.lastPointerDown = g.now()
	g.mu.Unlock()
}

// shouldHide reports whether the popup should hide after losing focus, given
// its state once focusCheckDelay has passed, and why (for debug logging).
func (g *focusGuard) shouldHide(visible, focused bool) (bool, string) {
	g.mu.Lock()
	sincePress := g.now().Sub(g.lastPointerDown)
	g.mu.Unlock()

	switch {
	case g.focusFollowsMouse:
		return false, "focus-follows-mouse window manager"
	case !visible:
		return false, "already hidden"
	case focused:
		return false, "focus regained"
	case sincePress < pointerGrace:
		return false, "caused by a press inside the popup"
	default:
		return true, "focus moved elsewhere"
	}
}

// focusFollowsMouse reports whether the Linux window manager moves focus with
// the pointer (Hyprland, Sway, i3). There, merely moving the mouse off the
// popup loses focus, so hiding on focus loss would make it unusable.
func focusFollowsMouse() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	for _, env := range []string{"HYPRLAND_INSTANCE_SIGNATURE", "SWAYSOCK", "I3SOCK"} {
		if os.Getenv(env) != "" {
			return true
		}
	}
	return false
}

// debugEnabled turns on focus/visibility logging with CAIRN_DESKTOP_DEBUG=1.
var debugEnabled = os.Getenv("CAIRN_DESKTOP_DEBUG") != ""

func debugf(format string, args ...any) {
	if debugEnabled {
		log.Printf("debug: "+format, args...)
	}
}

// WindowService exposes popup window hooks to the frontend.
type WindowService struct {
	guard *focusGuard
}

// PointerDown is called by the frontend on every mouse press in the popup.
func (w *WindowService) PointerDown() {
	debugf("pointer down")
	w.guard.pointerDown()
}
