---
id: TASK-145.16
title: T016 Bring the window forward or put it away
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:08:00Z'
labels:
  - desktop
  - us2
dependencies:
  - TASK-145.11
  - TASK-145.15
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 124000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T016 in specs/018-native-desktop-window/tasks.md.

In `desktop/launcher.go`, replace `Show`/`Toggle` with `BringForward(view string, token string)` (apply activation token via T002's helper, `Restore()` if minimised, `Show()`, `Focus()`, emit `launcher:show` with the pending-view logic kept) and `Toggle(token string)` that uses `shortcutAction` (`app.Hide()` on macOS, `window.Minimise()` or `window.Hide()` on Linux)
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Bringing Cairn forward restores, shows and focuses the window and opens the requested view
- [ ] #2 Toggling while in front hides the app on macOS and minimises (or hides in background mode) on Linux
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implemented in ndw-core, merged in 176f4f7. Code complete; awaiting manual QA (quickstart.md rows 1-10) before ticking criteria.
<!-- SECTION:NOTES:END -->
