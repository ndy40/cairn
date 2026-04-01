# Tasks: Edit Bookmark

**Input**: Design documents from `/specs/014-edit-bookmark/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

**Tests**: Not explicitly requested. Test tasks omitted.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: No new project setup needed — this feature extends existing code. Phase intentionally empty.

**Checkpoint**: No setup required — proceed to foundational phase.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Extend the store layer to support URL editing. All three interfaces (CLI, TUI, extension) depend on this.

**CRITICAL**: No user story work can begin until this phase is complete.

- [x] T001 Add `URL *string` field to `BookmarkPatch` struct and add `ErrDuplicateURL` sentinel error in `internal/store/bookmark.go`
- [x] T002 Update `UpdateFields` method to handle URL changes: validate non-empty URL, check for duplicate URL excluding self (`WHERE url = ? AND id != ?`), recalculate domain via `extractDomain()`, add `url` and `domain` to the dynamic SET clause in `internal/store/bookmark.go`

**Checkpoint**: Foundation ready — `store.UpdateFields()` now accepts URL patches with duplicate detection and domain recalculation. User story implementation can begin.

---

## Phase 3: User Story 1 — Edit Bookmark URL via CLI (Priority: P1) MVP

**Goal**: Users can edit a bookmark's URL (in addition to title and tags) from the command line.

**Independent Test**: Run `cairn edit <id> --url=<new-url>` and verify the bookmark's URL and domain update while other fields remain unchanged. Verify duplicate URL rejection with exit code 1.

### Implementation for User Story 1

- [x] T003 [US1] Add `--url` flag to `cmdEdit` function, track `urlSet` boolean via `fs.Visit`, and pass URL and urlSet to `runEdit` in `cmd/cairn/main.go`
- [x] T004 [US1] Update `runEdit` function signature to accept `url string, urlSet bool` parameters, add URL validation (non-empty when set), populate `BookmarkPatch.URL`, handle `ErrDuplicateURL` with "Duplicate URL" stderr and exit code 1, and update success messages to include URL in `cmd/cairn/main.go`
- [x] T005 [US1] Update usage string in `cmdEdit` to include `--url` flag and update help text file `cmd/cairn/help-edit.txt` (if it exists)

**Checkpoint**: `cairn edit <id> --url=<url> --title=<title> --tags=<tags>` fully works. URL editing, duplicate detection, and domain recalculation all functional via CLI.

---

## Phase 4: User Story 2 — Edit Bookmark URL and Tags via TUI (Priority: P2)

**Goal**: Users can edit a bookmark's URL and tags from the TUI edit panel (currently tags-only).

**Independent Test**: Launch TUI, press `e` on a bookmark, modify the URL field, press Enter, and verify the bookmark updates in the list.

### Implementation for User Story 2

- [x] T006 [US2] Extend `EditModel` struct in `internal/model/edit.go` to add a `urlInput textinput.Model` field, an `activeField int` for Tab navigation between URL and tags inputs, and update `newEditModel` to initialise the URL input pre-filled with the bookmark's current URL
- [x] T007 [US2] Update `EditModel.Update` method in `internal/model/edit.go` to handle Tab key for switching focus between URL and tags fields, and route key messages to the active input
- [x] T008 [US2] Update `EditModel.View` method in `internal/model/edit.go` to render both URL and tags input fields with labels, showing which field is focused
- [x] T009 [US2] Add a `URL() string` method to `EditModel` in `internal/model/edit.go` that returns the trimmed URL input value
- [x] T010 [US2] Update `updateEdit` handler in `internal/model/app.go` to build a `BookmarkPatch` with both URL and tags from the EditModel, call `store.UpdateFields` instead of `store.UpdateTags`, and handle `ErrDuplicateURL` by displaying an error message in the TUI status bar

**Checkpoint**: TUI edit panel shows URL and tags fields, Tab switches between them, Enter saves both, Escape cancels. Duplicate URL errors shown inline.

---

## Phase 5: User Story 3 — Edit Bookmark from Vicinae Extension (Priority: P3)

**Goal**: Users can edit a bookmark's URL and tags from the Vicinae browser extension.

**Independent Test**: Open extension, select a bookmark, trigger Edit action, modify URL and/or tags, submit, and verify changes persist via `cairn list --json`.

### Implementation for User Story 3

- [x] T011 [P] [US3] Add `bmEdit(id: number, url?: string, tags?: string)` function to `vicinae-extension/src/bm.ts` that calls `cairn edit <id> [--url=<url>] [--tags=<tags>]` and invalidates list cache on success
- [x] T012 [P] [US3] Create edit form component in `vicinae-extension/src/bm-edit.tsx` with pre-filled URL and tags fields, submit handler calling `bmEdit()`, duplicate URL error display (exit code 1), and success toast with navigation back to list
- [x] T013 [US3] Add "Edit Bookmark" action button to bookmark list items in `vicinae-extension/src/bm-list.tsx` that opens the edit form with the selected bookmark's data
- [x] T014 [US3] Add "Edit Bookmark" action button to search result items in `vicinae-extension/src/bm-search.tsx` that opens the edit form with the selected bookmark's data

**Checkpoint**: Extension shows Edit action on bookmarks in both list and search views. Edit form pre-fills URL and tags, submits via CLI, and refreshes the list.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and cleanup across all interfaces.

- [x] T015 Verify sync propagation: edit a bookmark URL via CLI, run `cairn sync push`, confirm `pending_sync` entry recorded with operation `'update'`
- [x] T016 Run `quickstart.md` validation — execute all manual test steps to confirm end-to-end functionality across CLI, TUI, and extension

---

## Dependencies & Execution Order

### Phase Dependencies

- **Foundational (Phase 2)**: No dependencies — can start immediately
- **User Story 1 (Phase 3)**: Depends on Phase 2 (T001, T002)
- **User Story 2 (Phase 4)**: Depends on Phase 2 (T001, T002) — can run in parallel with US1
- **User Story 3 (Phase 5)**: Depends on Phase 2 (T001, T002) and US1 completion (CLI `--url` flag must exist for extension to call it)
- **Polish (Phase 6)**: Depends on all user stories complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Phase 2 — no dependencies on other stories
- **User Story 2 (P2)**: Can start after Phase 2 — no dependencies on other stories (uses store directly)
- **User Story 3 (P3)**: Depends on US1 (the extension calls `cairn edit --url` which is implemented in US1)

### Within Each User Story

- Store layer (Phase 2) before any interface tasks
- CLI flags before CLI logic (T003 before T004)
- TUI model before TUI handler (T006-T009 before T010)
- Extension API wrapper and form before list/search integration (T011-T012 before T013-T014)

### Parallel Opportunities

- T001 has no internal parallelism (single file)
- T011 and T012 can run in parallel (different files)
- T013 and T014 can run in parallel (different files)
- US1 and US2 can run in parallel after Phase 2 (Go CLI vs Go TUI — different files)

---

## Parallel Example: User Story 3

```bash
# Launch extension API and form in parallel (different files):
Task: "Add bmEdit() function in vicinae-extension/src/bm.ts"
Task: "Create edit form component in vicinae-extension/src/bm-edit.tsx"

# Then launch list and search integration in parallel (different files):
Task: "Add Edit action to vicinae-extension/src/bm-list.tsx"
Task: "Add Edit action to vicinae-extension/src/bm-search.tsx"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 2: Foundational (T001-T002)
2. Complete Phase 3: User Story 1 (T003-T005)
3. **STOP and VALIDATE**: Test `cairn edit <id> --url=<new-url>` independently
4. URL editing works from CLI — ship it

### Incremental Delivery

1. Phase 2 → Foundation ready (BookmarkPatch supports URL)
2. Add User Story 1 → CLI edit with URL → Validate (MVP!)
3. Add User Story 2 → TUI edit panel with URL → Validate
4. Add User Story 3 → Vicinae extension edit → Validate
5. Phase 6 → Polish and sync verification

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- No new Go dependencies required
- No schema migrations — URL and domain columns already writable
- Sync propagation is handled automatically by existing `UpdateFields` → `pending_sync` logic
- Commit after each task or logical group
