# Data Model: Bookmark Expiry & Last-Visited Removal

**Feature**: 004-bookmark-expiry
**Date**: 2026-03-06

---

## Entity Changes

**No schema changes.** The `bookmarks` table and all columns remain exactly as defined in feature 002. The `last_visited_at` column stays in place — it will simply never be written to or read in the UI going forward.

---

## Updated Store Operation: ArchiveStale

The existing `ArchiveStale()` method in `internal/store/archive.go` is modified:

**Old logic** (183-day threshold, compound condition):
```
WHERE is_permanent = 0
  AND is_archived  = 0
  AND (
    (last_visited_at IS NOT NULL AND last_visited_at <= datetime('now', '-183 days'))
    OR
    (last_visited_at IS NULL AND created_at <= datetime('now', '-183 days'))
  )
```

**New logic** (30-day threshold, creation-date only):
```
WHERE is_permanent = 0
  AND is_archived  = 0
  AND created_at <= datetime('now', '-30 days')
```

**Behaviour**:
- Archives all active, non-pinned bookmarks whose `created_at` is 30 or more days before now.
- Pinned bookmarks (`is_permanent = 1`) are never archived.
- Already-archived bookmarks (`is_archived = 1`) are not evaluated.
- Returns the count of rows archived.

---

## Removed Store Operation: UpdateLastVisited

The `UpdateLastVisited(id int64) error` method is deleted from `internal/store/archive.go`.

No callers remain once the call is removed from `openBookmarkCmd` in `internal/model/app.go`.

---

## Expiry State Transition

The bookmark lifecycle is unchanged except for the archival trigger:

```
Active (is_archived=0)
  → [startup, age ≥ 30 days, is_permanent=0] → Archived (is_archived=1)
  → [user presses 'r' in archive view]        → Active (is_archived=0)
```

The permanent flag continues to gate expiry at the transition point.

---

## Browse Row Format (Updated)

The description line on each bookmark row in the browse list changes:

**Old format**: `domain · YYYY-MM-DD · Last: YYYY-MM-DD [· #tags]`
or `domain · YYYY-MM-DD · Never visited [· #tags]`

**New format**: `domain · YYYY-MM-DD [· #tags]`

The `LastVisitedAt` field on the `Bookmark` struct is no longer read in the UI (though it remains in the struct since the DB column is retained).
