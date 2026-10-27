package main

import (
	"errors"
	"fmt"
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

// Errors SyncNow reports instead of syncing.
var (
	errSyncNotSetUp = errors.New("sync is not set up; run `cairn sync setup`")
	errSyncBusy     = errors.New("another sync is running; try again shortly")
)

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

	// stateMu guards the outcome of the last sync, shown in the tray menu.
	stateMu sync.Mutex
	state   syncState

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

// SyncNow pulls, ignoring the throttle, then pushes, and returns once both
// are done. The error is what to report: sync not set up, another sync
// running, or the backend's warning. Used by "Sync Now".
func (y *Syncer) SyncNow() error {
	y.begin()
	defer y.finish()
	if err := y.run("pull", false); err != nil {
		return err
	}
	return y.run("push", false)
}

// State returns the outcome of the last sync, for the tray status line.
func (y *Syncer) State() syncState {
	y.stateMu.Lock()
	defer y.stateMu.Unlock()
	return y.state
}

func (y *Syncer) setState(update func(*syncState)) {
	y.stateMu.Lock()
	defer y.stateMu.Unlock()
	update(&y.state)
}

// start counts the sync as in flight before its goroutine starts, so a Wait
// that follows the call always sees it.
func (y *Syncer) start(kind string, throttle bool) {
	y.begin()
	go func() {
		defer y.finish()
		_ = y.run(kind, throttle)
	}()
}

// begin counts a sync as in flight; finish ends it.
func (y *Syncer) begin() {
	y.inflightMu.Lock()
	defer y.inflightMu.Unlock()
	if y.inflight == 0 {
		y.idle = make(chan struct{})
	}
	y.inflight++
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

// run does one pull or push. A throttled pull is skipped and reports nil;
// otherwise the error says why nothing synced, or carries the backend's warning.
func (y *Syncer) run(kind string, throttle bool) error {
	y.mu.Lock()
	defer y.mu.Unlock()

	if throttle && time.Since(y.lastPull) < minPullInterval {
		return nil
	}

	engine := y.engine()
	if engine == nil {
		y.setState(func(s *syncState) { s.NotSetUp = true })
		return errSyncNotSetUp
	}

	// Share the CLI's lockfile so a desktop sync never overlaps a CLI one.
	lockPath, lockFile, ok := csync.AcquireLock(kind)
	if !ok {
		return errSyncBusy
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
		y.setState(func(s *syncState) { s.NotSetUp, s.Warning = false, warning })
		if y.onWarning != nil {
			y.onWarning(warning)
		}
		return errors.New(warning)
	}
	y.setState(func(s *syncState) { s.NotSetUp, s.Warning, s.LastOK = false, "", time.Now() })
	if y.onOK != nil {
		y.onOK()
	}
	return nil
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

// syncState is the outcome of the last sync.
type syncState struct {
	// NotSetUp is true when the last sync found no backend set up.
	NotSetUp bool
	// LastOK is when a sync last succeeded; zero if none has this run.
	LastOK time.Time
	// Warning is the last sync's problem, "" once a sync succeeds.
	Warning string
}

// syncStatusText is the tray menu's sync status line.
func syncStatusText(s syncState, now time.Time) string {
	switch {
	case s.Warning != "":
		return "Sync problem: " + s.Warning
	case s.NotSetUp:
		return "Sync not set up"
	case !s.LastOK.IsZero():
		return "Synced " + ago(now.Sub(s.LastOK))
	default:
		return "Not synced yet"
	}
}

// ago describes a duration in the past in whole minutes or hours.
func ago(d time.Duration) string {
	switch m := int(d / time.Minute); {
	case m < 1:
		return "just now"
	case m == 1:
		return "1 minute ago"
	case m < 60:
		return fmt.Sprintf("%d minutes ago", m)
	case m < 120:
		return "1 hour ago"
	default:
		return fmt.Sprintf("%d hours ago", m/60)
	}
}
