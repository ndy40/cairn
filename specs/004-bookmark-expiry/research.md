# Research: Bookmark Expiry & Last-Visited Removal

**Feature**: 004-bookmark-expiry
**Date**: 2026-03-06

---

## Decision 1: Expiry Threshold Query

**Decision**: Replace the compound `last_visited_at`/`created_at` condition in `ArchiveStale()` with a single `created_at <= datetime('now', '-30 days')` condition.

**Rationale**: The existing SQL in `internal/store/archive.go` has a two-branch WHERE clause — one branch for bookmarks that have a `last_visited_at`, another for those that don't. Since last-visited tracking is being removed, `last_visited_at` will no longer be written to. The expiry rule is now purely based on creation date. Simplifying to a single `created_at` check is cleaner, correct, and still requires no schema changes.

**Alternatives considered**:
- Keep the two-branch query and just change the threshold numbers: Would work but is unnecessarily complex given that `last_visited_at` will never be populated going forward. Rejected.
- Add a new SQL column (e.g., `expires_at`): Over-engineering for a fixed 30-day rule. Rejected.

---

## Decision 2: Remove UpdateLastVisited

**Decision**: Delete the `UpdateLastVisited(id int64) error` method from `internal/store/archive.go` and remove its call from `openBookmarkCmd` in `internal/model/app.go`.

**Rationale**: `UpdateLastVisited` is only called from one location (`openBookmarkCmd`). With last-visited tracking removed, it becomes dead code. Leaving dead code in place violates the no-unnecessary-abstractions principle in the constitution. Removing both the call and the function keeps the codebase minimal.

**Alternatives considered**:
- Leave the function in place as a no-op for forward compatibility: No callers will remain; there is no API or external consumer. Rejected.
- Replace with a stub that immediately returns nil: Equally unnecessary. Rejected.

---

## Decision 3: Database Column Retention

**Decision**: Leave the `last_visited_at` column in the `bookmarks` table unchanged (no schema migration to drop it).

**Rationale**: Dropping a column in SQLite requires recreating the table. The spec explicitly documents this assumption to avoid a destructive migration. Existing data in the column does no harm — it will never be read again by the application. The constitution gate requires backward-compatible migrations (ALTER TABLE with DEFAULT only; no destructive changes).

**Alternatives considered**:
- DROP COLUMN via table recreation: Allowed in SQLite 3.35+ but requires copying all rows and risks data loss. Violates the constitution's backward-compatible migration gate. Rejected.

---

## Decision 4: Browse Row Display

**Decision**: Remove both the `"Last: " + date` branch and the `"Never visited"` branch from `BookmarkItem.Description()` in `internal/model/browse.go`. The description line becomes: `domain · YYYY-MM-DD[· #tag1 #tag2]`.

**Rationale**: Both conditional branches in `Description()` relate to last-visited data. Removing both gives a clean, consistent description format. No other view references `LastVisitedAt`.

**Alternatives considered**:
- Replace with a static placeholder like "Visit tracking disabled": Adds noise with no value. Rejected.

---

## Decision 5: openBookmarkCmd Simplification

**Decision**: Remove the `_ = s.UpdateLastVisited(b.ID)` call from `openBookmarkCmd`. The command becomes: open URL → on success, reload bookmarks. On error, return `openURLErrMsg`.

**Rationale**: The function is already documented (feature 003). Removing the update call makes the sequence simpler and eliminates the store dependency on `UpdateLastVisited`. The list reload still happens on success.

**Alternatives considered**:
- Keep the call but pass a no-op implementation: Needless indirection. Rejected.

---

## Files Changed

| File | Change |
|------|--------|
| `internal/store/archive.go` | Simplify `ArchiveStale()` SQL to `created_at <= datetime('now', '-30 days')`; remove `UpdateLastVisited()` |
| `internal/model/browse.go` | Remove last-visited branches from `BookmarkItem.Description()` |
| `internal/model/app.go` | Remove `_ = s.UpdateLastVisited(b.ID)` from `openBookmarkCmd`; update comment |

**No new files. No schema migrations. No new dependencies.**
