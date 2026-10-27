---
id: TASK-145.25
title: T025 Start in the background at login
status: To Do
assignee: []
created_date: '2026-10-27T03:57:00Z'
labels:
  - desktop
  - us3
dependencies:
  - TASK-145.20
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 133000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T025 in specs/018-native-desktop-window/tasks.md.

In `desktop/main.go` and `desktop/launcher.go`, add the `--background` flag (keep `--hidden` as an alias) that starts without a window and pulls sync on start; make the Linux tray only appear when `TrayAvailable`; make `SetSetting("runInBackground", true)` return an error when no tray is available
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 cairn-desktop --background (and --hidden) starts with no window and pulls sync
- [ ] #2 On Linux the tray icon appears only when a tray is available, and Run in background cannot be turned on without one
<!-- AC:END -->
