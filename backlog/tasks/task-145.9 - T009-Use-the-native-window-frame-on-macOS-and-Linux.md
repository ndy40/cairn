---
id: TASK-145.9
title: T009 Use the native window frame on macOS and Linux
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
labels:
  - desktop
  - us1
dependencies:
  - TASK-145.8
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 117000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T009 in specs/018-native-desktop-window/tasks.md.

In `desktop/launcher.go` `NewLauncher`, remove `Frameless`, `AlwaysOnTop` and `Windows.HiddenOnTaskbar`; set `MinWidth: 520`, `MinHeight: 360`, `Title: "Cairn"`; add `Mac: application.MacWindow{TitleBar: application.MacTitleBarHiddenInsetUnified, Backdrop: application.MacBackdropTranslucent}` and leave `Linux.Menu` unset so GTK draws its own title bar with no menu bar
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 macOS shows native traffic lights inset in a unified title band; Linux shows GTK's own title bar titled Cairn
- [ ] #2 The window is no longer always on top or hidden from the taskbar, and cannot shrink below 520x360
<!-- AC:END -->
