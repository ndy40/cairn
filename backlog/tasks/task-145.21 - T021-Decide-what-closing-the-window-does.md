---
id: TASK-145.21
title: T021 Decide what closing the window does
status: Done
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:18:00Z'
labels:
  - desktop
  - us3
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 129000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T021 in specs/018-native-desktop-window/tasks.md.

In `desktop/lifecycle.go`, add pure `closeAction(os string, runInBackground, trayAvailable bool) action` (macOS → hide window; Linux → hide only when both true, else quit) with table tests in `desktop/lifecycle_test.go`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The close decision function returns hide on macOS, hide on Linux only with background mode and a tray, and quit otherwise, and is unit tested
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
closeAction in desktop/lifecycle.go with table tests. ndw-lifecycle, merged into 018-native-desktop-window (0bedb24).
<!-- SECTION:FINAL_SUMMARY:END -->
