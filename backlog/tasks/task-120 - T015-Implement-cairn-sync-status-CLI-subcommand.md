---
id: TASK-120
title: 'T015: Implement cairn sync status CLI subcommand'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:45'
updated_date: '2026-03-11 11:40'
labels:
  - us1
  - phase-3
dependencies:
  - TASK-118
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement cairn sync status CLI subcommand in cmd/cairn/main.go: load sync config, query pending change count from store. Print formatted text output (Sync: configured/not configured, Backend, Device ID, Last sync, Pending changes). Support --json flag for JSON output. Exit code always 0.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 sync status shows configured/not configured state
- [x] #2 Shows Backend, Device ID, Last sync, Pending changes when configured
- [x] #3 Shows helpful message when not configured
- [x] #4 --json flag outputs JSON format
- [ ] #5 Exit code is always 0
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added 'cairn sync status' showing backend, device ID, last sync time, and pending change count.
<!-- SECTION:FINAL_SUMMARY:END -->
