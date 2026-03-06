---
id: TASK-96
title: >-
  T002 [005] Bootstrap vicinae-extension/ package with package.json and
  tsconfig.json
status: Done
assignee: []
created_date: '2026-03-06 16:02'
updated_date: '2026-03-06 17:10'
labels:
  - feature-005
  - setup
dependencies: []
documentation:
  - specs/005-vicinae-extension/contracts/extension-commands.md
  - specs/005-vicinae-extension/plan.md
priority: high
ordinal: 96000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create `vicinae-extension/` at the repository root containing:\n\n**package.json** — must declare:\n- `name: "bm-bookmarks"`, `title: "Bookmark Manager"`\n- Three commands: `search-bookmarks → src/search.tsx`, `list-bookmarks → src/list.tsx`, `add-bookmark → src/add.tsx`\n- Dependency: `@vicinae/api`\n- devDependency: `typescript`\n- Scripts: `build: vici build`, `dev: vici develop`\n\n**tsconfig.json** — ESNext target, JSX support, strict mode\n\n**src/** — empty directory\n\nPer `specs/005-vicinae-extension/contracts/extension-commands.md`
<!-- SECTION:DESCRIPTION:END -->
