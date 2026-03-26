# Tasks: Interactive Setup Configuration Prompts

**Input**: Design documents from `/specs/012-interactive-setup/`
**Prerequisites**: plan.md ✓, spec.md ✓, research.md ✓, data-model.md ✓, contracts/cli-contract.md ✓, quickstart.md ✓

**Tests**: Not explicitly requested — no test tasks generated.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Verify baseline and identify integration points before making changes. No new files needed.

- [x] T001 Confirm `go build ./cmd/cairn && go vet ./cmd/cairn` passes on the current branch (baseline check)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Propagate `cfgManager` into the sync setup call chain so the prompt helper can read and write config. This must complete before any user story can be implemented.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T002 Update `runSync` signature in `cmd/cairn/main.go` to add `cfgManager *config.Manager` parameter (pass through from caller; keep all other params unchanged)
- [x] T003 Update `runSyncSetup` signature in `cmd/cairn/main.go` to replace `appCfg *config.AppConfig` with `cfgManager *config.Manager` (derive `appKey` from `cfgManager.Get().DropboxAppKey`)
- [x] T004 Update the call in `main()` at `cmd/cairn/main.go` from `runSync(resolvedDB, appCfg, args[1:])` to `runSync(resolvedDB, cfgManager, args[1:])`

**Checkpoint**: `go build ./cmd/cairn && go vet ./cmd/cairn` must pass before user story work begins

---

## Phase 3: User Story 1 — Guided Dropbox App Key Prompt (Priority: P1) 🎯 MVP

**Goal**: When `CAIRN_DROPBOX_APP_KEY` is unset and `cairn.json` has no `dropbox_app_key`, `cairn sync setup` prompts the user interactively instead of printing a fatal error.

**Independent Test**: `unset CAIRN_DROPBOX_APP_KEY && rm -f ~/.config/cairn/cairn.json && cairn sync setup` — CLI must prompt "Enter your Dropbox App Key:", accept a non-empty value, and proceed to OAuth. Pressing Enter on an empty input must re-prompt with "Error: App Key cannot be empty. Please try again."

### Implementation for User Story 1

- [x] T005 [US1] Add `promptForSetupConfig(cfgManager *config.Manager)` function stub (empty body) in `cmd/cairn/main.go` — function must compile and be callable
- [x] T006 [US1] Implement App Key prompt loop inside `promptForSetupConfig` in `cmd/cairn/main.go`: print prompt, read stdin via `bufio.NewReader(os.Stdin).ReadString('\n')`, trim whitespace, re-prompt with error on empty input, call `cfgManager.Set("dropbox_app_key", key)` on valid input
- [x] T007 [US1] Add early-return guard in `promptForSetupConfig` in `cmd/cairn/main.go`: skip App Key prompt block when `cfgManager.Get().DropboxAppKey != ""`
- [x] T008 [US1] Call `promptForSetupConfig(cfgManager)` at the start of `runSyncSetup` in `cmd/cairn/main.go`, before the existing `appKey == ""` guard; remove the `fatalf` error for missing App Key (prompt now handles it)

**Checkpoint**: User Story 1 is independently testable — run setup with no env var and verify interactive prompt works

---

## Phase 4: User Story 2 — Optional Database Path Prompt (Priority: P2)

**Goal**: After the App Key is resolved, `cairn sync setup` optionally prompts for a custom database path. Pressing Enter keeps the default silently.

**Independent Test**: With `CAIRN_DROPBOX_APP_KEY` set in env (so App Key prompt is skipped) and no `CAIRN_DB_PATH` set and no `db_path` in `cairn.json`: running `cairn sync setup` must show the DB path prompt with the default path visible. Pressing Enter must proceed without writing `db_path` to `cairn.json`. Entering a custom path must write `db_path` to `cairn.json`.

### Implementation for User Story 2

- [x] T009 [US2] Add DB path prompt block to `promptForSetupConfig` in `cmd/cairn/main.go`: print prompt showing default path (`config.DefaultDBPath()`), read one line from stdin; if non-empty, call `cfgManager.Set("db_path", path)`
- [x] T010 [US2] Add skip condition to the DB path prompt block in `cmd/cairn/main.go`: skip when `os.Getenv("CAIRN_DB_PATH") != ""` OR when `cfgManager.Get().DBPath != config.DefaultDBPath()`

