---
id: TASK-85
title: 'T007 [003] [US2] Document last-visited update flow in app.go'
status: Done
assignee: []
created_date: '2026-03-06 10:17'
updated_date: '2026-03-06 17:10'
labels:
  - feature-003
  - US2
  - docs
dependencies:
  - TASK-84
documentation:
  - specs/003-edit-bookmark-help/data-model.md
  - specs/003-edit-bookmark-help/research.md
priority: medium
ordinal: 14000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add a comment block in `internal/model/app.go` above `openBookmarkCmd` documenting the last-visited update flow as confirmed in research decision 4. No functional code changes needed — the existing implementation is already correct.\n\nFlow to document:\n1. User presses Enter on bookmark\n2. openBookmarkCmd fires (tea.Cmd)\n3. openURLRaw(url) — starts browser process\n4. On success: s.UpdateLastVisited(b.ID) — writes datetime('now') to DB\n5. Returns loadBookmarks(s)() — triggers list reload showing updated 'Last: <date>'
<!-- SECTION:DESCRIPTION:END -->
