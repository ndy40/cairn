---
id: TASK-145.43
title: T043 Show sync state on the tray and menu bar icon
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:18:00Z'
labels:
  - desktop
  - us6
dependencies:
  - TASK-145.32
  - TASK-145.42
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 151000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T043 in specs/018-native-desktop-window/tasks.md.

In `desktop/menus.go`, reflect sync state on the tray/menu bar icon: the status line shows "Synced N minutes ago" or the warning text, and the macOS menu bar icon gets a warning label/badge (`SetLabel` or an alternate template icon) on `sync:warning`, cleared on `sync:ok`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The tray menu shows the last sync time or the warning
- [ ] #2 The macOS menu bar icon shows a warning badge while sync has a problem
<!-- AC:END -->
