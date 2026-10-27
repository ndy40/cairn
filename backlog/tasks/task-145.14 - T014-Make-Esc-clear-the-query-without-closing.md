---
id: TASK-145.14
title: T014 Make Esc clear the query without closing
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:08:00Z'
labels:
  - desktop
  - us1
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 122000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T014 in specs/018-native-desktop-window/tasks.md.

In `desktop/frontend/src/main.ts`, change Esc handling: non-empty query → clear and re-search; empty query → do nothing (no hide call); remove any code path that hides the window on Esc
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Esc clears a non-empty query; Esc on an empty query does nothing and the window stays open
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implemented in ndw-core, merged in 176f4f7. Code complete; awaiting manual QA (quickstart.md rows 1-10) before ticking criteria.
<!-- SECTION:NOTES:END -->
