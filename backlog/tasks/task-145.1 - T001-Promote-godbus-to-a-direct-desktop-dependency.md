---
id: TASK-145.1
title: T001 Promote godbus to a direct desktop dependency
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
ordinal: 109000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T001 in specs/018-native-desktop-window/tasks.md.

Promote `github.com/godbus/dbus/v5` from indirect to a direct dependency in `desktop/go.mod` (`go get github.com/godbus/dbus/v5@v5.2.2 && go mod tidy`), keeping the CLI module untouched
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 desktop/go.mod lists github.com/godbus/dbus/v5 as a direct dependency
- [ ] #2 The CLI module (root go.mod) is unchanged and still builds with CGO_ENABLED=0
<!-- AC:END -->
