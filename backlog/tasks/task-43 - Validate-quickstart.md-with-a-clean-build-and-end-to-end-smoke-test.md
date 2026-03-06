---
id: TASK-43
title: Validate quickstart.md with a clean build and end-to-end smoke test
status: Done
assignee: []
created_date: '2026-03-06 04:07'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:polish'
dependencies: []
priority: low
ordinal: 52000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Verify that the quickstart.md guide is accurate and the application works end-to-end before the feature is considered complete.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 CGO_ENABLED=0 go build -ldflags='-s -w' -o bm ./cmd/bm succeeds and produces a binary
- [ ] #2 bm launches and displays the bookmark list (or empty state) without errors
- [ ] #3 bm add https://example.com saves a bookmark non-interactively
- [ ] #4 bm list shows the saved bookmark
- [ ] #5 bm search 'example' returns the bookmark
- [ ] #6 quickstart.md is updated if any step in the guide is inaccurate or missing
<!-- AC:END -->
