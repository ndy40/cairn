---
id: TASK-145.7
title: T007 Add WindowService bindings for platform and settings
status: Done
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:08:00Z'
labels:
  - desktop
  - foundation
dependencies:
  - TASK-145.4
  - TASK-145.5
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 115000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T007 in specs/018-native-desktop-window/tasks.md.

Create `desktop/window_service.go`: move `WindowService` out of `focus.go`, give it fields for the app, launcher, settings path/settings and platform info, and add `Platform() PlatformInfo`, `Settings() DesktopSettings` (with `runInBackground` resolved) and `SetSetting(key string, value bool) error` (persist only; side effects added in US3/US5); register it in `desktop/main.go`; run `wails3 generate bindings` so `frontend/bindings/` includes the new methods
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 The frontend can read platform info and settings and change a setting through generated bindings
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented in ndw-core, merged in 176f4f7. WindowService exposes Platform, Settings and SetSetting; bindings regenerated.
<!-- SECTION:FINAL_SUMMARY:END -->
