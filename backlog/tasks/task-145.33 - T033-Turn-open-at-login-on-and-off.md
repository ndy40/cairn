---
id: TASK-145.33
title: T033 Turn open at login on and off
status: To Do
assignee: []
created_date: '2026-10-27T03:58:00Z'
labels:
  - desktop
  - us5
dependencies:
  - TASK-145.3
  - TASK-145.7
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 141000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T033 in specs/018-native-desktop-window/tasks.md.

In `desktop/window_service.go`, add `OpenAtLogin() (bool, error)` and `SetOpenAtLogin(enabled bool) error` using `app.Autostart` with identifier `io.github.ndy40.cairn.desktop` (macOS) / `cairn-desktop` (Linux) and arguments `--background` (plus the one-time plist check if T003 found it necessary); make `SetSetting("showInMenuBar", …)` show/hide the tray immediately
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Settings can turn open at login on and off, registering the app with --background
- [ ] #2 An existing installer-created login item shows as on
<!-- AC:END -->
