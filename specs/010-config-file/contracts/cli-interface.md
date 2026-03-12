# CLI Interface Contract: Configuration

**Feature**: 010-config-file

## `cairn config` Command

Displays the resolved configuration values. No changes to the command syntax — only the resolution logic changes.

### Output Format

```text
CAIRN_DB_PATH=<resolved path>
CAIRN_DROPBOX_APP_KEY=<resolved key>
```

### Resolution Order (highest to lowest priority)

1. Environment variable (`CAIRN_DB_PATH`, `CAIRN_DROPBOX_APP_KEY`)
2. CLI flag (`--db` for db path)
3. Config file (`cairn.json` → `db_path`, `dropbox_app_key`)
4. Default (OS-appropriate path for db, empty string for dropbox key)

### Error Behavior

| Condition | Behavior |
|-----------|----------|
| Config file missing | No error, no warning — use remaining sources |
| Config file invalid JSON | Print error to stderr, exit with code 1 |
| Config file wrong types | Print error to stderr, exit with code 1 |
| Config file empty (`{}`) | No error — treat as no overrides |
| Config file has unknown keys | Silently ignored |

## `cairn.json` File Format

```json
{
  "db_path": "string (optional)",
  "dropbox_app_key": "string (optional)"
}
```

All keys are optional. The file itself is optional.
