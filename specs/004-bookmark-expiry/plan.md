# Implementation Plan: Bookmark Expiry & Last-Visited Removal

**Branch**: `004-bookmark-expiry` | **Date**: 2026-03-06 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-bookmark-expiry/spec.md`

---

## Summary

Two focused changes: (1) simplify the startup archive check to expire bookmarks based on creation date alone (30-day threshold), replacing the previous 183-day last-visited-based rule; (2) remove all last-visited tracking — delete `UpdateLastVisited()`, remove its call from `openBookmarkCmd`, and strip the "Last:" / "Never visited" labels from browse rows. No schema migrations. No new dependencies.

---

## Technical Context

**Language/Version**: Go 1.22+
**Primary Dependencies**: charmbracelet/bubbletea, bubbles, lipgloss, modernc.org/sqlite — all unchanged
**Storage**: SQLite (WAL, FTS5). No schema changes; `last_visited_at` column retained but no longer written or displayed.
**Testing**: `go test ./...` (unchanged)
**Target Platform**: Linux, macOS, Windows (unchanged)
**Project Type**: CLI/TUI application (unchanged)
**Performance Goals**: Startup archive check completes before first TUI frame (<500 ms)
**Constraints**: Zero new CGO dependencies; no destructive schema migrations
**Scale/Scope**: Single user, local (unchanged)

---

## Constitution Check

| Gate | Status | Notes |
|------|--------|-------|
| No CGO dependencies | PASS | No new dependencies |
| Single binary deployment | PASS | No new external runtime |
| Task management | PASS | Tasks go through Backlog CLI after `/speckit.tasks` |
| Backward-compatible migrations | PASS | No schema changes at all |

No violations. No Complexity Tracking table required.

---

## Project Structure

### Documentation (this feature)

```text
specs/004-bookmark-expiry/
├── spec.md
├── plan.md              # This file
├── research.md
├── data-model.md
├── contracts/
│   └── browse-row.md
├── checklists/
│   └── requirements.md
└── tasks.md             # /speckit.tasks output (NOT created here)
```

### Source Code (additions/modifications)

```text
bookmark-manager/
├── internal/
│   ├── store/
│   │   └── archive.go   # MODIFIED: ArchiveStale() threshold 183→30 days, creation-date only;
│   │                    #           remove UpdateLastVisited()
│   └── model/
│       ├── app.go       # MODIFIED: remove s.UpdateLastVisited(b.ID) call from openBookmarkCmd;
│       │                #           update comment
│       └── browse.go    # MODIFIED: remove last-visited branches from BookmarkItem.Description()
```

**Structure Decision**: Three targeted edits to existing files. No new files. No new packages.

---

## Phase 0 Output

- [x] research.md

## Phase 1 Output

- [x] data-model.md
- [x] contracts/browse-row.md
- [ ] tasks.md — `/speckit.tasks` (next step)
