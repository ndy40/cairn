---
id: TASK-26
title: Implement search mode sub-model in internal/model/search.go
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us3'
dependencies: []
priority: medium
ordinal: 66000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create the SearchModel sub-model with a bubbles/textinput for the search bar and real-time fuzzy filtering of the bookmark list as the user types.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 SearchModel contains a bubbles/textinput.Model with placeholder 'Search bookmarks…'
- [ ] #2 On each keystroke, Search(term, allBookmarks) is called and filtered results are made available via Results()
- [ ] #3 Pressing Ctrl+A clears the search term; pressing Escape signals the parent to return to StateBrowse
- [ ] #4 Arrow key events are delegated to the parent browse list (not consumed by the text input)
<!-- AC:END -->
