---
id: TASK-9
title: Implement clipboard wrapper in internal/clipboard/clipboard.go
status: Done
assignee: []
created_date: '2026-03-06 04:04'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:foundation'
dependencies: []
priority: high
ordinal: 86000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Provide a simple interface for reading the system clipboard. This is used by the Ctrl+P shortcut to paste URLs.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Read() returns the clipboard text string and nil error when clipboard contains text
- [ ] #2 Read() returns a descriptive error when clipboard is empty
- [ ] #3 Read() returns a descriptive error when clipboard access fails (e.g., xclip not installed on Linux)
<!-- AC:END -->
