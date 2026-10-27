---
id: TASK-145.11
title: T011 Restore and save window geometry
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
labels:
  - desktop
  - us1
dependencies:
  - TASK-145.6
  - TASK-145.9
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 119000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T011 in specs/018-native-desktop-window/tasks.md.

In `desktop/launcher.go`, apply `restorePlan` (T006) when creating the window using `app.Screen` data and `PlatformInfo.Wayland`; centre only when there is no saved state; save size, position and screen (debounced, existing `saveSizeDelay`) on `WindowDidResize` and `WindowDidMove`, skipping position on Wayland; stop calling `window.Center()` in `Show`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 After resizing, moving and relaunching, the window reopens at the same size and, outside Wayland, the same position
- [ ] #2 The window is centred only on first launch
<!-- AC:END -->
