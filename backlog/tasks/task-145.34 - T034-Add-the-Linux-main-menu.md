---
id: TASK-145.34
title: T034 Add the Linux main menu
status: To Do
assignee: []
created_date: '2026-10-27T03:58:00Z'
labels:
  - desktop
  - us5
dependencies:
  - TASK-145.12
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 142000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T034 in specs/018-native-desktop-window/tasks.md.

In `desktop/frontend/index.html` and `desktop/frontend/src/main.ts`, add the Linux main menu popover for `#menu-btn` (`role="menu"`, items Add Bookmark…, Sync Now, Preferences, Keyboard Shortcuts, About Cairn, Quit) with arrow-key navigation and Esc to close
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 On Linux the main menu button opens Add Bookmark, Sync Now, Preferences, Keyboard Shortcuts, About Cairn and Quit, usable by keyboard
<!-- AC:END -->
