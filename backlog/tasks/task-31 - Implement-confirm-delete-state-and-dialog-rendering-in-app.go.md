---
id: TASK-31
title: Implement confirm-delete state and dialog rendering in app.go
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us4'
dependencies: []
priority: low
ordinal: 71000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add the StateConfirmDelete mode to the root app model. When triggered, store the target bookmark's ID and title, then render a confirmation dialog overlaid on the browse list.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 StateConfirmDelete stores the selected bookmark ID and title
- [ ] #2 View() for StateConfirmDelete renders the browse list in the background (dimmed)
- [ ] #3 A centred lipgloss modal shows: 'Delete «{title}»?' with footer [y/Enter] Confirm Delete  [n/Esc] Cancel
<!-- AC:END -->
