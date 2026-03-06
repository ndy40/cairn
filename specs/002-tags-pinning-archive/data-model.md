# Data Model: Tags, Pinning, Archive & Startup Checks

**Feature**: 002-tags-pinning-archive
**Date**: 2026-03-06
**Extends**: `specs/001-tui-bookmark-manager/data-model.md`

---

## Entities

### Bookmark (extended)

The `Bookmark` entity gains five new fields in schema migration v2. All existing fields (id, url, domain, title, description, created_at) are unchanged.

| Field | Type | Nullable | Default | Constraints |
|-------|------|----------|---------|-------------|
| `id` | integer | No | — | Primary key, auto-increment (unchanged) |
| `url` | text | No | — | Unique, non-empty (unchanged) |
| `domain` | text | No | — | Lowercase hostname (unchanged) |
| `title` | text | No | — | Non-empty (unchanged) |
| `description` | text | No | `''` | Meta description (unchanged) |
| `created_at` | text | No | — | UTC ISO-8601, set at insert (unchanged) |
| `tags` | text | No | `'[]'` | JSON array of 0–3 lowercase tag strings, each 1–32 chars |
| `last_visited_at` | text | Yes | `NULL` | UTC ISO-8601; NULL means never visited |
| `is_permanent` | integer | No | `0` | Boolean (0/1); 1 = exempt from auto-archiving |
| `is_archived` | integer | No | `0` | Boolean (0/1); 1 = moved to archive |
| `archived_at` | text | Yes | `NULL` | UTC ISO-8601; NULL when not archived |

**Validation Rules** (new fields):

- `tags`: must be a valid JSON array. Each element must be a non-empty, non-whitespace string of 1–32 characters, stored lowercase. Duplicates are silently removed before saving. Array length must not exceed 3 elements.
- `last_visited_at`: set to `datetime('now')` whenever the bookmark is opened via the TUI. Never set by CLI subcommands.
- `is_permanent`: toggled by the user at any time. Removing the flag takes effect on the next startup archive check with no grace period.
- `is_archived`: set to 1 by the startup archive check when the bookmark is eligible. Set back to 0 by a user restore action.
- `archived_at`: set to `datetime('now')` when `is_archived` is set to 1. Cleared (set to NULL) when restored.

**State Transitions** (updated):

```
[active] ──── open ──────────────────────────────── [active, last_visited_at updated]
[active] ──── p key (toggle) ──────────────────────  [active, is_permanent toggled]
[active] ──── startup archive check (eligible) ────  [archived]
[archived] ── restore (r key in archive view) ─────  [active, archived_at cleared]
[active] ──── delete ───────────────────────────────  [deleted]
[archived] ── delete ───────────────────────────────  [deleted]
```

**Archive Eligibility**:

A bookmark is eligible for archiving if ALL of the following are true:
1. `is_permanent = 0`
2. `is_archived = 0`
3. Either:
   - `last_visited_at IS NOT NULL AND last_visited_at <= datetime('now', '-183 days')`
   - `last_visited_at IS NULL AND created_at <= datetime('now', '-183 days')`

---

### Tag (value object — not a standalone entity)

Tags are stored inline on the `Bookmark` as a JSON array. They are not persisted as a separate entity.

| Attribute | Constraint |
|-----------|------------|
| Type | string |
| Minimum length | 1 character (after trimming whitespace) |
| Maximum length | 32 characters (truncated at input level) |
| Case | Always lowercase (normalised at save time) |
| Deduplication | Duplicate values silently removed before saving |
| Maximum count | 3 per bookmark |

**Example stored value**: `'["work","go","tools"]'`

---

### CheckResult (transient — not persisted)

Returned by `display.CheckPrerequisites()` at startup; used only to decide whether to block or warn.

