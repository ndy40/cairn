---
id: TASK-98
title: 'T004 [005] Create vicinae-extension/src/bm.ts shared CLI wrapper'
status: Done
assignee: []
created_date: '2026-03-06 16:02'
updated_date: '2026-03-06 17:10'
labels:
  - feature-005
  - foundational
  - extension
dependencies: []
documentation:
  - specs/005-vicinae-extension/data-model.md
  - specs/005-vicinae-extension/research.md
priority: high
ordinal: 104000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create `vicinae-extension/src/bm.ts` — shared utility module used by all three commands:\n\n1. Export `Bookmark` TypeScript interface matching JSON schema:\n   - `ID`, `URL`, `Domain`, `Title`, `Description`, `CreatedAt`, `Tags: string[]`, `IsPermanent`, `IsArchived`\n\n2. Export `bmAvailable(): boolean` — runs `which bm`, returns true if exit 0\n\n3. Export `bmList(): Bookmark[]` — runs `bm list --json`, parses JSON array, returns results (empty array on error)\n\n4. Export `bmSearch(query: string): Bookmark[]` — runs `bm search <query> --json --limit 20`, parses JSON array\n\n5. Export `bmAdd(url: string, tags?: string): { exitCode: number; stderr: string }` — runs `bm add <url> [--tags <tags>]`, returns exit code and stderr\n\nAll CLI calls use `child_process.spawnSync` with args as array (no shell interpolation).\n\nPer `specs/005-vicinae-extension/data-model.md`
<!-- SECTION:DESCRIPTION:END -->
