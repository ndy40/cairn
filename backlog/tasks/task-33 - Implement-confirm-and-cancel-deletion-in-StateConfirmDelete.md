---
id: TASK-33
title: Implement confirm and cancel deletion in StateConfirmDelete
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us4'
dependencies: []
priority: low
ordinal: 73000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Handle the user's response in the confirm-delete dialog: y/Enter deletes the bookmark, n/Esc cancels and returns to browsing.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Pressing y or Enter in StateConfirmDelete runs store.DeleteByID(id) as a tea.Cmd
- [ ] #2 On successful delete, the bookmark list is reloaded from the store and the app transitions to StateBrowse
- [ ] #3 Pressing n or Escape in StateConfirmDelete transitions to StateBrowse without any deletion
- [ ] #4 If store.DeleteByID fails, an error message is shown in the footer and the app returns to StateBrowse
<!-- AC:END -->
