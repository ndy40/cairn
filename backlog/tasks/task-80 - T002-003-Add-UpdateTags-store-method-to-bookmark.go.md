---
id: TASK-80
title: 'T002 [003] Add UpdateTags store method to bookmark.go'
status: Done
assignee: []
created_date: '2026-03-06 10:17'
updated_date: '2026-03-06 17:10'
labels:
  - feature-003
  - foundational
  - store
dependencies: []
documentation:
  - specs/003-edit-bookmark-help/data-model.md
  - specs/003-edit-bookmark-help/tasks.md
priority: high
ordinal: 9000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add `UpdateTags(id int64, tags []string) error` to `internal/store/bookmark.go`. Normalises via `NormaliseTags()`, JSON-encodes the result, executes `UPDATE bookmarks SET tags = ? WHERE id = ?`. Returns error on failure, nil on success. Does NOT modify any other column.
<!-- SECTION:DESCRIPTION:END -->
