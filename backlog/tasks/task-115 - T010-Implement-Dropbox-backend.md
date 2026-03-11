---
id: TASK-115
title: 'T010: Implement Dropbox backend'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:41'
updated_date: '2026-03-11 11:36'
labels:
  - us1
  - phase-3
dependencies:
  - TASK-111
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement Dropbox backend in internal/sync/backend/dropbox.go: create DropboxBackend struct implementing SyncBackend interface. Constructor takes oauth2 token and app key. Use dropbox-sdk-go-unofficial/v6 for Upload() (files/upload with WriteMode.Overwrite at /cairn/sync.json), Download() (files/download), and Exists() (files/get_metadata, return ErrNotFound on path_not_found). Wrap HTTP client with golang.org/x/oauth2.TokenSource for automatic token refresh. Map Dropbox API errors to sentinel errors (ErrAuthExpired, ErrNetworkFailure, ErrQuotaExceeded).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 DropboxBackend struct implements SyncBackend interface
- [x] #2 Upload writes file to /cairn/sync.json with overwrite mode
- [x] #3 Download retrieves file contents from /cairn/sync.json
- [x] #4 Exists checks file metadata and returns ErrNotFound when missing
- [x] #5 HTTP client wrapped with oauth2.TokenSource for auto-refresh
- [x] #6 Dropbox API errors mapped to sentinel errors
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Implemented in internal/sync/backend/dropbox.go
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Created internal/sync/backend/dropbox.go implementing SyncBackend. Uses dropbox-sdk-go-unofficial/v6 for Upload (overwrite mode), Download, and Exists. OAuth2 token auto-refresh via x/oauth2. Error mapping to sentinel errors.
<!-- SECTION:FINAL_SUMMARY:END -->
