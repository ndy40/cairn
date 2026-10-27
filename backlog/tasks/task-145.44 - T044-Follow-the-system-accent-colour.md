---
id: TASK-145.44
title: T044 Follow the system accent colour
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:18:00Z'
labels:
  - desktop
  - us6
dependencies:
  - TASK-145.4
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 152000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T044 in specs/018-native-desktop-window/tasks.md.

Accent colour: in `desktop/frontend/src/style.css` use `AccentColor` under `@supports (color: AccentColor)` on macOS; in `desktop/platform.go` read `org.freedesktop.appearance accent-color` from the Settings portal via godbus (RGB 0–1 tuple → hex, ignore when absent or out of range) into `PlatformInfo.AccentColor`; in `desktop/frontend/src/main.ts` apply it to `--accent` and darken it until white text reaches 4.5:1 (small contrast helper)
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The selection uses the system accent colour on macOS and from the Settings portal on Linux when available
- [ ] #2 An accent too light for white text is darkened to reach 4.5:1
<!-- AC:END -->
