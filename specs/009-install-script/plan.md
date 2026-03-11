# Implementation Plan: Installation Script

**Branch**: `009-install-script` | **Date**: 2026-03-11 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/009-install-script/spec.md`

## Summary

Provide a single shell script (`install.sh`) that users can run via `curl | sh` to install the Cairn CLI binary on Linux or macOS. The script detects OS/architecture, downloads the correct binary from GitHub Releases, verifies its checksum, and places it on the user's PATH. Optionally, it detects whether Vicinae is installed and offers to install the Cairn extension.

## Technical Context

**Language/Version**: POSIX shell (sh), no bash-specific features required for maximum portability
**Primary Dependencies**: curl or wget (for downloading), sha256sum or shasum (for checksum verification)
**Storage**: N/A (file-based installation only)
**Testing**: Manual testing on Linux (amd64/arm64) and macOS (amd64/arm64); automated via CI with Docker containers
**Target Platform**: Linux (amd64, arm64), macOS (amd64, arm64)
**Project Type**: CLI installer script
**Performance Goals**: Installation completes within 60 seconds on standard broadband
**Constraints**: Must work with only standard Unix utilities (no Python, no Go, no package managers required)
**Scale/Scope**: Single script file, ~200-400 lines

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate                           | Status | Notes                                                                                     |
|--------------------------------|--------|-------------------------------------------------------------------------------------------|
| No CGO                         | PASS   | This feature is a shell script, not Go code. No CGO involved.                             |
| Single binary                  | PASS   | The script installs the existing single binary. No new binaries introduced.               |
| Task management                | PASS   | Tasks will be created via Backlog CLI after `/speckit.tasks`.                              |
| Backward-compatible migrations | N/A    | No schema changes. This is an installer, not application code.                            |

## Project Structure

### Documentation (this feature)

```text
specs/009-install-script/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (minimal - no persistent data)
├── quickstart.md        # Phase 1 output
└── contracts/
    └── cli-interface.md # Script flags and environment variables
```

### Source Code (repository root)

```text
install.sh               # The installation script (repo root, downloadable via raw URL)
```

**Structure Decision**: A single `install.sh` file at the repository root. No new directories or packages needed. This is a standalone shell script that downloads and installs existing build artifacts.

## Complexity Tracking

No constitution violations. Table not needed.
