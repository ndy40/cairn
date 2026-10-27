---
id: TASK-145
title: Native app window for cairn-desktop
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
labels:
  - desktop
  - feature
dependencies: []
documentation:
  - specs/018-native-desktop-window/spec.md
  - specs/018-native-desktop-window/plan.md
  - specs/018-native-desktop-window/tasks.md
priority: high
ordinal: 108000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Turn cairn-desktop from a frameless Spotlight-style popup into a regular app window with native chrome on macOS and Linux (GNOME, KDE), per the design canvas https://claude.ai/artifact/Qvs5oe46kZEPWXYcJmhrd5.

Spec and plan: specs/018-native-desktop-window/. Subtasks map 1:1 to tasks.md (T001-T050).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The window uses each platform's own controls and can be moved, resized and closed like any app (User Story 1)
- [ ] #2 Shortcut, type, Enter still opens a bookmark in seconds (User Story 2)
- [ ] #3 Close, quit and background behaviour follow the platform rules (User Story 3)
- [ ] #4 Add Bookmark opens as a sheet or dialog attached to the window (User Story 4)
- [ ] #5 Menus, Settings and shortcuts are where each platform puts them (User Story 5)
- [ ] #6 Rows, messages and colours match the design (User Story 6)
<!-- AC:END -->
