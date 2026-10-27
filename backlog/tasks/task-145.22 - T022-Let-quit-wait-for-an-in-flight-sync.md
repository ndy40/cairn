---
id: TASK-145.22
title: T022 Let quit wait for an in-flight sync
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:08:00Z'
labels:
  - desktop
  - us3
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 130000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T022 in specs/018-native-desktop-window/tasks.md.

In `desktop/sync.go`, add `Running() bool` and `Wait(timeout time.Duration) bool` (track in-flight runs with a `sync.WaitGroup` or counter) and emit a new `sync:ok` event after a successful run; extend `desktop/sync_test.go`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The syncer reports whether a sync is running and can wait up to a timeout for it
- [ ] #2 A successful sync emits sync:ok
<!-- AC:END -->
