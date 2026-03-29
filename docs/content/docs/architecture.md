---
title: "Architecture"
weight: 80
---

# Architecture

## Overview

Cairn is a single Go binary built with `CGO_ENABLED=0` for fully static, dependency-free execution across Linux, macOS, and Windows.

```
cairn (binary)
├── cmd/cairn/main.go        CLI entry point, flag parsing, command routing
├── internal/
│   ├── model/               TUI — bubbletea MVU state machine
│   │   ├── app.go           AppState machine, Update/View routing
│   │   ├── browse.go        Bookmark list view and keybindings
│   │   └── edit.go          Bookmark edit form
│   ├── store/               SQLite persistence layer
│   │   └── bookmark.go      CRUD, FTS5, tags, NormaliseTags
│   ├── sync/                Dropbox sync engine
│   │   ├── engine.go        Push/pull orchestration
│   │   ├── backend.go       Backend interface + Dropbox implementation
│   │   └── oauth.go         OAuth2 PKCE flow
│   ├── config/              Configuration management (Viper)
│   ├── search/              Fuzzy search (client-side ranking)
│   ├── fetcher/             URL title/description fetcher
│   ├── display/             Terminal capability detection
│   └── clipboard/           Clipboard integration
```

## Key Design Decisions

### Static Binary

Built with `CGO_ENABLED=0` using `modernc.org/sqlite` — a pure-Go SQLite port. This avoids the need for system libraries and produces a single, portable binary.

### TUI Architecture: bubbletea MVU

The interactive TUI uses [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) following the Model-View-Update (MVU) pattern:

- **Model** (`internal/model/app.go`) — `AppState` enum drives all state transitions
- **Update** — handles keypress messages and produces new state
- **View** — renders the current state using [lipgloss](https://github.com/charmbracelet/lipgloss) for styling

### Database: SQLite with FTS5

- Schema versioned via `schema_version` migrations in `store.Open()`
- FTS5 virtual table for full-text search over title, description, and domain
- WAL mode enabled for concurrent read access
- Tags stored as comma-separated strings; normalised by `NormaliseTags`

### Sync: Background Process Model

Background sync (`backgroundSyncPull` / `backgroundSyncPush`) spawns a detached subprocess of the same binary to run `cairn sync pull/push`. This avoids blocking the TUI and keeps the sync logic isolated. A lockfile prevents concurrent syncs.

## Cross-Device Synchronisation

Cairn uses a **snapshot-based** sync model. Each device maintains its own local SQLite database and independently pushes or pulls a shared JSON snapshot to/from a cloud storage backend.

> **Note:** Dropbox is the only supported backend at this time. S3 support is planned for a future release.

### Sync Architecture

```
  Device A (laptop)                    Cloud Store                 Device B (desktop)
  ─────────────────                 ──────────────────             ──────────────────
  ┌─────────────┐                   ┌──────────────┐              ┌─────────────┐
  │  cairn TUI  │                   │              │              │  cairn TUI  │
  │  or CLI     │                   │  Dropbox     │              │  or CLI     │
  └──────┬──────┘                   │      or      │              └──────┬──────┘
         │                          │  S3 (future) │                     │
  ┌──────▼──────┐    push snapshot  │              │  pull snapshot ┌──────▼──────┐
  │  SQLite DB  │ ─────────────────►│  sync.json   │◄────────────── │  SQLite DB  │
  │  (local)    │ ◄─────────────────│  (snapshot)  │───────────────►│  (local)    │
  └─────────────┘    pull & merge   │              │   push & merge └─────────────┘
                                    └──────────────┘
```

### How Sync Works

**Push (upload):**

1. Cairn exports all bookmarks from the local SQLite database.
2. A `SyncRecord` JSON snapshot is assembled, tagged with the device's unique ID.
3. The snapshot is uploaded to `cairn/sync.json` on the cloud store, overwriting the previous snapshot.
4. Pending local changes are cleared; `last_sync_at` is recorded in config.

**Pull (download + merge):**

1. Cairn downloads `cairn/sync.json` from the cloud store.
2. The remote snapshot is compared against local bookmarks by UUID and URL.
3. A three-way merge is performed:
   - **New remote bookmarks** → inserted locally.
   - **Updated remote bookmarks** → tags and metadata applied locally.
   - **Deleted remote bookmarks** → removed locally.
4. The merged state is re-uploaded as the new canonical snapshot.

### Multi-Device Sync Flow

```
  ┌──────────────────────────────────────────────────────────────────┐
  │                        Typical Workflow                          │
  │                                                                  │
  │  1. Device A adds bookmarks → cairn auto-pushes snapshot         │
  │                                           │                      │
  │                                           ▼                      │
  │                                    [ sync.json ]                 │
  │                                           │                      │
  │  2. Device B opens cairn → auto-pull pulls snapshot              │
  │                          → merges new bookmarks into local DB    │
  │                          → re-uploads merged state               │
  │                                           │                      │
  │  3. Both devices now share the same bookmark set                 │
  └──────────────────────────────────────────────────────────────────┘
```

### Storage Backends

| Backend  | Status    | Snapshot Path       | Auth Method      |
|----------|-----------|---------------------|------------------|
| Dropbox  | Supported | `/cairn/sync.json`  | OAuth2 PKCE flow |
| S3       | Planned   | —                   | —                |

Each device is assigned a unique `device_id` (UUID) on first sync setup. This ID is stored in the local sync config and embedded in every snapshot, allowing the merge logic to identify the origin of changes.

### Auto-Sync

When enabled, cairn automatically triggers a **push** after any write operation (add, edit, delete, archive) and a **pull** on TUI startup. Both run as detached background subprocesses so they never block the foreground UI or CLI.

### Configuration: Viper with Precedence

`internal/config` uses [Viper](https://github.com/spf13/viper) to merge configuration from:
1. Environment variables (`CAIRN_*`)
2. CLI flags (`--db`)
3. `cairn.json` config file
4. Compiled-in defaults

## Data Model

### Bookmark

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER | Primary key (auto-increment) |
| `url` | TEXT | Full URL (unique) |
| `title` | TEXT | Page title (fetched automatically) |
| `description` | TEXT | Page description |
| `domain` | TEXT | Extracted hostname |
| `tags` | TEXT | Comma-separated tag list |
| `is_permanent` | BOOLEAN | Pin flag (true = never auto-archived) |
| `is_archived` | BOOLEAN | Archive flag |
| `created_at` | DATETIME | Insert timestamp |

### Auto-archiving

On every TUI startup, `ArchiveStale()` archives bookmarks where `created_at <= NOW - 30 days` and `is_permanent = false`.

## Dependencies

| Package | Purpose |
|---------|---------|
| `modernc.org/sqlite` | Pure-Go SQLite (no CGO) |
| `charmbracelet/bubbletea` | TUI framework (MVU) |
| `charmbracelet/bubbles` | TUI components (list, input) |
| `charmbracelet/lipgloss` | Terminal styling |
| `golang.org/x/oauth2` | OAuth2 for Dropbox auth |
| `dropbox-sdk-go-unofficial/v6` | Dropbox API client |
| `spf13/viper` | Configuration management |
| `google/uuid` | Device ID generation for sync |
