---
id: TASK-145.51
title: Show every bookmark and the date added in the desktop list
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T09:17:00Z'
labels:
  - desktop
  - us6
dependencies: []
parent_task_id: TASK-145
ordinal: 159000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Found while testing locally: the window listed at most 50 bookmarks (maxResults cap), so with 101 bookmarks half were unreachable. Users also want to see when each bookmark was added.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 With an empty query the list includes every bookmark, and it scrolls through all of them
- [ ] #2 Each row shows the date the bookmark was added on the right, with the exact time in a tooltip
- [ ] #3 PageUp and PageDown move the selection a screen at a time
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implemented in 408d29d: removed maxResults in desktop/service.go, added AddedAt (RFC 3339) to the Bookmark DTO, addedDate() in render.ts, PageUp/PageDown in main.ts, fixed reveal() scrolling to the wrong row when section headings are present. Tests: TestSearchListsEveryBookmark, TestSearchIncludesAddedDate. Awaiting a look in the running app.
<!-- SECTION:NOTES:END -->
