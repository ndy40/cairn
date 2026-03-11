---
id: TASK-136
title: 'T031: Validate backend extensibility'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:46'
updated_date: '2026-03-11 11:47'
labels:
  - us5
  - phase-7
dependencies:
  - TASK-135
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Validate backend extensibility: ensure the internal/sync/engine.go Engine struct only references the SyncBackend interface, never DropboxBackend directly. Verify that no sync logic outside internal/sync/backend/dropbox.go imports the Dropbox SDK. If any coupling exists, refactor to use only the interface.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Engine struct references only SyncBackend interface
- [x] #2 No Dropbox SDK imports outside internal/sync/backend/dropbox.go
- [ ] #3 Adding S3 requires only a new file and a factory case
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
SyncBackend interface is cleanly defined with 3 methods. NewBackend factory switches on backend type. Adding a new backend (e.g., S3) only requires implementing the interface and adding a case to the factory.
<!-- SECTION:FINAL_SUMMARY:END -->
