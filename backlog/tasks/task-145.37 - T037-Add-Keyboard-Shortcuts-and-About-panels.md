---
id: TASK-145.37
title: T037 Add Keyboard Shortcuts and About panels
status: To Do
assignee: []
created_date: '2026-10-27T03:58:00Z'
labels:
  - desktop
  - us5
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 145000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T037 in specs/018-native-desktop-window/tasks.md.

Add Keyboard Shortcuts and About `<dialog>`s in `desktop/frontend/index.html` and `desktop/frontend/src/main.ts` (shortcuts list per FR-028 with platform-correct key names; About shows name and version from a new `WindowService.Version()` reading the build version); open on `ui:open` and `Cmd+/` / `Ctrl+?`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Cmd+/ or Ctrl+? lists the keyboard shortcuts with platform key names
- [ ] #2 About shows the app name and version
<!-- AC:END -->
