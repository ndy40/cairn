# Implementation Plan: Edit Bookmark

**Branch**: `014-edit-bookmark` | **Date**: 2026-04-01 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/014-edit-bookmark/spec.md`

## Summary

Add URL editing capability to the bookmark edit flow across all three interfaces (CLI, TUI, Vicinae extension). The CLI already supports editing title and tags; this feature extends `BookmarkPatch` with an optional URL field, adds duplicate URL detection for edits, recalculates the domain on URL change, extends the TUI edit panel with a URL input, and adds a new edit command and form to the Vicinae extension.

## Technical Context

**Language/Version**: Go 1.25.0 (CLI + TUI), TypeScript 5.x (Vicinae extension)
**Primary Dependencies**: charmbracelet/bubbletea + bubbles + lipgloss (TUI), Raycast API (extension)
**Storage**: SQLite via modernc.org/sqlite (WAL mode, FTS5) — no schema changes
**Testing**: `go test` for store and CLI logic; manual testing for TUI and extension
**Target Platform**: Linux, macOS (CLI/TUI); macOS (Vicinae/Raycast extension)
**Project Type**: CLI tool with TUI and browser extension
**Performance Goals**: N/A (single-user bookmark manager, trivial query load)
**Constraints**: Zero CGO, single binary
**Scale/Scope**: Single user, <10k bookmarks typical

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate                          | Status | Notes                                                  |
| ----------------------------- | ------ | ------------------------------------------------------ |
| No CGO                        | PASS   | No new dependencies; all existing deps are pure Go     |
| Single binary                 | PASS   | No external runtime introduced                         |
| Task management               | PASS   | Tasks will be created via Backlog CLI after `/speckit.tasks` |
| Backward-compatible migrations| PASS   | No schema changes at all                               |

**Post-Phase 1 re-check**: All gates still pass. No new dependencies, no schema changes, no abstractions added.

## Project Structure

### Documentation (this feature)

```text
specs/014-edit-bookmark/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   ├── cli-edit.md      # CLI edit command contract
│   └── vicinae-edit.md  # Extension edit contract
└── tasks.md             # Phase 2 output (via /speckit.tasks)
```

### Source Code (repository root)

```text
cmd/cairn/
├── main.go              # Bootstrap, command registry, main()
├── commands.go          # cmdXxx dispatch functions
├── run_bookmarks.go     # Bookmark subcommand logic — runEdit receives --url flag here
├── run_sync.go          # Sync subcommand logic + background sync
├── run_update.go        # Update subcommand logic
├── output.go            # Table + help rendering
├── util.go              # openStore, domainFromURL, fatalf
└── text/                # Embedded help text files

internal/
├── store/
│   └── bookmark.go      # Add URL to BookmarkPatch, update UpdateFields
└── model/
    ├── edit.go           # Add URL text input to EditModel
    └── app.go            # Update updateEdit to pass URL patch

vicinae-extension/src/
├── bm.ts                # Add bmEdit() function
├── bm-edit.tsx          # New: Edit form component
├── bm-list.tsx          # Add Edit action button
└── bm-search.tsx        # Add Edit action button
```

**Structure Decision**: No new packages or directories in Go code. One new file (`bm-edit.tsx`) in the extension following the existing `bm-add.tsx` pattern.

## Complexity Tracking

No constitution violations. Table not needed.
