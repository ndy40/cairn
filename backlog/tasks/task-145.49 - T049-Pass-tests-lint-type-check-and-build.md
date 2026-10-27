---
id: TASK-145.49
title: T049 Pass tests, lint, type-check and build
status: Done
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:34:00Z'
labels:
  - desktop
  - polish
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 157000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T049 in specs/018-native-desktop-window/tasks.md.

Run `go test ./...`, the desktop lint job, `npx tsc --noEmit` in `desktop/frontend` and `wails3 build` (or `just test-desktop` / `just build-desktop`); fix failures
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 go test, the desktop lint job, tsc --noEmit and the desktop build all pass
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
On 018-native-desktop-window after all merges: go test ./..., go vet ./..., golangci-lint run ./... (0 issues after 16b9547), npx tsc --noEmit and go build all pass in desktop/. wails3 build passed in each feature worktree.
<!-- SECTION:FINAL_SUMMARY:END -->
