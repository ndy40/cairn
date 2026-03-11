# Implementation Plan: Async Autosync

**Branch**: `007-async-autosync` | **Date**: 2026-03-11 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/007-async-autosync/spec.md`

## Summary

Auto-push sync after bookmark add/delete/edit currently blocks the CLI until the network round-trip completes. This plan moves auto-push to a detached background subprocess so the CLI returns immediately after persisting the local change. Pending sync records are already written atomically in the same transaction as the bookmark operation, so no data is lost if the background sync fails or the process is interrupted.

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**: modernc.org/sqlite, charmbracelet/bubbletea, golang.org/x/oauth2, dropbox-sdk-go-unofficial/v6
**Storage**: SQLite (WAL mode, FTS5) via modernc.org/sqlite
**Testing**: `go test ./...`
**Target Platform**: Linux, macOS (single binary CLI)
**Project Type**: CLI + TUI
**Performance Goals**: Bookmark commands return in < 500ms regardless of sync
**Constraints**: Zero CGO, single binary, no external runtime
**Scale/Scope**: Single-user local tool

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate                           | Status | Notes                                                                 |
|--------------------------------|--------|-----------------------------------------------------------------------|
| No CGO                         | PASS   | No new dependencies; uses only Go stdlib `os/exec`, `syscall`         |
| Single binary                  | PASS   | Background sync re-invokes the same `cairn` binary                    |
| Task management                | PASS   | Tasks will be created via Backlog CLI after `/speckit.tasks`          |
| Backward-compatible migrations | PASS   | No schema changes required                                            |

## Project Structure

### Documentation (this feature)

```text
specs/007-async-autosync/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (minimal - no schema changes)
├── quickstart.md        # Phase 1 output
└── contracts/           # Phase 1 output (CLI contract)
```

### Source Code (repository root)

```text
cmd/cairn/main.go          # Modified: replace synchronous autoSyncPush calls with background subprocess spawn
internal/sync/engine.go    # No changes needed (AutoPush/Push remain the same, just called differently)
internal/store/bookmark.go # No changes needed (pending changes already recorded atomically)
```

**Structure Decision**: This is a minimal change scoped to `cmd/cairn/main.go`. The sync engine, store, and backend remain unchanged. The only new code is a helper function to spawn `cairn sync push` as a detached background process.

## Complexity Tracking

No constitution violations. No complexity justifications needed.
