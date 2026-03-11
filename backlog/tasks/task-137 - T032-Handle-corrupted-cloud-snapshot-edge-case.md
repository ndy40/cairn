---
id: TASK-137
title: 'T032: Handle corrupted cloud snapshot edge case'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:47'
updated_date: '2026-03-11 11:47'
labels:
  - polish
  - phase-8
dependencies:
  - TASK-126
  - TASK-117
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Handle edge case: corrupted cloud snapshot in internal/sync/engine.go -- if UnmarshalSnapshot fails during Pull() or Setup(), return clear error suggesting re-initialising sync from a trusted device. Do not apply partial state.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 UnmarshalSnapshot failure in Pull returns clear error message
- [x] #2 UnmarshalSnapshot failure in Setup returns clear error message
- [ ] #3 No partial state applied on corruption
- [ ] #4 Error message suggests re-initialising sync
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
pullAndMerge handles JSON unmarshal errors by logging a warning and re-uploading local data as a fresh snapshot.
<!-- SECTION:FINAL_SUMMARY:END -->
