# Data Model: Async Autosync

**Feature**: 007-async-autosync
**Date**: 2026-03-11

## Schema Changes

**None.** This feature does not modify the database schema. The existing `pending_sync` table and `bookmarks` table remain unchanged.

## Existing Entities (unchanged)

### Bookmark
- Already records pending sync changes atomically during insert/delete/update operations
- No field additions or modifications

### Pending Sync Record
- Table: `pending_sync`
- Fields: `id`, `bookmark_uuid`, `change_type`, `created_at`
- Already written in the same transaction as the bookmark operation
- Cleared by `ClearPendingChanges()` after successful push

## Data Flow Change

**Before (synchronous)**:
1. User runs `cairn add <url>`
2. Bookmark inserted + pending change recorded (single transaction)
3. `autoSyncPush()` called synchronously — blocks until network I/O completes
4. Command returns

**After (asynchronous)**:
1. User runs `cairn add <url>`
2. Bookmark inserted + pending change recorded (single transaction)
3. Background subprocess spawned: `cairn sync push`
4. Command returns immediately
5. Background subprocess runs push (opens own DB connection, uploads, clears pending changes)
