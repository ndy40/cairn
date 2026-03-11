# Research: Async Autosync

**Feature**: 007-async-autosync
**Date**: 2026-03-11

## Research Task 1: Background execution strategy for CLI tools

### Decision: Detached subprocess via `os/exec`

### Rationale

A CLI tool is a short-lived process — when `main()` returns, all goroutines are killed. The options considered:

1. **Goroutine (fire-and-forget)**: Cannot work. The goroutine dies when the parent process exits.
2. **Goroutine with WaitGroup/timeout**: Still blocks the user for the timeout duration. Defeats the purpose.
3. **Detached subprocess**: Spawn `cairn sync push` as a background child process using `os/exec.Command` + `cmd.Start()` (without `cmd.Wait()`). The child process is automatically reparented to init/systemd when the parent exits and continues running independently.

Option 3 is the standard Unix pattern for "fire and forget" work from short-lived CLI tools. It requires no new dependencies — only Go stdlib `os/exec` and `syscall`.

### Alternatives Considered

- **Daemon/service**: Overkill for a single-user bookmark tool. Violates the "single binary, no external runtime" constitution principle.
- **Lock file + cron**: Complex, requires user configuration, unreliable timing.
- **Named pipe / socket**: Requires a long-running listener process — same problem as daemon.

## Research Task 2: Preventing concurrent background sync

### Decision: Rely on existing `autoSyncPush` behavior + OS-level process semantics

### Rationale

The existing `autoSyncPush` function opens its own DB connection, loads config, and runs push. If two background processes run simultaneously, SQLite WAL mode handles concurrent reads safely. The Dropbox upload is idempotent (overwrites the same remote file). The only risk is a race where both processes export, then both upload slightly different snapshots — but since pending changes are cleared after upload, the worst case is a redundant upload with the same data.

For extra safety, an optional file lock (`flock`) on a sync lock file could be used, but given the single-user nature and the idempotency of the operations, this is an optimization that can be deferred.

### Alternatives Considered

- **File lock (flock)**: Adds complexity for minimal benefit in a single-user tool. Could be added later if issues arise.
- **PID file check**: Fragile, stale PID files cause problems.

## Research Task 3: Suppressing background process output

### Decision: Redirect stdout and stderr to `/dev/null` (or `os.DevNull`)

### Rationale

The background subprocess must not write to the parent's terminal after the parent has returned the prompt. Using `cmd.Stdout = nil` and `cmd.Stderr = nil` in Go's `os/exec` already suppresses output (they default to the parent's file descriptors only if explicitly set). Setting them to `os.DevNull` via `os.Open(os.DevNull)` is the explicit, safe approach.

## Research Task 4: Self-invocation pattern

### Decision: Use `os.Executable()` to get the path to the current binary

### Rationale

The background process needs to run `cairn sync push`. Using `os.Executable()` returns the absolute path to the currently running binary, avoiding PATH resolution issues. The command becomes: `exec.Command(self, "sync", "push")`.

This ensures the exact same binary version is used, avoiding version mismatch issues.

## Research Task 5: TUI autosync behavior

### Decision: No changes needed for TUI

### Rationale

The TUI (`internal/model/app.go`) does NOT call `autoSyncPush` or `autoSyncPull`. Bookmark operations in the TUI (delete, update tags) only interact with the store and reload the bookmark list. Auto-sync in the TUI would require a different pattern (tea.Cmd) and is out of scope for this feature.
