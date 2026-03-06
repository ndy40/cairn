---
id: TASK-29
title: Render search mode view with input bar and filtered list
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us3'
dependencies: []
priority: medium
ordinal: 69000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement the View() for StateSearch: the search text input appears at the top, the filtered results list fills the remaining space, and the footer shows search-mode shortcuts.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 StateSearch View() renders the search textinput in a lipgloss-styled bar at the top of the terminal
- [ ] #2 Filtered bookmark list is rendered below the search bar
- [ ] #3 Footer shows: [Esc] Clear  [Enter] Open  [Ctrl+P] Add  [Ctrl+C] Quit
<!-- AC:END -->
