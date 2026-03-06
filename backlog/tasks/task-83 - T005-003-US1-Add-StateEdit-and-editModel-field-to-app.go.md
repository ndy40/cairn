---
id: TASK-83
title: 'T005 [003] [US1] Add StateEdit and editModel field to app.go'
status: Done
assignee: []
created_date: '2026-03-06 10:17'
updated_date: '2026-03-06 17:10'
labels:
  - feature-003
  - US1
  - tui
dependencies:
  - TASK-81
  - TASK-82
documentation:
  - specs/003-edit-bookmark-help/tasks.md
priority: high
ordinal: 12000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
In `internal/model/app.go`:\n- Add `StateEdit` to the `AppState` const block\n- Add `editModel EditModel` field to the `Model` struct
<!-- SECTION:DESCRIPTION:END -->
