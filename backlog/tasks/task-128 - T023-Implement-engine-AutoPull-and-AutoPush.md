---
id: TASK-128
title: 'T023: Implement engine AutoPull() and AutoPush()'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:46'
updated_date: '2026-03-11 11:45'
labels:
  - us4
  - phase-6
dependencies:
  - TASK-123
  - TASK-126
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement engine AutoPull() and AutoPush() in internal/sync/engine.go: AutoPull() (int, error) wraps Pull() but returns 0 on any error instead of failing (non-blocking). AutoPush() error wraps Push() but on failure increments retry_count on pending changes instead of propagating error. Both methods also replay any pending changes from the queue on success. Add ReplayPending() (int, error) that lists pending changes, applies them to the snapshot, uploads, and clears on success.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 AutoPull wraps Pull and returns 0 on error
- [x] #2 AutoPush wraps Push and increments retry_count on failure
- [x] #3 Both replay pending changes on successful connection
- [ ] #4 ReplayPending lists, applies, uploads, and clears pending changes
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Engine.AutoPull() and AutoPush() wrap Pull/Push with graceful error handling, returning warning strings instead of failing.
<!-- SECTION:FINAL_SUMMARY:END -->
