# CLI Behavior Contract: Async Autosync

**Feature**: 007-async-autosync
**Date**: 2026-03-11

## Changed Commands

### `cairn add <url>`
- **Before**: Blocks until sync push completes (2-5s with network)
- **After**: Returns immediately after local save; sync runs in background
- **Output**: Unchanged (bookmark saved confirmation)
- **Exit code**: Unchanged (0 on success, non-zero on local save failure)

### `cairn delete`
- **Before**: Blocks until sync push completes
- **After**: Returns immediately after local delete; sync runs in background
- **Output**: Unchanged
- **Exit code**: Unchanged

### `cairn edit` (tag updates)
- **Before**: Blocks until sync push completes
- **After**: Returns immediately after local edit; sync runs in background
- **Output**: Unchanged
- **Exit code**: Unchanged

## Unchanged Commands

### `cairn sync push`
- Remains synchronous and interactive
- Shows errors to the user

### `cairn sync pull`
- Remains synchronous and interactive

### CLI startup auto-pull
- Remains synchronous (must complete before displaying data)

## Background Process Behavior

- Spawned as a detached child process
- No stdout/stderr output to user's terminal
- Failures are silent; pending changes remain queued
- Process continues running after parent CLI exits
