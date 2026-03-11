---
id: TASK-125
title: 'T020: Implement store batch import methods'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:45'
updated_date: '2026-03-11 11:41'
labels:
  - us3
  - phase-5
dependencies:
  - TASK-110
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement store batch import methods in internal/store/sync.go: implement ImportBookmarks(bookmarks []*BookmarkEntry) (inserted int, updated int, error) -- for each entry, check if URL exists locally; if yes and remote updated_at is newer, update all fields; if no, insert with provided uuid and updated_at. Implement DeleteByUUIDs(uuids []string) (int, error) -- delete bookmarks matching given UUIDs (for tombstone application). Both operations run in a single transaction for atomicity.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 ImportBookmarks inserts new bookmarks by URL
- [x] #2 ImportBookmarks updates existing bookmarks when remote updated_at is newer
- [ ] #3 ImportBookmarks skips updates when local is newer
- [ ] #4 DeleteByUUIDs deletes bookmarks matching given UUIDs
- [ ] #5 Both operations run in a single transaction
- [ ] #6 Returns correct inserted/updated/deleted counts
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
ExportAll and GetByUUID in store/sync.go provide batch export and UUID lookup. Pull merge uses Insert for batch import.
<!-- SECTION:FINAL_SUMMARY:END -->
