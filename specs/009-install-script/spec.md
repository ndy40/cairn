# Feature Specification: Installation Script

**Feature Branch**: `009-install-script`
**Created**: 2026-03-11
**Status**: Draft
**Input**: User description: "I would like to have an installation script that will install the cli tool and optionally detect if vicinae is installed and install the extension."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Install the CLI Tool (Priority: P1)

A user wants to install Cairn on their machine with a single command. They run the installation script, which detects their operating system and architecture, downloads the correct pre-built binary from the latest release, and places it in a standard location on their system path. After installation completes, the user can immediately run `cairn` from their terminal.

**Why this priority**: This is the core value of the feature. Without CLI installation, nothing else matters. Every user needs this regardless of whether they use the browser extension.

**Independent Test**: Can be fully tested by running the install script on a clean machine and verifying `cairn --help` works afterward.

**Acceptance Scenarios**:

1. **Given** a Linux or macOS machine without Cairn installed, **When** the user runs the install script, **Then** the correct binary for their OS and architecture is downloaded and placed on the system path.
2. **Given** a machine where Cairn is already installed, **When** the user runs the install script, **Then** the existing installation is upgraded to the latest version.
3. **Given** a machine with no internet connectivity, **When** the user runs the install script, **Then** a clear error message is displayed explaining the failure.
4. **Given** a successful installation, **When** the user runs `cairn --help`, **Then** the CLI responds with usage information.

---

### User Story 2 - Automatic Vicinae Detection and Extension Installation (Priority: P2)

After installing the CLI, the script checks whether the Vicinae browser is installed on the user's system. If Vicinae is detected, the script offers to install the Cairn browser extension automatically. The user can accept or decline. If Vicinae is not detected, the script skips this step silently or with a brief informational message.

**Why this priority**: This adds convenience for users who use Vicinae, but the CLI is fully functional without it. This is an enhancement that reduces friction for a subset of users.

**Independent Test**: Can be tested by running the install script on a machine with Vicinae installed and verifying the extension appears in the browser's extension list.

**Acceptance Scenarios**:

1. **Given** a machine with Vicinae installed, **When** the CLI installation completes, **Then** the script detects Vicinae and prompts the user to install the extension.
2. **Given** the user accepts the extension install prompt, **When** the extension is installed, **Then** the Cairn extension is available in Vicinae.
3. **Given** the user declines the extension install prompt, **When** the script finishes, **Then** only the CLI is installed and no extension changes are made.
4. **Given** a machine without Vicinae installed, **When** the CLI installation completes, **Then** the script skips extension installation without errors.

---

### User Story 3 - Unattended Installation (Priority: P3)

A user or automated system wants to install Cairn without interactive prompts. They pass a flag or environment variable to the install script to skip all prompts, installing just the CLI (and optionally the extension via a separate flag). This supports scripted provisioning and CI environments.

**Why this priority**: Important for power users and automation but not required for the primary use case. Most users will use the interactive flow.

**Independent Test**: Can be tested by running the install script with the non-interactive flag in a CI-like environment and verifying the CLI is installed without any prompts.

**Acceptance Scenarios**:

1. **Given** the install script is run with a non-interactive flag, **When** installation completes, **Then** the CLI is installed without any user prompts.
2. **Given** the install script is run with flags for both CLI and extension in non-interactive mode, **When** Vicinae is present, **Then** both CLI and extension are installed without prompts.

---

### Edge Cases

- What happens when the user's OS/architecture combination is not supported (e.g., FreeBSD, 32-bit ARM)?
- What happens when the user does not have write permissions to the installation directory?
- What happens when a download is interrupted or the downloaded file is corrupted?
- What happens when the GitHub release page is unreachable or rate-limited?
- What happens when the Vicinae installation is in a non-standard location?
- What happens when a newer version of the extension is incompatible with the installed Vicinae version?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The install script MUST detect the user's operating system (Linux, macOS, Windows) and processor architecture (amd64, arm64).
- **FR-002**: The install script MUST download the correct pre-built binary for the detected OS and architecture from the project's published releases.
- **FR-003**: The install script MUST place the binary in a directory on the user's system path so it is immediately usable.
- **FR-004**: The install script MUST verify the integrity of the downloaded binary using published checksums before installation.
- **FR-005**: The install script MUST detect whether Vicinae browser is installed on the system.
- **FR-006**: The install script MUST prompt the user to install the Cairn browser extension if Vicinae is detected, and respect the user's choice.
- **FR-007**: The install script MUST support a non-interactive mode that skips all prompts and installs the CLI only by default.
- **FR-008**: The install script MUST display clear, actionable error messages when installation fails for any reason (unsupported platform, network failure, permission denied).
- **FR-009**: The install script MUST support upgrading an existing Cairn installation to the latest version.
- **FR-010**: The install script MUST support Linux and macOS only. Windows support is out of scope and will be addressed in a follow-up feature.

## Assumptions

- Pre-built binaries are already published as GitHub releases with checksums (this exists per the CI/release workflow).
- The Vicinae browser stores extensions in a well-known, discoverable directory on each supported OS.
- The install script will be hosted in the project repository and downloadable via a short URL (e.g., `curl -sSL https://... | sh`).
- Standard installation directories will be used: `~/.local/bin` on Linux, `/usr/local/bin` on macOS (with fallback to `~/.local/bin` if the user lacks sudo).
- The Vicinae extension is distributed as a directory that can be copied into the browser's extension folder.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A new user can install Cairn and run their first command within 60 seconds on a supported platform.
- **SC-002**: 95% of installation attempts on supported platforms succeed without manual intervention.
- **SC-003**: Users with Vicinae installed are prompted about the extension, and can complete extension installation within 30 seconds of accepting.
- **SC-004**: The install script provides a clear, human-readable error message for 100% of known failure modes (unsupported platform, network error, permission denied, corrupted download).
- **SC-005**: Unattended installations complete without any interactive prompts when the appropriate flag is provided.
