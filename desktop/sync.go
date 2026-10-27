package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/ndy40/cairn/internal/store"
	csync "github.com/ndy40/cairn/internal/sync"
)

// minPullInterval throttles pulls triggered by showing the popup, which can
// happen many times a minute.
const minPullInterval = 2 * time.Minute

// Syncer runs background sync against the configured backend, mirroring the
// CLI's auto-sync: pull on startup and when the popup is shown, push after
// writes. It is a no-op until sync is configured with `cairn sync setup`.
type Syncer struct {
	store *store.Store
	// onSynced is called after a successful pull or push so the UI can refresh;
	// a push merges the cloud snapshot first, so it can bring in bookmarks too.
	onSynced func()
	// onWarning reports a non-fatal sync problem (e.g. expired auth).
	onWarning func(msg string)

	// mu serialises syncs in this process; both push and pull rewrite sync.json.
	mu       sync.Mutex
	lastPull time.Time
}

// NewSyncer returns a Syncer for s.
func NewSyncer(s *store.Store, onSynced func(), onWarning func(msg string)) *Syncer {
	return &Syncer{store: s, onSynced: onSynced, onWarning: onWarning}
}

// Pull starts a background pull unless one ran within minPullInterval.
func (y *Syncer) Pull() {
	go y.run("pull", true)
}

// Push starts a background push.
func (y *Syncer) Push() {
	go y.run("push", false)
}

func (y *Syncer) run(kind string, throttle bool) {
	y.mu.Lock()
	defer y.mu.Unlock()

	if throttle && time.Since(y.lastPull) < minPullInterval {
		return
	}

	engine := y.engine()
	if engine == nil {
		return
	}

	// Share the CLI's lockfile so a desktop sync never overlaps a CLI one.
	lockPath, lockFile, ok := csync.AcquireLock(kind)
	if !ok {
		return
	}
	defer func() {
		_ = lockFile.Close()
		_ = os.Remove(lockPath)
	}()

	var warning string
	if kind == "pull" {
		y.lastPull = time.Now()
		_, warning = engine.AutoPull()
	} else {
		warning = engine.AutoPush()
	}

	if warning == "" && y.onSynced != nil {
		y.onSynced()
	}

	if warning != "" {
		log.Print(warning)
		if y.onWarning != nil {
			y.onWarning(warning)
		}
	}
}

// engine loads the sync config fresh each time, so running `cairn sync setup`
// or `cairn sync unlink` takes effect without restarting the app.
func (y *Syncer) engine() *csync.Engine {
	cfgPath := csync.DefaultConfigPath()
	cfg, err := csync.LoadConfig(cfgPath)
	if err != nil || !csync.IsConfigured(cfg) {
		return nil
	}
	b, err := csync.NewBackend(cfg)
	if err != nil {
		log.Printf("sync backend unavailable: %v", err)
		return nil
	}
	return csync.NewEngine(y.store, b, cfg, cfgPath)
}
