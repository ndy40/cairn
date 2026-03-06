---
id: TASK-92
title: 'T004 [004] [US2] Remove last-visited display from BookmarkItem.Description()'
status: Done
assignee: []
created_date: '2026-03-06 15:32'
updated_date: '2026-03-06 17:10'
labels:
  - feature-004
  - US2
  - tui
dependencies: []
documentation:
  - specs/004-bookmark-expiry/contracts/browse-row.md
priority: high
ordinal: 5000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
In `internal/model/browse.go`, delete the last-visited conditional block from `BookmarkItem.Description()`:\n\nRemove:\n```go\nif i.b.LastVisitedAt != nil {\n    desc += \" · Last: \" + i.b.LastVisitedAt.Format(\"2006-01-02\")\n} else {\n    desc += \" · Never visited\"\n}\n```\n\nThe resulting description line format becomes:\n`domain · YYYY-MM-DD [· #tag1 #tag2]`\n\nPer contract: `specs/004-bookmark-expiry/contracts/browse-row.md`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 No 'Last:' text visible on any bookmark row
- [ ] #2 No 'Never visited' text visible on any bookmark row
- [ ] #3 Domain and creation date still shown correctly
- [ ] #4 Tags still shown when present
<!-- AC:END -->
