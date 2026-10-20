package main

import (
	"testing"
)

func TestSyncerIsNoOpWhenNotConfigured(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)

	called := false
	y := NewSyncer(openTestStore(t),
		func() { called = true },
		func(msg string) { t.Errorf("unexpected warning: %s", msg) },
	)

	// Call run directly (Pull/Push only add a goroutine) so the test is synchronous.
	y.run("pull", true)
	y.run("push", false)

	if called {
		t.Error("onPulled called although sync is not configured")
	}
	if !y.lastPull.IsZero() {
		t.Error("lastPull set although no pull ran")
	}
}
