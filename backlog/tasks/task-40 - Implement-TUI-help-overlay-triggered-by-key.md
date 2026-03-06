---
id: TASK-40
title: Implement TUI help overlay triggered by ? key
status: Done
assignee: []
created_date: '2026-03-06 04:07'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:polish'
dependencies: []
priority: low
ordinal: 80000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Provide users with an in-app reference for all keyboard shortcuts, accessible at any time by pressing ?.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Pressing ? in any mode renders a full-screen lipgloss overlay listing all keyboard shortcuts from contracts/keyboard-shortcuts.md
- [ ] #2 The help overlay lists shortcuts organized by mode: Browse, Search, Add, Confirm-Delete
- [ ] #3 Any keypress closes the help overlay and returns to the previous mode
<!-- AC:END -->
