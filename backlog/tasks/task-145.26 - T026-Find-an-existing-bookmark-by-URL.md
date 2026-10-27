---
id: TASK-145.26
title: T026 Find an existing bookmark by URL
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:08:00Z'
labels:
  - desktop
  - us4
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 134000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T026 in specs/018-native-desktop-window/tasks.md.

In `desktop/service.go`, add `FindByURL(rawURL string) (*Bookmark, error)` that normalises like `Add` and returns the matching bookmark or nil; test in `desktop/service_test.go`; regenerate bindings
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 FindByURL returns the bookmark matching a URL after the same normalisation as Add, or nothing
- [ ] #2 Unit tested
<!-- AC:END -->
