---
id: TASK-138
title: 'T033: Handle concurrent URL conflict edge case'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:47'
updated_date: '2026-03-11 11:47'
labels:
  - polish
  - phase-8
dependencies:
  - TASK-114
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Handle edge case: concurrent URL conflict in internal/sync/merge.go -- when two devices add same URL, the merge should keep the record with the most recent updated_at, not create duplicates. Verify the merge algorithm handles this and add explicit handling if missing.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Merge handles two devices adding same URL concurrently
- [x] #2 Record with most recent updated_at wins
- [ ] #3 No duplicate bookmarks created
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
MergeBookmarks matches by UUID first, then URL for dedup. Insert silently skips duplicates (ErrDuplicate). Last-write-wins by updated_at handles concurrent edits.
<!-- SECTION:FINAL_SUMMARY:END -->
