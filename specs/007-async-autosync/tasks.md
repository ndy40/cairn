# Tasks: Async Autosync

**Input**: Design documents from `/specs/007-async-autosync/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

**Tests**: Not explicitly requested in spec. No test tasks generated.

**Organization**: Tasks are grouped by user story. Since US1-US3 share the same mechanism (background subprocess helper), they are sequenced within a single phase after the foundational helper is built.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Foundational (Background Subprocess Helper)

**Purpose**: Create the core helper function that all user stories depend on

- [x] T001 Implement `backgroundSyncPush` helper function in cmd/cairn/main.go that uses `os.Executable()` to resolve the current binary path, creates an `exec.Command` for `cairn sync push`, redirects stdout/stderr to `os.DevNull`, and calls `cmd.Start()` without `cmd.Wait()`. The function should accept `dbPath string` as parameter (to match the current `autoSyncPush` signature) but only use it to verify sync is configured before spawning. If sync is not configured or `os.Executable()` fails, the function should return silently (no error output). The existing `autoSyncPush` function must be preserved unchanged for now.

**Checkpoint**: Helper function exists and compiles. `go build ./cmd/cairn/...` succeeds.

---

## Phase 2: User Story 1 - Non-blocking bookmark creation (Priority: P1)

**Goal**: `cairn add` returns immediately after local save; sync runs in background

**Independent Test**: Run `time cairn add <url>` with sync configured — should return in < 500ms

- [x] T002 [US1] Replace the synchronous `autoSyncPush(dbPath)` call in `runAdd` function in cmd/cairn/main.go (currently called after successful bookmark save, around line 180) with `backgroundSyncPush(dbPath)`. Also replace the `autoSyncPush(dbPath)` call on the fetch-failure path (around line 176) with `backgroundSyncPush(dbPath)`. Verify `go build ./cmd/cairn/...` succeeds.

**Checkpoint**: `cairn add` returns immediately. Background sync push is spawned.

---

## Phase 3: User Story 2 - Non-blocking bookmark deletion (Priority: P1)

**Goal**: `cairn delete` returns immediately after local delete; sync runs in background

**Independent Test**: Run `cairn delete` with sync configured — should return in < 500ms

- [x] T003 [US2] Replace the synchronous `autoSyncPush(dbPath)` call in `runDelete` function in cmd/cairn/main.go (around line 264) with `backgroundSyncPush(dbPath)`. Verify `go build ./cmd/cairn/...` succeeds.

**Checkpoint**: `cairn delete` returns immediately. Background sync push is spawned.

---

## Phase 4: User Story 3 - Non-blocking bookmark editing (Priority: P1)

**Goal**: `cairn edit` (tag updates) returns immediately; sync runs in background

**Independent Test**: Edit bookmark tags with sync configured — should return in < 500ms

- [x] T004 [US3] Locate any `autoSyncPush(dbPath)` call after tag-editing operations in cmd/cairn/main.go and replace with `backgroundSyncPush(dbPath)`. If no autosync call exists after edit operations yet, add a `backgroundSyncPush(dbPath)` call after successful tag update in the appropriate edit handler. Verify `go build ./cmd/cairn/...` succeeds.

**Checkpoint**: Tag editing returns immediately. Background sync push is spawned if applicable.

---

## Phase 5: User Story 4 - Non-blocking startup auto-pull (Priority: P1)

**Goal**: CLI startup no longer blocks on remote sync pull; pull runs in background

**Independent Test**: Run `time cairn list` with sync configured — should return in < 500ms

- [x] T008 [US4] Replace synchronous `autoSyncPull(resolvedDB)` call in main() with `backgroundSyncPull()` that spawns a detached `cairn sync pull` subprocess. Remove old `autoSyncPull` function. Verify `go build` and `go vet` pass.

**Checkpoint**: All CLI commands start immediately without waiting for remote pull.

---

## Phase 6: User Story 5 - Background sync failure handling (Priority: P2)

**Goal**: Background sync failures are silent; pending changes remain queued; explicit sync still reports errors

**Independent Test**: Disconnect network, add a bookmark, verify no error output. Then run `cairn sync push` and verify it reports the error interactively.

- [x] T005 [US4] Verify that the `backgroundSyncPush` helper in cmd/cairn/main.go properly suppresses all output by ensuring stdout and stderr of the spawned subprocess are set to `os.DevNull`. Verify that the existing explicit `cairn sync push` and `cairn sync pull` commands remain synchronous and interactive (they should not be affected since they go through `runSyncPush`/`runSyncPull` handlers, not through `autoSyncPush`). No code changes expected if T001 was implemented correctly — this is a verification task.

**Checkpoint**: Background failures are silent. Explicit sync commands show errors normally.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Cleanup and final validation

- [x] T006 Remove the old synchronous `autoSyncPush` and `autoSyncPull` functions from cmd/cairn/main.go since they are no longer called anywhere. Run `go vet ./...` and `go build ./cmd/cairn/...` to verify no compilation errors or unused code warnings.
- [x] T007 Run quickstart.md validation: execute the manual test steps from specs/007-async-autosync/quickstart.md to verify non-blocking add, sync eventual consistency, offline behavior, and explicit sync unchanged behavior.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Foundational)**: No dependencies - start immediately
- **Phase 2 (US1)**: Depends on T001
- **Phase 3 (US2)**: Depends on T001 (can run in parallel with Phase 2)
- **Phase 4 (US3)**: Depends on T001 (can run in parallel with Phase 2-3)
- **Phase 5 (US4)**: Depends on T001-T004 (verification of all call sites)
- **Phase 6 (Polish)**: Depends on T001-T005

### User Story Dependencies

- **US1 (P1)**: Depends only on T001 (foundational helper)
- **US2 (P1)**: Depends only on T001 — independent of US1
- **US3 (P1)**: Depends only on T001 — independent of US1/US2
- **US4 (P2)**: Depends on T001-T004 (verifies behavior across all call sites)

### Within Each User Story

All user stories are single-task (one call site replacement each), so no intra-story ordering needed.

### Parallel Opportunities

- T002, T003, T004 modify different functions in the same file. They CAN be done sequentially in one pass but are logically independent.
- T005 is a verification task and must come after T001-T004.

---

## Parallel Example: User Stories 1-3

```bash
# After T001 (helper) is complete, these can be done in a single editing pass:
Task T002: "Replace autoSyncPush in runAdd (cmd/cairn/main.go)"
Task T003: "Replace autoSyncPush in runDelete (cmd/cairn/main.go)"
Task T004: "Replace/add backgroundSyncPush after edit (cmd/cairn/main.go)"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete T001: Create backgroundSyncPush helper
2. Complete T002: Replace autoSyncPush in runAdd
3. **STOP and VALIDATE**: `time cairn add <url>` returns in < 500ms
4. Deploy if ready

### Incremental Delivery

1. T001 → Foundation ready
2. T002 → Non-blocking add (MVP!)
3. T003 → Non-blocking delete
4. T004 → Non-blocking edit
5. T005 → Verify failure handling
6. T006-T007 → Cleanup and validation

---

## Notes

- All changes are in a single file: `cmd/cairn/main.go`
- No new dependencies, no schema changes, no new files
- The sync engine, store, and backend are completely untouched
- Auto-pull on startup remains synchronous (not in scope)
- TUI does not call autosync, so no TUI changes needed
- Commit after each task or after T002-T004 as a logical group