**Checkpoint**: User Stories 1 AND 2 are independently testable

---

## Phase 5: User Story 3 — Config File Written with Confirmation (Priority: P3)

**Goal**: After prompts are resolved, `cairn.json` is written to the OS config directory and the CLI prints the full path. Existing keys are preserved.

**Independent Test**: After completing the setup prompt flow (entering an App Key), verify that `cairn.json` exists at `~/.config/cairn/cairn.json` (Linux) and contains `dropbox_app_key`. Verify the CLI printed `Config written to <path>`. Verify a pre-existing `db_path` in the file was not removed.

### Implementation for User Story 3

- [x] T011 [US3] Add config write block at the end of `promptForSetupConfig` in `cmd/cairn/main.go`: call `cfgManager.WriteConfig()` and `fatalf` on error; print `"Config written to %s\n"` using `config.DefaultConfigPath()` as the path

**Checkpoint**: All three user stories are independently functional. `go build ./cmd/cairn && go vet ./cmd/cairn` must pass.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and cleanup

- [x] T012 Run `go build ./... && go vet ./...` from repo root and fix any compilation or vet errors in `cmd/cairn/main.go`
- [x] T013 [P] Perform manual validation per `specs/012-interactive-setup/quickstart.md`: run all three manual test cases (no env var, env var set, empty input) and confirm expected output
- [x] T014 [P] Verify Ctrl-D/EOF at the App Key prompt exits cleanly without writing a partial `cairn.json` (check `bufio.Reader.ReadString` EOF behaviour and add guard if needed in `cmd/cairn/main.go`)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 — **BLOCKS all user stories**
- **User Story phases (3, 4, 5)**: All depend on Phase 2 completion; US2 and US3 extend the same `promptForSetupConfig` function — implement sequentially
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1 (P1)**: Depends on Phase 2 (foundational) — no story dependencies
- **US2 (P2)**: Depends on US1 (extends `promptForSetupConfig`)
- **US3 (P3)**: Depends on US1 (extends `promptForSetupConfig`); logically after US2

### Within Each User Story

- T005 must precede T006, T007, T008 (stub before implementation)
- T006 must precede T007 (implement then guard)
- T007 must precede T008 (guard before wiring into runSyncSetup)

### Parallel Opportunities

- T009 and T010 within US2 can be written together (same function, adjacent lines)
- T013 and T014 in Polish can run in parallel (different test scenarios)

---

## Parallel Example: Polish Phase

```bash
# Run simultaneously:
Task T013: Manual test — no env var scenario
Task T014: Manual test — Ctrl-D/EOF scenario
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (baseline verification)
2. Complete Phase 2: Foundational (cfgManager propagation)
3. Complete Phase 3: User Story 1 (App Key prompt)
4. **STOP and VALIDATE**: `unset CAIRN_DROPBOX_APP_KEY && cairn sync setup` — confirm prompt appears and accepts input
5. US1 alone is a shippable improvement — eliminates the "CAIRN_DROPBOX_APP_KEY is required" fatal error

### Incremental Delivery

1. Setup + Foundational → build passes with refactored signatures
2. US1 → interactive App Key prompt working → **MVP**
3. US2 → optional DB path prompt added
4. US3 → config write confirmation printed
5. Polish → EOF handling verified, all manual tests pass

---

## Notes

- All changes are in `cmd/cairn/main.go` only — no new files, no new packages
- `config.Manager.WriteConfig()` already handles `MkdirAll` and file creation — no extra directory logic needed
- `bufio.NewReader(os.Stdin)` should be instantiated once in `promptForSetupConfig` and reused for both prompts to avoid consuming characters
- The `runSyncAuth` function also uses `appCfg.DropboxAppKey` — it is **not** affected by this feature (only `runSyncSetup` gets the interactive prompt)
- Total tasks: 14 | US1: 4 | US2: 2 | US3: 1 | Foundational: 3 | Setup: 1 | Polish: 3