| Field | Type | Description |
|-------|------|-------------|
| `DisplayType` | enum (Wayland, X11, Unknown) | Detected display environment |
| `ToolFound` | bool | Whether the required clipboard tool is present |
| `MissingTool` | string | Name of the missing tool (empty if found) |
| `InstallHint` | string | Human-readable install instruction (empty if no missing tool) |
| `ShouldBlock` | bool | True when display detected but tool absent; false for Unknown display |

---

## Storage Schema

### Migration v2: Extend `bookmarks` Table

Migration v2 uses `ALTER TABLE ... ADD COLUMN` for each new column. This is backward-compatible — existing rows receive the column defaults automatically.

```sql
-- Migration v2
ALTER TABLE bookmarks ADD COLUMN tags            TEXT    NOT NULL DEFAULT '[]';
ALTER TABLE bookmarks ADD COLUMN last_visited_at TEXT;
ALTER TABLE bookmarks ADD COLUMN is_permanent    INTEGER NOT NULL DEFAULT 0;
ALTER TABLE bookmarks ADD COLUMN is_archived     INTEGER NOT NULL DEFAULT 0;
ALTER TABLE bookmarks ADD COLUMN archived_at     TEXT;
```

No changes to existing indexes, FTS5 table, or triggers from migration v1.

### New Indexes

```sql
CREATE INDEX IF NOT EXISTS idx_bookmarks_is_archived  ON bookmarks(is_archived);
CREATE INDEX IF NOT EXISTS idx_bookmarks_archived_at  ON bookmarks(archived_at DESC);
CREATE INDEX IF NOT EXISTS idx_bookmarks_is_permanent ON bookmarks(is_permanent);
```

### Archive Eligibility Query

Used by `store.ArchiveStale()`:

```sql
UPDATE bookmarks
SET    is_archived = 1,
       archived_at = datetime('now')
WHERE  is_permanent = 0
  AND  is_archived  = 0
  AND  (
         (last_visited_at IS NOT NULL AND last_visited_at <= datetime('now', '-183 days'))
         OR
         (last_visited_at IS NULL AND created_at <= datetime('now', '-183 days'))
       );
```

### Active Bookmark List Query

The `List()` method gains an `is_archived = 0` filter:

```sql
SELECT id, url, domain, title, description, created_at,
       tags, last_visited_at, is_permanent, is_archived, archived_at
FROM   bookmarks
WHERE  is_archived = 0
ORDER BY created_at DESC;
```

### Archived Bookmark List Query

Used by `ListArchived()`:

```sql
SELECT id, url, domain, title, description, created_at,
       tags, last_visited_at, is_permanent, is_archived, archived_at
FROM   bookmarks
WHERE  is_archived = 1
ORDER BY archived_at DESC;
```

---

## Go Struct (updated)

```go
type Bookmark struct {
    ID            int64
    URL           string
    Domain        string
    Title         string
    Description   string
    CreatedAt     time.Time
    Tags          []string   // decoded from JSON; empty slice when no tags
    LastVisitedAt *time.Time // nil = never visited
    IsPermanent   bool
    IsArchived    bool
    ArchivedAt    *time.Time // nil when not archived
}
```

---

## In-Memory Tag Filtering

Tag filtering is applied in the `App` model, not in SQL. The filter pipeline is:

```
allBookmarks (active only, loaded at startup)
  → tagFilter(activeTagFilter []string)  [OR logic: include if bookmark has any selected tag]
  → twoStageSearch(term string)          [FTS5 pre-filter + fuzzy rank]
  → browse.load(filtered)
```

Tag filter predicate (OR logic):
```go
func tagFilter(bookmarks []*store.Bookmark, tags []string) []*store.Bookmark {
    if len(tags) == 0 {
        return bookmarks
    }
    tagSet := make(map[string]bool, len(tags))
    for _, t := range tags { tagSet[t] = true }
    var out []*store.Bookmark
    for _, b := range bookmarks {
        for _, bt := range b.Tags {
            if tagSet[bt] {
                out = append(out, b)
                break
            }
        }
    }
    return out
}
```
