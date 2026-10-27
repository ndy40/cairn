---
id: TASK-145.23
title: T023 Apply close and quit rules
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:18:00Z'
labels:
  - desktop
  - us3
dependencies:
  - TASK-145.21
  - TASK-145.22
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 131000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T023 in specs/018-native-desktop-window/tasks.md.

In `desktop/launcher.go`, change the `WindowClosing` hook to use `closeAction` (hide + cancel, or cancel and call the quit path); in `desktop/window_service.go` add `CloseWindow()` (same path) and `Quit()` which, if `Syncer.Running()`, emits `status` "Finishing sync…", then calls `app.Quit()`; in `desktop/main.go` `OnShutdown` call `syncer.Wait(5*time.Second)` before saving window state
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Closing the window follows the close rule on each platform; Cmd/Ctrl+W behaves like the close button
- [ ] #2 Quit shows Finishing sync... during a sync and quits within 5 seconds, without a confirmation
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
ndw-lifecycle, merged into 018-native-desktop-window (0bedb24). Code complete; awaiting manual QA (quickstart rows 11-14, 19). Quit waits for sync in a goroutine before app.Quit so the status can show; menus must call WindowService.Quit (round 3).
<!-- SECTION:NOTES:END -->
