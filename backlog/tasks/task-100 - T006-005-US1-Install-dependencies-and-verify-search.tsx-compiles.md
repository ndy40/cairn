---
id: TASK-100
title: 'T006 [005] [US1] Install dependencies and verify search.tsx compiles'
status: Done
assignee: []
created_date: '2026-03-06 16:02'
updated_date: '2026-03-06 17:10'
labels:
  - feature-005
  - US1
  - extension
dependencies:
  - TASK-99
documentation:
  - specs/005-vicinae-extension/tasks.md
priority: medium
ordinal: 98000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
In `vicinae-extension/`:\n1. Run `npm install` to install `@vicinae/api` and TypeScript\n2. Confirm `src/search.tsx` has zero TypeScript compilation errors (run `npx tsc --noEmit`)\n3. Optionally run `vici develop` to confirm the search command loads in Vicinae without runtime errors
<!-- SECTION:DESCRIPTION:END -->
