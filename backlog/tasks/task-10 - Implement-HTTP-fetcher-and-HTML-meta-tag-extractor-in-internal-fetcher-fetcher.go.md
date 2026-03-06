---
id: TASK-10
title: >-
  Implement HTTP fetcher and HTML meta tag extractor in
  internal/fetcher/fetcher.go
status: Done
assignee: []
created_date: '2026-03-06 04:05'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:foundation'
dependencies: []
priority: high
ordinal: 87000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Fetch a web page and extract its title and description from HTML meta tags. This runs when a bookmark is saved to automatically populate the bookmark's display name and description.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Fetch(url) performs HTTP GET with 8-second timeout and BookmarkManager/1.0 User-Agent
- [ ] #2 Response body is limited to 512 KB via io.LimitReader before parsing
- [ ] #3 Response body charset is detected and converted to UTF-8 using golang.org/x/net/html/charset before parsing
- [ ] #4 Title is extracted in priority order: og:title property > <title> tag > URL hostname as fallback
- [ ] #5 Description is extracted in priority order: og:description property > meta[name=description] content > empty string
- [ ] #6 Fetch returns (title, description, nil) on success and (hostname, "", err) on any fetch/parse error
<!-- AC:END -->
