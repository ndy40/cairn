# Quickstart: Bookmark Cloud Sync

**Feature**: 001-bookmark-sync
**Date**: 2026-03-11

## Prerequisites

- Go 1.22+ (project uses Go 1.25.0)
- Dropbox account (for testing)
- Dropbox App created at https://www.dropbox.com/developers/apps with:
  - App type: Scoped access
  - Access type: App folder (or Full Dropbox if preferred)
  - Permissions: `files.content.write`, `files.content.read`
  - Note the **App key** (no App secret needed for PKCE flow)

## New Dependencies

```
go get golang.org/x/oauth2
go get github.com/dropbox/dropbox-sdk-go-unofficial/v6
go get github.com/google/uuid
```

All three are pure Go (zero CGO).

> Note: `google/uuid` is already an indirect dependency of the project. Adding it as a direct dependency for bookmark UUID generation.

## Key File Locations

### Source Code (new)

```
internal/
├── sync/
│   ├── config.go          # SyncConfig load/save, config file path resolution
│   ├── engine.go          # Core sync orchestration: merge, push, pull, pending replay
│   ├── merge.go           # Merge algorithm: dedup by URL, last-write-wins by updated_at
│   ├── snapshot.go        # SyncRecord JSON serialization/deserialization
│   └── backend/
│       ├── backend.go     # SyncBackend interface definition
│       └── dropbox.go     # Dropbox implementation using SDK
├── store/
│   (existing files, plus:)
│   ├── sync.go            # New: pending_sync CRUD, bulk bookmark export for snapshot
│   └── store.go           # Modified: migration V3 added
└── model/
    (existing files, modified:)
    ├── app.go             # Modified: auto-pull on Init(), auto-push after modifying ops
```

### CLI (modified)

```
cmd/cairn/main.go          # Add sync subcommand routing
```

### Config (runtime)

```
$XDG_CONFIG_HOME/cairn/sync.json   # Linux
~/Library/Application Support/cairn/sync.json  # macOS
```

## Build & Run

```bash
# Build
go build -o cairn ./cmd/cairn

# Test sync setup
./cairn sync setup

# Test manual sync
./cairn sync push
./cairn sync pull
./cairn sync status

# Test auto-sync (just use the app normally — sync fires automatically)
./cairn add https://example.com
./cairn                        # TUI auto-pulls on startup
```

## Testing Strategy

### Unit Tests

- `internal/sync/merge_test.go` — merge algorithm: dedup, conflict resolution, tombstone application
- `internal/sync/config_test.go` — config load/save, default paths, file permissions
- `internal/sync/snapshot_test.go` — JSON round-trip serialization
- `internal/store/sync_test.go` — pending_sync table CRUD, migration V3

### Integration Tests

- `internal/sync/engine_test.go` — full push/pull cycle with a mock backend
- `internal/sync/backend/dropbox_test.go` — Dropbox upload/download with mock HTTP server

### Manual Testing

1. Set up sync on Device A → verify cloud snapshot created
2. Set up sync on Device B → verify bookmarks pulled
3. Add bookmark on Device A → verify auto-push
4. Open cairn on Device B → verify auto-pull at startup
5. Kill network on Device A → add bookmark → verify pending queue → restore network → verify auto-reconcile on next operation

## Architecture Notes

- **Sync engine** is decoupled from the store — it takes a `SyncBackend` interface and a `*store.Store`, orchestrating the merge.
- **Auto-sync hooks** are injected into the CLI flow (not the TUI model directly) — keeping the TUI model focused on UI state.
- **Pending changes** are written in the same transaction as the bookmark modification — atomic or nothing.
- **Token refresh** is handled by `golang.org/x/oauth2.TokenSource` — the Dropbox backend wraps the HTTP client with auto-refresh.
