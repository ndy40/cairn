---
id: TASK-145.19
title: T019 Focus search on show and honour hide-after-open
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
labels:
  - desktop
  - us2
dependencies:
  - TASK-145.16
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 127000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T019 in specs/018-native-desktop-window/tasks.md.

In `desktop/frontend/src/main.ts`, on every `launcher:show` focus `#query` and call `select()`; keep arrow keys moving the selection while focus stays in the field; after `BookmarkService.Open` succeeds call `WindowService.AfterOpen()`; add `AfterOpen()` to `desktop/window_service.go` that hides the app (macOS) or window (Linux) only when `hideAfterOpen` is on
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Every time the window comes forward the search field is focused with its text selected
- [ ] #2 Opening a link keeps the window open unless Hide Cairn after opening a link is on
<!-- AC:END -->
