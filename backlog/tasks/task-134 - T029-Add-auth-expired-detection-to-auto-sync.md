---
id: TASK-134
title: 'T029: Add auth-expired detection to auto-sync'
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
Add auth-expired detection to auto-sync in internal/sync/engine.go: in AutoPull() and AutoPush(), if the backend returns ErrAuthExpired, print warning 'Sync auth expired -- run cairn sync auth to reconnect' instead of the generic failure message. Queue changes locally as normal.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 AutoPull detects ErrAuthExpired and prints specific message
- [x] #2 AutoPush detects ErrAuthExpired and prints specific message
- [ ] #3 Changes queued locally on auth expiry
- [ ] #4 Generic failure message not shown for auth expiry
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
isAuthError() helper in engine.go detects auth expiry. AutoPull/AutoPush show specific re-auth message.
<!-- SECTION:FINAL_SUMMARY:END -->
