---
id: TASK-144
title: Connect Dropbox sync from the desktop app
status: To Do
assignee: []
created_date: '2026-10-20 11:44'
labels:
  - desktop
  - sync
  - feature
dependencies:
  - TASK-143
priority: medium
ordinal: 107000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Dropbox sync setup currently only works in a terminal (cairn sync setup): RunOAuth2Flow prints the auth URL and reads the code from stdin, so cairn-desktop can use an existing sync config but cannot create one. Let desktop users connect Dropbox without a terminal.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 From the desktop app (popup or tray menu) the user can start Dropbox setup, which opens the Dropbox authorization page in the default browser
- [ ] #2 User can paste the authorization code into the desktop app to complete setup; the initial sync runs and its result is shown
- [ ] #3 If dropbox_app_key is not configured, the desktop app asks for it and saves it to cairn.json
- [ ] #4 Invalid or expired codes and network failures show a clear error and leave sync unconfigured
- [ ] #5 cairn sync setup in the CLI behaves exactly as before
- [ ] #6 User guide docs describe connecting Dropbox from the desktop app
<!-- AC:END -->
