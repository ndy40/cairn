# Tasks: Delete Bookmarks from Vicinae Extension

**Input**: Design documents from `/specs/006-vicinae-delete-bookmark/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

**Tests**: No automated tests requested. Manual testing via `vici develop`.

**Organization**: Tasks grouped by user story. US1 is the MVP.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: No setup needed — existing project, no new dependencies.

*(No tasks — the extension project and all dependencies already exist.)*

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Add the `bmDelete()` CLI wrapper function that both user stories depend on.

- [x] T001 Add `bmDelete(id: number)` function to `vicinae-extension/src/bm.ts` — call `cairn delete <id>` via `spawnSync`, return `{ exitCode: number; stderr: string }` following the `bmAdd()` pattern

**Checkpoint**: CLI delete wrapper ready — user story implementation can now begin.

---

## Phase 3: User Story 1 - Delete a Bookmark from the List (Priority: P1) 🎯 MVP

**Goal**: Users can delete a bookmark from the List Bookmarks view with confirmation dialog, success/error toast, and automatic list refresh.

**Independent Test**: Open List Bookmarks → select bookmark → action panel shows "Delete Bookmark" → confirm → bookmark removed, success toast shown, list refreshes. Cancel confirmation → bookmark remains.

### Implementation for User Story 1

- [x] T002 [US1] Update `BookmarkListItem` in `vicinae-extension/src/bm-list.tsx` to accept an `onDelete` callback prop of type `(bookmark: Bookmark) => void`
- [x] T003 [US1] Add "Delete Bookmark" action to the `ActionPanel` in `BookmarkListItem` in `vicinae-extension/src/bm-list.tsx` — position after "Copy URL", use `confirmAlert()` with `Alert.ActionStyle.Destructive`, call `bmDelete(bookmark.id)`, show `showToast()` with `Toast.Style.Success` or `Toast.Style.Failure`, then call `onDelete` callback on success
- [x] T004 [US1] Update `ListBookmarks` component in `vicinae-extension/src/bm-list.tsx` to define a `handleDelete` function that re-calls `bmList()` and updates state via `setBookmarks`, pass it as `onDelete` prop to each `BookmarkListItem`
- [x] T005 [US1] Add imports for `confirmAlert`, `Alert`, `showToast`, `Toast` from `@vicinae/api` and `bmDelete` from `./bm` in `vicinae-extension/src/bm-list.tsx`

**Checkpoint**: User Story 1 fully functional — delete from list view works end-to-end.

---

## Phase 4: User Story 2 - Delete a Bookmark from Search Results (Priority: P2)

**Goal**: Users can delete a bookmark from the Search Bookmarks view with the same confirmation/toast/refresh behavior.

**Independent Test**: Open Search Bookmarks → search for a bookmark → action panel shows "Delete Bookmark" → confirm → bookmark removed, search results refresh.

### Implementation for User Story 2

- [x] T006 [US2] Update `BookmarkListItem` in `vicinae-extension/src/bm-search.tsx` to accept an `onDelete` callback prop of type `(bookmark: Bookmark) => void`
- [x] T007 [US2] Add "Delete Bookmark" action to the `ActionPanel` in `BookmarkListItem` in `vicinae-extension/src/bm-search.tsx` — same pattern as T003 (confirmAlert, bmDelete, showToast, onDelete callback)
- [x] T008 [US2] Update `SearchBookmarks` component in `vicinae-extension/src/bm-search.tsx` to define a `handleDelete` function that re-calls `bmSearch(query)` (or `bmList()` if query is empty) and updates state, pass it as `onDelete` prop to each `BookmarkListItem`
- [x] T009 [US2] Add imports for `confirmAlert`, `Alert`, `showToast`, `Toast` from `@vicinae/api` and `bmDelete` from `./bm` in `vicinae-extension/src/bm-search.tsx`

**Checkpoint**: Both user stories functional — delete works in both list and search views.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Formatting and validation.

- [x] T010 Run `npx biome format --write src` in `vicinae-extension/` to format all modified files
- [x] T011 Run `vici lint` in `vicinae-extension/` to verify no lint errors
- [x] T012 Run quickstart.md manual validation — verify all test scenarios pass in `vici develop`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Foundational (Phase 2)**: No dependencies — can start immediately
- **User Story 1 (Phase 3)**: Depends on T001 (bmDelete wrapper)
- **User Story 2 (Phase 4)**: Depends on T001 (bmDelete wrapper). Independent of US1.
- **Polish (Phase 5)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Depends only on T001. No dependencies on US2.
- **User Story 2 (P2)**: Depends only on T001. No dependencies on US1. Can be implemented in parallel with US1.

### Parallel Opportunities

- T002–T005 (US1) and T006–T009 (US2) can run in parallel after T001 completes (different files)
- Within US1: T005 (imports) can run in parallel with T002 (prop change) — different sections of same file, but safer to do sequentially
- T010 and T011 can run in parallel

---

## Parallel Example: User Stories 1 & 2

```bash
# After T001 completes, launch both stories in parallel:
# Agent A: US1 tasks in bm-list.tsx (T002, T003, T004, T005)
# Agent B: US2 tasks in bm-search.tsx (T006, T007, T008, T009)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete T001: bmDelete wrapper
2. Complete T002–T005: List view delete
3. **STOP and VALIDATE**: Test delete from list view in `vici develop`
4. Ship MVP

### Incremental Delivery

1. T001 → bmDelete ready
2. T002–T005 → US1 complete → Test list delete (MVP!)
3. T006–T009 → US2 complete → Test search delete
4. T010–T012 → Polish → Final validation

---

## Notes

- Total tasks: 12
- US1 tasks: 4 (T002–T005)
- US2 tasks: 4 (T006–T009)
- Foundational: 1 (T001)
- Polish: 3 (T010–T012)
- US1 and US2 can run in parallel (different files)
- No automated tests — manual testing via `vici develop`
- Suggested MVP: US1 only (delete from list view)
