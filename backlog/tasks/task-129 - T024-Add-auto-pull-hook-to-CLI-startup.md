---
id: TASK-129
title: 'T024: Add auto-pull hook to CLI startup'
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
Add auto-pull hook to CLI startup in cmd/cairn/main.go: after the first-run prompt check and before executing the main command, if sync is configured, call engine.AutoPull(). On success with changes > 0, print down-arrow N new bookmarks synced. On failure, print warning Sync pull failed: reason. Continue with the original command regardless.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Auto-pull fires on CLI startup when sync is configured
- [x] #2 Prints count of new bookmarks when changes pulled
- [ ] #3 Prints warning on sync failure
- [ ] #4 Original command always continues regardless of sync result
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
autoSyncPull() called at CLI startup, non-fatal with stderr warnings.
<!-- SECTION:FINAL_SUMMARY:END -->
