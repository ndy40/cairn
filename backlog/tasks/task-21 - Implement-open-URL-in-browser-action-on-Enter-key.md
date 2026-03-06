---
id: TASK-21
title: Implement open-URL-in-browser action on Enter key
status: Done
assignee: []
created_date: '2026-03-06 04:05'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us2'
dependencies: []
priority: medium
ordinal: 61000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Allow users to open a selected bookmark in their default browser by pressing Enter. Uses the OS-appropriate open command.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Pressing Enter on a selected bookmark runs xdg-open (Linux), open (macOS), or start (Windows) with the bookmark URL
- [ ] #2 The command runs as an async tea.Cmd so the TUI remains responsive
- [ ] #3 If the open command fails, an error message is shown in the footer
<!-- AC:END -->
