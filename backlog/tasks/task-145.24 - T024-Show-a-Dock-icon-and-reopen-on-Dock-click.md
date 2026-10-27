---
id: TASK-145.24
title: T024 Show a Dock icon and reopen on Dock click
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:08:00Z'
labels:
  - desktop
  - us3
dependencies:
  - TASK-145.16
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 132000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T024 in specs/018-native-desktop-window/tasks.md.

In `desktop/main.go`, set `Mac.ActivationPolicy` to `application.ActivationPolicyRegular` and register `events.Mac.ApplicationShouldHandleReopen` to call `launcher.BringForward("search", "")` when the window is not visible
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Cairn has a Dock icon on macOS
- [ ] #2 Clicking the Dock icon with no window open reopens the window
<!-- AC:END -->
