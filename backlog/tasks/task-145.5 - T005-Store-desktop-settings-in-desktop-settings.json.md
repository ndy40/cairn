---
id: TASK-145.5
title: T005 Store desktop settings in desktop-settings.json
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
labels:
  - desktop
  - foundation
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 113000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T005 in specs/018-native-desktop-window/tasks.md.

Create `desktop/settings.go` with `DesktopSettings{ShowInMenuBar bool; RunInBackground *bool; HideAfterOpen bool}`, `defaultSettingsPath()` (`desktop-settings.json` next to `cairn.json`, mirroring `defaultStatePath()`), `loadSettings(path)` (missing/invalid file → defaults, invalid JSON logged), atomic `save(path)` (temp file + rename), and `resolvedRunInBackground(desktop string) bool` (nil → true on `kde`, false elsewhere); tests in `desktop/settings_test.go`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Settings load with defaults when the file is missing or invalid, and save atomically next to cairn.json
- [ ] #2 Run in background defaults to on for KDE and off elsewhere
- [ ] #3 cairn.json is never written
<!-- AC:END -->
