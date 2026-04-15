# Implementation Plan: Edit Bookmark Title

**Branch**: `017-edit-bookmark-title` | **Date**: 2026-04-15 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/017-edit-bookmark-title/spec.md`

## Summary

Add a title text-input field to two existing editing surfaces — the TUI edit panel (`internal/model/edit.go`) and the Vicinae extension edit form (`vicinae-extension/src/bm-edit.tsx` + `bm.ts`) — so users can manually correct bookmark titles that were set to fallbacks by bot-protection pages. The CLI backend (`cairn edit --title`) and store layer (`UpdateFields`) already handle title writes; no schema or CLI changes are needed.

---

## Technical Context

**Language/Version**: Go 1.25.0 (TUI + CLI), TypeScript 5.x (Vicinae extension)
**Primary Dependencies**: charmbracelet/bubbletea + bubbles + lipgloss (TUI, existing), `@vicinae/api` (extension, existing) — **no new dependencies**
**Storage**: modernc.org/sqlite — no schema changes; `Bookmark.Title` column already exists and is writable via `UpdateFields`
**Testing**: `go test` (Go); manual smoke-test via quickstart.md for both surfaces
**Target Platform**: Linux/macOS/Windows terminal (TUI), Vicinae launcher (extension)
**Project Type**: CLI + TUI (Go) + launcher extension (TypeScript/React)
**Performance Goals**: Immediate field response; save completes within existing `UpdateFields` latency (<50ms local SQLite write)
**Constraints**: `CGO_ENABLED=0`; pure-Go; no new external libraries
**Scale/Scope**: Single-user local application; only affects the edit panel on two surfaces

---

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate | Status | Notes |
|------|--------|-------|
| No CGO | ✅ PASS | No new dependencies; existing stack is pure Go |
| Single binary | ✅ PASS | No new runtime or external process introduced |
| Task management | ✅ PASS | Tasks will be created via Backlog CLI after `/speckit.tasks` |
| Backward-compatible migrations | ✅ PASS | No schema changes required |

**Post-design re-check**: All gates still pass. Changes are confined to UI layer only.

---

## Project Structure

### Documentation (this feature)

```text
specs/017-edit-bookmark-title/
├── plan.md              ← this file
├── research.md          ← Phase 0 output
├── data-model.md        ← Phase 1 output
├── quickstart.md        ← Phase 1 output
├── contracts/
│   ├── tui-edit-panel.md
│   └── extension-edit-form.md
└── tasks.md             ← Phase 2 output (/speckit.tasks — NOT created here)
```

### Source Code (affected files only)

```text
internal/
  model/
    edit.go                    ← add titleInput textinput.Model + tab navigation + validation

vicinae-extension/
  src/
    bm.ts                      ← add optional title param to bmEdit()
    bm-edit.tsx                ← add Form.TextField for title, wire validation + change detection
```

**Structure Decision**: Changes are isolated to the model layer of the TUI and the form component + helper function of the extension. No new files needed in the Go codebase; no new files needed in the extension. All data persistence routes through the existing `cairn edit --title` CLI path.

---

## Complexity Tracking

No constitution violations. Table omitted.
