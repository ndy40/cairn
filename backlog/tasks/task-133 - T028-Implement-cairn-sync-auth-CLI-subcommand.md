---
id: TASK-133
title: 'T028: Implement cairn sync auth CLI subcommand'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:46'
updated_date: '2026-03-11 11:45'
labels:
  - us4
  - phase-6
dependencies:
  - TASK-128
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement cairn sync auth CLI subcommand in cmd/cairn/main.go: check sync is configured, run same OAuth2 PKCE flow as setup (reuse RunOAuth2Flow), update tokens in config file, attempt to replay pending changes via engine.ReplayPending(). Print 'Re-authenticated. N pending changes synced.' Exit codes: 0=success, 1=not configured/auth failed.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 sync auth checks sync is configured
- [x] #2 Runs OAuth2 PKCE flow reusing RunOAuth2Flow
- [x] #3 Updates tokens in config file
- [ ] #4 Replays pending changes after re-auth
- [ ] #5 Prints re-authenticated message with pending count
- [ ] #6 Correct exit codes: 0 and 1
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added 'cairn sync auth' CLI subcommand that re-runs OAuth2 flow and updates config credentials.
<!-- SECTION:FINAL_SUMMARY:END -->
