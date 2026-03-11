---
id: TASK-110
title: 'T005: Create pending sync store methods'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:41'
updated_date: '2026-03-11 11:31'
labels:
  - foundational
  - phase-2
dependencies:
  - TASK-108
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create pending sync store methods in internal/store/sync.go: implement AddPendingChange(bookmarkUUID, operation, payload string) error (insert into pending_sync within same transaction as bookmark write), ListPendingChanges() ([]*PendingChange, error), ClearPendingChanges() error, PendingChangeCount() (int, error), and ExportAllBookmarks() ([]*Bookmark, error) (all non-archived bookmarks with uuid and updated_at for snapshot building). Define PendingChange struct matching pending_sync columns.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 PendingChange struct defined matching pending_sync schema
- [x] #2 AddPendingChange inserts into pending_sync table
- [x] #3 ListPendingChanges returns all pending changes ordered by created_at
- [x] #4 ClearPendingChanges deletes all rows from pending_sync
- [x] #5 PendingChangeCount returns count of pending changes
- [x] #6 ExportAllBookmarks returns all non-archived bookmarks
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Implemented in internal/store/sync.go
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Created internal/store/sync.go with PendingChange struct, InsertPendingChange (tx-aware), ListPendingChanges, ClearPendingChanges, DeletePendingChange, IncrementRetryCount, PendingChangeCount, ExportAll, GetByUUID, and BeginTx methods.
<!-- SECTION:FINAL_SUMMARY:END -->
