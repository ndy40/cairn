---
id: TASK-145.38
title: T038 Implement the keyboard map
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:34:00Z'
labels:
  - desktop
  - us5
dependencies:
  - TASK-145.31
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 146000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T038 in specs/018-native-desktop-window/tasks.md.

Implement the FR-028 keyboard map in `desktop/frontend/src/main.ts`: `Cmd/Ctrl+F` focus+select search, `Shift+Cmd/Ctrl+C` copy link, `Cmd/Ctrl+N` add, `Cmd/Ctrl+W` `CloseWindow()`, `Cmd/Ctrl+Q` `Quit()`, `Cmd/Ctrl+,` settings, `Cmd+/` / `Ctrl+?` shortcuts; on macOS skip keys the app menu already handles so actions do not fire twice
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 All shortcuts from FR-028 work on both platforms
- [ ] #2 On macOS no action fires twice when the menu bar also handles the key
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
ndw-menus, merged into 018-native-desktop-window. Code complete; awaiting manual QA (quickstart rows 14, 18-24; Story 5). Deviations: Settings is an in-window dialog on macOS; sync warning on macOS is a ⚠︎ label beside the menu bar icon; Wails has no Dock menu.
<!-- SECTION:NOTES:END -->
