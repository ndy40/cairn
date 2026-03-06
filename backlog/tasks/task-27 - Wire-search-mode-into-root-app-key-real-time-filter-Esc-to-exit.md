---
id: TASK-27
title: 'Wire search mode into root app: / key, real-time filter, Esc to exit'
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us3'
dependencies: []
priority: medium
ordinal: 67000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Connect the search mode lifecycle in app.go: / enters search, each keystroke updates the filtered list displayed by BrowseModel, Esc restores the full list.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Pressing / in StateBrowse transitions to StateSearch and focuses the search input
- [ ] #2 Each search keystroke calls search.Search() and BrowseModel.Load(filteredResults)
- [ ] #3 Pressing Escape in StateSearch clears the search term, reloads the full bookmark list, and transitions to StateBrowse
- [ ] #4 Enter in StateSearch opens the currently selected bookmark in the browser
<!-- AC:END -->
