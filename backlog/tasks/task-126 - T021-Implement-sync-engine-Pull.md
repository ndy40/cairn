---
id: TASK-126
title: 'T021: Implement sync engine Pull()'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:45'
updated_date: '2026-03-11 11:41'
labels:
  - us3
  - phase-5
dependencies:
  - TASK-125
  - TASK-114
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement sync engine Pull() in internal/sync/engine.go: implement Pull() (inserted int, deleted int, error). Download cloud snapshot via backend.Download(). Unmarshal. Call merge algorithm with local bookmarks and remote snapshot. Apply inserts/updates via store.ImportBookmarks(). Apply deletes via store.DeleteByUUIDs(). Update last_sync_at in config. Return counts.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Pull downloads cloud snapshot
- [x] #2 Unmarshals snapshot correctly
- [x] #3 Calls merge algorithm with local and remote data
- [ ] #4 Applies inserts and updates via ImportBookmarks
- [ ] #5 Applies deletes via DeleteByUUIDs
- [ ] #6 Updates last_sync_at in config
- [ ] #7 Returns correct inserted and deleted counts
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Engine.Pull() downloads snapshot, merges via MergeBookmarks, applies changes, uploads merged result, updates last_sync_at.
<!-- SECTION:FINAL_SUMMARY:END -->
