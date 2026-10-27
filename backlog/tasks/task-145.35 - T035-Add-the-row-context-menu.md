---
id: TASK-145.35
title: T035 Add the row context menu
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:18:00Z'
labels:
  - desktop
  - us5
dependencies:
  - TASK-145.30
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 143000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T035 in specs/018-native-desktop-window/tasks.md.

In `desktop/frontend/index.html` and `desktop/frontend/src/main.ts`, add the row context menu (`contextmenu` event and Shift+F10 / Menu key) with Open Link, Copy Link and Pin/Unpin (label follows `pinned`), and handle `menu:action` (`open`, `copy`, `pin`, `find`) against the selected row
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Right-click, Shift+F10 or the Menu key on a row offers Open Link, Copy Link and Pin/Unpin
- [ ] #2 macOS File menu actions act on the selected row
<!-- AC:END -->
