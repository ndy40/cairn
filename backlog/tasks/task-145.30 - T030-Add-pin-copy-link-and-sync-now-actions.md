---
id: TASK-145.30
title: T030 Add pin, copy link and sync-now actions
status: Done
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:34:00Z'
labels:
  - desktop
  - us5
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 138000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T030 in specs/018-native-desktop-window/tasks.md.

In `desktop/service.go`, add `SetPinned(id int64, pinned bool) error` (wraps `store.SetPermanent`, emits `bookmarks:changed` "pin" via an injected callback), `CopyLink(id int64) error` (injected clipboard writer using `app.Clipboard.SetText`) and `SyncNow()` (calls a new `Syncer.PullNow()` that skips the throttle, then `Push`, in `desktop/sync.go`); tests in `desktop/service_test.go`; regenerate bindings
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Bookmarks can be pinned and unpinned, their link copied, and a sync run on demand from the desktop app
- [x] #2 Unit tested
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
SetPinned, CopyLink and SyncNow on BookmarkService, Syncer.SyncNow without throttle; tested. ndw-menus, merged into 018-native-desktop-window.
<!-- SECTION:FINAL_SUMMARY:END -->
