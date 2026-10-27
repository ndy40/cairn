---
id: TASK-145.32
title: T032 Rebuild the tray and menu bar icon menu
status: To Do
assignee: []
created_date: '2026-10-27T03:58:00Z'
labels:
  - desktop
  - us5
dependencies:
  - TASK-145.25
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 140000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T032 in specs/018-native-desktop-window/tasks.md.

Move tray setup from `desktop/launcher.go` into `desktop/menus.go` `buildTray`: items Search Bookmarks…, Add Bookmark…, a disabled sync status line, Settings…, Quit Cairn (via `WindowService.Quit`); macOS uses the template icon and respects `showInMenuBar` (`tray.Show()/Hide()`); add a Dock menu with New Bookmark if Wails exposes one, otherwise note it in `research.md`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The tray/menu bar menu offers Search Bookmarks, Add Bookmark, sync status, Settings and Quit
- [ ] #2 On macOS the icon can be hidden from Settings
<!-- AC:END -->
