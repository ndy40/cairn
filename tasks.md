# Tasks: Makefile Build & Install

**Input**: User request — add a Makefile with build and install targets for the cairn CLI.

**Context**: Go 1.25.0 project, module `github.com/ndy40/cairn`, entry point `cmd/cairn/main.go`. Binary name: `cairn`. Install location: `$HOME/.local/bin`.

## Phase 1: Setup

**Purpose**: Create the Makefile with build and install targets

- [x] T001 Create Makefile at project root (`Makefile`) with a `build` target that compiles `cmd/cairn/main.go` into a binary named `cairn` (output to project root or `bin/` directory) using `go build -o cairn ./cmd/cairn`
- [x] T002 Add `install` target to `Makefile` that depends on `build` and copies the `cairn` binary to `$HOME/.local/bin/cairn`, creating the directory if it doesn't exist (`mkdir -p $(HOME)/.local/bin`)

---

## Phase 2: Polish

**Purpose**: Ensure completeness

- [x] T003 Add a `clean` target to remove the built binary
- [x] T004 Verify `make`, `make build`, `make install`, and `make clean` all work correctly
- [x] T005 Add an `extension` target to `Makefile` that runs `npm run build` in `vicinae-extension/` directory to build the Vicinae extension
- [x] T006 Verify `make extension` works correctly

---

## Dependencies & Execution Order

- **T001**: No dependencies — create the Makefile with build target
- **T002**: Depends on T001 — adds install target to same file
- **T003**: Depends on T001 — adds clean target
- **T004**: Depends on T001–T003 — validation

### Parallel Opportunities

- T002 and T003 can be written in parallel (different targets, same file — but simple enough to do sequentially)

---

## Implementation Strategy

### MVP First

1. T001: Create Makefile with `build` target
2. T002: Add `install` target
3. **STOP and VALIDATE**: Run `make build` and `make install`

### Summary

- **Total tasks**: 4
- **Core tasks**: 2 (build + install)
- **Polish tasks**: 2 (clean + validation)
- **No tests requested**
