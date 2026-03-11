---
id: TASK-123
title: 'T018: Implement sync engine Push()'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:45'
updated_date: '2026-03-11 11:41'
labels:
  - us2
  - phase-4
dependencies:
  - TASK-117
  - TASK-122
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement sync engine Push() in internal/sync/engine.go: implement Push() (int, error). Check if cloud snapshot exists. If exists: download, unmarshal. Export all local bookmarks. Build new SyncRecord from local bookmarks, preserving existing tombstones from cloud and adding any new tombstones from pending_sync delete operations. Marshal and upload. Clear pending_sync table on success. Update last_sync_at in config. Return count of changes pushed.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Push downloads existing cloud snapshot if it exists
- [x] #2 Builds SyncRecord from local bookmarks
- [x] #3 Preserves existing tombstones from cloud
- [ ] #4 Adds new tombstones from pending delete operations
- [ ] #5 Uploads marshalled snapshot to cloud
- [ ] #6 Clears pending_sync on success
- [ ] #7 Updates last_sync_at in config
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Engine.Push() uploads snapshot, clears pending changes, updates last_sync_at.
<!-- SECTION:FINAL_SUMMARY:END -->
