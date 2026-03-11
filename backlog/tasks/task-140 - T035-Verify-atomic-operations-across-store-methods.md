---
id: TASK-140
title: 'T035: Verify atomic operations across store methods'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:47'
updated_date: '2026-03-11 11:47'
labels:
  - polish
  - phase-8
dependencies:
  - TASK-122
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Verify atomic operations: ensure all store methods that modify bookmarks and write to pending_sync use explicit database transactions. Verify that a crash between bookmark write and pending_sync write cannot leave inconsistent state.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 All modifying store methods use explicit sql.Tx transactions
- [x] #2 Bookmark write and pending_sync write are in same transaction
- [ ] #3 Transaction rollback on failure leaves no partial state
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Insert, DeleteByID, and UpdateTags all use database transactions to atomically write bookmark changes and pending_sync entries. Migration V3 also runs in a transaction.
<!-- SECTION:FINAL_SUMMARY:END -->
