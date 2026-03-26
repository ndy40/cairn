# Data Model: Interactive Setup Configuration Prompts

## Entities

This feature does not introduce new data entities or schema changes. It extends the existing configuration write path.

---

## Affected Existing Entities

### AppConfig (`internal/config/config.go`)

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| `DropboxAppKey` | `string` | env `CAIRN_DROPBOX_APP_KEY` → `cairn.json` `dropbox_app_key` → interactive prompt | Prompted when empty after all other sources |
| `DBPath` | `string` | env `CAIRN_DB_PATH` → CLI `--db` → `cairn.json` `db_path` → default | Prompted when value equals the built-in default |

### cairn.json (config file on disk)

Fields written by setup prompt flow:

| Key | Type | Written when |
|-----|------|-------------|
| `dropbox_app_key` | string | User provides value at prompt (always in this flow) |
| `db_path` | string | User provides a non-empty custom path |

**Preservation rule**: All pre-existing keys in `cairn.json` are preserved. The config Manager (viper) loads existing values before writing, so unrelated keys survive the write.

---

## State Transitions

```
cairn sync setup invoked
        │
        ▼
AppKey already resolved? ──YES──► skip prompt
        │ NO
        ▼
Prompt: "Enter Dropbox App Key:"
        │
        ├─ empty input ──► re-prompt with error message
        │
        └─ non-empty input
              │
              ▼
        Write dropbox_app_key to cairn.json
              │
              ▼
        DB path already overridden? ──YES──► skip prompt
              │ NO
              ▼
        Prompt: "Enter database path (press Enter for default: <path>):"
              │
              ├─ empty input ──► use default, do not write db_path to file
              │
              └─ non-empty input ──► write db_path to cairn.json
                    │
                    ▼
              Print: "Config written to <path>"
                    │
                    ▼
              Proceed with OAuth (existing Setup engine)
```

---

## No Schema Changes

No SQLite schema changes. No new migration version required.
