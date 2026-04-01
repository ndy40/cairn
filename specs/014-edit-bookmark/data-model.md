# Data Model: Edit Bookmark

**Feature**: 014-edit-bookmark
**Date**: 2026-04-01

## Entities

### Bookmark (existing — no schema changes)

| Field         | Type      | Editable via this feature? | Notes                                    |
| ------------- | --------- | -------------------------- | ---------------------------------------- |
| ID            | int64     | No                         | Primary key, immutable                   |
| UUID          | string    | No                         | Sync identifier, immutable               |
| URL           | string    | **Yes (new)**              | UNIQUE constraint; domain recalculated   |
| Domain        | string    | Derived                    | Auto-recalculated when URL changes       |
| Title         | string    | Yes (existing)             | Already editable via CLI                 |
| Description   | string    | No                         | Set at insert time from page fetch       |
| CreatedAt     | time.Time | No                         | Immutable                                |
| UpdatedAt     | time.Time | Auto                       | Refreshed on every edit                  |
| Tags          | []string  | Yes (existing)             | JSON array, max 3, normalised            |
| LastVisitedAt | *time.Time| No                         | Retained but not written                 |
| IsPermanent   | bool      | No (via pin command)       | Separate command                         |
| IsArchived    | bool      | No                         | Managed by expiry system                 |
| ArchivedAt    | *time.Time| No                         | Managed by expiry system                 |

### BookmarkPatch (existing — modified)

| Field | Type      | Before  | After             |
| ----- | --------- | ------- | ----------------- |
| Title | *string   | Present | Unchanged         |
| Tags  | *[]string | Present | Unchanged         |
| URL   | *string   | —       | **New (optional)** |

When `URL` is non-nil in BookmarkPatch:
1. Validate non-empty
2. Check no other bookmark has the same URL (exclude self by ID)
3. Recalculate domain via `extractDomain()`
4. Update `url`, `domain`, and `updated_at` columns

## Validation Rules

| Rule                  | Applies to | Behaviour                                           |
| --------------------- | ---------- | --------------------------------------------------- |
| Non-empty URL         | URL edit   | Reject with error before DB query                   |
| No duplicate URL      | URL edit   | Query `url = ? AND id != ?`; reject with ErrDuplicate |
| Tag normalisation     | Tags edit  | Lowercase, trim, dedup, truncate 32 runes, max 3   |
| Non-empty title       | Title edit | Reject with error (existing behaviour)              |

## State Transitions

No new states. Edit operations update field values but do not change bookmark lifecycle state (active/archived/pinned).

## Sync Impact

All edit operations (URL, title, tags) record a `pending_sync` entry with operation `'update'`. This is already handled by `UpdateFields` — no changes needed to the sync recording logic.
