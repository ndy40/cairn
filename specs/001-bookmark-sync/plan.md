# Implementation Plan: Bookmark Cloud Sync

**Branch**: `001-bookmark-sync` | **Date**: 2026-03-11 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-bookmark-sync/spec.md`

## Summary

Add bidirectional bookmark sync across devices using cloud storage (Dropbox first). The sync engine uses a JSON snapshot stored as a single file in Dropbox, with local pending-change tracking for offline resilience. Sync fires automatically (pull on startup, push after modifications) and is also available via manual CLI commands. The backend is abstracted behind an interface to support future providers (S3, etc.).

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**: modernc.org/sqlite (existing), golang.org/x/oauth2 (new), dropbox-sdk-go-unofficial/v6 (new), google/uuid (promote from indirect)
**Storage**: SQLite (existing, WAL mode, FTS5) + local JSON config file for sync credentials
**Testing**: Go standard `testing` package, table-driven tests, mock backend for integration
**Target Platform**: Linux, macOS (cross-platform CLI)
**Project Type**: CLI + TUI (single binary)
**Performance Goals**: Sync of 10,000 bookmarks < 30 seconds; auto-sync adds < 2 seconds to CLI startup
**Constraints**: Zero CGO, single binary, offline-capable (pending queue), atomic operations
**Scale/Scope**: Up to 10,000 bookmarks per user; single-file snapshot < 10MB

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate | Requirement | Status | Notes |
|------|-------------|--------|-------|
| No CGO | All dependencies pure Go | ✅ PASS | `dropbox-sdk-go-unofficial/v6` is pure Go; `golang.org/x/oauth2` is pure Go; `google/uuid` is pure Go |
| Single binary | No external runtime | ✅ PASS | All new code compiles into the existing cairn binary. Sync config is a JSON file, not a separate process. |
| Task management | Tasks via Backlog CLI | ✅ PASS | Will use `backlog task create` after `/speckit.tasks` |
| Backward-compatible migrations | ALTER TABLE with DEFAULT only | ✅ PASS | V3 migration uses `ALTER TABLE bookmarks ADD COLUMN uuid TEXT NOT NULL DEFAULT ''` and `ADD COLUMN updated_at TEXT NOT NULL DEFAULT ''`. Backfill runs as UPDATE statements. New table `pending_sync` created fresh. No destructive changes. |

**Post-Phase 1 Re-check**: All gates still pass. No new dependencies or schema changes beyond what's listed.

## Project Structure

### Documentation (this feature)

```text
specs/001-bookmark-sync/
├── plan.md              # This file
├── research.md          # Phase 0: technology decisions
├── data-model.md        # Phase 1: schema, entities, state transitions
├── quickstart.md        # Phase 1: build, run, test instructions
├── contracts/
│   ├── cli-commands.md          # CLI subcommand contracts
│   └── sync-backend-interface.md  # SyncBackend interface contract
└── tasks.md             # Phase 2 output (via /speckit.tasks)
```

### Source Code (repository root)

```text
cmd/cairn/
└── main.go              # Modified: add sync subcommand routing + first-run prompt + auto-sync hooks

internal/
├── sync/                # NEW PACKAGE: sync engine
│   ├── config.go        # SyncConfig: load/save JSON, OS-specific paths, mode 0600
│   ├── engine.go        # Orchestration: Push(), Pull(), Setup(), AutoPull(), AutoPush()
│   ├── merge.go         # Merge algorithm: dedup by URL, last-write-wins by updated_at
│   ├── snapshot.go      # SyncRecord JSON marshal/unmarshal
│   └── backend/         # NEW SUBPACKAGE: backend implementations
│       ├── backend.go   # SyncBackend interface + error types
│       └── dropbox.go   # Dropbox implementation via SDK + oauth2.TokenSource
├── store/               # MODIFIED: sync-related store methods
│   ├── store.go         # Modified: add migration V3
│   ├── bookmark.go      # Modified: set uuid + updated_at on Insert; update updated_at on tag edit
│   └── sync.go          # NEW: pending_sync CRUD, ExportAll(), ImportBatch()
└── model/               # MODIFIED: TUI auto-sync integration
    └── app.go           # Modified: auto-pull on Init() via tea.Cmd
```

**Structure Decision**: Follows the existing single-project layout. New `internal/sync/` package owns sync orchestration and backend abstraction. Store modifications are minimal (new file + migration). CLI routing extends the existing subcommand pattern in `main.go`.

## Complexity Tracking

No constitution violations to justify. All gates pass.
