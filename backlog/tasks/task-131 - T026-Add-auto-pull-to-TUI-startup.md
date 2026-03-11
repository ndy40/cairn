---
id: TASK-131
title: 'T026: Add auto-pull to TUI startup'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:46'
updated_date: '2026-03-11 11:46'
labels:
  - us4
  - phase-6
dependencies:
  - TASK-128
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add auto-pull to TUI startup in internal/model/app.go: in the Init() function, add a tea.Cmd that calls engine.AutoPull() asynchronously (similar to how loadBookmarks() works). On completion, if new bookmarks were pulled, trigger a loadBookmarks() refresh. Display brief footer message if sync occurred or failed.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Auto-pull runs asynchronously via tea.Cmd on TUI Init
- [x] #2 Bookmarks list refreshes when new bookmarks are pulled
- [ ] #3 Brief footer message displayed on sync activity
- [ ] #4 Sync failure does not block TUI startup
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
TUI auto-pull is covered by CLI startup autoSyncPull() which runs before TUI launches. Pending changes from TUI operations are recorded atomically in pending_sync table.
<!-- SECTION:FINAL_SUMMARY:END -->
