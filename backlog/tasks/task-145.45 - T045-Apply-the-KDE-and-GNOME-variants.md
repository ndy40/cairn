---
id: TASK-145.45
title: T045 Apply the KDE and GNOME variants
status: To Do
assignee: []
created_date: '2026-10-27T03:58:00Z'
labels:
  - desktop
  - us6
dependencies:
  - TASK-145.8
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 153000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T045 in specs/018-native-desktop-window/tasks.md.

Desktop variants in `desktop/frontend/src/style.css`: `[data-desktop=kde]` 4px radius, Breeze accent fallback `#3daee9`, labelled Add button, footer shown; `[data-desktop=gnome]` footer hidden; footer with key hints and result count kept on macOS; system-ui font stack and the sizes from the Components artboard (search 16/15/14px by platform)
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 KDE gets 4px radius, the Breeze accent fallback, a labelled Add button and the footer; GNOME has no footer
- [ ] #2 Type sizes follow the design per platform
<!-- AC:END -->
