---
title: "Dropbox Sync"
weight: 50
---

# Dropbox Sync

Cairn can sync bookmarks to Dropbox, keeping them in sync across multiple devices. Sync runs automatically in the background on startup — it will not block the TUI or any CLI command.

## Prerequisites

You need a Dropbox app key. Either:

- Use an existing Cairn-compatible Dropbox app key provided by your team/setup, or
- Create your own at [dropbox.com/developers/apps](https://www.dropbox.com/developers/apps)

## Setup

### 1. Provide the App Key

Add the app key to `cairn.json` (recommended):

```json
{
  "dropbox_app_key": "your-app-key-here"
}
```

Or export it temporarily for setup:

```sh
export CAIRN_DROPBOX_APP_KEY=your-app-key-here
```

### 2. Run Setup

```sh
cairn sync setup
```

This will:
1. Open a browser for OAuth2 authentication with Dropbox
2. Save credentials to `~/.config/cairn/sync.json` (Linux) or the OS equivalent
3. Perform an initial sync of all existing bookmarks

### 3. Verify

```sh
cairn sync status
```

**Example output:**

```
Backend:         dropbox
Device ID:       a1b2c3d4-...
Last sync:       2025-07-15 09:30:00 UTC
Pending changes: 0
```

## Sync Commands

| Command | Description |
|---------|-------------|
| `cairn sync setup` | First-time setup and initial sync |
| `cairn sync push` | Push local changes to Dropbox |
| `cairn sync pull` | Pull changes from Dropbox |
| `cairn sync status` | Show sync state and pending count |
| `cairn sync auth` | Re-authenticate (refresh OAuth2 token) |
| `cairn sync unlink` | Disconnect sync; local bookmarks are preserved |

## Automatic Background Sync

When sync is configured, Cairn automatically:

- **Pulls** on startup — picks up bookmarks added on other devices
- **Pushes** after `add` and `delete` — keeps remote in sync immediately

Background sync runs as a detached subprocess and will not block the TUI or produce any terminal output.

A lockfile (`cairn-sync-pull.lock` / `cairn-sync-push.lock`) prevents concurrent syncs. Stale locks older than 10 minutes are cleaned up automatically.

## Re-authentication

If your Dropbox token expires:

```sh
cairn sync auth
```

This re-runs the OAuth2 flow and updates the token in `sync.json`.

## Unlinking Sync

To disconnect a device from Dropbox sync:

```sh
cairn sync unlink
```

Your local bookmarks are **not deleted**. Only the sync configuration is removed.

## Credential Storage

OAuth2 tokens are stored in:

| OS | Path |
|----|------|
| Linux | `~/.config/cairn/sync.json` |
| macOS | `~/Library/Application Support/cairn/sync.json` |
| Windows | `%APPDATA%\cairn\sync.json` |

The file is created with `0600` (owner-only) permissions. See [Security]({{< relref "/docs/security" >}}) for best practices.
