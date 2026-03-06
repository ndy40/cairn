---
id: TASK-15
title: Implement fetch-and-save flow when user confirms add modal
status: Done
assignee: []
created_date: '2026-03-06 04:05'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us1'
dependencies: []
priority: high
ordinal: 55000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
When the user presses Enter in the add modal, run the page fetch and store insert as an async tea.Cmd so the TUI remains responsive. Handle all outcomes: success, duplicate, and fetch failure.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 On Enter in AddModel, a tea.Cmd is dispatched that calls fetcher.Fetch(url) followed by store.Insert(...)
- [ ] #2 On successful save, the bookmark list is reloaded and the app transitions to StateBrowse
- [ ] #3 On duplicate URL error, the modal stays open and shows 'Already bookmarked' in the modal status
- [ ] #4 On fetch failure, the bookmark is still saved with the URL hostname as title; modal shows 'Saved (title unavailable)' then transitions to StateBrowse
- [ ] #5 On invalid URL (no http/https scheme), the modal shows 'Invalid URL' and does not attempt to save
<!-- AC:END -->
