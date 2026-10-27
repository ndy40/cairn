---
id: TASK-145.13
title: T013 Style the title band and drag regions
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:08:00Z'
labels:
  - desktop
  - us1
dependencies:
  - TASK-145.12
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 121000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T013 in specs/018-native-desktop-window/tasks.md.

In `desktop/frontend/src/style.css`: remove the `.launcher` border; add `--titleband-h` (52px macOS, measured on a build) and make `.titleband` `--wails-draggable: drag` with the input and buttons `no-drag`; on `[data-os=darwin]` pad the band left for the traffic lights and give only the band a translucent background while `.results` and the footer are opaque; on `[data-os=linux]` style it as a search bar under the GTK title bar (title-bar fill, hairline below); remove the old drag rules on `.results`, `.hints` and `.add`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Dragging empty title-band space moves the window; dragging in the search field selects text
- [ ] #2 Only the macOS title band is translucent; the list stays opaque
- [ ] #3 The CSS window border is removed
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implemented in ndw-core, merged in 176f4f7. Code complete; awaiting manual QA (quickstart.md rows 1-10) before ticking criteria.
<!-- SECTION:NOTES:END -->
