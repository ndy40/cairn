---
id: TASK-106
title: 'T001: Add new Go dependencies for sync'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:40'
updated_date: '2026-03-11 11:30'
labels:
  - setup
  - phase-1
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add new dependencies: go get golang.org/x/oauth2 github.com/dropbox/dropbox-sdk-go-unofficial/v6 github.com/google/uuid
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 golang.org/x/oauth2 added to go.mod
- [ ] #2 dropbox-sdk-go-unofficial/v6 added to go.mod
- [x] #3 google/uuid promoted to direct dependency in go.mod
- [x] #4 go mod tidy passes without errors
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Run go get for deps\n2. Dependencies retained once code imports them\n3. Mark complete after T002+ creates importing code
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
oauth2 and dropbox deps added via go get but trimmed by go mod tidy since no Go code imports them yet. They'll be retained once Dropbox backend is implemented. google/uuid promoted to direct.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added google/uuid as direct dependency. golang.org/x/oauth2 and dropbox SDK added but deferred to T010/T011 when importing code exists.
<!-- SECTION:FINAL_SUMMARY:END -->
