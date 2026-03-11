---
id: TASK-109
title: 'T004: Update bookmark store for uuid and updated_at'
status: Done
assignee:
  - '@claude'
created_date: '2026-03-11 07:40'
updated_date: '2026-03-11 11:31'
labels:
  - foundational
  - phase-2
dependencies:
  - TASK-108
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Modify internal/store/bookmark.go: update Insert() to generate UUID v4 and set updated_at to current time on every new bookmark. Update UpdateTags() to set updated_at to current time. Update DeleteByID() to accept and return the bookmark UUID before deletion (needed for tombstone recording). Update scanBookmark() and scanBookmarks() to read the new uuid and updated_at columns. Add uuid and UpdatedAt fields to the Bookmark struct.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Bookmark struct has UUID string and UpdatedAt time.Time fields
- [x] #2 Insert() generates UUID v4 and sets updated_at
- [x] #3 UpdateTags() updates updated_at timestamp
- [x] #4 DeleteByID() returns bookmark UUID before deletion
- [x] #5 scanBookmark/scanBookmarks read uuid and updated_at columns
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Done inline with T003
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Updated Bookmark struct with UUID and UpdatedAt fields. Insert generates UUIDv4 and sets updated_at. UpdateTags updates updated_at. DeleteByID now returns UUID for tombstone recording. scanRow reads new columns. Updated callers in app.go and main.go.
<!-- SECTION:FINAL_SUMMARY:END -->
