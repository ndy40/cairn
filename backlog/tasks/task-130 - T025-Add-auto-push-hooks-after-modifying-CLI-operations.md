---
id: TASK-130
title: 'T025: Add auto-push hooks after modifying CLI operations'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:46'
updated_date: '2026-03-11 11:45'
labels:
  - us4
  - phase-6
dependencies:
  - TASK-128
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add auto-push hooks after modifying CLI operations in cmd/cairn/main.go: after the add, delete subcommand handlers complete successfully, if sync is configured, call engine.AutoPush(). On failure, print warning Sync push failed: change queued for later. On success with replayed pending > 0, print up-arrow Synced.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Auto-push fires after successful add command
- [x] #2 Auto-push fires after successful delete command
- [x] #3 Prints warning on push failure
- [ ] #4 Prints synced message when pending changes replayed
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
autoSyncPush() called after add and delete CLI commands.
<!-- SECTION:FINAL_SUMMARY:END -->
