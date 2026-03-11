---
id: TASK-139
title: 'T034: Add sync help text to CLI output'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:47'
updated_date: '2026-03-11 11:45'
labels:
  - polish
  - phase-8
dependencies:
  - TASK-118
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add sync help text to cmd/cairn/main.go help output: update the help subcommand to include sync subcommands (sync setup, sync push, sync pull, sync status, sync auth, sync unlink) with brief descriptions.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Help output includes all six sync subcommands
- [x] #2 Each subcommand has a brief description
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Updated main printHelp() with sync commands section and added printSyncHelp() for 'cairn sync' subcommand.
<!-- SECTION:FINAL_SUMMARY:END -->
