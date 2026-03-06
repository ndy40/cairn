---
id: TASK-90
title: 'T002 [004] [US1] Simplify ArchiveStale() to 30-day created_at threshold'
status: Done
assignee: []
created_date: '2026-03-06 15:31'
updated_date: '2026-03-06 17:10'
labels:
  - feature-004
  - US1
  - store
dependencies: []
documentation:
  - specs/004-bookmark-expiry/data-model.md
  - specs/004-bookmark-expiry/research.md
priority: high
ordinal: 3000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
In `internal/store/archive.go`, replace the compound WHERE clause in `ArchiveStale()` with a single creation-date condition:\n\nOld:\n```\nAND (\n  (last_visited_at IS NOT NULL AND last_visited_at <= datetime('now', '-183 days'))\n  OR\n  (last_visited_at IS NULL AND created_at <= datetime('now', '-183 days'))\n)\n```\n\nNew:\n```\nAND created_at <= datetime('now', '-30 days')\n```\n\nKeep `is_permanent = 0` and `is_archived = 0` conditions unchanged. Update the function comment to describe the new rule.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Active non-pinned bookmarks with created_at >= 30 days old are archived at startup
- [ ] #2 Pinned bookmarks of any age are NOT archived
- [ ] #3 Active bookmarks created < 30 days ago are NOT archived
- [ ] #4 Startup count message shows correct number
<!-- AC:END -->
