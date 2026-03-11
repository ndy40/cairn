# Quickstart: Async Autosync

**Feature**: 007-async-autosync
**Date**: 2026-03-11

## What Changed

Bookmark add, delete, and edit commands no longer wait for cloud sync to complete before returning. The sync happens in a background process.

## How to Test

### Prerequisites
- Sync must be configured: `cairn sync setup`
- A valid Dropbox connection with `CAIRN_DROPBOX_APP_KEY` set

### Manual Testing

1. **Verify non-blocking add**:
   ```bash
   time cairn add https://example.com
   # Should return in < 500ms (previously could take 2-5 seconds)
   ```

2. **Verify sync still happens**:
   ```bash
   cairn add https://test-async.com
   sleep 5  # Wait for background sync
   cairn sync status
   # Should show no pending changes (background sync completed)
   ```

3. **Verify offline behavior**:
   ```bash
   # Disconnect network
   cairn add https://offline-test.com
   # Should return immediately, no error
   # Reconnect network
   cairn sync push
   # Should push the pending change
   ```

4. **Verify explicit sync unchanged**:
   ```bash
   cairn sync push
   # Should still block and show progress/errors as before
   ```

### Automated Testing

```bash
go test ./...
```

## Key Implementation Detail

The `autoSyncPush` function in `cmd/cairn/main.go` now spawns `cairn sync push` as a detached background process using `os/exec.Command` + `cmd.Start()` instead of calling the sync engine directly. The background process inherits no stdout/stderr, so it cannot interfere with the user's terminal.
