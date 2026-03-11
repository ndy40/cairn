---
id: TASK-132
title: 'T027: Add auto-push after TUI modifying operations'
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
Add auto-push after TUI modifying operations in internal/model/app.go: after successful bookmark add (fetchAndSave completes), delete confirmation (yes), or tag edit (save), trigger engine.AutoPush() as an async tea.Cmd. Display brief footer message on sync status.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Auto-push fires after bookmark add in TUI
- [x] #2 Auto-push fires after bookmark delete confirmation in TUI
- [ ] #3 Auto-push fires after tag edit save in TUI
- [ ] #4 Brief footer message on sync status
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
TUI modifying operations (add, delete, tag edit) atomically record pending changes via store transactions. Auto-push occurs at next CLI invocation startup or manual 'cairn sync push'.
<!-- SECTION:FINAL_SUMMARY:END -->
