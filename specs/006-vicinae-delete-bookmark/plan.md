# Implementation Plan: Delete Bookmarks from Vicinae Extension

**Branch**: `006-vicinae-delete-bookmark` | **Date**: 2026-03-11 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/006-vicinae-delete-bookmark/spec.md`

## Summary

Add a "Delete Bookmark" action to the Vicinae extension's list and search views. Uses a confirmation dialog (`confirmAlert`) before calling `cairn delete <id>` via the existing CLI wrapper pattern. Shows success/failure toast notifications and refreshes the list after deletion.

## Technical Context

**Language/Version**: TypeScript 5.9.2 (vicinae-extension only — no Go changes needed)
**Primary Dependencies**: @vicinae/api ^0.20.3 (existing), React (existing), Node.js child_process (existing)
**Storage**: N/A (delegates to cairn CLI which manages SQLite)
**Testing**: Manual testing via `vici develop` (no automated test framework in the extension)
**Target Platform**: Vicinae launcher (Linux)
**Project Type**: Launcher extension (TypeScript/React)
**Performance Goals**: N/A (single CLI call per delete)
**Constraints**: Must use `cairn delete <id>` CLI — no direct database access
**Scale/Scope**: 3 files modified, ~50 lines added

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate                            | Status | Notes                                                                          |
|---------------------------------|--------|--------------------------------------------------------------------------------|
| No CGO                          | PASS   | No Go changes. Extension is pure TypeScript.                                   |
| Single binary                   | PASS   | No changes to Go binary.                                                       |
| Task management                 | PASS   | Tasks will be created via Backlog CLI after `/speckit.tasks`.                   |
| Backward-compatible migrations  | PASS   | No schema changes.                                                             |

All gates pass. No violations.

## Project Structure

### Documentation (this feature)

```text
specs/006-vicinae-delete-bookmark/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (minimal — no new entities)
├── quickstart.md        # Phase 1 output
└── contracts/           # Phase 1 output (CLI contract reference)
```

### Source Code (repository root)

```text
vicinae-extension/
└── src/
    ├── bm.ts            # Add bmDelete() function
    ├── bm-list.tsx       # Add Delete action to BookmarkListItem, add list refresh
    └── bm-search.tsx     # Add Delete action to BookmarkListItem, add list refresh
```

**Structure Decision**: No new files needed. All changes are additions to existing files in the vicinae-extension package.
