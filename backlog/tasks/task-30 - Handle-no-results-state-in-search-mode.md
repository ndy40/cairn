---
id: TASK-30
title: Handle no-results state in search mode
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us3'
dependencies: []
priority: medium
ordinal: 70000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
When a search returns no matches, show a clear message to the user instead of an empty list.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 When search.Search() returns an empty slice, the list area shows a centred message: 'No results for «{term}»'
- [ ] #2 When the search term is cleared, the full bookmark list is restored immediately
<!-- AC:END -->
