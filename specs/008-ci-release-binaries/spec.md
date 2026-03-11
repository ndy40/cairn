# Feature Specification: CI Release Binaries

**Feature Branch**: `008-ci-release-binaries`
**Created**: 2026-03-11
**Status**: Draft
**Input**: User description: "Implement GitHub workflow actions that compile and generate Go binaries for Linux, macOS, and Windows, and make them available for download. Support release tags."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Tag-triggered release with downloadable binaries (Priority: P1)

A maintainer pushes a version tag (e.g., `v1.2.0`) to the repository. The CI system automatically compiles the application for Linux, macOS, and Windows, creates a release on the repository's release page, and attaches the compiled binaries as downloadable assets. Users can then visit the release page and download the binary for their operating system.

**Why this priority**: This is the core value proposition — without automated release builds, users cannot easily install the application. Tag-triggered releases are the standard distribution mechanism for CLI tools.

**Independent Test**: Can be tested by pushing a tag to the repository and verifying that a release is created with downloadable binaries for all three platforms.

**Acceptance Scenarios**:

1. **Given** a maintainer pushes a tag matching the pattern `v*` (e.g., `v1.0.0`), **When** the CI pipeline runs, **Then** binaries are compiled for Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64), and a release is created on the repository's release page with all binaries attached.
2. **Given** a release is created, **When** a user visits the release page, **Then** they can download the binary for their operating system and architecture, and the binary is immediately executable (no additional installation steps required).
3. **Given** a tag is pushed, **When** the CI pipeline compiles the binaries, **Then** the version string embedded in the binary matches the tag name (e.g., running `cairn version` outputs `v1.2.0`).

---

### User Story 2 - Build validation on pull requests (Priority: P2)

A contributor opens a pull request. The CI system automatically compiles the application for all target platforms to verify the code compiles successfully. This provides early feedback before merging, without creating a release.

**Why this priority**: Catching compilation failures before merge prevents broken releases. This is a standard CI best practice that complements the release workflow.

**Independent Test**: Can be tested by opening a pull request and verifying the CI pipeline runs and reports build success or failure.

**Acceptance Scenarios**:

1. **Given** a contributor opens or updates a pull request, **When** the CI pipeline runs, **Then** the application is compiled for all target platforms and the build status is reported on the pull request.
2. **Given** the code has a compilation error, **When** the CI pipeline runs on a pull request, **Then** the build fails and the contributor sees the error in the CI output.
3. **Given** a pull request build succeeds, **When** the contributor views the pull request, **Then** no release is created and no binaries are published.

---

### Edge Cases

- What happens if a tag is pushed that does not match the version pattern (e.g., `test-tag`)? The release workflow should not trigger.
- What happens if a tag is deleted? The associated release should remain (manual cleanup is acceptable).
- What happens if the same tag is pushed twice? The workflow should not create duplicate releases; it should update or skip if a release already exists.
- What happens if compilation fails for one platform but succeeds for others? The entire release should fail — partial releases with missing platforms are confusing for users.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The CI system MUST compile the application for Linux (amd64 and arm64), macOS (amd64 and arm64), and Windows (amd64) — producing 5 binaries total.
- **FR-002**: The CI system MUST create a release on the repository's release page when a version tag (matching `v*`) is pushed.
- **FR-003**: The CI system MUST attach all compiled binaries to the release as downloadable assets with clear platform/architecture naming (e.g., `cairn-linux-amd64`, `cairn-darwin-arm64`, `cairn-windows-amd64.exe`).
- **FR-004**: The CI system MUST embed the tag version into the compiled binary so that `cairn version` outputs the correct version string.
- **FR-005**: The CI system MUST compile the application on every pull request to validate build correctness, without creating a release.
- **FR-006**: The CI system MUST fail the entire release if any platform compilation fails, preventing partial releases.
- **FR-007**: The CI system MUST generate SHA256 checksums for all release binaries and include them in the release.
- **FR-008**: The release MUST only trigger on tags matching the `v*` pattern (e.g., `v1.0.0`, `v2.3.1-beta`), not on arbitrary tags.

### Key Entities

- **Release**: A versioned distribution point containing compiled binaries for all supported platforms, associated with a git tag.
- **Binary Asset**: A platform-specific compiled executable, named with platform and architecture identifiers, attached to a release.
- **Checksum File**: A file containing SHA256 hashes for all binary assets, enabling users to verify download integrity.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of version tags pushed to the repository result in a release with downloadable binaries for all 5 platform/architecture combinations within 10 minutes.
- **SC-002**: Users can download and run the binary for their platform without any additional installation steps — the binary is self-contained and immediately executable.
- **SC-003**: Every pull request receives build validation feedback (pass/fail) within 5 minutes of being opened or updated.
- **SC-004**: The version output of every released binary exactly matches the git tag that triggered the release.
- **SC-005**: Every release includes a checksum file that users can use to verify download integrity.

## Assumptions

- The project uses semantic versioning tags (e.g., `v1.0.0`). Pre-release tags like `v1.0.0-beta` are also supported.
- The binary is statically compiled with `CGO_ENABLED=0` (consistent with the project's zero-CGO constraint), so no dynamic library dependencies exist on any platform.
- Release notes are auto-generated from commit history between tags. Custom release notes can be edited manually after creation.
- The CI system has permission to create releases and upload assets to the repository.
- ARM64 support for Windows is not included (low demand for CLI tools on Windows ARM).
