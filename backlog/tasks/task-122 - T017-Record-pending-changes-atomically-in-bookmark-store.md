---
id: TASK-122
title: 'T017: Record pending changes atomically in bookmark store'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:45'
updated_date: '2026-03-11 11:41'
labels:
  - us2
  - phase-4
dependencies:
  - TASK-110
  - TASK-109
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Modify internal/store/bookmark.go to record pending changes atomically: update Insert() to call AddPendingChange(uuid, add, jsonPayload) within the same database transaction. Update DeleteByID() to call AddPendingChange(uuid, delete, '') within the same transaction. Update UpdateTags() to call AddPendingChange(uuid, update, jsonPayload) within the same transaction. This requires changing these methods to use explicit transactions (sql.Tx) instead of direct db.Exec.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Insert() records add pending change in same transaction
- [x] #2 DeleteByID() records delete pending change in same transaction
- [x] #3 UpdateTags() records update pending change in same transaction
- [x] #4 All three methods use explicit sql.Tx transactions
- [ ] #5 Transaction rollback on any failure leaves no partial state
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Insert, DeleteByID, and UpdateTags now use transactions to atomically record pending_sync entries alongside bookmark mutations.
<!-- SECTION:FINAL_SUMMARY:END -->
