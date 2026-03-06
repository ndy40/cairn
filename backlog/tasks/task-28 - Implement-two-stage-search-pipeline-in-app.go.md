---
id: TASK-28
title: Implement two-stage search pipeline in app.go
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us3'
dependencies: []
priority: medium
ordinal: 68000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Combine FTS5 pre-filtering with fuzzy ranking for efficient search at scale. For short queries, skip FTS and fuzzy-rank everything directly.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 When search term length >= 3, FTSSearch() is called first to get candidate IDs, then only those bookmarks are passed to search.Search()
- [ ] #2 When search term length < 3, all bookmarks are passed directly to search.Search() without FTS pre-filtering
- [ ] #3 Search results update in real time on each keystroke
<!-- AC:END -->
