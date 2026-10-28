---
id: TASK-145.46
title: T046 Rewrite the Desktop App user guide
status: Done
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:34:00Z'
labels:
  - desktop
  - polish
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 154000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T046 in specs/018-native-desktop-window/tasks.md.

Rewrite `docs/content/docs/desktop-app.md` for the window app: opening, window behaviour, close/quit per platform, menus, Add Bookmark, Settings (start at login, menu bar, background, hide after open), global shortcut and the Wayland fallback, keyboard map, `--background`/`--hidden`, `desktop-settings.json`, updated uninstall steps
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 docs/content/docs/desktop-app.md describes the window, close/quit per platform, menus, Add Bookmark, Settings, the global shortcut and Wayland fallback, the keyboard map and flags
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Rewrote docs/content/docs/desktop-app.md for the native window (install, window, close/quit per platform, shortcut and Wayland fallback, add, menus, settings and open at login, keyboard map, sync messages, flags, files, update, uninstall) and replaced popup wording in _index, quickstart, sync, configuration, installation and architecture pages. Checked against the merged menus code. ndw-docs, merged into 018-native-desktop-window.
<!-- SECTION:FINAL_SUMMARY:END -->
