---
id: TASK-11
title: Implement CLI entry point and subcommand dispatcher in cmd/bm/main.go
status: Done
assignee: []
created_date: '2026-03-06 04:05'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:foundation'
dependencies: []
priority: high
ordinal: 88000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create the main binary entry point. When run with no arguments it launches the TUI; subcommands (add, list, search, delete, version, help) run non-interactively. Also resolves the database path from --db flag or BM_DB_PATH env var.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Running ./bm with no arguments opens the store and launches the bubbletea TUI
- [ ] #2 --db flag and BM_DB_PATH env var are both recognized for overriding the default database path; --db takes precedence
- [ ] #3 Unrecognized subcommands print a usage error and exit with code 3
- [ ] #4 Store open failure is reported to stderr and exits with code 3
<!-- AC:END -->
