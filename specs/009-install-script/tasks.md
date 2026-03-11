# Tasks: Installation Script

**Input**: Design documents from `/specs/009-install-script/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: No test tasks included (not explicitly requested in feature specification).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: Create the install script skeleton with argument parsing and utility functions

- [x] T001 Create install.sh at repository root with POSIX sh shebang, license header, and `set -eu` error handling
- [x] T002 Implement argument parsing in install.sh supporting `--install-dir`/`-d`, `--version`/`-v`, `--non-interactive`/`-y`, `--with-extension`, and `--help`/`-h` flags per contracts/cli-interface.md
- [x] T003 Implement utility functions in install.sh: `log_info`, `log_error`, `log_success` for consistent output formatting, and `has_command` to check for available commands (curl/wget, sha256sum/shasum)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: OS/architecture detection and download infrastructure that ALL user stories depend on

**Warning**: No user story work can begin until this phase is complete

- [x] T004 Implement OS detection in install.sh using `uname -s` to identify Linux and Darwin (macOS), exiting with code 2 and clear error message for unsupported platforms (FreeBSD, Windows/MSYS, etc.)
- [x] T005 Implement architecture detection in install.sh using `uname -m` to map to `amd64` (x86_64) and `arm64` (aarch64), exiting with code 2 for unsupported architectures
- [x] T006 Implement download function in install.sh that uses curl (preferred) or wget (fallback) to fetch files from GitHub Releases, with proper error handling for network failures (exit code 1) and HTTP errors
- [x] T007 Implement GitHub Release URL resolution in install.sh: construct download URLs for binary (`cairn-{os}-{arch}`) and checksums.txt from `https://github.com/ndy40/cairn/releases/download/{version}/`, defaulting to latest release via GitHub API when `--version` is not specified

**Checkpoint**: Foundation ready - OS/arch detection works, download infrastructure functional

---

## Phase 3: User Story 1 - Install the CLI Tool (Priority: P1) MVP

**Goal**: Users can install the Cairn CLI binary with a single command

**Independent Test**: Run `sh install.sh --install-dir /tmp/cairn-test` on a clean machine and verify `/tmp/cairn-test/cairn --help` works

### Implementation for User Story 1

- [x] T008 [US1] Implement checksum verification in install.sh: download checksums.txt, extract the expected SHA-256 for the target binary, compute checksum of downloaded file using sha256sum or shasum, and exit with code 3 if verification fails
- [x] T009 [US1] Implement binary installation in install.sh: create install directory (`~/.local/bin` default or `--install-dir` / `CAIRN_INSTALL_DIR`), copy verified binary as `cairn`, set executable permissions (`chmod +x`), exit with code 4 if directory creation or write fails due to permissions
- [x] T010 [US1] Implement PATH check in install.sh: after installation, verify the install directory is on the user's PATH; if not, print instructions for adding it (with shell-specific advice for bash/zsh profile files)
- [x] T011 [US1] Implement upgrade handling in install.sh: if `cairn` already exists in the target directory, back up the existing binary to `cairn.bak` before overwriting, and remove the backup on success
- [x] T012 [US1] Implement the main installation flow in install.sh: wire together detection, download, verification, and installation steps with progress messages; print version installed and success message on completion

**Checkpoint**: User Story 1 complete - `curl -sSL .../install.sh | sh` installs a working `cairn` binary

---

## Phase 4: User Story 2 - Vicinae Detection and Extension Installation (Priority: P2)

**Goal**: After CLI installation, detect Vicinae and offer to install the Cairn browser extension

**Independent Test**: Run install.sh on a machine with `vici` on PATH and verify the extension prompt appears; accept and verify extension is installed

### Implementation for User Story 2

- [x] T013 [US2] Implement Vicinae detection in install.sh: use `command -v vici` to check if Vicinae CLI is available on PATH; if not found, skip extension installation with an informational message
- [x] T014 [US2] Implement interactive extension prompt in install.sh: when Vicinae is detected and script is running interactively (TTY check via `[ -t 0 ]`), prompt the user "Vicinae detected. Install Cairn extension? [y/N]" and process their response
- [x] T015 [US2] Implement extension download and installation in install.sh: download pre-built extension archive (`cairn-vicinae-extension.tar.gz`) from the same GitHub Release, verify its checksum, extract to a temporary directory, and copy to Vicinae extensions directory (discovered via `vici` CLI or standard paths)
- [x] T016 [US2] Update .github/workflows/release.yml to build the Vicinae extension (`cd vicinae-extension && npm ci && npm run build`) and include the built extension as `cairn-vicinae-extension.tar.gz` in the release artifacts alongside the CLI binaries and checksums

