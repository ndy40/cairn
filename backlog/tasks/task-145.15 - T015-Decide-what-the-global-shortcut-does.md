---
id: TASK-145.15
title: T015 Decide what the global shortcut does
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
labels:
  - desktop
  - us2
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 123000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T015 in specs/018-native-desktop-window/tasks.md.

Create `desktop/lifecycle.go` with pure `shortcutAction(os string, visible, focused, background bool) action` returning `actionShow`, `actionHideApp` (macOS), `actionMinimise` or `actionHideWindow` (Linux background mode); table tests in `desktop/lifecycle_test.go`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The shortcut decision function covers macOS and Linux, in front or not, with and without background mode, and is unit tested
<!-- AC:END -->
