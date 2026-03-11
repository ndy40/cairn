# CLI Contract: Sync Commands

**Feature**: 001-bookmark-sync
**Date**: 2026-03-11

## New Subcommands

All sync subcommands are grouped under `cairn sync`.

---

### `cairn sync setup`

Set up cloud sync for this device.

**Arguments**: None
**Flags**:
- `--backend <type>` — Backend type (default: `dropbox`). Future: `s3`.
- `--db <path>` — Override database path (inherited global flag).

**Interactive flow**:
1. If sync is already configured, print "Sync already configured (backend: dropbox). Reconfigure? (y/N)" and exit on decline.
2. Print "Visit this URL to authorize cairn:" followed by the OAuth2 authorization URL.
3. Prompt "Enter the authorization code:" and wait for user input.
4. Exchange code for tokens, store in sync config file.
5. Generate device ID (UUID v4), store in sync config.
6. If cloud snapshot exists: download, merge with local DB (cloud wins on URL conflict).
7. If cloud snapshot does not exist: upload local bookmarks as initial snapshot.
8. Print "Sync configured. N bookmarks synced."

**Exit codes**:
- 0: Success
- 1: Auth failed or user cancelled
- 3: Error (network, DB, or file I/O)

---

### `cairn sync push`

Manually push local changes to cloud.

**Arguments**: None
**Flags**:
- `--db <path>` — Override database path.

**Behavior**:
1. Check sync is configured; error if not.
2. Replay any pending changes from `pending_sync` table.
3. Build full snapshot from local DB.
4. Upload snapshot to cloud backend.
5. Clear pending changes on success.
6. Update `last_sync_at` in config.

**Output**:
- "Pushed N changes. Up to date." on success.
- "Already up to date." if no changes since last sync.
- "Sync not configured. Run `cairn sync setup` first." if unconfigured.

**Exit codes**:
- 0: Success (including already up to date)
- 1: Sync not configured
- 3: Error

---

### `cairn sync pull`

Manually pull changes from cloud.

**Arguments**: None
**Flags**:
- `--db <path>` — Override database path.

**Behavior**:
1. Check sync is configured; error if not.
2. Download snapshot from cloud backend.
3. Merge cloud bookmarks with local DB (cloud wins on URL conflict by `updated_at`).
4. Apply tombstones: delete locally any bookmarks that appear in cloud tombstones.
5. Update `last_sync_at` in config.

**Output**:
- "Pulled N new, M deleted. Up to date." on success.
- "Already up to date." if no changes.
- "Sync not configured. Run `cairn sync setup` first." if unconfigured.

**Exit codes**:
- 0: Success
- 1: Sync not configured
- 3: Error

---

### `cairn sync status`

Show current sync configuration and status.

**Arguments**: None
**Flags**:
- `--db <path>` — Override database path.
- `--json` — Output as JSON.

**Output (text)**:
```
Sync: configured
Backend: dropbox
Device ID: a1b2c3d4-...
Last sync: 2026-03-11 10:30:00 UTC
Pending changes: 3
```

Or if not configured:
```
Sync: not configured
Run `cairn sync setup` to get started.
```

**Output (JSON)**:
```json
{
  "configured": true,
  "backend": "dropbox",
  "device_id": "a1b2c3d4-...",
  "last_sync_at": "2026-03-11T10:30:00Z",
  "pending_changes": 3
}
```

**Exit codes**:
- 0: Always (status is informational)

---

### `cairn sync auth`

Re-authenticate with the cloud backend (when refresh token expires or is revoked).

**Arguments**: None
**Flags**:
- `--db <path>` — Override database path.

**Behavior**:
1. Check sync is configured; error if not.
2. Run same OAuth2 flow as setup (generate URL, prompt for code).
3. Update tokens in sync config file.
4. Attempt to replay pending changes.
5. Print "Re-authenticated. N pending changes synced."

**Exit codes**:
- 0: Success
- 1: Sync not configured or auth failed
- 3: Error

---

### `cairn sync unlink`

Remove sync configuration from this device without deleting local bookmarks.

**Arguments**: None
**Flags**:
- `--db <path>` — Override database path.

**Behavior**:
1. Check sync is configured; error if not.
2. Confirm: "Unlink this device from sync? Local bookmarks will be kept. (y/N)"
3. On confirm: delete sync config file, clear `pending_sync` table.
4. Print "Sync unlinked. Local bookmarks preserved."

**Exit codes**:
- 0: Success
- 1: Sync not configured or user cancelled
- 3: Error

---

## First-Run Prompt Behavior

When **any** `cairn` command is run and no sync config file exists and `sync_declined` is not set:

```
No sync configured — connect to Dropbox? (y/N):
```

- **y**: Run the `sync setup` flow inline, then continue with the original command.
- **N** (or Enter): Set `sync_declined: true` in sync config file. Continue with the original command. Never prompt again.

This prompt fires **once** — on the very first CLI invocation with no sync state.

---

## Auto-Sync Behavior (non-interactive)

### On CLI startup (before main operation):

1. If sync is configured → attempt pull.
2. On success with changes → print `↓ N new bookmarks synced`.
3. On success with no changes → silent.
4. On failure → print `⚠ Sync pull failed: <reason>`. Continue normally.

### After modifying operations (add, delete, edit tags):

1. If sync is configured → record change in `pending_sync`, then attempt push.
2. On success → silent (or print `↑ Synced` if pending queue was > 0).
3. On failure → print `⚠ Sync push failed: change queued for later`.

### Pending replay:

On every successful connection (push or pull), also replay any pending changes from `pending_sync`.
