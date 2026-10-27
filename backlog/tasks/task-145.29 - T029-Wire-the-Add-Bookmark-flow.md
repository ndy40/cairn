---
id: TASK-145.29
title: T029 Wire the Add Bookmark flow
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:08:00Z'
labels:
  - desktop
  - us4
dependencies:
  - TASK-145.26
  - TASK-145.28
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 137000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T029 in specs/018-native-desktop-window/tasks.md.

In `desktop/frontend/src/main.ts`, replace the add view swap with `openAdd(prefill?)` using `showModal()`: prefill from `ClipboardURL()` and select it, or from `+text` typed in search with the cursor at the end; Enter submits, Esc/Cancel closes and refocuses search; on success close, refresh, select the new row and set status "Added. Fetching title…"; on `ErrAlreadyBookmarked` close, `FindByURL`, select that row and show an "Already bookmarked" toast (add a minimal `showToast(text)` helper, styled in US6); wire the Add button, `Ctrl/Cmd+N` and `launcher:show "add"`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The URL is prefilled and selected from the clipboard, or from +text typed in search
- [ ] #2 After adding, the dialog closes, the new row is selected and Added. Fetching title... is shown
- [ ] #3 A duplicate selects the existing row and shows an Already bookmarked toast
- [ ] #4 Esc or Cancel closes the dialog and refocuses search
<!-- AC:END -->
