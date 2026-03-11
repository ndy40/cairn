---
id: TASK-124
title: 'T019: Add cairn sync push CLI subcommand'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:45'
updated_date: '2026-03-11 11:41'
labels:
  - us2
  - phase-4
dependencies:
  - TASK-123
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add cairn sync push CLI subcommand in cmd/cairn/main.go: load config, check IsConfigured, construct Dropbox backend and engine. Call engine.Push(). Print 'Pushed N changes. Up to date.' or 'Already up to date.' on success. Print 'Sync not configured. Run cairn sync setup first.' if unconfigured. Exit codes: 0=success, 1=not configured, 3=error.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 sync push subcommand implemented
- [x] #2 Prints change count on success
- [ ] #3 Prints already up to date when no changes
- [ ] #4 Error message when sync not configured
- [ ] #5 Correct exit codes: 0, 1, 3
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added 'cairn sync push' CLI subcommand using openSyncEngine helper.
<!-- SECTION:FINAL_SUMMARY:END -->
