---
id: TASK-25
title: Implement multi-field fuzzy search wrapper in internal/search/fuzzy.go
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us3'
dependencies: []
priority: medium
ordinal: 65000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Wrap sahilm/fuzzy to search across title, domain, and description fields with field-weight scoring. This produces intuitive fuzzy ranking for the interactive search mode.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Search(query, bookmarks) runs sahilm/fuzzy FindFrom separately against title (weight 3), domain (weight 2), description (weight 1)
- [ ] #2 Results are merged by bookmark ID, taking the highest weighted score per bookmark
- [ ] #3 Returned slice is sorted by composite score descending
- [ ] #4 Empty query returns the full bookmarks slice unchanged (no filtering)
- [ ] #5 Handles case-insensitive matching
<!-- AC:END -->
