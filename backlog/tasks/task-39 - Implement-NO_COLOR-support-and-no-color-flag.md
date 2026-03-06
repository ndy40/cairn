---
id: TASK-39
title: Implement NO_COLOR support and --no-color flag
status: Done
assignee: []
created_date: '2026-03-06 04:07'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:polish'
dependencies: []
priority: low
ordinal: 79000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Respect the no-color.org standard so the CLI output is usable in environments without ANSI color support.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 When NO_COLOR env var is set (any value), ANSI colors are disabled in bm list and bm search output
- [ ] #2 --no-color flag also disables colors regardless of NO_COLOR
<!-- AC:END -->
