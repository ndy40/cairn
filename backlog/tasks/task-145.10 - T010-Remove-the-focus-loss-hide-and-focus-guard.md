---
id: TASK-145.10
title: T010 Remove the focus-loss hide and focus guard
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:08:00Z'
labels:
  - desktop
  - us1
dependencies:
  - TASK-145.9
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 118000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T010 in specs/018-native-desktop-window/tasks.md.

Delete `desktop/focus.go` and `desktop/focus_test.go`; remove the `WindowLostFocus`/`WindowFocus` handlers, `hideIfUnfocused`, the `focus` field and `newFocusGuard()` from `desktop/launcher.go` and `desktop/main.go`; move `debugf` into `desktop/launcher.go` if still used; remove the `WindowService.PointerDown()` call from `desktop/frontend/src/main.ts`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The window stays open when another app takes focus
- [ ] #2 focus.go and its pointer-down workaround are gone
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implemented in ndw-core, merged in 176f4f7. Code complete; awaiting manual QA (quickstart.md rows 1-10) before ticking criteria.
<!-- SECTION:NOTES:END -->
