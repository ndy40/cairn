# Tasks: Configuration File Support

**Input**: Design documents from `/specs/010-config-file/`
**Prerequisites**: plan.md, spec.md, data-model.md, contracts/cli-interface.md, research.md, quickstart.md

**Tests**: Not explicitly requested — no test tasks included.

**Organization**: Tasks grouped by user story. Precedence: env vars > CLI flags > config file > defaults.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: Create the new config package with core types

- [x] T001 Create `internal/config/config.go` with `FileConfig` struct (pointer fields: `DBPath *string`, `DropboxAppKey *string`), `DefaultConfigPath()` function returning OS-appropriate path to `cairn.json`, and `LoadFile(path string) (*FileConfig, error)` function that reads and parses the JSON file (returns nil if file doesn't exist, error if invalid JSON or wrong types)

---

## Phase 2: Foundational

**Purpose**: Config resolution logic that all user stories depend on

**CRITICAL**: Must complete before user story work begins

- [x] T002 Add `Resolve(fileConfig *FileConfig, cliDBFlag string, defaultDBPath string) AppConfig` function to `internal/config/config.go` that merges sources in precedence order: start with defaults, overlay config file values (if non-nil pointers), overlay CLI flag values (if non-empty), overlay environment variables (`CAIRN_DB_PATH`, `CAIRN_DROPBOX_APP_KEY`) if set — returning the final `AppConfig` struct

**Checkpoint**: Config package complete — can load, parse, and resolve configuration from all sources

---

## Phase 3: User Story 1 - Load Configuration from File (Priority: P1) MVP

**Goal**: cairn reads settings from optional `cairn.json` and uses them when no higher-precedence source overrides

**Independent Test**: Create `~/.config/cairn/cairn.json` with `{"db_path": "/tmp/test.db"}`, unset `CAIRN_DB_PATH`, run `cairn config` — output should show `CAIRN_DB_PATH=/tmp/test.db`

### Implementation for User Story 1

- [x] T003 [US1] Update `cmd/cairn/main.go` to call `config.DefaultConfigPath()` and `config.LoadFile()` at startup, storing the loaded `FileConfig`; if `LoadFile` returns an error (invalid JSON / wrong types), print the error to stderr and exit with code 1
- [x] T004 [US1] Update `resolveDBPath` in `cmd/cairn/main.go` to accept the loaded `FileConfig` and call `config.Resolve()` instead of directly checking env var and default — wire the resolved `AppConfig.DBPath` into the store initialization
- [x] T005 [US1] Update the `config` subcommand handler in `cmd/cairn/main.go` to use `config.Resolve()` for both `CAIRN_DB_PATH` and `CAIRN_DROPBOX_APP_KEY` output values, replacing the current direct `os.Getenv` calls

**Checkpoint**: Config file loading works end-to-end. `cairn config` reflects values from `cairn.json`.

---

## Phase 4: User Story 2 - Configuration Precedence (Priority: P2)

**Goal**: Precedence order is correctly enforced: env vars > CLI flags > config file > defaults

**Independent Test**: Set `CAIRN_DB_PATH=/env/path.db`, create `cairn.json` with `{"db_path": "/file/path.db"}`, run `cairn --db /flag/path.db config` — output should show `CAIRN_DB_PATH=/env/path.db`

### Implementation for User Story 2

- [x] T006 [US2] Update all call sites in `cmd/cairn/main.go` that read `CAIRN_DROPBOX_APP_KEY` via `os.Getenv` (sync setup, sync auth commands) to use the resolved `AppConfig.DropboxAppKey` instead, ensuring env var override works consistently across all commands

**Checkpoint**: All 4 precedence scenarios from spec acceptance criteria pass correctly.

---

## Phase 5: User Story 3 - README Documentation (Priority: P3)

**Goal**: Users can find configuration docs in the README

**Independent Test**: Open `README.md` and find a "Configuration" section with config file paths, JSON keys, env vars, and precedence order

### Implementation for User Story 3

- [x] T007 [US3] Add a "Configuration" section to `README.md` documenting: config file location per OS (Linux/macOS/Windows), supported JSON keys (`db_path`, `dropbox_app_key`), supported environment variables (`CAIRN_DB_PATH`, `CAIRN_DROPBOX_APP_KEY`), CLI flags (`--db`), and precedence order (env vars > CLI flags > config file > defaults)

**Checkpoint**: README contains complete configuration documentation.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Validation and cleanup

- [x] T008 Run `go build ./...` and `go vet ./...` to verify no compilation or lint errors
- [x] T009 Run quickstart.md scenarios 1–6 manually to validate all precedence and error cases

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — create config package
- **Foundational (Phase 2)**: Depends on Phase 1 — adds Resolve function
- **US1 (Phase 3)**: Depends on Phase 2 — integrates config into CLI
- **US2 (Phase 4)**: Depends on Phase 3 — extends config to all Dropbox call sites
- **US3 (Phase 5)**: Can start after Phase 2 (independent of US1/US2 code)
- **Polish (Phase 6)**: Depends on all phases complete

### Within Each User Story

- T003 → T004 → T005 (sequential, same file)
- T006 depends on T005 (needs resolved AppConfig wired in)
- T007 independent of code tasks

### Parallel Opportunities

- T007 (README) can be written in parallel with Phase 3/4 code tasks
- T001 and T007 could start in parallel (different files)

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Create config package (T001)
2. Complete Phase 2: Add Resolve function (T002)
3. Complete Phase 3: Wire into CLI (T003–T005)
4. **STOP and VALIDATE**: Run `cairn config` with and without `cairn.json`
5. If working, continue to US2 and US3

### Incremental Delivery

1. T001–T002: Config package ready
2. T003–T005: Config file loading works (MVP!)
3. T006: Full precedence across all commands
4. T007: Documentation complete
5. T008–T009: Final validation

---

## Notes

- No test tasks generated (not requested)
- All code changes touch 2 files: `internal/config/config.go` (new) and `cmd/cairn/main.go` (modified)
- `README.md` is the only documentation file modified
- No database schema changes, no migrations
- Commit after each phase for clean history
