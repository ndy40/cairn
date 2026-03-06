# Implementation Plan: Tags, Pinning, Archive & Startup Checks

**Branch**: `002-tags-pinning-archive` | **Date**: 2026-03-06 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-tags-pinning-archive/spec.md`

---

## Summary

Extends the existing TUI bookmark manager with: (1) a startup prerequisite check that detects the display environment (Wayland/X11) and verifies the clipboard tool is installed; (2) tags (up to 3 per bookmark, lowercase, deduplicated, 32-char max) with in-memory filtering composable with the existing fuzzy search; (3) last-visited timestamp recorded each time a bookmark is opened; (4) automatic startup archiving of bookmarks unvisited for 183+ days; and (5) a permanent flag that exempts bookmarks from archiving. This is a pure extension of the existing codebase — schema migration v2 adds new columns; no breaking changes to existing behaviour.

---

## Technical Context

**Language/Version**: Go 1.22+ (unchanged from feature 001)

**Primary Dependencies** (unchanged from feature 001):
- `charmbracelet/bubbletea` + `charmbracelet/bubbles` + `charmbracelet/lipgloss` — TUI framework
- `modernc.org/sqlite` — Pure Go embedded SQLite with FTS5 and JSON functions
- `sahilm/fuzzy` — Fuzzy search scoring
- `atotto/clipboard` + `wl-paste` (system binary, optional) — Clipboard read

**No new library dependencies** are required. All new functionality is built on the existing dependency set.

**Storage**: SQLite migration v2 adds five new columns to the `bookmarks` table:
- `tags` TEXT (JSON array, default `'[]'`)
- `last_visited_at` TEXT (nullable ISO-8601 datetime)
- `is_permanent` INTEGER (0/1, default 0)
- `is_archived` INTEGER (0/1, default 0)
- `archived_at` TEXT (nullable ISO-8601 datetime)

**Testing**: `go test ./...` (unchanged)

**Target Platform**: Linux (primary), macOS, Windows (unchanged)

**Project Type**: CLI/TUI application (unchanged)

**Performance Goals**:
- Tag filtering updates list within 200ms for ≤1,000 bookmarks (per SC-003) — in-memory filter, no DB round-trip
- Startup archive check completes within 2 seconds for ≤1,000 bookmarks (per SC-005)
- Prerequisite check displays result within 1 second of launch (per SC-001)

**Constraints**:
- Zero new CGO dependencies
- Backward-compatible schema migration (existing bookmarks gain new columns with safe defaults)
- Archive check is silent on zero eligible bookmarks; shows count only when ≥1 bookmark is archived

**Scale/Scope**: Single user, local, ≤1,000 bookmarks (unchanged)

---

## Constitution Check

No `.specify/memory/constitution.md` exists; same lightweight gates as feature 001 apply.

| Gate | Status | Notes |
|------|--------|-------|
| No CGO dependencies | PASS | No new dependencies added |
| Single binary deployment | PASS | No new external runtime required |
| No unnecessary abstractions | PASS | New `internal/display/` package has a single clear responsibility |
| No external runtime dependencies | PASS | wl-paste is checked for presence, not bundled |
| Backward-compatible migration | PASS | ALTER TABLE with DEFAULT values; existing data unaffected |

No violations. No Complexity Tracking table required.

---

## Project Structure

### Documentation (this feature)

```text
specs/002-tags-pinning-archive/
├── spec.md              # Feature specification (/speckit.specify output)
├── plan.md              # This file (/speckit.plan output)
├── research.md          # Phase 0 research findings
├── data-model.md        # Extended entity and schema migration v2
├── quickstart.md        # Updated developer guide for new features
├── contracts/
│   ├── keyboard-shortcuts.md   # Updated TUI keyboard contract (new shortcuts)
│   └── startup-check.md        # Startup prerequisite check behaviour contract
├── checklists/
│   └── requirements.md  # Spec quality checklist (already complete)
└── tasks.md             # Phase 2 output (/speckit.tasks — NOT created here)
```

### Source Code (additions to existing structure)

```text
bookmark-manager/
├── cmd/
│   └── bm/
│       └── main.go              # MODIFIED: add display/prereq check before TUI launch; archive check on startup
├── internal/
│   ├── display/
│   │   └── check.go             # NEW: detect Wayland/X11, verify clipboard tool, return check result
│   ├── model/
│   │   ├── app.go               # MODIFIED: StateArchive added; tag filter state; archive count footer msg
│   │   ├── browse.go            # MODIFIED: show tags, permanent indicator; 'p' toggle; 'a' open archive
│   │   ├── archive.go           # NEW: archive list model (mirrors browse.go pattern)
│   │   ├── search.go            # UNCHANGED
│   │   └── add.go               # MODIFIED: tag input field added
│   ├── store/
│   │   ├── store.go             # MODIFIED: migration v2 added to migrations slice
│   │   ├── bookmark.go          # MODIFIED: Insert/List updated for new fields; new methods below
│   │   ├── archive.go           # NEW: ArchiveStale(), ListArchived(), RestoreByID(), UpdateLastVisited()
│   │   └── search.go            # UNCHANGED
│   ├── fetcher/
│   │   └── fetcher.go           # UNCHANGED
│   ├── search/
│   │   └── fuzzy.go             # UNCHANGED (tag filter is applied before fuzzy, in app.go)
│   └── clipboard/
│       └── clipboard.go         # UNCHANGED
├── go.mod                       # UNCHANGED
└── go.sum                       # UNCHANGED
```

**Structure Decision**: Minimal extension of the existing single-project layout. Two new files (`internal/display/check.go`, `internal/model/archive.go`, `internal/store/archive.go`), targeted modifications to existing files. No new packages except `internal/display/`.

---

## Phase 0 Output

- [x] research.md — All technical decisions resolved

## Phase 1 Output

- [x] data-model.md — Extended Bookmark entity, migration v2 schema
- [x] contracts/keyboard-shortcuts.md — Updated keyboard contract with new shortcuts
- [x] contracts/startup-check.md — Prerequisite check behaviour and exit codes
- [x] quickstart.md — Updated developer guide
- [ ] tasks.md — Generated by `/speckit.tasks` (next step)
