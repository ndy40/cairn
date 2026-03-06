---
id: TASK-41
title: Implement bm version subcommand
status: Done
assignee: []
created_date: '2026-03-06 04:07'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:polish'
dependencies: []
priority: low
ordinal: 81000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Provide a version command so users can check which version of the binary they are running.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 bm version prints the version string and exits with code 0
- [ ] #2 Version is embedded at build time via go build -ldflags='-X main.version=0.1.0'
- [ ] #3 When built without the ldflags, version prints 'dev'
<!-- AC:END -->
