---
id: TASK-32
title: Wire d/Delete key in browse mode to trigger confirm-delete
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us4'
dependencies: []
priority: low
ordinal: 72000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Allow users to initiate deletion of a selected bookmark using the d or Delete key. Transitions to the confirmation dialog without immediately deleting.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Pressing d or Delete with a bookmark selected in StateBrowse transitions to StateConfirmDelete with the selected bookmark's ID and title
- [ ] #2 Pressing d or Delete when the list is empty has no effect
<!-- AC:END -->
