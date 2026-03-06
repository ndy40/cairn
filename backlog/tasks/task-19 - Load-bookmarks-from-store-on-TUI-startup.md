---
id: TASK-19
title: Load bookmarks from store on TUI startup
status: Done
assignee: []
created_date: '2026-03-06 04:05'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us1'
dependencies: []
priority: high
ordinal: 59000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Ensure the bookmark list is populated when the app launches by loading all bookmarks from the database as an Init() command.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 app.go Init() returns a tea.Cmd that calls store.List() asynchronously
- [ ] #2 Update() receives the list result message and calls BrowseModel.Load(bookmarks)
- [ ] #3 App launches showing existing bookmarks without requiring any user interaction
<!-- AC:END -->
