# Data Model: Bookmark Cloud Sync

**Feature**: 001-bookmark-sync
**Date**: 2026-03-11

## Schema Migration V3

### Altered Table: `bookmarks`

New columns added to existing `bookmarks` table:

| Column | Type | Default | Nullable | Description |
|--------|------|---------|----------|-------------|
| `uuid` | TEXT | `''` | NOT NULL | Globally unique identifier (UUID v4). Populated on insert. Backfilled for existing rows during migration. |
| `updated_at` | TEXT | `''` | NOT NULL | RFC 3339 timestamp of last modification. Set on insert, updated on tag edit, delete (tombstone). Backfilled from `created_at` during migration. |

New indexes:

| Index | Columns | Purpose |
|-------|---------|---------|
| `idx_bookmarks_uuid` | `uuid` | Unique lookup for sync merge |
| `idx_bookmarks_updated_at` | `updated_at DESC` | Efficient "changes since" queries |

### New Table: `pending_sync`

| Column | Type | Default | Nullable | Description |
|--------|------|---------|----------|-------------|
| `id` | INTEGER | AUTO | NOT NULL | Primary key |
| `bookmark_uuid` | TEXT | — | NOT NULL | UUID of the affected bookmark |
| `operation` | TEXT | — | NOT NULL | One of: `add`, `update`, `delete` |
| `payload` | TEXT | `'{}'` | NOT NULL | JSON snapshot of the bookmark at time of change (for add/update) or empty (for delete) |
| `created_at` | TEXT | — | NOT NULL | RFC 3339 timestamp when the pending change was recorded |
| `retry_count` | INTEGER | `0` | NOT NULL | Number of failed sync attempts for this entry |

Index:

| Index | Columns | Purpose |
|-------|---------|---------|
| `idx_pending_sync_created_at` | `created_at ASC` | Process pending changes in chronological order |

### Migration V3 Backfill Logic

1. `ALTER TABLE bookmarks ADD COLUMN uuid TEXT NOT NULL DEFAULT ''`
2. `ALTER TABLE bookmarks ADD COLUMN updated_at TEXT NOT NULL DEFAULT ''`
3. Backfill `uuid` for all existing rows: generate a UUID v4 per row
4. Backfill `updated_at` for all existing rows: copy from `created_at`
5. Create unique index on `uuid`
6. Create index on `updated_at`
7. `CREATE TABLE pending_sync (...)`
8. Create index on `pending_sync.created_at`
9. Record version 3 in `schema_version`

## Entity: SyncConfig (file-based, not in SQLite)

Stored at: `$XDG_CONFIG_HOME/cairn/sync.json` (Linux) or `~/Library/Application Support/cairn/sync.json` (macOS)

```json
{
  "backend": "dropbox",
  "device_id": "a1b2c3d4-...",
  "last_sync_at": "2026-03-11T10:30:00Z",
  "sync_declined": false,
  "dropbox": {
    "access_token": "sl.xxx...",
    "refresh_token": "xxx...",
    "token_expiry": "2026-03-11T14:30:00Z",
    "app_key": "your-app-key"
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `backend` | string | Backend type: `"dropbox"`, `"s3"` (future) |
| `device_id` | string | UUID v4 generated at setup time |
| `last_sync_at` | string (RFC 3339) | Timestamp of last successful sync |
| `sync_declined` | bool | True if user declined first-run prompt |
| `dropbox` | object | Dropbox-specific credentials (present only when backend is "dropbox") |
| `dropbox.access_token` | string | Current short-lived access token |
| `dropbox.refresh_token` | string | Long-lived refresh token |
| `dropbox.token_expiry` | string (RFC 3339) | When the current access token expires |
| `dropbox.app_key` | string | Dropbox app key for PKCE flow |

File permissions: mode 0600 (owner read/write only).

## Entity: SyncRecord (cloud snapshot format)

Stored in Dropbox at: `/cairn/sync.json`

```json
{
  "version": 1,
  "last_updated_by": "a1b2c3d4-...",
  "last_updated_at": "2026-03-11T10:30:00Z",
  "bookmarks": [
    {
      "uuid": "e5f6g7h8-...",
      "url": "https://example.com",
      "domain": "example.com",
      "title": "Example",
      "description": "An example site",
      "tags": ["dev", "tools"],
      "is_permanent": false,
      "is_archived": false,
      "archived_at": null,
      "created_at": "2026-03-01T08:00:00Z",
      "updated_at": "2026-03-10T12:00:00Z",
      "deleted": false
    }
  ],
  "tombstones": [
    {
      "uuid": "z9y8x7w6-...",
      "url": "https://removed.com",
      "deleted_at": "2026-03-09T15:00:00Z",
      "deleted_by": "a1b2c3d4-..."
    }
  ]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `version` | int | Schema version for forward compatibility |
| `last_updated_by` | string | Device ID that last wrote this snapshot |
| `last_updated_at` | string (RFC 3339) | When this snapshot was last written |
| `bookmarks` | array | All active bookmarks across all devices |
| `tombstones` | array | Deleted bookmark markers; retained for propagation |

### Bookmark entry fields (within SyncRecord)

All fields from the local `bookmarks` table, plus:
- `deleted` (bool) — false for active bookmarks (redundant with tombstones but simplifies parsing)

### Tombstone entry fields

| Field | Type | Description |
|-------|------|-------------|
| `uuid` | string | UUID of the deleted bookmark |
| `url` | string | URL for reference/logging |
| `deleted_at` | string (RFC 3339) | When the deletion occurred |
| `deleted_by` | string | Device ID that performed the deletion |

## Entity Relationships

```
SyncConfig (local file)
  └── references → SyncBackend type ("dropbox" | "s3")

Bookmark (SQLite)
  ├── uuid → used as merge key in SyncRecord
  ├── updated_at → used for conflict resolution (last-write-wins)
  └── url → secondary dedup key (UNIQUE constraint)

PendingChange (SQLite, pending_sync table)
  └── bookmark_uuid → references Bookmark.uuid

SyncRecord (cloud JSON file)
  ├── bookmarks[] → each maps to a Bookmark by uuid
  └── tombstones[] → each maps to a deleted Bookmark by uuid
```

## State Transitions

### Bookmark Sync Lifecycle

```
[Local Only] ──(sync push)──→ [Synced]
[Cloud Only] ──(sync pull)──→ [Synced]
[Synced] ──(local edit)──→ [Modified, Pending Push]
[Modified, Pending Push] ──(sync push)──→ [Synced]
[Synced] ──(local delete)──→ [Tombstone, Pending Push]
[Tombstone, Pending Push] ──(sync push)──→ [Tombstone in Cloud]
[Tombstone in Cloud] ──(sync pull on other device)──→ [Deleted Locally]
```

### Pending Change Lifecycle

```
[Created] ──(sync attempt succeeds)──→ [Removed from pending_sync]
[Created] ──(sync attempt fails)──→ [retry_count++, stays in pending_sync]
[retry_count > 0] ──(next sync attempt succeeds)──→ [Removed from pending_sync]
```
