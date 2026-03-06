---
id: TASK-35
title: Implement bm list subcommand
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:polish'
dependencies: []
priority: low
ordinal: 75000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Allow users to list all bookmarks from the command line, suitable for piping to other tools.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 bm list prints all bookmarks as tab-separated rows: id, title, url, domain, created_at
- [ ] #2 bm list --json prints all bookmarks as a JSON array
- [ ] #3 bm list prints nothing and exits cleanly when no bookmarks exist
<!-- AC:END -->
