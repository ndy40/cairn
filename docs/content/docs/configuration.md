---
title: "Configuration"
weight: 40
---

# Configuration

Cairn can be configured via an optional JSON config file, environment variables, and CLI flags.

## Precedence (highest to lowest)

1. **Environment variables** — always win
2. **CLI flags** (e.g. `--db`)
3. **Config file** (`cairn.json`)
4. **Defaults** (OS-appropriate paths)

## Config File

Create a `cairn.json` file in your OS config directory:

| OS | Path |
|----|------|
| Linux | `$XDG_CONFIG_HOME/cairn/cairn.json` (default: `~/.config/cairn/cairn.json`) |
| macOS | `~/Library/Application Support/cairn/cairn.json` |
| Windows | `%APPDATA%\cairn\cairn.json` |

All keys are optional. Example:

```json
{
  "db_path": "/path/to/bookmarks.db",
  "dropbox_app_key": "your-app-key",
  "disable_auto_archive": false
}
```

## Supported Settings

| JSON key | Environment variable | CLI flag | Description |
|----------|---------------------|----------|-------------|
| `db_path` | `CAIRN_DB_PATH` | `--db` | Path to the SQLite bookmark database |
| `dropbox_app_key` | `CAIRN_DROPBOX_APP_KEY` | — | Dropbox app key for sync |
| `disable_auto_archive` | `CAIRN_DISABLE_AUTO_ARCHIVE` | — | Set to `true` to stop auto-archiving stale bookmarks (default `false`) |

## Disable Auto-Archiving

By default, Cairn keeps your active list tidy: on every TUI startup it
auto-archives any bookmark older than 30 days that isn't pinned. Archived
bookmarks are never deleted — you can review and [restore]({{< relref "/docs/quickstart" >}}#8-review-and-restore-archived-bookmarks)
them from the archive view (press `a` in the TUI).

If you'd rather keep every bookmark in the main list indefinitely, disable
auto-archiving.

**Via the config file** — add to `cairn.json`:

```json
{
  "disable_auto_archive": true
}
```

**Via an environment variable:**

```sh
export CAIRN_DISABLE_AUTO_ARCHIVE=true
```

With this enabled, Cairn skips the archive check on startup and no bookmarks are
moved to the archive. Bookmarks already in the archive stay there until you
restore them; disabling auto-archiving does not un-archive anything, and you can
still restore them manually. Pinning individual bookmarks (`p` in the TUI, or
`cairn pin <id>`) remains an option if you only want to protect a few rather
than turning the feature off entirely.

Verify the setting took effect with `cairn config` (see below); it prints
`CAIRN_DISABLE_AUTO_ARCHIVE=true` when enabled.

## Inspect Resolved Config

```sh
cairn config
```

Prints the effective values that Cairn will use, after applying all precedence rules.

## Database Location

If `db_path` is not set, Cairn uses a platform-appropriate default:

| OS | Default path |
|----|-------------|
| Linux | `$XDG_DATA_HOME/cairn/bookmarks.db` (default: `~/.local/share/cairn/bookmarks.db`) |
| macOS | `~/Library/Application Support/cairn/bookmarks.db` |
| Windows | `%APPDATA%\cairn\bookmarks.db` |

The database is created automatically on first run if it does not exist.

## Security Notes

- Store `dropbox_app_key` in `cairn.json`, not in environment variables — env vars are visible in process listings (`ps aux`) and shell history
- Ensure `cairn.json` has owner-only read permissions: `chmod 600 ~/.config/cairn/cairn.json`
- Never commit `cairn.json` or `sync.json` to version control

See [Security]({{< relref "/docs/security" >}}) for full guidance.
