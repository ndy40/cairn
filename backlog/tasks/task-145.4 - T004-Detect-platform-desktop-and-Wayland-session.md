---
id: TASK-145.4
title: T004 Detect platform, desktop and Wayland session
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
labels:
  - desktop
  - foundation
dependencies:
  - TASK-145.1
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 112000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T004 in specs/018-native-desktop-window/tasks.md.

Create `desktop/platform.go` with `type PlatformInfo struct{ OS, Desktop string; Wayland, TrayAvailable bool; AccentColor string }` (JSON tags per `data-model.md`) and a pure `detectPlatform(goos string, getenv func(string) string) PlatformInfo` mapping `XDG_CURRENT_DESKTOP` (case-insensitive, colon-separated) to `gnome`/`kde`/`other` and `XDG_SESSION_TYPE=wayland` to `Wayland`; macOS sets `TrayAvailable=true`; add table tests in `desktop/platform_test.go`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 PlatformInfo reports os, desktop (gnome/kde/other), wayland and tray availability from the environment
- [ ] #2 Table tests cover GNOME, KDE, unknown desktops and Wayland/X11
<!-- AC:END -->
