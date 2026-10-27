---
id: TASK-145.3
title: 'T003 Spike: detect installer-created login items'
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
labels:
  - desktop
  - setup
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 111000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T003 in specs/018-native-desktop-window/tasks.md.

Spike open-at-login detection (research R9): on macOS 13+ with the installer's `~/Library/LaunchAgents/io.github.ndy40.cairn.desktop.plist` present, and on Linux with `~/.config/autostart/cairn-desktop.desktop` present, confirm `app.Autostart.Status()` with identifiers `io.github.ndy40.cairn.desktop` / `cairn-desktop` reports `Enabled: true`; record the outcome under R9 in `specs/018-native-desktop-window/research.md` and, if macOS does not detect it, note that T037 must add the one-time plist check
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Settings' open-at-login status reads as on when the installer-created LaunchAgent (macOS 13+) or XDG autostart entry (Linux) exists, or the need for a fallback check is recorded
- [ ] #2 research.md R9 records the outcome
<!-- AC:END -->
