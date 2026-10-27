---
id: TASK-145.18
title: T018 Report global shortcut registration status
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
labels:
  - desktop
  - us2
dependencies:
  - TASK-145.16
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 126000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T018 in specs/018-native-desktop-window/tasks.md.

In `desktop/launcher.go` `RegisterHotkey`, store a `ShortcutStatus{Hotkey, Registered, Error, ToggleCommand: "cairn-desktop --toggle"}` (per `data-model.md`) instead of only logging; the hotkey calls `Toggle("")`; expose `ShortcutStatus()` on `WindowService` in `desktop/window_service.go`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 When the shortcut cannot be registered the app keeps running and the failure reason is available to Settings
<!-- AC:END -->
