---
id: TASK-14
title: Wire Ctrl+P shortcut to open add bookmark modal
status: Done
assignee: []
created_date: '2026-03-06 04:05'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us1'
dependencies: []
priority: high
ordinal: 54000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Connect the global Ctrl+P keypress to the add-bookmark flow: read the clipboard, pre-fill the modal input, and transition to StateAdd. Handle the empty-clipboard edge case gracefully.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Pressing Ctrl+P in StateBrowse or StateSearch reads the clipboard and transitions to StateAdd with the URL pre-filled in AddModel
- [ ] #2 If clipboard is empty, the application stays in the current state and shows an inline footer message: 'Clipboard is empty'
- [ ] #3 If clipboard read fails (e.g., xclip not installed), a descriptive error is shown in the footer
<!-- AC:END -->
