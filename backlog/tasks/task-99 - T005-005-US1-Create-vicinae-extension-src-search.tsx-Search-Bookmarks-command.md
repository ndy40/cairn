---
id: TASK-99
title: >-
  T005 [005] [US1] Create vicinae-extension/src/search.tsx Search Bookmarks
  command
status: Done
assignee: []
created_date: '2026-03-06 16:02'
updated_date: '2026-03-06 17:10'
labels:
  - feature-005
  - US1
  - extension
dependencies:
  - TASK-98
documentation:
  - specs/005-vicinae-extension/contracts/extension-commands.md
priority: high
ordinal: 103000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create `vicinae-extension/src/search.tsx` implementing the Search Bookmarks command:\n\n- Use `List` component from `@vicinae/api`\n- On mount check `bmAvailable()`; if false, show error panel with install hint\n- When search query is empty: call `bmList()` to show all bookmarks\n- When query changes: call `bmSearch(query)` and update list\n- Each result is a `List.Item` with:\n  - Title: `bookmark.Title || bookmark.URL`\n  - Subtitle: `bookmark.Domain`\n  - Accessories: formatted CreatedAt (YYYY-MM-DD), tags as `#tag1 #tag2`, pin badge if `IsPermanent`\n- Primary action (Enter): open `bookmark.URL` in default browser\n- Secondary action: copy URL to clipboard\n- Empty state when no results: show \"No bookmarks found\"\n\nPer `specs/005-vicinae-extension/contracts/extension-commands.md`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Typing a query filters results in real time
- [ ] #2 Empty query shows full bookmark list
- [ ] #3 Enter on a result opens URL in browser
- [ ] #4 No results shows 'No bookmarks found' empty state
- [ ] #5 Missing bm CLI shows error message
<!-- AC:END -->
