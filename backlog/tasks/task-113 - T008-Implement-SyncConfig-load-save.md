---
id: TASK-113
title: 'T008: Implement SyncConfig load/save'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:41'
updated_date: '2026-03-11 11:33'
labels:
  - foundational
  - phase-2
dependencies:
  - TASK-107
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement SyncConfig load/save in internal/sync/config.go: define SyncConfig struct matching data-model.md (Backend, DeviceID, LastSyncAt, SyncDeclined, Dropbox credentials). Implement DefaultConfigPath() string (XDG_CONFIG_HOME/cairn/sync.json on Linux, ~/Library/Application Support/cairn/sync.json on macOS). Implement LoadConfig() (*SyncConfig, error) (returns nil if file does not exist), SaveConfig(cfg *SyncConfig) error (write JSON with mode 0600), IsConfigured(cfg *SyncConfig) bool (true if backend is set and not just declined).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 SyncConfig struct defined with Backend, DeviceID, LastSyncAt, SyncDeclined, Dropbox fields
- [x] #2 DefaultConfigPath returns OS-appropriate path
- [x] #3 LoadConfig returns nil when file does not exist
- [x] #4 LoadConfig correctly parses existing config file
- [x] #5 SaveConfig writes JSON with mode 0600
- [x] #6 IsConfigured returns true only when backend is set and not just declined
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Implemented in internal/sync/config.go
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Created internal/sync/config.go with SyncConfig and DropboxConfig structs, DefaultConfigPath (Linux/macOS), LoadConfig (nil on missing file), SaveConfig (0600 perms), and IsConfigured helper.
<!-- SECTION:FINAL_SUMMARY:END -->
