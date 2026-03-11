---
id: TASK-108
title: 'T003: Implement schema migration V3'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:40'
updated_date: '2026-03-11 11:30'
labels:
  - foundational
  - phase-2
dependencies:
  - TASK-106
  - TASK-107
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement schema migration V3 in internal/store/store.go: add uuid and updated_at columns to bookmarks table with DEFAULT values, backfill existing rows (uuid from google/uuid, updated_at from created_at), create unique index on uuid, create index on updated_at, create pending_sync table with id/bookmark_uuid/operation/payload/created_at/retry_count columns, create index on pending_sync.created_at, record version 3 in schema_version. Follow existing migration pattern in the migrations slice.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 uuid TEXT NOT NULL DEFAULT '' column added to bookmarks
- [x] #2 updated_at TEXT NOT NULL DEFAULT '' column added to bookmarks
- [x] #3 Existing rows backfilled with UUID v4 and updated_at from created_at
- [x] #4 Unique index idx_bookmarks_uuid created
- [x] #5 Index idx_bookmarks_updated_at created
- [x] #6 pending_sync table created with all columns
- [x] #7 Index idx_pending_sync_created_at created
- [x] #8 Version 3 recorded in schema_version
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add migrateV3 to store.go\n2. Add uuid/updated_at columns\n3. Backfill existing rows\n4. Create pending_sync table\n5. Create indexes
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented migration V3 in store.go: adds uuid and updated_at columns to bookmarks, backfills existing rows with UUIDv4 and created_at, creates pending_sync table with indexes. Migration runs in a transaction for atomicity.
<!-- SECTION:FINAL_SUMMARY:END -->
