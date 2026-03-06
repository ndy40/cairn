# Data Model: Edit Bookmark Tags, Last-Visited Visibility & CLI Help

**Feature**: 003-edit-bookmark-help
**Date**: 2026-03-06

---

## Entity Changes

**No entity changes.** The `Bookmark` entity and all its fields (including `tags`) are unchanged from feature 002. No schema migration is required.

---

## New Store Operation: UpdateTags

A single new method on the `Store` is added to `internal/store/bookmark.go`:

```
UpdateTags(id int64, tags []string) error
```

**Behaviour**:
- Normalises `tags` using the existing `NormaliseTags()` helper (lowercase, dedup, truncate, max 3).
- JSON-encodes the normalised slice.
- Executes: `UPDATE bookmarks SET tags = ? WHERE id = ?`
- Returns an error if the update fails; returns nil on success.
- Does NOT modify any other column.

**Validation**:
- `id` must refer to an existing bookmark; no validation is performed (callers obtain IDs from the list).
- `tags` follows the same rules as `Insert()`: comma-split by the caller, then `NormaliseTags()` applied inside `UpdateTags`.

---

## No State Transition Changes

The bookmark state machine (active → archived → active via restore; permanent flag toggle) is unchanged. Tag editing is orthogonal to all lifecycle states.

---

## Confirmed: Last-Visited Update Flow

The existing `UpdateLastVisited(id int64) error` method in `internal/store/archive.go` is called by `openBookmarkCmd` in `internal/model/app.go` after a successful browser launch. No changes needed. The flow is:

```
User presses Enter on bookmark
  → openBookmarkCmd fires (tea.Cmd)
    → openURLRaw(url) — starts browser process
    → on success: s.UpdateLastVisited(b.ID) — writes datetime('now') to DB
    → returns loadBookmarks(s)() — triggers list reload
      → browse list shows updated "Last: <date>"
```

This flow covers both browse mode and search mode (`Enter` in both routes to `openBookmarkCmd`).
