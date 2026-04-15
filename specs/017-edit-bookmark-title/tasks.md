# Tasks: Edit Bookmark Title (017)

**Input**: `specs/017-edit-bookmark-title/`
**Prerequisites**: plan.md ✓, spec.md ✓, research.md ✓, data-model.md ✓, contracts/ ✓

**No new dependencies.** Changes touch exactly 3 source files:
- `internal/model/edit.go` (TUI model)
- `internal/model/app.go` (TUI save handler)
- `vicinae-extension/src/bm.ts` (CLI helper)
- `vicinae-extension/src/bm-edit.tsx` (extension form)

US1 (TUI) and US2 (Extension) are fully independent — they can be implemented in parallel by different agents.

---

## Phase 1: Setup

- [ ] T001 Run `go build ./... && go vet ./...` to confirm a clean baseline in the repo root

---

## Phase 2: User Story 1 — Title Field in the TUI (Priority: P1) 🎯 MVP

**Goal**: The TUI edit panel exposes an editable title field (first in tab order), validates it on save, and persists the new title.

**Independent Test**: Open TUI → select a bookmark → press edit key → edit title → Enter → list shows updated title. Empty title rejected with inline error.

### Implementation for User Story 1

- [ ] T002 [US1] Rewrite `EditModel` in `internal/model/edit.go`:
  - Rename field constants: `editFieldTitle = 0`, `editFieldURL = 1`, `editFieldTags = 2`
  - Add `titleInput textinput.Model` (CharLimit 500, Width 58, pre-populated with `b.Title`, focused on open) and `titleErr string` to the struct
  - Update `newEditModel()`: initialise `titleInput` as described; `urlInput` and `tagsInput` start blurred (Title is focused)
  - Update `Update()`: Tab/Shift+Tab cycles across all 3 fields; dispatch key events to whichever `textinput` is active
  - Update `View()`: replace static bold title header with labelled `titleInput.View()` + error line when `titleErr != ""` (rendered in red/warning colour)
  - Add `Title() string` method returning `strings.TrimSpace(titleInput.Value())`

- [ ] T003 [US1] Update `updateEdit()` in `internal/model/app.go` to wire validation and the title patch:
  - On `"enter"`: call `a.editModel.Title()`; if empty, set `a.editModel.titleErr = "Title cannot be empty"` and return without state transition
  - If title is valid: clear `titleErr`, set `a.state = StateBrowse`, build `BookmarkPatch`
  - Include title in patch: if `title != origTitle`, set `patch.Title = &title`
  - Preserve existing URL-change + auto-fetch behaviour: if URL changed AND title was NOT manually edited (title == origTitle), auto-fetch still sets `patch.Title`; if title WAS manually edited, skip auto-fetch (user's value wins)
  - Store `origTitle` alongside `editOrigURL` when entering edit state (line ~254 in app.go)

**Checkpoint**: `go build ./... && go vet ./...` passes. TUI edit panel shows title as first editable field; empty-title save is rejected; valid edit persists.

---

## Phase 3: User Story 2 — Title Field in the Vicinae Extension (Priority: P2)

**Goal**: The extension edit form exposes a pre-populated title field, validates it, detects changes, and passes `--title` to the CLI only when the value changed.

**Independent Test**: Open extension → edit form → change title → submit → `cairn list` shows updated title. Empty title blocked with field error. Unchanged title does not trigger a CLI call.

### Implementation for User Story 2

- [ ] T004 [P] [US2] Add optional `title?: string` parameter to `bmEdit()` in `vicinae-extension/src/bm.ts`:
  - Append `"--title", title` to CLI args when `title !== undefined && title.trim() !== ""`
  - Parameter position: after `tags` — signature becomes `bmEdit(id, url?, tags?, title?)`

- [ ] T005 [US2] Update `EditBookmarkForm` in `vicinae-extension/src/bm-edit.tsx`:
  - Add `titleValue` / `setTitleValue` state initialised from `bookmark.Title`
  - Add `titleError` / `setTitleError` state
  - Add `validateTitle(value: string): boolean` — returns false and sets error if trimmed value is empty
  - Add `Form.TextField` for title above the URL field: `id="title"`, `title="Title"`, `maxLength={500}`, pre-populated, `error={titleError}`, `onBlur` validates
  - In `handleSubmit`: validate title first (return early on empty); compute `titleChanged = values.title.trim() !== bookmark.Title`; update the no-op guard to include `titleChanged`; pass `titleChanged ? values.title.trim() : undefined` as the `title` argument to `bmEdit()`

**Checkpoint**: Extension edit form shows title field; empty title shows inline error; title-only change calls CLI with `--title`; unchanged title does not call CLI.

---

## Phase 4: Polish & Cross-Cutting Concerns

- [ ] T006 [P] Run `go build ./... && go vet ./...` and confirm clean build
- [ ] T007 [P] Run `go test ./...` and confirm no regressions
- [ ] T008 [P] Run quickstart.md smoke tests manually for both surfaces (TUI + extension)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies.
- **Phase 2 (US1 — TUI)**: Depends on Phase 1. T003 depends on T002 (needs `Title()` method and `titleErr` field).
- **Phase 3 (US2 — Extension)**: Depends on Phase 1 only. Fully independent of Phase 2 — different language, different files.
- **Phase 4 (Polish)**: Depends on Phases 2 and 3.

### Within User Story 1

```
T002 (edit.go — model + Title() method)
  → T003 (app.go — validation + patch assembly)
```

### Within User Story 2

```
T004 (bm.ts — bmEdit() signature)
  → T005 (bm-edit.tsx — form field + change detection)
```

### Parallel Opportunities

- T002 and T004 touch different languages and files — can be implemented in parallel by two agents.
- Once T002 is done, T003 can start while T004→T005 continues in parallel.
- T006, T007, T008 can all run in parallel.

---

## Parallel Example: Both User Stories Together

```bash
# Agent A — TUI:
Task T002: Rewrite EditModel in internal/model/edit.go
Task T003: Update updateEdit() in internal/model/app.go

# Agent B — Extension (can start at same time as Agent A):
Task T004: Add title param to bmEdit() in vicinae-extension/src/bm.ts
Task T005: Add title field to EditBookmarkForm in vicinae-extension/src/bm-edit.tsx
```

---

## Implementation Strategy

### MVP First (US1 — TUI only)

1. T001 — baseline check.
2. T002 → T003 — TUI title field, validation, patch assembly.
3. **STOP and VALIDATE** via quickstart.md TUI tests.
4. Ship US1 independently — users can already fix titles in the terminal.

### Full Feature (US1 + US2)

5. T004 → T005 — extension title field (can overlap with US1 or follow after).
6. T006 + T007 + T008 — polish and verify.

---

## Notes

- Total tasks: **8** (T001–T008)
- US1 tasks: 2 (T002–T003)
- US2 tasks: 2 (T004–T005)
- Parallel opportunities: T002‖T004, T006‖T007‖T008
- Files changed: `internal/model/edit.go`, `internal/model/app.go`, `vicinae-extension/src/bm.ts`, `vicinae-extension/src/bm-edit.tsx`
- No schema changes, no new dependencies, no new files
