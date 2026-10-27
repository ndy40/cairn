---
id: TASK-145.31
title: T031 Build the macOS menu bar
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:18:00Z'
labels:
  - desktop
  - us5
dependencies:
  - TASK-145.30
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 139000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T031 in specs/018-native-desktop-window/tasks.md.

Create `desktop/menus.go` with `buildAppMenu(app, launcher, svc)` (darwin only): Cairn menu (About Cairn, Settings… `CmdOrCtrl+,`, Services role, Hide/Hide Others/Show All roles, Quit `CmdOrCtrl+Q`), File (New Bookmark… `CmdOrCtrl+N`, Open Link, Copy Link `Shift+CmdOrCtrl+C`, Pin, Sync Now, Close Window `CmdOrCtrl+W`), Edit role plus Find Bookmarks `CmdOrCtrl+F`, Window role, Help (Cairn Help → docs URL, Keyboard Shortcuts `CmdOrCtrl+/`); items emit `menu:action` / `ui:open` or call `WindowService` methods; set it with `app.Menu.Set` from `desktop/main.go`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 macOS shows Cairn, File, Edit, Window and Help menus with the items and shortcuts from the spec, using system roles for Edit and Window
<!-- AC:END -->
