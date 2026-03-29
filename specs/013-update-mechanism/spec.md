# Feature Specification: Cairn Self-Update Mechanism

**Feature Branch**: `013-update-mechanism`
**Created**: 2026-03-26
**Status**: Draft
**Input**: User description: "I would like to implement a simple update mechanism for cairn. We need a way to update the install version of cairn and extension."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Check for and Apply CLI Update (Priority: P1)

A user who has cairn installed wants to update to the latest version without manually downloading and replacing the binary. They run a single command and cairn downloads and installs the newest release.

**Why this priority**: The CLI binary is the core tool. Keeping it up to date is the most critical and frequent user need. This is a complete, independently valuable capability on its own.

**Independent Test**: Can be tested by installing an older version of cairn, running the update command, and verifying the new version is reported by `cairn version`.

**Acceptance Scenarios**:

1. **Given** cairn is installed and a newer version is available, **When** the user runs `cairn update`, **Then** cairn downloads the new binary, replaces itself, and reports the version it updated to.
2. **Given** cairn is installed and already at the latest version, **When** the user runs `cairn update`, **Then** cairn reports that no update is needed and exits cleanly.
3. **Given** cairn is installed but cannot reach the release source (no network), **When** the user runs `cairn update`, **Then** cairn reports a clear error and leaves the existing installation untouched.
4. **Given** cairn downloads a new binary, **When** the download completes, **Then** cairn verifies the binary integrity before replacing the existing installation.

---

### User Story 2 - Check for Available Update Without Applying (Priority: P2)

A user wants to know if an update is available before deciding to apply it. They run a check command that reports version information without making any changes.

**Why this priority**: Gives users awareness and control before committing to an update. Especially useful in scripted or automated environments where silent updates are undesirable.

**Independent Test**: Can be tested by running `cairn update --check` and verifying output describes available version without modifying any files.

**Acceptance Scenarios**:

1. **Given** a newer version is available, **When** the user runs `cairn update --check`, **Then** cairn reports the current version, the available version, and exits without making changes.
2. **Given** cairn is already at the latest version, **When** the user runs `cairn update --check`, **Then** cairn reports that the current version is up to date.

---

### User Story 3 - Update the Vicinae Extension (Priority: P3)

A user who has the Vicinae browser extension installed alongside cairn wants to update it to the latest compatible version. They run a command that updates the extension package.

**Why this priority**: The extension depends on the CLI and must stay in sync. However, it is a secondary concern after the core CLI update is working, and not all users have the extension installed.

**Independent Test**: Can be tested by installing an older extension version, running the update command with the extension flag, and verifying the extension version matches the latest release.

**Acceptance Scenarios**:

1. **Given** the Vicinae extension is installed and a newer version is available, **When** the user runs `cairn update --extension`, **Then** cairn downloads and installs the updated extension package.
2. **Given** the Vicinae extension is not installed, **When** the user runs `cairn update --extension`, **Then** cairn reports that the extension is not found and provides guidance on how to install it.
3. **Given** an update is in progress, **When** the download fails partway through, **Then** the existing extension installation is preserved and an error is reported.

---

### Edge Cases

- What happens when the user lacks write permission to the cairn installation directory?
- How does the system handle a corrupted or incomplete download before replacement?
- What happens if the update is interrupted mid-way (e.g., power loss, process kill)?
- How does the system behave if the current version string cannot be determined?
- What happens when a pre-release or development build is installed — does it offer to update to stable?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Users MUST be able to trigger an update of the cairn CLI binary via a single command.
- **FR-002**: The update command MUST check the currently installed version against the latest available release before downloading anything.
- **FR-003**: The update command MUST support a `--check` flag that reports version status without applying any changes.
- **FR-004**: The update command MUST verify the integrity of a downloaded binary before replacing the installed version.
- **FR-005**: If an update fails at any point, the system MUST leave the existing installation intact and report a descriptive error message.
- **FR-006**: The update command MUST display the version it is updating from and to when an update is applied.
- **FR-007**: Users MUST be able to update the Vicinae extension via the update command using a dedicated flag (e.g., `--extension`).
- **FR-008**: If the extension is not installed, the system MUST report this clearly and not treat it as an error.
- **FR-009**: The update mechanism MUST work on all supported platforms (Linux, macOS) without requiring elevated privileges, provided the installation directory is writable by the current user.
- **FR-010**: The update mechanism MUST only perform network activity when the user explicitly runs the update command — cairn MUST NOT check for updates automatically in the background during normal usage.

### Key Entities

- **Installed Version**: The version identifier embedded in the currently running cairn binary.
- **Available Release**: The latest published release, including its version identifier and downloadable artifacts.
- **Binary Artifact**: The platform-specific executable to be downloaded and installed.
- **Integrity Checksum**: A hash provided alongside each release artifact used to verify download correctness.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can update cairn to the latest version in under 60 seconds on a standard broadband connection.
- **SC-002**: The update command completes without error in 100% of cases where the user has write access to the install location and a network connection is available.
- **SC-003**: A failed or interrupted update never results in a broken cairn installation — the previous binary is always preserved.
- **SC-004**: Users can determine whether an update is available without any files being modified on their system.
- **SC-005**: The extension update path works independently of the CLI update — each can be updated without requiring the other.

## Assumptions

- Cairn releases are published to a publicly accessible location (consistent with feature 008 — GitHub Releases).
- Each release provides platform-specific binary artifacts and associated integrity checksums.
- The Vicinae extension is distributed as a file-based package that can be replaced on disk (consistent with the existing install-script feature 009).
- The extension update only applies to users who have previously installed the extension; it is not a mandatory step.
- The update command replaces binaries in-place at the path where cairn is currently installed — it does not change the installation directory.
- Users run cairn with sufficient file-system permissions to overwrite their installation; if not, a clear permission error is shown.
