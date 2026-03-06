---
id: TASK-102
title: 'T008 [005] [US3] Create vicinae-extension/src/add.tsx Add Bookmark command'
status: Done
assignee: []
created_date: '2026-03-06 16:03'
updated_date: '2026-03-06 17:10'
labels:
  - feature-005
  - US3
  - extension
dependencies:
  - TASK-101
documentation:
  - specs/005-vicinae-extension/contracts/extension-commands.md
priority: medium
ordinal: 101000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create `vicinae-extension/src/add.tsx` implementing the Add Bookmark command:\n\n- Use `Form` component from `@vicinae/api`\n- On mount: read clipboard; if content starts with `http://` or `https://`, pre-fill URL field\n- URL field: required, placeholder `https://example.com`\n- Tags field: optional, placeholder `work, go, tools  (comma-separated, max 3)`\n- On submit:\n  1. Validate URL is non-empty and starts with `http://` or `https://`; show validation error if not\n  2. Call `bmAdd(url, tags)`\n  3. Exit 0 or 2: show success toast \"Saved\", close form\n  4. Exit 1: show inline error \"Already bookmarked\"\n  5. Exit 3: show inline error with stderr content\n\nPer `specs/005-vicinae-extension/contracts/extension-commands.md`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 URL field pre-filled from clipboard if valid URL
- [ ] #2 Empty URL blocked with validation error
- [ ] #3 Successful save shows toast and closes form
- [ ] #4 Duplicate URL shows 'Already bookmarked'
- [ ] #5 Tags saved when provided
<!-- AC:END -->
