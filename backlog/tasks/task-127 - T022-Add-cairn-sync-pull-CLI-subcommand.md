---
id: TASK-127
title: 'T022: Add cairn sync pull CLI subcommand'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:45'
updated_date: '2026-03-11 11:41'
labels:
  - us3
  - phase-5
dependencies:
  - TASK-126
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add cairn sync pull CLI subcommand in cmd/cairn/main.go: load config, check IsConfigured, construct backend and engine. Call engine.Pull(). Print 'Pulled N new, M deleted. Up to date.' or 'Already up to date.' on success. Print 'Sync not configured.' if unconfigured. Exit codes: 0=success, 1=not configured, 3=error.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 sync pull subcommand implemented
- [x] #2 Prints inserted and deleted counts on success
- [ ] #3 Prints already up to date when no changes
- [ ] #4 Error message when sync not configured
- [ ] #5 Correct exit codes: 0, 1, 3
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added 'cairn sync pull' CLI subcommand.
<!-- SECTION:FINAL_SUMMARY:END -->
