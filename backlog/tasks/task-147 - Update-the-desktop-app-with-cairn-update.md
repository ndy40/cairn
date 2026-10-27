---
id: TASK-147
title: Update the desktop app with cairn update
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T13:30:00Z'
updated_date: '2026-10-27T13:35:00Z'
labels:
  - cli
  - desktop
dependencies: []
ordinal: 108000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
cairn update replaces the CLI and Vicinae extension but not cairn-desktop, so desktop installs drift behind the CLI they share a database with. Releases already ship matching desktop archives.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 cairn update updates an installed desktop app to the latest release after updating the CLI
- [x] #2 cairn update --desktop updates only the desktop app; --check reports desktop version status without changes
- [x] #3 Downloaded desktop archive is verified against checksums.txt; a mismatch leaves the installed app untouched
- [x] #4 A running desktop app is restarted on the new version
- [x] #5 install.sh records the installed desktop version so cairn update can compare it
- [x] #6 CLI help and docs describe desktop updates
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. internal/updater/desktop.go: DetectDesktop, CheckDesktopVersion (marker file desktop-version.txt in the cairn data dir), UpdateDesktop (download + required checksum, Linux binary/icon swap via same-dir rename, macOS .app swap via ditto, restart if running)
2. cmd/cairn/run_update.go: default update runs CLI then desktop (if installed); --desktop for desktop only
3. install.sh writes the desktop version marker
4. Tests for update/mismatch/version; help text, README and docs
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- Implemented updater/desktop.go, --desktop flag, installer version marker, docs
- Smoke test against live v0.7.2 found that name-based pgrep/pkill restarted an unrelated cairn-desktop install; now matches processes by executable path (/proc/<pid>/exe on Linux, ps comm on macOS)
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
`cairn update` now updates the desktop app alongside the CLI so both stay on the same release.

Changes:
- internal/updater/desktop.go: DetectDesktop (binary beside cairn or on PATH on Linux; ~/Applications/cairn-desktop.app on macOS), CheckDesktopVersion (reads desktop-version.txt in the cairn data dir; missing = unknown = update), UpdateDesktop (required checksum, atomic same-dir binary swap + icon on Linux, ditto-unpacked .app swap with rollback on macOS, restart only if this install was running).
- cmd/cairn: `cairn update` runs CLI then desktop (if installed); `--desktop` updates only the desktop app; `--check` reports both. Exit codes 3/4 for checksum/permission as for the CLI.
- install.sh writes desktop-version.txt on install.
- Help text, README, cli-reference and desktop-app docs updated.

Tests:
- go test ./... ; golangci-lint run; cross-built for windows/darwin
- New updater tests: replace + restart, not-running, checksum mismatch, missing checksums, version check, detection, process matching
- Manual: sandboxed update against the live v0.7.2 release (CLI + desktop), real running app left untouched

Not tested: macOS path end-to-end (compiled only).
<!-- SECTION:FINAL_SUMMARY:END -->
