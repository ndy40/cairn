---
id: TASK-145.47
title: T047 Describe the desktop app as a native window in the README
status: Done
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T03:58:00Z'
labels:
  - desktop
  - polish
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 155000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T047 in specs/018-native-desktop-window/tasks.md.

Update the desktop app description in `README.md` ("a hotkey popup" → a native desktop app with a global shortcut) and the `--with-desktop`/`--no-autostart` notes
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 README.md describes cairn-desktop as a native app window with a global shortcut and explains --no-autostart
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
README describes cairn-desktop as a native app window with a global shortcut, opened from the app menu (Linux) or the Applications folder (macOS), and explains --no-autostart.

Commits: 0da40a6, 68d16af (merged in 3ea32f6).
<!-- SECTION:FINAL_SUMMARY:END -->
