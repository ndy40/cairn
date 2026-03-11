# Data Model: CI Release Binaries

**Feature**: 008-ci-release-binaries
**Date**: 2026-03-11

## Schema Changes

**None.** This is a CI/CD-only feature. No application data model changes.

## CI Entities

### Release (GitHub)
- **Tag**: Git tag matching `v*` pattern (e.g., `v1.2.0`)
- **Title**: Tag name
- **Body**: Auto-generated release notes from commit history
- **Assets**: 5 compiled binaries + 1 checksum file (6 total)

### Binary Asset
- **Name**: `cairn-{os}-{arch}[.exe]`
- **Platforms**: linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64
- **Build flags**: `CGO_ENABLED=0`, version injected via ldflags

### Checksum File
- **Name**: `checksums.txt`
- **Format**: Standard sha256sum output (`<hash>  <filename>`)
- **Contents**: SHA256 hash for each of the 5 binary assets
