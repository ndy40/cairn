---
id: TASK-145.20
title: T020 Watch for a system tray on Linux
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:08:00Z'
labels:
  - desktop
  - us3
dependencies:
  - TASK-145.4
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 128000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T020 in specs/018-native-desktop-window/tasks.md.

In `desktop/platform.go`, add `watchTray(ctx, onChange func(bool))` using godbus: `NameHasOwner("org.kde.StatusNotifierWatcher")` on the session bus at start, then subscribe to `NameOwnerChanged` for that name; on change update `PlatformInfo.TrayAvailable` and emit `platform:changed`; no-op on macOS; a missing session bus means `false`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Tray availability reflects whether a StatusNotifier host is running and updates live
- [ ] #2 A missing session bus means no tray
<!-- AC:END -->
