---
id: TASK-145.1
title: T001 Promote godbus to a direct desktop dependency
status: Done
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:18:00Z'
labels:
  - desktop
  - setup
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 109000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T001 in specs/018-native-desktop-window/tasks.md.

Promote `github.com/godbus/dbus/v5` from indirect to a direct dependency in `desktop/go.mod` (`go get github.com/godbus/dbus/v5@v5.2.2 && go mod tidy`), keeping the CLI module untouched
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 desktop/go.mod lists github.com/godbus/dbus/v5 as a direct dependency
- [x] #2 The CLI module (root go.mod) is unchanged and still builds with CGO_ENABLED=0
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Not done in ndw-core: godbus stays indirect until something imports it (go mod tidy would undo it). Lands with T020 or T044.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
godbus/dbus/v5 is a direct desktop dependency (imported by platform.go for tray detection); CLI module untouched. Done in ndw-lifecycle, merged into 018-native-desktop-window (0bedb24).
<!-- SECTION:FINAL_SUMMARY:END -->
