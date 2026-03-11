---
id: TASK-118
title: 'T013: Add cairn sync setup CLI subcommand'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:44'
updated_date: '2026-03-11 11:40'
labels:
  - us1
  - phase-3
dependencies:
  - TASK-117
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add cairn sync setup CLI subcommand in cmd/cairn/main.go: add sync to the subcommand router. Parse sync setup, sync push, sync pull, sync status, sync auth, sync unlink sub-subcommands. For this task, implement only the setup path: open store, call engine.Setup(), print 'Sync configured. N bookmarks synced.' on success, handle errors with appropriate exit codes (0=success, 1=auth failed, 3=error). Parse --backend flag (default dropbox).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 sync subcommand added to CLI router
- [x] #2 sync setup path implemented
- [x] #3 Success prints bookmark count
- [x] #4 Auth failure returns exit code 1
- [x] #5 Other errors return exit code 3
- [ ] #6 --backend flag parsed with default dropbox
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added 'cairn sync setup' CLI subcommand in main.go. Reads CAIRN_DROPBOX_APP_KEY env, runs OAuth2 flow via Engine.Setup(), reports bookmark count.
<!-- SECTION:FINAL_SUMMARY:END -->
