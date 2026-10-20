package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// lockTTL is the age after which a sync lockfile is considered stale.
const lockTTL = 10 * time.Minute

// AcquireLock creates the lockfile for the given sync kind ("push" or "pull")
// next to the sync config, so that only one sync of each kind runs at a time
// across the CLI, TUI and desktop app. It returns the lock path, the open lock
// file (nil when not acquired) and whether the lock was acquired. Stale locks
// (older than lockTTL or held by a dead process) are cleaned up first.
// The caller removes the lock path once the sync finishes.
func AcquireLock(kind string) (string, *os.File, bool) {
	lockPath := filepath.Join(filepath.Dir(DefaultConfigPath()), fmt.Sprintf("cairn-sync-%s.lock", kind))
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if os.IsExist(err) {
			if cleanupStaleLock(lockPath) {
				lockFile, err = os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
				if err == nil {
					_, _ = fmt.Fprintf(lockFile, "%d\n", os.Getpid())
					return lockPath, lockFile, true
				}
			}
			return lockPath, nil, false
		}
		return lockPath, nil, false
	}
	_, _ = fmt.Fprintf(lockFile, "%d\n", os.Getpid())
	return lockPath, lockFile, true
}

func cleanupStaleLock(lockPath string) bool {
	info, err := os.Stat(lockPath)
	if err != nil {
		return false
	}
	if time.Since(info.ModTime()) > lockTTL {
		_ = os.Remove(lockPath)
		return true
	}
	pid, err := readLockPID(lockPath)
	if err != nil || pid <= 0 {
		_ = os.Remove(lockPath)
		return true
	}
	if !processAlive(pid) {
		_ = os.Remove(lockPath)
		return true
	}
	return false
}

func readLockPID(lockPath string) (int, error) {
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}
