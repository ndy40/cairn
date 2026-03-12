# Data Model: Configuration File Support

**Feature**: 010-config-file
**Date**: 2026-03-12

## Entities

### FileConfig

Represents the raw configuration loaded from `cairn.json`. All fields are pointers to distinguish "not set" from "set to zero value".

| Field | Type | JSON Key | Description |
|-------|------|----------|-------------|
| DBPath | *string | `db_path` | Path to the SQLite database file |
| DropboxAppKey | *string | `dropbox_app_key` | Dropbox OAuth app key for sync |

### AppConfig

Represents the final resolved configuration after merging all sources (env vars, CLI flags, config file, defaults).

| Field | Type | Description |
|-------|------|-------------|
| DBPath | string | Resolved database path |
| DropboxAppKey | string | Resolved Dropbox app key |

## Relationships

- `FileConfig` is an intermediate representation loaded from disk.
- `AppConfig` is produced by merging: defaults → FileConfig overrides → CLI flag overrides → env var overrides.
- No database entities. No schema changes. No migrations.

## File Format

```json
{
  "db_path": "/custom/path/to/bookmarks.db",
  "dropbox_app_key": "my-app-key"
}
```

All keys are optional. Missing keys are not overridden (fall through to lower-precedence sources).

## File Location

| OS | Path |
|----|------|
| Linux | `$XDG_CONFIG_HOME/cairn/cairn.json` (default: `~/.config/cairn/cairn.json`) |
| macOS | `~/Library/Application Support/cairn/cairn.json` |
| Windows | `%APPDATA%/cairn/cairn.json` |
