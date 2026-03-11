---
id: TASK-116
title: 'T011: Implement OAuth2 PKCE flow for Dropbox'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:41'
updated_date: '2026-03-11 11:36'
labels:
  - us1
  - phase-3
dependencies:
  - TASK-113
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement OAuth2 PKCE flow for Dropbox in internal/sync/auth.go: implement RunOAuth2Flow(appKey string) (*oauth2.Token, error) using PKCE with no-redirect. Generate code_verifier and code_challenge (S256). Build authorization URL with Dropbox OAuth2 endpoint, print to stdout. Read authorization code from stdin. Exchange code + verifier for token via golang.org/x/oauth2.Config.Exchange(). Return token containing access_token, refresh_token, and expiry.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 RunOAuth2Flow function implemented
- [x] #2 PKCE code_verifier and code_challenge (S256) generated correctly
- [x] #3 Authorization URL printed to stdout
- [x] #4 Authorization code read from stdin
- [x] #5 Code exchanged for token with access_token, refresh_token, and expiry
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Implemented in internal/sync/auth.go
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Created internal/sync/auth.go with RunOAuth2Flow using PKCE S256. Prints auth URL to stdout, reads code from stdin, exchanges for token with offline access.
<!-- SECTION:FINAL_SUMMARY:END -->
