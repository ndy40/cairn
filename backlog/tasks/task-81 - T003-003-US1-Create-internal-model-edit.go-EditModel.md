---
id: TASK-81
title: 'T003 [003] [US1] Create internal/model/edit.go EditModel'
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
ordinal: 10000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create `internal/model/edit.go` with `EditModel` struct containing a `textinput.Model` for tags and a read-only title string. Implement:\n- `New(b Bookmark) EditModel` — pre-fills tags field comma-separated\n- `Update(msg tea.Msg) (EditModel, tea.Cmd)`\n- `View() string` — shows read-only title heading + tags input\n- `Tags() []string` — splits input on comma, returns slice
<!-- SECTION:DESCRIPTION:END -->
