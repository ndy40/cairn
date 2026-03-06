---
id: TASK-20
title: Add full keyboard navigation to browse list
status: Done
assignee: []
created_date: '2026-03-06 04:05'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us2'
dependencies: []
priority: medium
ordinal: 60000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Extend BrowseModel with vim-style navigation keys and jump-to-top/bottom shortcuts so users can navigate large bookmark lists efficiently.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Arrow keys ↑/↓ move selection up and down in the list
- [ ] #2 j and k keys also move selection (vim-style)
- [ ] #3 g jumps to the first item in the list
- [ ] #4 G jumps to the last item in the list
- [ ] #5 All navigation keys are shown in the list's help text
<!-- AC:END -->
