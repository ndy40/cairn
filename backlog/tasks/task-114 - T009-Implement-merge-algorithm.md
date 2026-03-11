---
id: TASK-114
title: 'T009: Implement merge algorithm'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:41'
updated_date: '2026-03-11 11:33'
labels:
  - foundational
  - phase-2
dependencies:
  - TASK-109
  - TASK-112
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement merge algorithm in internal/sync/merge.go: implement MergeBookmarks(local []*store.Bookmark, remote *SyncRecord) (toInsert []*BookmarkEntry, toUpdate []*BookmarkEntry, toDelete []string, error). Logic: iterate remote bookmarks, match by URL against local set; if URL exists locally and remote updated_at is newer, mark for update; if URL not in local, mark for insert; iterate remote tombstones, if UUID exists locally mark for delete. Return three slices for the caller to apply atomically.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 MergeBookmarks function implemented
- [x] #2 Remote bookmarks not in local are marked for insert
- [x] #3 Remote bookmarks with newer updated_at override local
- [x] #4 Remote tombstones mark matching local bookmarks for delete
- [x] #5 URL deduplication works correctly
- [x] #6 Merge returns three slices: toInsert, toUpdate, toDelete
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Implemented in internal/sync/merge.go
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Created internal/sync/merge.go with MergeBookmarks function. Matches by UUID first, then URL for dedup. Remote-newer wins for updates. Tombstones trigger local deletes. Returns MergeResult with three slices.
<!-- SECTION:FINAL_SUMMARY:END -->
