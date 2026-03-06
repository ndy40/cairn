---
id: TASK-17
title: Implement minimal bookmark list display in internal/model/browse.go
status: Done
assignee: []
created_date: '2026-03-06 04:05'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us1'
dependencies: []
priority: high
ordinal: 57000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create the BrowseModel with a bubbles/list that displays saved bookmarks. This is required for US1 so the user can see the bookmark was saved after pressing Ctrl+P.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 BrowseModel wraps a bubbles/list.Model
- [ ] #2 BookmarkItem satisfies list.Item interface with Title() returning the page title and Description() returning domain + ' · ' + formatted created_at date
- [ ] #3 Load(bookmarks []Bookmark) replaces the list items with the provided slice
- [ ] #4 Built-in list filtering is disabled (fuzzy search is handled separately in US3)
<!-- AC:END -->
