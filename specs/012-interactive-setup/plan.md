# Implementation Plan: Interactive Setup Configuration Prompts

**Branch**: `012-interactive-setup` | **Date**: 2026-03-26 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/012-interactive-setup/spec.md`

## Summary

Add interactive prompts to `cairn sync setup` so users who have not set `CAIRN_DROPBOX_APP_KEY` are guided to enter their Dropbox App Key at the terminal rather than receiving a hard error. Optionally prompt for a custom database path. Persist supplied values to `cairn.json` using the existing `config.Manager` before proceeding to OAuth. No new dependencies, no schema changes, no new packages.

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**: stdlib only (`bufio`, `fmt`, `os`, `strings`) + existing `internal/config`, `internal/sync`
**Storage**: `cairn.json` (existing, via `config.Manager` / viper) — no SQLite changes
**Testing**: `go test ./...`
**Target Platform**: Linux, macOS, Windows (existing cross-platform support)
**Project Type**: CLI tool
**Performance Goals**: N/A — setup is a one-time interactive flow
**Constraints**: No new external dependencies; pure Go; single binary
**Scale/Scope**: Two new functions in `cmd/cairn/main.go`; no new files

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate | Status | Notes |
|------|--------|-------|
| No CGO | PASS | stdlib + existing pure-Go deps only |
| Single binary | PASS | No new runtime or external process |
| Task management | PASS | Tasks created via Backlog CLI after `/speckit.tasks` |
| Backward-compatible migrations | PASS | No schema changes |

All gates pass. No violations to justify.

## Project Structure

### Documentation (this feature)

```text
specs/012-interactive-setup/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   └── cli-contract.md  # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit.tasks — not yet created)
```

### Source Code (repository root)

```text
cmd/cairn/
└── main.go              # Modify: add promptForSetupConfig(); update runSyncSetup() signature

internal/config/
└── config.go            # No change (existing Manager.Set, WriteConfig, DefaultConfigPath reused)

internal/sync/
└── (all files)          # No change
```

**Structure Decision**: Single-project layout. All prompt logic lives in `cmd/cairn/main.go` as a private helper. No new packages — three new functions (prompt helper, stdin reader, path confirmation printer) in the existing CLI entry point.

## Phase 0: Research

See [research.md](./research.md) for full findings. Summary:

- **Input method**: `bufio.NewReader(os.Stdin).ReadString('\n')` — handles paths with spaces, arbitrary key characters.
- **Config write**: Reuse `config.Manager.Set()` + `WriteConfig()`. Viper serialises all loaded settings, so existing keys are preserved automatically.
- **Prompt conditions**: App Key prompt when `cfgManager.Get().DropboxAppKey == ""`; DB path prompt when `CAIRN_DB_PATH` unset AND resolved path equals `config.DefaultDBPath()`.
- **No new dependencies**: stdlib only.

## Phase 1: Design

### data-model.md

See [data-model.md](./data-model.md). No new entities. Existing `AppConfig` fields (`DropboxAppKey`, `DBPath`) and `cairn.json` keys (`dropbox_app_key`, `db_path`) are the only data involved.

### contracts/cli-contract.md

See [contracts/cli-contract.md](./contracts/cli-contract.md). Documents exact prompt text, error messages, and conditions under which each prompt is shown or skipped.

### quickstart.md

See [quickstart.md](./quickstart.md). Lists affected files, pseudocode for `promptForSetupConfig`, and manual test steps.

## Implementation Phases

### Phase A — Core Prompt Helper (P1)

1. In `cmd/cairn/main.go`, add `promptForSetupConfig(cfgManager *config.Manager)`:
   - If `cfgManager.Get().DropboxAppKey == ""`: loop reading stdin until non-empty trimmed input; call `cfgManager.Set("dropbox_app_key", key)`.
   - If `os.Getenv("CAIRN_DB_PATH") == ""` and `cfgManager.Get().DBPath == config.DefaultDBPath()`: prompt once for db path; if non-empty, call `cfgManager.Set("db_path", path)`.
   - If any Set was called: call `cfgManager.WriteConfig()`; print "Config written to `<path>`".
2. Update `runSyncSetup` to accept `cfgManager *config.Manager` and call `promptForSetupConfig(cfgManager)` before the existing `appKey == ""` guard.
3. Pass `cfgManager` from `main()` when calling `runSyncSetup`.

### Phase B — Unit Tests (P1)

Add tests in `cmd/cairn/` (or a new `internal/setup/` if extracted):
- Test prompt is skipped when App Key already set.
- Test prompt is skipped when `CAIRN_DB_PATH` is set.
- Test empty input re-prompts (simulate via `strings.NewReader`).
- Test valid input writes correct keys to config.

### Phase C — Integration Test (P2)

Manual or scripted: run `cairn sync setup` with no env vars, exercise prompts, verify `cairn.json` content.
