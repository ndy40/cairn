---
id: TASK-91
title: 'T003 [004] [US1] Remove UpdateLastVisited() from archive.go'
status: Done
assignee: []
created_date: '2026-03-06 15:31'
updated_date: '2026-03-06 17:10'
labels:
  - feature-004
  - US1
  - store
dependencies:
  - TASK-90
documentation:
  - specs/004-bookmark-expiry/research.md
priority: high
ordinal: 4000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Delete the entire `UpdateLastVisited(id int64) error` function (including its comment) from `internal/store/archive.go`. This method is no longer called anywhere once T005 removes the call from `openBookmarkCmd`.
<!-- SECTION:DESCRIPTION:END -->
