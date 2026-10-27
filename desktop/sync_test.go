package main

import (
	"testing"
	"time"
)

func TestSyncerIsNoOpWhenNotConfigured(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)

	called := false
	y := NewSyncer(openTestStore(t),
		func() { called = true },
		func(msg string) { t.Errorf("unexpected warning: %s", msg) },
		func() { t.Error("onOK called although sync is not configured") },
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

func TestSyncerWait(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)

	y := NewSyncer(openTestStore(t), nil, nil, nil)
	if y.Running() {
		t.Error("Running before any sync")
	}
	if !y.Wait(time.Millisecond) {
		t.Error("Wait with nothing running = false, want true")
	}

	// Hold the sync lock so the started push stays in flight.
	y.mu.Lock()
	y.Push()
	if !y.Running() {
		t.Error("Running = false while a push is in flight")
	}
	if y.Wait(20 * time.Millisecond) {
		t.Error("Wait = true while a push is in flight")
	}
	y.mu.Unlock()

	if !y.Wait(5 * time.Second) {
		t.Fatal("Wait timed out after the push could finish")
	}
	if y.Running() {
		t.Error("Running after the push finished")
	}

	// A second round reuses the counter after it dropped to zero.
	y.Pull()
	if !y.Wait(5 * time.Second) {
		t.Error("Wait timed out on a second sync")
	}
}
