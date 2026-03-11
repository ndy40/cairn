---
id: TASK-117
title: 'T012: Implement sync engine Setup()'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:41'
updated_date: '2026-03-11 11:36'
labels:
  - us1
  - phase-3
dependencies:
  - TASK-114
  - TASK-115
  - TASK-116
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement sync engine Setup() in internal/sync/engine.go: create Engine struct holding *store.Store, SyncBackend, *SyncConfig. Implement Setup(appKey string) error: run OAuth2 flow, generate device UUID, save config. Check if cloud snapshot exists via backend.Exists(). If exists: download, unmarshal, call merge, apply inserts/updates/deletes to store atomically. If not exists: export all local bookmarks via store.ExportAllBookmarks(), build SyncRecord, marshal, upload via backend.Upload(). Update last_sync_at. Return bookmark count for display.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Engine struct created with Store, SyncBackend, SyncConfig fields
- [x] #2 Setup runs OAuth2 flow and generates device UUID
- [x] #3 Setup saves config with tokens and device ID
- [x] #4 When cloud snapshot exists: downloads, merges, and applies changes atomically
- [x] #5 When no cloud snapshot: exports local bookmarks and uploads initial snapshot
- [x] #6 last_sync_at updated after successful setup
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Implemented in internal/sync/engine.go
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Created internal/sync/engine.go with Engine struct. Setup() runs OAuth2 flow, generates device UUID, saves config, then either pulls+merges (if cloud exists) or uploads initial snapshot. Updates last_sync_at.
<!-- SECTION:FINAL_SUMMARY:END -->
