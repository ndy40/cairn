# Research: Installation Script

**Feature**: 009-install-script
**Date**: 2026-03-11

## R1: Installation Directory Strategy

**Decision**: Use `~/.local/bin` as the default install directory on both Linux and macOS, with an option to specify a custom directory via `--install-dir`.

**Rationale**:
- `~/.local/bin` does not require sudo/root privileges
- It follows the XDG Base Directory Specification on Linux
- Modern macOS versions include `~/.local/bin` in default PATH (or it can be added easily)
- Avoids permission issues that plague `/usr/local/bin` on macOS with SIP
- Many popular installers (rustup, deno, bun) use a user-local directory

**Alternatives considered**:
- `/usr/local/bin`: Requires sudo on macOS, may conflict with Homebrew. Rejected for permission complexity.
- `~/bin`: Less standard, not always on PATH. Rejected.

## R2: Download and Checksum Verification

**Decision**: Use the GitHub Releases API to find the latest release, download the correct binary, and verify against the `checksums.txt` file already published by the CI workflow.

**Rationale**:
- The existing `release.yml` workflow already generates `checksums.txt` with SHA-256 sums for all binaries
- Binary names follow a predictable pattern: `cairn-{os}-{arch}` (e.g., `cairn-linux-amd64`, `cairn-darwin-arm64`)
- GitHub Releases API provides a stable endpoint for fetching latest release metadata

**Alternatives considered**:
- GPG signatures: More secure but requires users to have GPG installed and import a key. Overkill for current project size. Can be added later.
- No verification: Unacceptable for security.

## R3: Vicinae Extension Detection

**Decision**: Detect Vicinae by checking for the `vici` CLI command on PATH (using `command -v vici`).

**Rationale**:
- Vicinae provides a `vici` CLI tool that is installed alongside the browser
- Checking for the CLI is more reliable than scanning filesystem paths that may vary across versions
- If `vici` is available, the extension can potentially be installed via `vici install` or by copying to the extensions directory

**Alternatives considered**:
- Check for application bundles (`/Applications/Vicinae.app` on macOS, desktop files on Linux): Fragile, varies by install method.
- Check for config directories: May exist without the app being functional.

## R4: Vicinae Extension Installation Method

**Decision**: Build the extension from source in the repo's `vicinae-extension/` directory using `npm run build`, then use standard Vicinae extension installation conventions.

**Rationale**:
- The extension source lives in the project repo at `vicinae-extension/`
- Building requires Node.js and the `vici` CLI
- However, requiring Node.js for the install script is a heavy dependency

**Revised Decision**: Pre-build the extension and publish it as a release artifact alongside the CLI binaries. The install script downloads the pre-built extension archive and extracts it to the Vicinae extensions directory.

**Alternatives considered**:
- Build from source during install: Requires Node.js and npm, too heavy for an install script.
- `vici install <extension-name>` from a registry: Only viable if the extension is published to a Vicinae registry (not currently the case).

## R5: Non-Interactive Mode

**Decision**: Support `--non-interactive` flag (or `-y` for short) that skips all prompts. Extension installation in non-interactive mode requires explicit `--with-extension` flag.

**Rationale**:
- Non-interactive mode is essential for CI/CD and scripted provisioning
- Defaulting to CLI-only in non-interactive mode is safest (principle of least surprise)
- Explicit `--with-extension` makes the intent clear in automation scripts

**Alternatives considered**:
- Environment variable `CAIRN_NONINTERACTIVE=1`: Less discoverable but could be supported as secondary option.
- Auto-detect TTY: Good enhancement but not sufficient as the sole mechanism.

## R6: Upgrade Strategy

**Decision**: The install script overwrites the existing binary. No version comparison is performed; it always installs the latest release. If a specific version is desired, support `--version vX.Y.Z` flag.

**Rationale**:
- Simple and predictable behavior
- Users who want to pin versions can use `--version`
- No need to maintain a version database or registry

**Alternatives considered**:
- Skip if already at latest: Requires version comparison logic, adds complexity for little benefit.
- Side-by-side versions: Overkill for a single CLI tool.

## R7: Shell Compatibility

**Decision**: Write the script targeting POSIX sh, avoiding bash-isms.

**Rationale**:
- macOS ships with zsh by default, not bash (since Catalina)
- Some minimal Linux containers may not have bash
- POSIX sh is the lowest common denominator guaranteed on all target platforms
- `#!/bin/sh` ensures broadest compatibility

**Alternatives considered**:
- Bash-only: Would exclude some environments unnecessarily.
- Zsh-only: Would exclude Linux servers.
