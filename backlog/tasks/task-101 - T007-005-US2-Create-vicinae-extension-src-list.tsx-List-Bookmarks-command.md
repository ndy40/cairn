---
id: TASK-101
title: 'T007 [005] [US2] Create vicinae-extension/src/list.tsx List Bookmarks command'
status: Done
assignee: []
created_date: '2026-03-06 16:02'
updated_date: '2026-03-06 17:10'
labels:
  - feature-005
  - US2
  - extension
dependencies:
  - TASK-100
documentation:
  - specs/005-vicinae-extension/contracts/extension-commands.md
priority: medium
ordinal: 102000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create `vicinae-extension/src/list.tsx` implementing the List Bookmarks command:\n\n- Use `List` component from `@vicinae/api`\n- On mount check `bmAvailable()`; if false, show error panel\n- Call `bmList()` once on mount to load all active bookmarks\n- Render using same `List.Item` structure as search.tsx (title, domain, date, tags, pin badge)\n- Primary action (Enter): open URL in browser; secondary: copy URL\n- Empty state: "No bookmarks saved yet"\n- Vicinae's built-in search box handles client-side filtering — no additional CLI calls on query change\n\nPer `specs/005-vicinae-extension/contracts/extension-commands.md`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Full bookmark list loads on command open
- [ ] #2 Typing in launcher filters list without CLI call
- [ ] #3 Enter opens URL in browser
- [ ] #4 Empty list shows 'No bookmarks saved yet'
<!-- AC:END -->
