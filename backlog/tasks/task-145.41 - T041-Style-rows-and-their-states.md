---
id: TASK-145.41
title: T041 Style rows and their states
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:18:00Z'
labels:
  - desktop
  - us6
dependencies:
  - TASK-145.40
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 149000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T041 in specs/018-native-desktop-window/tasks.md.

Style rows, sections, states (default, hover, selected, window inactive via `blur`/`focus` on `window`) and empty states in `desktop/frontend/src/style.css` per the Components artboard: 52px rows, 7px radius, 8px inset; selection = accent fill with white text on macOS, accent tint (16% GNOME, 18% + 1px outline KDE) with dark text on Linux
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Rows match the design in default, hover, selected and inactive-window states in light and dark
- [ ] #2 Selection text meets 4.5:1 contrast on every platform
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
ndw-rows, merged into 018-native-desktop-window (0bedb24). Code complete; awaiting visual/contrast QA (rows 22-24). Deviations: pure white text and outlined chips on macOS selection, --text-2-strong on Linux tints, for 4.5:1.
<!-- SECTION:NOTES:END -->
