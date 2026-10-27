---
id: TASK-145.28
title: T028 Style Add Bookmark as a sheet or dialog
status: To Do
assignee: []
created_date: '2026-10-27T03:58:00Z'
labels:
  - desktop
  - us4
dependencies:
  - TASK-145.27
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 136000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T028 in specs/018-native-desktop-window/tasks.md.

In `desktop/frontend/src/style.css`, style `#add-dialog` as a macOS sheet under the title band (`[data-os=darwin]`, slide-down, attached to the band, list visible behind) and as a centred libadwaita-style dialog with a dimmed `::backdrop` and Cancel/Add in its header on Linux
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 On macOS it appears as a sheet under the title band; on Linux as a centred dialog over a dimmed window
<!-- AC:END -->
