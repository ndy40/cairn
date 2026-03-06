---
id: TASK-36
title: Implement bm search subcommand
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:polish'
dependencies: []
priority: low
ordinal: 76000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Allow users to run a fuzzy search from the command line and get matching bookmarks in a pipeable format.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 bm search <query> prints matching bookmarks in the same tab-separated format as bm list
- [ ] #2 bm search --json outputs matching bookmarks as a JSON array
- [ ] #3 bm search --limit N returns at most N results (default 10)
- [ ] #4 bm search prints nothing and exits cleanly when no matches are found
<!-- AC:END -->
