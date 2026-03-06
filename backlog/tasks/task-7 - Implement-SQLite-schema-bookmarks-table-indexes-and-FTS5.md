---
id: TASK-7
title: 'Implement SQLite schema: bookmarks table, indexes, and FTS5'
status: Done
assignee: []
created_date: '2026-03-06 04:04'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:foundation'
dependencies: []
priority: high
ordinal: 84000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Define the full database schema: the bookmarks table, domain and date indexes for fast sorting, and the FTS5 virtual table with sync triggers. This schema underpins all bookmark storage and full-text search.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 bookmarks table exists with columns: id (PK autoincrement), url (TEXT UNIQUE NOT NULL), domain (TEXT NOT NULL), title (TEXT NOT NULL), description (TEXT NOT NULL DEFAULT ''), created_at (TEXT NOT NULL)
- [ ] #2 Index on domain column exists
- [ ] #3 Index on created_at DESC exists
- [ ] #4 FTS5 virtual table bookmarks_fts exists with columns: title, description, domain; content='bookmarks'; content_rowid='id'
- [ ] #5 After-insert trigger populates bookmarks_fts when a bookmark is inserted
- [ ] #6 After-delete trigger removes entry from bookmarks_fts when a bookmark is deleted
<!-- AC:END -->
