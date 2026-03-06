---
id: TASK-82
title: 'T004 [003] [US1] Add editKey binding to browse.go'
status: Done
assignee: []
created_date: '2026-03-06 10:17'
updated_date: '2026-03-06 17:10'
labels:
  - feature-003
  - US1
  - tui
dependencies:
  - TASK-80
documentation:
  - specs/003-edit-bookmark-help/contracts/keyboard-shortcuts.md
  - specs/003-edit-bookmark-help/tasks.md
priority: high
ordinal: 11000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add `editKey` key binding (key `e`) to `internal/model/browse.go`. This is the trigger key for entering edit mode from browse mode. No conflict with existing bindings.
<!-- SECTION:DESCRIPTION:END -->
