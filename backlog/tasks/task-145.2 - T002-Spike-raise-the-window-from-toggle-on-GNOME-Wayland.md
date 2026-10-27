---
id: TASK-145.2
title: 'T002 Spike: raise the window from --toggle on GNOME Wayland'
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:08:00Z'
labels:
  - desktop
  - setup
dependencies:
  - TASK-145.1
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 110000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T002 in specs/018-native-desktop-window/tasks.md.

Spike the Wayland activation token (research R8): add `desktop/activation_linux.go` (cgo, `//go:build linux`, `#cgo pkg-config: gtk4`) with `applyActivationToken(w *application.WebviewWindow, token string)` that calls `gtk_window_set_startup_id` on the window's native `GtkWindow*` (via `w.NativeWindow()`), and `desktop/activation_other.go` (`//go:build !linux`) as a no-op; on GNOME Wayland, bind a custom shortcut to `cairn-desktop --toggle` and confirm the window is raised and focused; record the outcome (works / falls back) under R8 in `specs/018-native-desktop-window/research.md`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 On GNOME Wayland, a custom shortcut running cairn-desktop --toggle raises and focuses the window, or the fallback is documented
- [ ] #2 research.md R8 records the outcome and the manual steps used
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Code merged in 176f4f7 (desktop/activation_linux.go). NOT YET VERIFIED: manual GNOME Wayland check per research.md R8.
<!-- SECTION:NOTES:END -->
