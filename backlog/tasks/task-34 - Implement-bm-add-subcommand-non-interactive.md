---
id: TASK-34
title: Implement bm add subcommand (non-interactive)
status: Done
assignee: []
created_date: '2026-03-06 04:06'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:polish'
dependencies: []
priority: low
ordinal: 74000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Allow users to save a bookmark from the command line without launching the TUI, useful for scripting and quick saves.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 bm add <url> fetches the page, saves the bookmark, and prints 'Saved: "<title>" (<domain>)' on success
- [ ] #2 bm add exits with code 1 and prints 'Already bookmarked' if the URL already exists
- [ ] #3 bm add saves with fallback title and exits with code 2 if the page fetch fails
- [ ] #4 bm add exits with code 3 and an error message for unexpected errors
<!-- AC:END -->
