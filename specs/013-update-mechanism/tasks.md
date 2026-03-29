# Tasks: Cairn Self-Update Mechanism

**Input**: Design documents from `/specs/013-update-mechanism/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli-contract.md, quickstart.md

**Tests**: Included — network-dependent code requires httptest-based unit tests for reliable verification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: Create the new package skeleton that all subsequent tasks will build on.

- [x] T001 Create `internal/updater/` package with empty `updater.go` (package declaration only) and `updater_test.go` stub

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared HTTP client, version check, and checksum helpers required by every user story.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T002 Implement HTTP client helper (8s timeout, `cairn/<version>` User-Agent) and `CheckLatestVersion(currentVersion string) (latest string, available bool, err error)` in `internal/updater/updater.go` — queries GitHub Releases API, parses `tag_name`, treats `"dev"` as always out of date
- [x] T003 [P] Implement `downloadFile(url, destPath string) error` and `verifyChecksum(filePath, expectedHash string) error` helpers in `internal/updater/updater.go` — stream to temp file, SHA256 via `crypto/sha256`, hex compare

**Checkpoint**: Foundation ready — user story implementation can now begin.

---

## Phase 3: User Story 1 — Check for and Apply CLI Update (Priority: P1) 🎯 MVP

**Goal**: `cairn update` checks for a newer binary, downloads it, verifies SHA256, backs up the current binary, and atomically replaces it. On no-update or error, nothing changes on disk.

**Independent Test**: Install an older cairn binary, run `cairn update`, verify `cairn version` reports the new version and the `.bak` file is absent.

### Implementation for User Story 1

- [x] T004 [US1] Implement `UpdateBinary(currentVersion, latestVersion string) error` in `internal/updater/updater.go` — resolve binary path via `os.Executable()+filepath.EvalSymlinks`, build platform asset name (`cairn-{GOOS}-{GOARCH}`), download binary + `checksums.txt`, verify SHA256, backup to `.bak`, attempt `os.Rename` (fall back to copy+delete on EXDEV), remove `.bak` on success, restore `.bak` on failure
- [x] T005 [US1] Add `update` subcommand entry to `cmd/cairn/main.go` — parse `--check` and `--extension` flags, call `updater.CheckLatestVersion(version)`, branch on `--check` vs apply, print messages per `contracts/cli-contract.md`, map errors to exit codes 0/1/3/4
- [x] T006 [P] [US1] Write unit tests for `UpdateBinary` in `internal/updater/updater_test.go` using `net/http/httptest` — cover: update applied, already up to date, network error, checksum mismatch, permission error (unwritable target dir)

**Checkpoint**: `cairn update` and `cairn update --check` (binary path) are fully functional and independently testable.

---

## Phase 4: User Story 2 — Check for Available Update Without Applying (Priority: P2)

**Goal**: `cairn update --check` prints current vs latest version and exits without downloading, modifying, or creating any files.

**Independent Test**: Run `cairn update --check`; verify stdout contains version information and that no temp files, `.bak` files, or binary changes occur on disk.

### Implementation for User Story 2

- [x] T007 [US2] Verify the `--check` branch in the `update` handler in `cmd/cairn/main.go` exits after printing version status — ensure no call to `UpdateBinary` is made and no files are touched when `--check` is set
- [x] T008 [P] [US2] Write unit tests for `--check` path in `internal/updater/updater_test.go` — assert no files are created or modified, correct stdout output for update-available and already-up-to-date cases

**Checkpoint**: `cairn update --check` is independently testable and verifiably makes no file-system changes.

---

## Phase 5: User Story 3 — Update the Vicinae Extension (Priority: P3)

**Goal**: `cairn update --extension` detects whether the extension is installed, compares its version, and replaces the extension files if a newer version exists. Reports "not installed" clearly when the extension directory is absent.

**Independent Test**: Create a mock extension directory with a `version.txt` set to an older version, run `cairn update --extension`, verify `version.txt` is updated and extension files are replaced.

### Implementation for User Story 3

- [x] T009 [US3] Implement `DetectExtension() (dir string, installed bool)` in `internal/updater/updater.go` — mirror `install.sh` XDG logic: macOS `~/Library/Application Support/vicinae/extensions/cairn`, Linux `$XDG_DATA_HOME/vicinae/extensions/cairn` (fallback `~/.local/share/...`)
- [x] T010 [P] [US3] Implement `CheckExtensionVersion(dir string) (current, latest string, available bool, err error)` in `internal/updater/updater.go` — read `dir/version.txt`, call `CheckLatestVersion`, return comparison
- [x] T011 [US3] Implement `UpdateExtension(dir, latestVersion string) error` in `internal/updater/updater.go` — build archive asset name (`vicinae-extension-{version}.tar.gz`), download + verify SHA256, extract to `dir`, write `dir/version.txt`
- [x] T012 [US3] Wire `--extension` flag in `cmd/cairn/main.go` — when set, call `DetectExtension`, print "not installed" if absent, otherwise call `CheckExtensionVersion` and `UpdateExtension`; support `--extension --check` combination
- [x] T013 [P] [US3] Write unit tests for extension update path in `internal/updater/updater_test.go` — cover: not installed, already up to date, update applied, checksum mismatch, `--extension --check` (no file changes)

**Checkpoint**: All three user stories are independently functional and testable.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Help text, platform guard, and final build verification.

- [x] T014 Add `update` help text to `cmd/cairn/main.go` — integrate with existing `cairn help` and `cairn update --help` per `contracts/cli-contract.md` (flags, exit codes)
- [x] T015 [P] Add Windows guard in `internal/updater/updater.go` — when `runtime.GOOS == "windows"`, `UpdateBinary` and `UpdateExtension` print "Windows in-process update is not supported; re-run the install script to update" and return nil
- [x] T016 Run `go build ./... && go vet ./...` and `go test ./internal/updater/...`; fix any issues

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 — **blocks all user stories**
- **User Stories (Phase 3–5)**: All depend on Phase 2; US1 and US3 share updater.go but in distinct functions, so they can proceed in parallel once Phase 2 is done
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **US1 (P1)**: Depends on Phase 2 only — no dependency on US2 or US3
- **US2 (P2)**: Depends on US1 (T005 handler must exist to add the `--check` branch verification)
- **US3 (P3)**: Depends on Phase 2 only — independent of US1 and US2

### Within Each User Story

- Shared helpers (T002, T003) before any story implementation
- `UpdateBinary` (T004) before subcommand wiring (T005)
- Implementation before tests where tests verify the implementation (T006, T008, T013)
- All US3 functions (T009, T010) before `UpdateExtension` (T011)

### Parallel Opportunities

- T003 (checksum helper) can be written in parallel with T002 (HTTP client + version check) — different function groups in the same file
- T006, T008, T013 (tests) can be written in parallel with their respective implementations once function signatures are defined
- T009 and T010 (extension detection + version check) can be written in parallel
- T014 and T015 (help text and Windows guard) can be written in parallel

---

## Parallel Example: User Story 1

```bash
# After T002 + T003 (Foundational) complete, these can run in parallel:
Task T004: "Implement UpdateBinary in internal/updater/updater.go"
Task T006: "Write unit tests (skeleton) in internal/updater/updater_test.go"
# Then T005 wires everything into main.go
```

## Parallel Example: User Story 3

```bash
# T009 and T010 can run in parallel (different functions, same file):
Task T009: "Implement DetectExtension in internal/updater/updater.go"
Task T010: "Implement CheckExtensionVersion in internal/updater/updater.go"
# Then T011 (UpdateExtension) — depends on T009 + T010
# Then T012 (wire in main.go) + T013 (tests) in parallel
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001)
2. Complete Phase 2: Foundational (T002, T003)
3. Complete Phase 3: US1 (T004, T005, T006)
4. **STOP and VALIDATE**: Run `cairn update --check` and `cairn update` against a real or mock older version
5. Ship MVP — users can now self-update the binary

### Incremental Delivery

1. Setup + Foundational → shared infrastructure ready
2. US1 → binary update works → **MVP**
3. US2 → `--check` flag verified and tested independently
4. US3 → extension update works
5. Polish → help text, Windows guard, final build pass

### Parallel Team Strategy

With two developers after Phase 2 is complete:
- Developer A: US1 (T004 → T005 → T006) + US2 (T007 → T008)
- Developer B: US3 (T009, T010 in parallel → T011 → T012 → T013)
- Both finish at Phase 6 together

---

## Notes

- [P] tasks = different functions/files, no blocking dependencies
- Tests use `net/http/httptest` to serve mock GitHub API and binary download responses — no real network calls in tests
- The `version.txt` approach for extension versioning is new; write it during `UpdateExtension` and handle the case where it is absent (treat as `"unknown"` / always update)
- `os.Rename` is atomic on POSIX same-filesystem; the temp file is created in the same directory as the target binary to guarantee same-filesystem placement
- Commit after each phase checkpoint to keep git history clean
