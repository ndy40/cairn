# Implementation Plan: Configuration File Support

**Branch**: `010-config-file` | **Date**: 2026-03-12 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/010-config-file/spec.md`

## Summary

Add an optional `cairn.json` configuration file that provides a persistent way to configure cairn without environment variables. The config is loaded from the OS-appropriate config directory and merged with other sources using the precedence: env vars > CLI flags > config file > defaults. A new `internal/config` package handles loading, parsing, and merging. The `cairn config` command is updated to reflect the resolved values. README is updated with configuration documentation.

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**: `encoding/json` (stdlib) — no new external dependencies
**Storage**: JSON file (`cairn.json`) in OS config directory; no schema changes
**Testing**: `go test`
**Target Platform**: Linux, macOS, Windows (cross-platform)
**Project Type**: CLI tool
**Performance Goals**: N/A (config loaded once at startup)
**Constraints**: Zero CGO, single binary, no external runtime
**Scale/Scope**: 2 config keys (`db_path`, `dropbox_app_key`), 1 new internal package

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate | Requirement | Status |
|------|-------------|--------|
| No CGO | All dependencies pure Go | PASS — uses only stdlib `encoding/json` |
| Single binary | No external runtime | PASS — reads a local JSON file |
| Task management | Tasks must go through Backlog CLI after `/speckit.tasks` | PASS — will use Backlog CLI |
| Backward-compatible migrations | ALTER TABLE with DEFAULT values only | PASS — no schema changes |

## Project Structure

### Documentation (this feature)

```text
specs/010-config-file/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── cli-interface.md
└── tasks.md
```

### Source Code (repository root)

```text
internal/
├── config/
│   └── config.go          # NEW: AppConfig struct, Load(), DefaultConfigPath(), Resolve()
├── store/
│   └── store.go           # MODIFIED: resolveDBPath uses config.Resolve()
├── sync/
│   └── config.go          # UNCHANGED: sync-specific config stays separate
cmd/
└── cairn/
    └── main.go            # MODIFIED: load config file, update resolveDBPath, update config command
README.md                  # MODIFIED: add Configuration section
```

**Structure Decision**: New `internal/config` package with a single file. This package has a single clear responsibility: load and merge configuration from multiple sources. It does not duplicate sync config — that remains in `internal/sync/config.go`.

## Complexity Tracking

No constitution violations. No complexity tracking needed.
