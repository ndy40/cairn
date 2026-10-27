---
id: TASK-145.36
title: T036 Add the Settings panel
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:34:00Z'
labels:
  - desktop
  - us5
dependencies:
  - TASK-145.33
  - TASK-145.18
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 144000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T036 in specs/018-native-desktop-window/tasks.md.

Add a Settings `<dialog id="settings-dialog">` in `desktop/frontend/index.html` and its logic in `desktop/frontend/src/main.ts`: switches (`role="switch"`) for Open Cairn at login, Show in menu bar (macOS only), Run in background (Linux only; disabled with "Requires a system tray…" when `!trayAvailable`), Hide Cairn after opening a link; a shortcut row that shows the accelerator when registered, or the error and the `cairn-desktop --toggle` command with a Copy Command button when not; the closing note about Quit on macOS; open it on `ui:open "settings"` and `Cmd/Ctrl+,`; style per the Settings artboard in `desktop/frontend/src/style.css`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Settings offers open at login, show in menu bar (macOS), run in background (Linux, disabled without a tray) and hide after opening a link, applied without restart
- [ ] #2 The shortcut row shows the shortcut, or the failure and a Copy Command button for cairn-desktop --toggle
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
ndw-menus, merged into 018-native-desktop-window. Code complete; awaiting manual QA (quickstart rows 14, 18-24; Story 5). Deviations: Settings is an in-window dialog on macOS; sync warning on macOS is a ⚠︎ label beside the menu bar icon; Wails has no Dock menu.
<!-- SECTION:NOTES:END -->
