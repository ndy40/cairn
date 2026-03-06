---
id: TASK-84
title: 'T006 [003] [US1] Wire edit flow in app.go'
status: Done
assignee: []
created_date: '2026-03-06 10:17'
updated_date: '2026-03-06 17:10'
labels:
  - feature-003
  - US1
  - tui
dependencies:
  - TASK-83
documentation:
  - specs/003-edit-bookmark-help/contracts/keyboard-shortcuts.md
  - specs/003-edit-bookmark-help/tasks.md
priority: high
ordinal: 13000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
In `internal/model/app.go` wire the complete edit flow:\n- Handle `editKey` press in `updateBrowse()`: initialise `m.editModel` with selected bookmark, set state to `StateEdit`\n- Add `updateEdit()`: route `Enter` (call `s.UpdateTags` with `m.editModel.Tags()`, then `loadBookmarks`) and `Esc` (return to `StateBrowse`)\n- Add `editView()` returning `m.editModel.View()`\n- Route `StateEdit` in top-level `Update()` and `View()`\n- `e` on empty list does nothing
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Pressing 'e' in browse mode opens the edit panel pre-filled with current tags
- [ ] #2 Pressing Enter saves tags and returns to browse with list reloaded
- [ ] #3 Pressing Esc discards changes and returns to browse
- [ ] #4 Pressing 'e' on an empty list does nothing
- [ ] #5 More than 3 tags: only first 3 saved
<!-- AC:END -->
