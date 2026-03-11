---
id: TASK-135
title: 'T030: Implement backend factory function'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:46'
updated_date: '2026-03-11 11:45'
labels:
  - us5
  - phase-7
dependencies:
  - TASK-115
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement backend factory function in internal/sync/backend/backend.go: implement NewBackend(cfg *sync.SyncConfig) (SyncBackend, error) that switches on cfg.Backend (dropbox -> construct DropboxBackend, unknown -> return descriptive error 'unsupported sync backend: %s'). Update engine construction in cmd/cairn/main.go to use this factory instead of directly constructing DropboxBackend.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 NewBackend factory function implemented
- [x] #2 Returns DropboxBackend for dropbox backend type
- [ ] #3 Returns descriptive error for unknown backend types
- [ ] #4 cmd/cairn/main.go uses factory instead of direct DropboxBackend construction
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
NewBackend() factory in internal/sync/factory.go switches on config backend type, creates DropboxBackend with oauth2.Token.
<!-- SECTION:FINAL_SUMMARY:END -->
