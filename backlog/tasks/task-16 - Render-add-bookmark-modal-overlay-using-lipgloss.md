---
id: TASK-16
title: Render add bookmark modal overlay using lipgloss
status: Done
assignee: []
created_date: '2026-03-06 04:05'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us1'
dependencies: []
priority: high
ordinal: 56000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement the View() rendering for StateAdd: the background browse list appears dimmed, and a centred lipgloss-bordered modal is overlaid with the URL input and status message.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 StateAdd View() renders the browse list in the background (dimmed or at reduced opacity)
- [ ] #2 Modal overlay is centred, at least 60 characters wide, with a lipgloss border
- [ ] #3 Modal contains: 'Add Bookmark' label, URL textinput, status message line
- [ ] #4 Modal footer shows: [Enter] Save  [Esc] Cancel
<!-- AC:END -->
