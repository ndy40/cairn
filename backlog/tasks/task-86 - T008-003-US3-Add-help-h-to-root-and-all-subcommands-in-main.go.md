---
id: TASK-86
title: 'T008 [003] [US3] Add --help/-h to root and all subcommands in main.go'
status: Done
assignee: []
created_date: '2026-03-06 10:17'
updated_date: '2026-03-06 17:10'
labels:
  - feature-003
  - US3
  - cli
dependencies:
  - TASK-84
documentation:
  - specs/003-edit-bookmark-help/contracts/cli-interface.md
  - specs/003-edit-bookmark-help/research.md
priority: medium
ordinal: 15000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update `cmd/bm/main.go` to support --help/-h flags exiting with code 0:\n\n1. Root command: check `os.Args` for `-h`/`--help` before `flag.Parse()`, call `printHelp()`, exit 0\n2. For each subcommand FlagSet (add, list, search, delete, version, help):\n   - Set `fs.Usage` to print per-subcommand help text (per contracts/cli-interface.md) and exit 0\n   - Add `fs.Bool(\"help\", false, \"show help\")` flag to each FlagSet\n   - After `fs.Parse()`, check the help flag and exit 0 if set\n\nAll help responses must exit with code 0 (not the default code 2 from the flag package).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 bm --help prints usage and exits 0
- [ ] #2 bm -h prints usage and exits 0
- [ ] #3 bm add --help prints add usage and exits 0
- [ ] #4 bm list --help prints list usage and exits 0
- [ ] #5 bm search --help prints search usage and exits 0
- [ ] #6 bm delete --help prints delete usage and exits 0
- [ ] #7 bm version --help prints version usage and exits 0
- [ ] #8 bm help --help prints help usage and exits 0
<!-- AC:END -->
