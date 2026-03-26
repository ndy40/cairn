---
id: TASK-142
title: Add "Toggle Pin" action to Vicinae extension list view
status: Done
assignee:
  - '@claude'
created_date: '2026-03-26 07:59'
updated_date: '2026-03-26 08:06'
labels:
  - vicinae-extension
  - cli
  - go
  - typescript
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The Vicinae extension list view displays a 📌 indicator for pinned bookmarks but provides no way to pin or unpin a bookmark from within the extension. Users must open the TUI to change a bookmark's pin state.

This task adds a "Toggle Pin" action to the list view. It requires:
- A new `cairn pin <id>` CLI subcommand (Go) that toggles `is_permanent` in SQLite.
- A `TogglePin(id int64) error` store method (Go).
- A `bmPin(id)` helper function in `bm.ts` (TypeScript).
- A "Pin Bookmark" / "Unpin Bookmark" action in `bm-list.tsx`'s `ActionPanel`.

No schema changes are required — `is_permanent` already exists from feature 002.

**Spec**: `specs/011-vicinae-pin-bookmark/spec.md`
**Plan**: `specs/011-vicinae-pin-bookmark/plan.md`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Invoking 'Toggle Pin' on an unpinned bookmark sets IsPermanent=true; 📌 appears on the item after list refresh
- [x] #2 Invoking 'Toggle Pin' on a pinned bookmark sets IsPermanent=false; 📌 is removed after list refresh
- [x] #3 Pin state persists: subsequent `cairn list --json` reflects the updated IsPermanent value
- [x] #4 A failure toast is shown when `cairn pin` exits non-zero; list state is unchanged
- [x] #5 `cairn pin <id>` exits 0 on success, 1 if bookmark not found, 3 on unexpected error
- [x] #6 `go build ./... && go vet ./...` pass with no new warnings
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add `case "pin"` + `runPin()` to cmd/cairn/main.go — reuse GetByID + SetPermanent
2. Add `bmPin(id)` export to vicinae-extension/src/bm.ts — same pattern as bmDelete
3. Add `onPin` prop + Pin/Unpin action to BookmarkListItem in bm-list.tsx
4. Add `handlePin` callback to ListBookmarks and pass to each BookmarkListItem
5. Run go build ./... && go vet ./... to verify
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added "Toggle Pin" action to the Vicinae extension list view.

Changes:
- cmd/cairn/main.go: Added `case "pin"` + `runPin()` — fetches bookmark via GetByID, inverts IsPermanent, calls SetPermanent. Exits 0 on success, 1 if not found, 3 on error.
- vicinae-extension/src/bm.ts: Added `bmPin(id)` export — same cache-invalidation pattern as bmDelete.
- vicinae-extension/src/bm-list.tsx: Added `onPin` prop to BookmarkListItem, "Pin/Unpin Bookmark" action in ActionPanel, and `handlePin` callback in ListBookmarks.

No schema changes. No new dependencies. `go build ./... && go vet ./...` pass.
<!-- SECTION:FINAL_SUMMARY:END -->
