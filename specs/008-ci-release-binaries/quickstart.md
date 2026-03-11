# Quickstart: CI Release Binaries

**Feature**: 008-ci-release-binaries
**Date**: 2026-03-11

## What Changed

Two GitHub Actions workflows added:
1. **build.yml** — Validates compilation on every pull request
2. **release.yml** — Creates releases with downloadable binaries when a version tag is pushed

## How to Test

### Test PR Build Validation

1. Create a branch and open a pull request
2. Verify the "Build" workflow runs and shows pass/fail status on the PR
3. Introduce a compilation error, push, and verify the build fails

### Test Release Creation

1. Create and push a version tag:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. Go to the repository's Releases page on GitHub

3. Verify:
   - A release named `v0.1.0` exists
   - 5 binary assets are attached (linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64.exe)
   - A `checksums.txt` file is attached
   - Release notes are auto-generated

4. Download a binary for your platform and verify:
   ```bash
   chmod +x cairn-linux-amd64
   ./cairn-linux-amd64 version
   # Should output: bm version v0.1.0
   ```

### Test Edge Cases

- Push a non-version tag (e.g., `test-tag`) — release workflow should NOT trigger
- Delete a tag — release should remain (manual cleanup acceptable)

## Binary Naming Convention

| Platform | Architecture | Binary Name              |
|----------|-------------|--------------------------|
| Linux    | amd64       | cairn-linux-amd64        |
| Linux    | arm64       | cairn-linux-arm64        |
| macOS    | amd64       | cairn-darwin-amd64       |
| macOS    | arm64       | cairn-darwin-arm64       |
| Windows  | amd64       | cairn-windows-amd64.exe  |
