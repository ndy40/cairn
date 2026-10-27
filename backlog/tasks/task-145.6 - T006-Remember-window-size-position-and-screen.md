---
id: TASK-145.6
title: T006 Remember window size, position and screen
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
labels:
  - desktop
  - foundation
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 114000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T006 in specs/018-native-desktop-window/tasks.md.

Extend `desktop/windowstate.go`: defaults 720×520, minimum 520×360; add optional `X, Y *int` and `Screen string` fields; add pure `restorePlan(state windowState, wayland bool, screens []screenRect) (w, h int, x, y int, hasPos bool)` that clamps size to the target screen work area, restores position only when not Wayland and the saved screen exists, and otherwise reports no position; keep old `{width,height}` files loading; extend `desktop/windowstate_test.go` (old file, missing screen, Wayland, clamp)
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Default size is 720x520 and the minimum is 520x360
- [ ] #2 Position is restored on macOS and X11 only when the saved screen exists; Wayland restores size only
- [ ] #3 Older desktop-state.json files with only width and height still load
<!-- AC:END -->
