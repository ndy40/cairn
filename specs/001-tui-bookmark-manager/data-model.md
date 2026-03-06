# Data Model: TUI Bookmark Manager

**Feature**: 001-tui-bookmark-manager
**Date**: 2026-03-06

---

## Entities

### Bookmark

The central entity. Represents a single saved web page.

| Field | Type | Nullable | Constraints |
|-------|------|----------|-------------|
| `id` | integer | No | Primary key, auto-increment |
| `url` | text | No | Unique, non-empty, valid URL format |
| `domain` | text | No | Extracted from URL hostname, lowercase |
| `title` | text | No | Non-empty; defaults to URL if fetch fails |
| `description` | text | Yes | Meta description or og:description content |
| `created_at` | datetime | No | UTC timestamp, set at insert time |

**Derived Field**:
- `domain` is always extracted from the `url` at save time (e.g., `https://www.example.com/path` → `example.com`). It is stored separately to enable efficient domain-based search without URL parsing at query time.

**Validation Rules**:
- `url` must be a syntactically valid URL with a scheme (`http` or `https`) and a non-empty host.
- `url` must be unique across all bookmarks (duplicate prevention per FR-005).
- `title` defaults to `url` if the page fetch fails or returns no title.
- `description` is stored as empty string if not available (not NULL, to simplify search queries).
- `created_at` is assigned by the application at save time, not by the user.

**State Transitions**:
Bookmarks are created and deleted. There is no explicit edit or archive state in this version.

```
[created] ──── delete ──── [deleted]
```

---

### SearchQuery (transient — not persisted)

Represents the current user-entered search term during an interactive session.

| Field | Type | Description |
|-------|------|-------------|
| `term` | string | The raw text typed by the user |
| `results` | []Bookmark | Ranked list of matching bookmarks |

This entity exists only in application memory. It is reset each time the user clears the search field or exits the application.

---

## Storage Schema

### Table: `bookmarks`

```sql
CREATE TABLE IF NOT EXISTS bookmarks (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    url         TEXT    NOT NULL UNIQUE,
    domain      TEXT    NOT NULL,
    title       TEXT    NOT NULL,
    description TEXT    NOT NULL DEFAULT '',
    created_at  TEXT    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_bookmarks_domain ON bookmarks(domain);
CREATE INDEX IF NOT EXISTS idx_bookmarks_created_at ON bookmarks(created_at DESC);
```

### Table: `schema_version`

Tracks schema migrations applied at startup.

```sql
CREATE TABLE IF NOT EXISTS schema_version (
    version     INTEGER PRIMARY KEY,
    applied_at  TEXT NOT NULL
);
```

Initial version: `1`.

### FTS5 Virtual Table: `bookmarks_fts`

Enables fast full-text search across title, description, and domain as a pre-filter layer before fuzzy ranking.

```sql
CREATE VIRTUAL TABLE IF NOT EXISTS bookmarks_fts USING fts5(
    title,
    description,
    domain,
    content='bookmarks',
    content_rowid='id'
);
```

FTS triggers keep the virtual table in sync:

```sql
CREATE TRIGGER bookmarks_ai AFTER INSERT ON bookmarks BEGIN
    INSERT INTO bookmarks_fts(rowid, title, description, domain)
    VALUES (new.id, new.title, new.description, new.domain);
END;

CREATE TRIGGER bookmarks_ad AFTER DELETE ON bookmarks BEGIN
    INSERT INTO bookmarks_fts(bookmarks_fts, rowid, title, description, domain)
    VALUES ('delete', old.id, old.title, old.description, old.domain);
END;
```

---

## Relationships

This is a single-entity model. There are no foreign keys or inter-entity relationships in the initial version.

---

## Storage Location

The SQLite database file is stored in the user's OS-appropriate data directory:

| Platform | Path |
|----------|------|
| Linux | `$XDG_DATA_HOME/bookmark-manager/bookmarks.db` or `~/.local/share/bookmark-manager/bookmarks.db` |
| macOS | `~/Library/Application Support/bookmark-manager/bookmarks.db` |
| Windows | `%APPDATA%\bookmark-manager\bookmarks.db` |

The directory is created automatically on first launch if it does not exist.

---

## Search Architecture

Search operates in two stages:

1. **Pre-filter (SQL FTS5)**: When the search term is 3+ characters, query `bookmarks_fts` for candidates matching any of the three fields. Returns a candidate set (typically much smaller than the full collection).
2. **Fuzzy rank (sahilm/fuzzy)**: Apply fuzzy scoring to the candidate set across title (3×), domain (2×), and description (1×). Sort by composite score descending.

For terms under 3 characters, skip FTS5 and fuzzy-rank all bookmarks directly (fast at ≤1,000 records).

This two-stage approach keeps interactive search fast (under 100ms) even as the collection grows beyond 1,000 entries.
