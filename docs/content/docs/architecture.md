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
