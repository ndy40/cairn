---
id: TASK-37
title: Implement bm delete subcommand
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:polish'
dependencies: []
priority: low
ordinal: 77000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Allow users to delete a bookmark by ID from the command line without launching the TUI.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 bm delete <id> removes the bookmark and prints 'Deleted' on success
- [ ] #2 bm delete <id> exits with code 1 and prints 'Not found' if the ID does not exist
<!-- AC:END -->
