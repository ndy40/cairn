# Research: CI Release Binaries

**Feature**: 008-ci-release-binaries
**Date**: 2026-03-11

## Research Task 1: GitHub Actions for Go cross-compilation

### Decision: Use `GOOS` and `GOARCH` environment variables with `go build`

### Rationale

Go natively supports cross-compilation via `GOOS` and `GOARCH` environment variables. Combined with `CGO_ENABLED=0`, this produces fully static binaries for any supported platform from a single Linux runner. No need for platform-specific runners (macOS/Windows runners are slower and more expensive).

### Alternatives Considered

- **Platform-specific runners**: Running each build on its native OS (macOS for darwin, Windows for windows). More expensive, slower, and unnecessary since Go cross-compiles cleanly with CGO disabled.
- **GoReleaser**: Full-featured release tool. Adds a dependency and configuration complexity. Overkill for a simple 5-target build matrix.
- **Makefile-driven builds**: Would work but duplicates logic that belongs in the workflow definition.

## Research Task 2: GitHub release creation action

### Decision: Use `softprops/action-gh-release`

### Rationale

`softprops/action-gh-release` is the most widely used GitHub Action for creating releases. It supports:
- Automatic release creation from tags
- File glob patterns for attaching assets
- Auto-generated release notes
- Idempotent behavior (won't create duplicates)

### Alternatives Considered

- **`gh release create` via CLI**: Requires manual scripting for asset upload. More verbose but no external action dependency.
- **`actions/create-release` (official)**: Deprecated by GitHub in favor of community actions.
- **`ncipollo/release-action`**: Similar capabilities but less widely adopted.

## Research Task 3: Version embedding via ldflags

### Decision: Use `-ldflags "-X main.version=$TAG"` to embed version at compile time

### Rationale

The cairn binary already has `var version = "dev"` in `cmd/cairn/main.go`. Using Go's `-ldflags -X` linker flag replaces this value at compile time with the git tag. This is the standard Go pattern for version injection — no code changes needed.

The tag name is available in GitHub Actions via `${{ github.ref_name }}` when triggered by a tag push.

## Research Task 4: SHA256 checksum generation

### Decision: Use `sha256sum` on Linux runner to generate checksums

### Rationale

`sha256sum` is available on all GitHub-hosted Ubuntu runners. Generate a single `checksums.txt` file containing SHA256 hashes for all binaries, then attach it to the release alongside the binaries.

Format: `<hash>  <filename>` (standard sha256sum output).

## Research Task 5: Build matrix strategy

### Decision: Use a matrix strategy with GOOS/GOARCH pairs

### Rationale

GitHub Actions matrix strategy allows defining all 5 target combinations and running them in parallel jobs. This is cleaner than a single job with a loop and provides per-target failure visibility.

Matrix entries:
| GOOS    | GOARCH | Binary Name              |
|---------|--------|--------------------------|
| linux   | amd64  | cairn-linux-amd64        |
| linux   | arm64  | cairn-linux-arm64        |
| darwin  | amd64  | cairn-darwin-amd64       |
| darwin  | arm64  | cairn-darwin-arm64       |
| windows | amd64  | cairn-windows-amd64.exe  |
