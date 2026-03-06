---
id: TASK-6
title: Implement SQLite store open and schema migration runner
status: Done
assignee: []
created_date: '2026-03-06 04:04'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:foundation'
dependencies: []
priority: high
ordinal: 83000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create the core database lifecycle logic: open the SQLite file at the correct OS-specific path, enable WAL mode for concurrent reads, and run schema migrations on startup. This is the foundation every other store operation depends on.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Store opens SQLite at XDG_DATA_HOME/bookmark-manager/bookmarks.db on Linux, ~/Library/Application Support/bookmark-manager/bookmarks.db on macOS
- [ ] #2 Parent directory is created automatically with os.MkdirAll if it does not exist
- [ ] #3 WAL mode is enabled: PRAGMA journal_mode=WAL
- [ ] #4 schema_version table is created on first open
- [ ] #5 Migration runner checks schema_version and applies new migrations in order
<!-- AC:END -->
