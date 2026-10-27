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
	// onPulled is called after a successful pull so the UI can refresh.
	onPulled func()
	// onWarning reports a non-fatal sync problem (e.g. expired auth).
	onWarning func(msg string)
	// onOK is called after a sync completes without a warning, so the UI
	// can clear a warning shown earlier.
	onOK func()

	// mu serialises syncs in this process; both push and pull rewrite sync.json.
	mu       sync.Mutex
	lastPull time.Time

	// inflightMu guards inflight and idle, which let quitting wait for
	// started syncs to finish.
	inflightMu sync.Mutex
	inflight   int
	// idle is closed when inflight drops to zero; nil before the first sync.
	idle chan struct{}
}

// NewSyncer returns a Syncer for s.
func NewSyncer(s *store.Store, onPulled func(), onWarning func(msg string), onOK func()) *Syncer {
	return &Syncer{store: s, onPulled: onPulled, onWarning: onWarning, onOK: onOK}
}

// Pull starts a background pull unless one ran within minPullInterval.
func (y *Syncer) Pull() {
	y.start("pull", true)
}

// Push starts a background push.
func (y *Syncer) Push() {
	y.start("push", false)
}

// start counts the sync as in flight before its goroutine starts, so a Wait
// that follows the call always sees it.
func (y *Syncer) start(kind string, throttle bool) {
	y.inflightMu.Lock()
	if y.inflight == 0 {
		y.idle = make(chan struct{})
	}
	y.inflight++
	y.inflightMu.Unlock()

	go func() {
		defer y.finish()
		y.run(kind, throttle)
	}()
}

func (y *Syncer) finish() {
	y.inflightMu.Lock()
	defer y.inflightMu.Unlock()
	y.inflight--
	if y.inflight == 0 {
		close(y.idle)
	}
}

// Running reports whether a sync started by Pull or Push is still going.
func (y *Syncer) Running() bool {
	y.inflightMu.Lock()
	defer y.inflightMu.Unlock()
	return y.inflight > 0
}

// Wait blocks until no sync is running or timeout passes, and reports
// whether syncing finished. Quitting uses it so a push isn't cut short.
func (y *Syncer) Wait(timeout time.Duration) bool {
	y.inflightMu.Lock()
	if y.inflight == 0 {
		y.inflightMu.Unlock()
		return true
	}
	idle := y.idle
	y.inflightMu.Unlock()

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-idle:
		return true
	case <-timer.C:
		return false
	}
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
		if warning == "" && y.onPulled != nil {
			y.onPulled()
		}
	} else {
		warning = engine.AutoPush()
	}

	if warning != "" {
		log.Print(warning)
		if y.onWarning != nil {
			y.onWarning(warning)
		}
		return
	}
	if y.onOK != nil {
		y.onOK()
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
