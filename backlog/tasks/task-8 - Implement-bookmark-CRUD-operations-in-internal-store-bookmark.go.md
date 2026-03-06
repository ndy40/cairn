---
id: TASK-8
title: Implement bookmark CRUD operations in internal/store/bookmark.go
status: Done
assignee: []
created_date: '2026-03-06 04:04'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:foundation'
dependencies: []
priority: high
ordinal: 85000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Provide all bookmark data operations needed by the rest of the application: insert, list, delete, duplicate check, and get-by-id.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Insert(url, domain, title, description, createdAt) inserts a bookmark and returns a specific error type when url already exists (duplicate)
- [ ] #2 List() returns all bookmarks ordered by created_at DESC
- [ ] #3 DeleteByID(id) removes a bookmark by its integer ID
- [ ] #4 ExistsByURL(url) returns true if a bookmark with that URL is already stored
- [ ] #5 GetByID(id) returns the bookmark or an error if not found
<!-- AC:END -->
