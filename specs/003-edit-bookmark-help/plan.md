# Implementation Plan: Edit Bookmark Tags, Last-Visited Visibility & CLI Help

**Branch**: `003-edit-bookmark-help` | **Date**: 2026-03-06 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-edit-bookmark-help/spec.md`

---

## Summary

Three focused improvements to the existing TUI bookmark manager: (1) an edit panel triggered by `e` in browse mode that lets users add or change tags on any existing active bookmark; (2) confirm and verify last-visited date updates immediately in the list after pressing Enter to open a bookmark (list reloads automatically); (3) standard `--help`/`-h` flag support on the root command and all subcommands. No schema changes. No new dependencies.

---

## Technical Context

**Language/Version**: Go 1.22+ (unchanged)

**Primary Dependencies**: All unchanged from features 001–002. No new libraries required.

**Storage**: No schema changes. Tags field already exists from migration v2 (feature 002). The edit operation adds one new store method `UpdateTags(id int64, tags []string) error` — no DDL.

**Testing**: `go test ./...` (unchanged)

**Target Platform**: Linux, macOS, Windows (unchanged)

**Project Type**: CLI/TUI application (unchanged)

**Performance Goals**:
- Edit panel opens within one render frame (<100ms)
- `bm --help` exits within 100ms

**Constraints**:
- Zero new CGO dependencies
- No schema migrations required
- `e` key is currently unbound in browse mode — safe to assign

**Scale/Scope**: Single user, local (unchanged)

---

## Constitution Check

| Gate | Status | Notes |
|------|--------|-------|
| No CGO dependencies | PASS | No new dependencies |
| Single binary deployment | PASS | No new external runtime |
| Task management | PASS | Tasks will go through Backlog CLI after `/speckit.tasks` |
| Backward-compatible migrations | PASS | No schema changes needed |

No violations. No Complexity Tracking table required.

---

## Project Structure

### Documentation (this feature)

```text
specs/003-edit-bookmark-help/
├── spec.md
├── plan.md              # This file
├── research.md
├── data-model.md
├── contracts/
│   ├── keyboard-shortcuts.md
│   └── cli-interface.md
├── checklists/
│   └── requirements.md
└── tasks.md             # /speckit.tasks output (NOT created here)
```

### Source Code (additions/modifications)

```text
bookmark-manager/
├── cmd/
│   └── bm/
│       └── main.go              # MODIFIED: --help/-h on root and each subcommand
├── internal/
│   ├── model/
│   │   ├── app.go               # MODIFIED: StateEdit; 'e' key handler; editView()
│   │   ├── browse.go            # MODIFIED: editKey binding
│   │   └── edit.go              # NEW: EditModel with read-only title + editable tags input
│   └── store/
│       └── bookmark.go          # MODIFIED: UpdateTags(id int64, tags []string) error
```

**Structure Decision**: One new file (`internal/model/edit.go`), targeted edits to four existing files. No new packages.

---

## Phase 0 Output

- [x] research.md

## Phase 1 Output

- [x] data-model.md
- [x] contracts/keyboard-shortcuts.md
- [x] contracts/cli-interface.md
- [ ] tasks.md — `/speckit.tasks` (next step)
