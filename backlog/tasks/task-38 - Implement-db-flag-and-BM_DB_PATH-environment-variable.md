---
id: TASK-38
title: Implement --db flag and BM_DB_PATH environment variable
status: Done
assignee: []
created_date: '2026-03-06 04:07'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:polish'
dependencies: []
priority: low
ordinal: 78000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Allow users to override the default database location for all commands and the TUI, enabling use of multiple bookmark databases or non-standard paths.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 --db <path> flag overrides the default database path for all commands
- [ ] #2 BM_DB_PATH environment variable overrides the default path when --db is not set
- [ ] #3 --db takes precedence over BM_DB_PATH when both are set
- [ ] #4 Both flags work consistently for bm (TUI), bm add, bm list, bm search, and bm delete
<!-- AC:END -->
