---
id: TASK-97
title: 'T003 [005] Add --tags flag to bm add subcommand in main.go'
status: Done
assignee: []
created_date: '2026-03-06 16:02'
updated_date: '2026-03-06 17:10'
labels:
  - feature-005
  - foundational
  - cli
dependencies: []
documentation:
  - specs/005-vicinae-extension/contracts/cli-interface.md
  - specs/005-vicinae-extension/data-model.md
priority: high
ordinal: 97000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update `cmd/bm/main.go` to add an optional `--tags` flag to the `bm add` subcommand:\n\n1. Change the `"add"` case to use a FlagSet instead of direct positional args:\n   - `fs := flag.NewFlagSet("add", flag.ContinueOnError)`\n   - Add `fs.Usage` to call `printAddHelp()` and exit 0\n   - `tagsFlag := fs.String("tags", "", "comma-separated tags")`\n   - Parse remaining args after detecting help flags\n   - First non-flag arg is the URL\n2. Pass `store.NormaliseTagsFromString(*tagsFlag)` to `s.Insert()` instead of `nil`\n3. Update `printAddHelp()` to document `--tags` flag\n\nPer `specs/005-vicinae-extension/contracts/cli-interface.md`
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 bm add https://example.com --tags "go, tools" saves bookmark with tags
- [ ] #2 bm add https://example.com (no --tags) behaves identically to before
- [ ] #3 bm add --help shows --tags in usage text
<!-- AC:END -->
