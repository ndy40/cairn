# Tasks: CI Release Binaries

**Input**: Design documents from `/specs/008-ci-release-binaries/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, contracts/

**Tests**: Not explicitly requested. No test tasks generated.

**Organization**: Two user stories — US1 (tag-triggered release) is the MVP, US2 (PR build validation) adds CI safety net. Both produce separate workflow files so they can be implemented in parallel.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Phase 1: Setup

**Purpose**: Create the directory structure for GitHub Actions workflows

- [x] T001 Create the `.github/workflows/` directory at the repository root if it does not already exist.

**Checkpoint**: Directory `.github/workflows/` exists.

---

## Phase 2: User Story 1 - Tag-triggered release with downloadable binaries (Priority: P1)

**Goal**: Pushing a `v*` tag compiles cairn for 5 platforms, creates a GitHub Release, and attaches all binaries + checksums

**Independent Test**: Push a tag `v0.0.1-test` and verify a release is created with 5 binaries + checksums.txt

- [x] T002 [US1] Create `.github/workflows/release.yml` with the following configuration: trigger on push of tags matching `v*`; use a matrix strategy to cross-compile for 5 targets (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64); each matrix job checks out code, sets up Go 1.25.0, and runs `CGO_ENABLED=0 GOOS=$matrix.goos GOARCH=$matrix.goarch go build -ldflags "-X main.version=${{ github.ref_name }}" -o cairn-$matrix.goos-$matrix.goarch ./cmd/cairn/`; include `.exe` extension for the windows target; upload each binary as a build artifact using `actions/upload-artifact`.

- [x] T003 [US1] Add a release job to `.github/workflows/release.yml` that depends on the build matrix jobs (`needs: build`), downloads all build artifacts using `actions/download-artifact`, generates a `checksums.txt` file using `sha256sum` on all binaries, and creates a GitHub Release using `softprops/action-gh-release@v2` with `generate_release_notes: true` and attaches all binaries plus `checksums.txt` as release assets. The release job must use `permissions: contents: write` to allow creating releases.

**Checkpoint**: Pushing a `v*` tag triggers the workflow, builds all 5 binaries, and creates a release with all assets attached.

---

## Phase 3: User Story 2 - Build validation on pull requests (Priority: P2)

**Goal**: Every PR gets compilation validated for all 5 platforms before merge

**Independent Test**: Open a PR and verify the build workflow runs and reports status

- [x] T004 [P] [US2] Create `.github/workflows/build.yml` with the following configuration: trigger on `pull_request` events (opened, synchronize, reopened) targeting `main` branch; use the same matrix strategy as release.yml to cross-compile for all 5 targets (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64); each matrix job checks out code, sets up Go 1.25.0, and runs `CGO_ENABLED=0 GOOS=$matrix.goos GOARCH=$matrix.goarch go build -o /dev/null ./cmd/cairn/` (output to /dev/null since no artifacts need to be retained); no release creation or artifact upload.

**Checkpoint**: Opening or updating a PR triggers the build workflow and reports pass/fail status.

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Validation and documentation

- [x] T005 Verify both workflow files are valid YAML by running `python3 -c "import yaml; yaml.safe_load(open('.github/workflows/release.yml')); yaml.safe_load(open('.github/workflows/build.yml')); print('Valid')"` or equivalent YAML validation. Fix any syntax errors.
- [x] T006 Run quickstart.md validation: review the manual test steps from specs/008-ci-release-binaries/quickstart.md and confirm the workflow files implement all described behaviors (tag-triggered release, PR validation, binary naming, version embedding, checksums).

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (US1)**: Depends on T001 (directory exists)
- **Phase 3 (US2)**: Depends on T001 (directory exists) — can run in parallel with Phase 2
- **Phase 4 (Polish)**: Depends on T002-T004

### User Story Dependencies

- **US1 (P1)**: Independent — only needs the directory from T001
- **US2 (P2)**: Independent — only needs the directory from T001. Can be implemented in parallel with US1

### Parallel Opportunities

- T002/T003 and T004 work on different files (`release.yml` vs `build.yml`) and can be done in parallel
- T002 and T003 work on the same file sequentially (T003 adds the release job after T002 sets up the build jobs)

---

## Parallel Example: User Stories 1 and 2

```bash
# After T001 (setup) is complete, these can run in parallel:
Task T002+T003: "Create release.yml (build matrix + release job)"
Task T004: "Create build.yml (PR validation)"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete T001: Create directory
2. Complete T002-T003: Create release.yml
3. **STOP and VALIDATE**: Push a test tag and verify release creation
4. Deploy if ready

### Incremental Delivery

1. T001 → Directory ready
2. T002-T003 → Tag-triggered releases work (MVP!)
3. T004 → PR build validation added
4. T005-T006 → Validation and polish

---

## Notes

- All workflow files are new (no existing CI configuration to modify)
- No Go code changes needed — `var version = "dev"` already exists and is overridden by ldflags
- The `.exe` extension for Windows must be handled in the matrix (include block or conditional)
- `softprops/action-gh-release` handles idempotency — won't create duplicate releases
- The `permissions: contents: write` is required for the release job to create releases
