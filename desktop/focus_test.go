package main

import (
	"testing"
	"time"
)

// newTestGuard returns a guard with a controllable clock.
func newTestGuard(t *testing.T) (*focusGuard, *time.Time) {
	t.Helper()
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	g := &focusGuard{now: func() time.Time { return now }}
	return g, &now
}

func TestShouldHideWhenFocusMovesElsewhere(t *testing.T) {
	g, _ := newTestGuard(t)
	if hide, reason := g.shouldHide(true, false); !hide {
		t.Errorf("shouldHide = false (%s), want true when focus moved to another app", reason)
	}
}

func TestShouldNotHideAfterPressInsidePopup(t *testing.T) {
	g, now := newTestGuard(t)
	g.pointerDown()

	// Focus loss right after a click/drag/resize inside the popup.
	*now = now.Add(focusCheckDelay)
	if hide, _ := g.shouldHide(true, false); hide {
		t.Error("shouldHide = true right after a press inside the popup, want false")
	}

	// Long after the press, a focus loss means the user went elsewhere.
	*now = now.Add(pointerGrace)
	if hide, _ := g.shouldHide(true, false); !hide {
		t.Error("shouldHide = false long after the last press, want true")
	}
}

func TestShouldNotHideWhenFocusRegainedOrHidden(t *testing.T) {
	g, _ := newTestGuard(t)
	if hide, _ := g.shouldHide(true, true); hide {
		t.Error("shouldHide = true although the window has focus again")
	}
	if hide, _ := g.shouldHide(false, false); hide {
		t.Error("shouldHide = true although the window is already hidden")
	}
}

func TestShouldNotHideOnFocusFollowsMouse(t *testing.T) {
	g, _ := newTestGuard(t)
	g.focusFollowsMouse = true
	if hide, _ := g.shouldHide(true, false); hide {
		t.Error("shouldHide = true on a focus-follows-mouse window manager")
	}
}
