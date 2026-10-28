---
id: TASK-145.42
title: T042 Add toasts, the sync banner and the status area
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:18:00Z'
labels:
  - desktop
  - us6
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 150000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T042 in specs/018-native-desktop-window/tasks.md.

Add toast and banner elements to `desktop/frontend/index.html` (`role="status"`) and helpers in `desktop/frontend/src/main.ts`: `showToast(text, action?)` for ~4s; `showBanner(text)` with Dismiss that persists until dismissed or `sync:ok`; route `sync:warning` → banner, `sync:ok` → clear, `status` → footer status area on macOS/KDE and a toast on GNOME; style them in `desktop/frontend/src/style.css` with the warning icon and `--warn-*` tokens
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Confirmations show as toasts for about 4 seconds
- [ ] #2 Sync problems show a banner until dismissed or a later sync succeeds
- [ ] #3 Background progress shows in the footer on macOS and KDE and as a toast on GNOME
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
ndw-rows, merged into 018-native-desktop-window (0bedb24). Code complete; awaiting visual/contrast QA (rows 22-24). Deviations: pure white text and outlined chips on macOS selection, --text-2-strong on Linux tints, for 4.5:1.
<!-- SECTION:NOTES:END -->
