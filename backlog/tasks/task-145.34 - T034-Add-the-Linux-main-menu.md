---
id: TASK-145.34
title: T034 Add the Linux main menu
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:34:00Z'
labels:
  - desktop
  - us5
dependencies:
  - TASK-145.12
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 142000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T034 in specs/018-native-desktop-window/tasks.md.

In `desktop/frontend/index.html` and `desktop/frontend/src/main.ts`, add the Linux main menu popover for `#menu-btn` (`role="menu"`, items Add Bookmark…, Sync Now, Preferences, Keyboard Shortcuts, About Cairn, Quit) with arrow-key navigation and Esc to close
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 On Linux the main menu button opens Add Bookmark, Sync Now, Preferences, Keyboard Shortcuts, About Cairn and Quit, usable by keyboard
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
ndw-menus, merged into 018-native-desktop-window. Code complete; awaiting manual QA (quickstart rows 14, 18-24; Story 5). Deviations: Settings is an in-window dialog on macOS; sync warning on macOS is a ⚠︎ label beside the menu bar icon; Wails has no Dock menu.
<!-- SECTION:NOTES:END -->
