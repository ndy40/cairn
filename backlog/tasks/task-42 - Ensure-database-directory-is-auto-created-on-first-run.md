---
id: TASK-42
title: Ensure database directory is auto-created on first run
status: Done
assignee: []
created_date: '2026-03-06 04:07'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:polish'
dependencies: []
priority: low
ordinal: 82000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Guarantee the app works on a fresh install without requiring users to manually create the data directory.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 On store open, parent directories for the database path are created with os.MkdirAll if they do not exist
- [ ] #2 App starts successfully on a system where the data directory has never existed
- [ ] #3 Appropriate permissions are set on the created directory (0700)
<!-- AC:END -->
