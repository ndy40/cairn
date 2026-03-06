---
id: TASK-22
title: Implement empty state display when no bookmarks exist
status: Done
assignee: []
created_date: '2026-03-06 04:05'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us2'
dependencies: []
priority: medium
ordinal: 62000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Show a helpful message when the bookmark list is empty so new users know what to do.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 When List() returns zero bookmarks, the list area shows a centred message: 'No bookmarks yet. Press Ctrl+P to add your first bookmark.'
- [ ] #2 The empty state message is rendered with lipgloss centring within the available terminal height
<!-- AC:END -->
