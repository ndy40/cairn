# Tasks: `cairn show` Command (016)

**Input**: `specs/016-show-command/`  
**Prerequisites**: plan.md ✓, spec.md ✓  
**No new dependencies** — all changes use existing packages.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: User story label (US1, US2)

---

## Phase 1: Setup

**Purpose**: Confirm baseline before adding anything.

- [ ] T001 Run `go build ./... && go vet ./...` to confirm clean baseline in repo root

---

## Phase 2: Foundational (Blocking Prerequisite)

**Purpose**: `Store.Count()` is required by both user stories — must land first.

- [ ] T002 Add `Count() (int64, error)` method to `internal/store/bookmark.go` using `SELECT COUNT(1) FROM bookmarks`

**Checkpoint**: `go build ./...` passes; `Store` now exposes `Count()`.

---

## Phase 3: User Story 1 — Database Path and Bookmark Count (Priority: P1) 🎯 MVP

**Goal**: `cairn show` prints the resolved database path and total bookmark count.

**Independent Test**: Run `cairn show`; output contains lines starting with `Database:` and `Bookmarks:` with a non-negative integer.

### Implementation for User Story 1

- [ ] T003 [US1] Create `cmd/cairn/run_show.go` with `runShow(db string, cfgManager *config.Manager)` that opens the store, calls `Count()`, and prints `Database:` and `Bookmarks:` lines
- [ ] T004 [P] [US1] Add `cmdShow(ctx cmdContext)` handler to `cmd/cairn/commands.go` (checks `-h`/`--help`, calls `runShow`)
- [ ] T005 [P] [US1] Register `"show": {run: cmdShow, autoSync: false}` in the `commands` map in `cmd/cairn/main.go`

**Checkpoint**: `cairn show` prints database path and bookmark count; `go build ./... && go vet ./...` passes.

---

## Phase 4: User Story 2 — Sync Status and Last Sync Date (Priority: P2)

**Goal**: `cairn show` also reports whether sync is configured and the timestamp of the last sync.

**Independent Test**:
- With no sync config: output contains `Sync:        not configured`.
- With sync configured + past sync: output contains `Sync:        configured (dropbox)` and a `Last sync:` timestamp line.
- With sync configured but never run: output contains `Last sync:   never`.

### Implementation for User Story 2

- [ ] T006 [US2] Extend `runShow` in `cmd/cairn/run_show.go` to load sync config via `csync.LoadConfig`, call `csync.IsConfigured`, and append `Sync:` and `Last sync:` lines to the output

**Checkpoint**: `cairn show` prints all four fields correctly under all sync states.

---

## Phase 5: Polish & Cross-Cutting Concerns

- [ ] T007 [P] Run `go build ./... && go vet ./...` and confirm clean build
- [ ] T008 [P] Run `go test ./...` to confirm no regressions across the whole project

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies.
- **Phase 2 (Foundational)**: Depends on Phase 1. Blocks both user stories.
- **Phase 3 (US1)**: Depends on Phase 2. T004 and T005 can run in parallel with T003 once Phase 2 is done.
- **Phase 4 (US2)**: Depends on T003 (extends `runShow`). T006 is a single targeted edit.
- **Phase 5 (Polish)**: Depends on Phases 3 and 4.

### Parallel Opportunities

- T004 (`commands.go`) and T005 (`main.go`) can be written in parallel with T003 (`run_show.go`) since they touch different files.
- T007 and T008 can run in parallel.

---

## Parallel Example: User Story 1

```
# After T002 completes, launch in parallel:
Task T003: Create cmd/cairn/run_show.go
Task T004: Add cmdShow to cmd/cairn/commands.go
Task T005: Register "show" in cmd/cairn/main.go
```

---

## Implementation Strategy

### MVP First (US1 Only)

1. T001 — baseline check.
2. T002 — add `Store.Count()`.
3. T003 + T004 + T005 (parallel) — wire up `cairn show` with DB path + count.
4. **STOP and VALIDATE**: `cairn show` prints `Database:` and `Bookmarks:`.

### Full Feature (US1 + US2)

5. T006 — extend `runShow` with sync fields.
6. T007 + T008 — build and test.

---

## Notes

- Total tasks: **8** (T001–T008)
- US1 tasks: 3 (T003–T005)
- US2 tasks: 1 (T006)
- Foundational tasks: 1 (T002)
- Parallel opportunities: T003‖T004‖T005, T007‖T008
- All changes confined to `internal/store/bookmark.go`, `cmd/cairn/run_show.go` (new), `cmd/cairn/commands.go`, `cmd/cairn/main.go`
