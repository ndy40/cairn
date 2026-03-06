---
id: TASK-18
title: Wire browse list into root app View() with footer
status: Done
assignee: []
created_date: '2026-03-06 04:05'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us1'
dependencies: []
priority: high
ordinal: 58000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Connect BrowseModel to the root app view so the bookmark list fills the terminal and the footer shows available shortcuts.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 StateBrowse View() renders BrowseModel.View() occupying full terminal height minus one footer row
- [ ] #2 Footer shows: [/] Search  [Ctrl+P] Add  [Enter] Open  [d] Delete  [?] Help  [Ctrl+C] Quit
- [ ] #3 List title is set to 'Bookmarks'
<!-- AC:END -->
