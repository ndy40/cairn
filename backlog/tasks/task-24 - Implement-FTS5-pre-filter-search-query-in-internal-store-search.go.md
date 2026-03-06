---
id: TASK-24
title: Implement FTS5 pre-filter search query in internal/store/search.go
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us3'
dependencies: []
priority: medium
ordinal: 64000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Provide a fast SQL full-text search that returns candidate bookmark IDs before fuzzy ranking. Used as the first stage of the two-stage search pipeline for terms of 3+ characters.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 FTSSearch(db, term) queries bookmarks_fts using FTS5 MATCH syntax
- [ ] #2 Returns a slice of int64 bookmark IDs matching the term in any of title, description, or domain
- [ ] #3 Returns all bookmark IDs (no filter) when term length is less than 3 characters
- [ ] #4 Returns empty slice (not error) when no FTS matches are found
<!-- AC:END -->
