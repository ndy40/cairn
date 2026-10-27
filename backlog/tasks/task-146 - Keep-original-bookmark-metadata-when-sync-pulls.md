---
id: TASK-146
title: Keep original bookmark metadata when sync pulls
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T09:31:00Z'
updated_date: '2026-10-27T09:36:00Z'
labels:
  - sync
  - bug
dependencies: []
references:
  - internal/sync/engine.go
  - internal/sync/merge.go
  - internal/store/bookmark.go
priority: high
ordinal: 160000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Pulling from Dropbox inserts remote bookmarks with Store.Insert, which stamps created_at with the pull time and a new UUID and drops pinned/archived state. Other machines therefore show the import date, deletes (matched by UUID) never reach them, and the merged snapshot they upload overwrites the cloud's original dates. Found while adding the date-added column to cairn-desktop (TASK-145.51).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 A bookmark pulled onto another machine keeps its original UUID, created_at, updated_at, pinned and archived state
- [x] #2 Pulled bookmarks are not queued as local changes to push back
- [x] #3 When a pulled bookmark matches a local one (by UUID or URL), the earlier created_at is kept on both, so dates already overwritten heal on the next sync
- [x] #4 A bookmark deleted on one machine is deleted on another after push and pull
- [x] #5 Unit tests cover the merge rule and a pull round trip
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- fix-sync-metadata (from main): 50a5ad8 InsertSynced + created_at healing; 2801cbd tombstones for local deletes, remote tombstones carried forward on pull.
- Verified by orchestrator: go test ./..., go vet, golangci-lint (0 issues), CGO_ENABLED=0 build.
- Not merged yet. Open: diverged UUIDs (agent recommends adopting the UUID of the earliest copy), Push uploads without pulling first (can undo another device's delete or drop its adds), tombstones never expire, cloud updates only change tags.
<!-- SECTION:NOTES:END -->