**Checkpoint**: User Story 2 complete - Vicinae users are prompted and can install the extension automatically

---

## Phase 5: User Story 3 - Unattended Installation (Priority: P3)

**Goal**: Support non-interactive installation for scripted provisioning and CI environments

**Independent Test**: Run `sh install.sh --non-interactive --install-dir /tmp/cairn-test` and verify CLI installs without any prompts; run `sh install.sh --non-interactive --with-extension --install-dir /tmp/cairn-test` with `vici` available and verify both CLI and extension install without prompts

### Implementation for User Story 3

- [x] T017 [US3] Wire non-interactive mode in install.sh: when `--non-interactive`/`-y` flag is set, skip all prompts (Vicinae extension prompt) and install CLI only by default; also auto-detect piped input (`! [ -t 0 ]`) and treat as non-interactive
- [x] T018 [US3] Wire `--with-extension` flag in install.sh: when combined with `--non-interactive`, automatically install the Vicinae extension without prompting (if Vicinae is detected); if Vicinae is not detected, print a warning to stderr and continue with CLI-only installation

**Checkpoint**: All user stories independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final quality improvements across the entire script

- [x] T019 Add comprehensive error messages to install.sh for all exit codes (1-4) per contracts/cli-interface.md, ensuring each error message includes actionable guidance (e.g., "Permission denied: try running with --install-dir ~/bin or use sudo")
- [x] T020 Add `--help` output to install.sh displaying all flags, environment variables, examples, and exit codes per contracts/cli-interface.md
- [x] T021 Implement cleanup in install.sh: use a trap handler to remove temporary files (downloaded binary, checksums.txt, extension archive) on exit regardless of success or failure
- [x] T022 Run end-to-end validation of install.sh on Linux and macOS per quickstart.md test scenarios: fresh install, upgrade, non-interactive mode, specific version, custom install directory

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - US1 (Phase 3): No dependencies on other stories
  - US2 (Phase 4): No dependency on US1 for implementation, but logically extends the install flow
  - US3 (Phase 5): Depends on US1 and US2 being implemented (wires flags into existing prompts)
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Independent of US1 but builds on same script
- **User Story 3 (P3)**: Depends on US1 and US2 being implemented (adds flags to control their behavior)

### Within Each User Story

- Tasks within a story are sequential (all modify the same file: install.sh)
- No [P] markers within stories because all tasks modify the same file

### Parallel Opportunities

- T004 and T005 (OS and arch detection) could be developed as independent functions, but both target install.sh
- US1 and US2 are logically independent but share install.sh, so sequential implementation is recommended
- T016 (release workflow update) is independent of all install.sh tasks and can be done in parallel

---

## Parallel Example: Foundational Phase

```bash
# T016 (release workflow) can run in parallel with all install.sh work:
Task: "Update .github/workflows/release.yml to build Vicinae extension"
# runs independently of install.sh development
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T003)
2. Complete Phase 2: Foundational (T004-T007)
3. Complete Phase 3: User Story 1 (T008-T012)
4. **STOP and VALIDATE**: Run `sh install.sh` on Linux and macOS, verify `cairn --help` works
5. Deploy/merge if ready - users can install the CLI

### Incremental Delivery

1. Setup + Foundational → Script skeleton ready
2. Add User Story 1 → CLI installation works → **MVP!**
3. Add User Story 2 → Vicinae extension installation works
4. Add User Story 3 → Non-interactive mode works
5. Polish → Error handling, cleanup, help text polished

---

## Notes

- All tasks modify a single file (install.sh) except T016 (release workflow) — sequential execution within each phase is required
- No test tasks generated (not requested in spec)
- T016 (release workflow for extension artifact) is a prerequisite for US2 runtime but can be developed in parallel
- Commit after each completed phase for clean history
- Stop at any checkpoint to validate the story independently
