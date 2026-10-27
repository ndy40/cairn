---
id: TASK-145.8
title: T008 Apply platform attributes and design tokens in the frontend
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
labels:
  - desktop
  - foundation
dependencies:
  - TASK-145.7
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 116000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T008 in specs/018-native-desktop-window/tasks.md.

At startup in `desktop/frontend/src/main.ts`, call `WindowService.Platform()` and set `data-os` and `data-desktop` on `document.documentElement`; in `desktop/frontend/src/style.css` replace the `:root` palette with the token set from the design (`--accent`, `--text`, `--text-2`, `--list`, `--hover`, `--warn-bg`, `--warn-fg`, light and dark, fallback accent `#1766d1` / `#2f6fd6`)
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The page root carries data-os and data-desktop for the running platform
- [ ] #2 Light and dark colour tokens from the design replace the old palette
<!-- AC:END -->
