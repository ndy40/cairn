# Workflow Trigger Contracts

**Feature**: 008-ci-release-binaries
**Date**: 2026-03-11

## build.yml — PR Build Validation

### Trigger
- **Event**: `pull_request` (opened, synchronize, reopened)
- **Branches**: `main`

### Behavior
- Compiles for all 5 platform/architecture targets
- Reports pass/fail on the pull request
- Does NOT create a release or upload artifacts to releases

### Outputs
- Build status (pass/fail) per platform

---

## release.yml — Tag-Triggered Release

### Trigger
- **Event**: `push`
- **Filter**: Tags matching `v*` (e.g., `v1.0.0`, `v2.3.1-beta`)

### Behavior
- Compiles for all 5 platform/architecture targets with version embedded
- Generates SHA256 checksums
- Creates a GitHub Release with auto-generated notes
- Attaches all binaries and checksum file to the release

### Outputs
- GitHub Release with:
  - `cairn-linux-amd64`
  - `cairn-linux-arm64`
  - `cairn-darwin-amd64`
  - `cairn-darwin-arm64`
  - `cairn-windows-amd64.exe`
  - `checksums.txt`

### Failure Behavior
- If any platform build fails, the entire release workflow fails
- No partial release is created
