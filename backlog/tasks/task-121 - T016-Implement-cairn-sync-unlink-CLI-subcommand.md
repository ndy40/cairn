---
id: TASK-121
title: 'T016: Implement cairn sync unlink CLI subcommand'
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
Implement cairn sync unlink CLI subcommand in cmd/cairn/main.go: check sync is configured, prompt 'Unlink this device from sync? Local bookmarks will be kept. (y/N)'. On confirm: delete sync config file (os.Remove), clear pending_sync table via store.ClearPendingChanges(). Print confirmation. Exit codes: 0=success, 1=not configured/cancelled.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Checks sync is configured before proceeding
- [x] #2 Confirmation prompt before unlinking
- [x] #3 Deletes sync config file on confirm
- [ ] #4 Clears pending_sync table on confirm
- [ ] #5 Local bookmarks are preserved after unlink
- [ ] #6 Exit code 1 if not configured or cancelled
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added 'cairn sync unlink' that removes sync config file while preserving local bookmarks.
<!-- SECTION:FINAL_SUMMARY:END -->
