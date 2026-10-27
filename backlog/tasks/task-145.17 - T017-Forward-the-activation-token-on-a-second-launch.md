---
id: TASK-145.17
title: T017 Forward the activation token on a second launch
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:08:00Z'
labels:
  - desktop
  - us2
dependencies:
  - TASK-145.2
  - TASK-145.16
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 125000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T017 in specs/018-native-desktop-window/tasks.md.

In `desktop/main.go`: set `SingleInstanceOptions.AdditionalData` to `{"activation_token": os.Getenv("XDG_ACTIVATION_TOKEN"), "startup_id": os.Getenv("DESKTOP_STARTUP_ID")}`; add pure `activationToken(data map[string]string) string` (prefers activation_token) with a test in `desktop/lifecycle_test.go`; in `Launcher.HandleArgs` route no args → `BringForward("search")`, `--toggle` → `Toggle`, `--add` → `BringForward("add")`, all passing the token
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 A second launch forwards its activation token and its flags reach the running instance (none: bring forward, --toggle: toggle, --add: add)
- [x] #2 Token extraction is unit tested
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implemented in ndw-core, merged in 176f4f7. Token extraction unit tested. Forwarding and flag routing await manual QA (quickstart rows 8-10).
<!-- SECTION:NOTES:END -->
