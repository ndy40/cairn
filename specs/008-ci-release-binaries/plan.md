# Implementation Plan: CI Release Binaries

**Branch**: `008-ci-release-binaries` | **Date**: 2026-03-11 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/008-ci-release-binaries/spec.md`

## Summary

Add GitHub Actions workflows to automatically compile cairn for Linux, macOS, and Windows on every pull request (build validation) and on version tags (release with downloadable binaries). Releases include SHA256 checksums and embed the tag version into the binary via `-ldflags`.

## Technical Context

**Language/Version**: Go 1.25.0 (compiled with `CGO_ENABLED=0` for static binaries)
**Primary Dependencies**: GitHub Actions (runs-on: ubuntu-latest), `actions/checkout`, `actions/setup-go`, `actions/upload-artifact`, `softprops/action-gh-release`
**Storage**: N/A (CI-only feature, no application storage changes)
**Testing**: Workflow validated by pushing a tag and verifying release creation
**Target Platform**: GitHub-hosted runners (ubuntu-latest)
**Project Type**: CI/CD configuration
**Performance Goals**: Release pipeline completes within 10 minutes; PR builds within 5 minutes
**Constraints**: Zero CGO (already enforced), single binary output per platform
**Scale/Scope**: 5 target binaries (linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate                           | Status | Notes                                                              |
|--------------------------------|--------|--------------------------------------------------------------------|
| No CGO                         | PASS   | All binaries built with `CGO_ENABLED=0`                            |
| Single binary                  | PASS   | Each platform produces a single static executable                  |
| Task management                | PASS   | Tasks will be created via Backlog CLI after `/speckit.tasks`       |
| Backward-compatible migrations | N/A    | No schema changes — CI-only feature                                |

## Project Structure

### Documentation (this feature)

```text
specs/008-ci-release-binaries/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (minimal — no application data changes)
├── quickstart.md        # Phase 1 output
└── contracts/           # Phase 1 output (workflow trigger contracts)
```

### Source Code (repository root)

```text
.github/
└── workflows/
    ├── build.yml        # PR build validation (compile all platforms, no release)
    └── release.yml      # Tag-triggered release (compile + create GitHub release with assets)
```

**Structure Decision**: Two separate workflow files for clarity — `build.yml` for PR validation, `release.yml` for tag-triggered releases. This keeps concerns separated and makes it easy to modify triggers independently.

## Complexity Tracking

No constitution violations. No complexity justifications needed.
