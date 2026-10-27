---
id: TASK-146
title: Keep original bookmark metadata when sync pulls
status: Done
assignee:
  - '@claude'
created_date: '2026-10-27T09:31:00Z'
updated_date: '2026-10-27T10:41:00Z'
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

- Merged to main as f499637 (PR #10, released in v0.7.2). The open items were all fixed in the same PR: diverged UUIDs converge on the earliest copy, push merges before uploading, tombstones expire after 180 days, and edits and pins sync. A 30-day first_seen_at guard stops restored old dates from auto-archiving.
- Merged into 018-native-desktop-window (5c3a772); the desktop now pushes after a pin (0df2932).
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Sync now preserves bookmark identity and metadata across devices (PR #10, v0.7.2).

Changes:
- Pull stores remote bookmarks with their original UUID, created/updated dates and pin/archive state (Store.InsertSynced), without queuing them to push back.
- The earliest created_at wins on a UUID or URL match, so dates already overwritten heal on the next sync; diverged copies converge on the earliest copy's UUID.
- Local deletes write tombstones, pulls carry remote tombstones forward, and tombstones expire after 180 days.
- Push merges the cloud snapshot before uploading, so it no longer undoes other devices' deletes or drops their adds.
- Newer remote edits (URL, title, description, tags, pin) apply locally; archive state stays per-device.
- Auto-archive also waits 30 days after a bookmark arrives on a device (migration v5, first_seen_at).

Tests: store, merge and two-device engine round trips; go test, vet, golangci-lint and the CGO_ENABLED=0 CLI build pass.
<!-- SECTION:FINAL_SUMMARY:END -->
