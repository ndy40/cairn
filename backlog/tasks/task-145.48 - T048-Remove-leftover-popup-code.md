---
id: TASK-145.48
title: T048 Remove leftover popup code
status: To Do
assignee: []
created_date: '2026-10-27T03:58:00Z'
labels:
  - desktop
  - polish
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 156000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T048 in specs/018-native-desktop-window/tasks.md.

Remove dead code and constants left from the popup in `desktop/launcher.go`, `desktop/main.go` and `desktop/frontend/src/main.ts` (old `viewAdd` view swap, `#hints-add`, focus-guard leftovers); regenerate bindings
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 No popup-only code, constants or markup remain; bindings are regenerated
<!-- AC:END -->
