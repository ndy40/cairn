---
id: TASK-145.26
title: T026 Find an existing bookmark by URL
status: Done
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:18:00Z'
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
- [x] #1 FindByURL returns the bookmark matching a URL after the same normalisation as Add, or nothing
- [x] #2 Unit tested
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
BookmarkService.FindByURL normalises like Add; tested. Known limit: scans store.List(), so archived duplicates are not found (the toast still shows). ndw-add, merged into 018-native-desktop-window (0bedb24).
<!-- SECTION:FINAL_SUMMARY:END -->
