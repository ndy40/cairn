---
id: TASK-145.39
title: T039 Leave start-at-login to the desktop app in install.sh
status: Done
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T03:58:00Z'
labels:
  - desktop
  - us5
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 147000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T039 in specs/018-native-desktop-window/tasks.md.

In `install.sh`: stop writing `~/.config/autostart/cairn-desktop.desktop` and `~/Library/LaunchAgents/io.github.ndy40.cairn.desktop.plist` for new installs (leave existing ones untouched); keep `--no-autostart` accepted as a no-op that prints "Start at login is now set in Cairn's Settings"; restart a running app with `--background` instead of `--hidden`; update the usage text
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 New installs get no login item; existing ones are left in place
- [x] #2 --no-autostart is accepted as a no-op with a notice
- [x] #3 A running app is restarted with --background
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
The installer no longer creates start-at-login items; the desktop app owns them (FR-042).

Changes:
- install.sh stops writing ~/.config/autostart/cairn-desktop.desktop and ~/Library/LaunchAgents/io.github.ndy40.cairn.desktop.plist; existing ones are left in place.
- --no-autostart is accepted as a no-op and prints a notice.
- A running app is restarted with --background; popup wording removed from prompts and hints.

Tests: sh -n install.sh; sh install.sh --help. shellcheck not available.

Risk: --background needs the ndw-core change; do not release install.sh before it is merged.

Commits: 6c4a222 (merged in 3ea32f6).
<!-- SECTION:FINAL_SUMMARY:END -->
