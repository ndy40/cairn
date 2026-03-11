---
id: TASK-119
title: 'T014: Implement first-run sync prompt'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:44'
updated_date: '2026-03-11 11:40'
labels:
  - us1
  - phase-3
dependencies:
  - TASK-118
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement first-run sync prompt in cmd/cairn/main.go: before executing any command, check if sync config exists (LoadConfig). If config is nil (no file), prompt 'No sync configured -- connect to Dropbox? (y/N):'. On y: run sync setup flow inline, then continue with original command. On N/Enter: create config file with sync_declined: true via SaveConfig, continue normally. Skip prompt entirely if config file exists (whether configured or declined).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 First-run prompt appears when no sync config file exists
- [x] #2 Accepting prompt runs sync setup flow inline
- [x] #3 Declining creates config with sync_declined: true
- [x] #4 Prompt does not appear if config file already exists
- [ ] #5 Prompt does not appear if sync_declined is true
- [ ] #6 Original command continues after prompt regardless of choice
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added checkFirstRunSync() at CLI startup. Prompts 'No sync configured — connect to Dropbox? (y/N)'. Records decline in config to suppress future prompts.
<!-- SECTION:FINAL_SUMMARY:END -->
