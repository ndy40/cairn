---
id: TASK-93
title: 'T005 [004] [US2] Remove UpdateLastVisited call from openBookmarkCmd in app.go'
status: Done
assignee: []
created_date: '2026-03-06 15:32'
updated_date: '2026-03-06 17:10'
labels:
  - feature-004
  - US2
  - tui
dependencies: []
documentation:
  - specs/004-bookmark-expiry/research.md
priority: high
ordinal: 6000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
In `internal/model/app.go`, remove the `_ = s.UpdateLastVisited(b.ID)` line from `openBookmarkCmd`. Update the function comment to remove references to last-visited recording. The command should simply: open the URL, and on success return `loadBookmarks(s)()`.
<!-- SECTION:DESCRIPTION:END -->
