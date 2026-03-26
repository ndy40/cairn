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
  "dropbox_app_key": "your-app-key"
}
```

## Supported Settings

| JSON key | Environment variable | CLI flag | Description |
|----------|---------------------|----------|-------------|
| `db_path` | `CAIRN_DB_PATH` | `--db` | Path to the SQLite bookmark database |
| `dropbox_app_key` | `CAIRN_DROPBOX_APP_KEY` | — | Dropbox app key for sync |

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
