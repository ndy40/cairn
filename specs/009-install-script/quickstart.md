# Quickstart: Installation Script

**Feature**: 009-install-script
**Date**: 2026-03-11

## What This Feature Does

Provides a one-line installation command for Cairn. Users run a single `curl | sh` command to download and install the Cairn CLI. If Vicinae is detected, it optionally installs the browser extension too.

## Developer Setup

No special setup needed. The install script is a standalone POSIX shell script with no build step.

### Prerequisites

- A GitHub personal access token is NOT required (releases are public)
- For testing: access to Linux and/or macOS machines (or Docker)

### File Location

```
/install.sh    # The script lives at the repository root
```

### Testing Locally

```sh
# Test the script locally (it will install to ~/.local/bin by default)
sh install.sh

# Test with a custom directory
sh install.sh --install-dir /tmp/cairn-test

# Test non-interactive mode
sh install.sh --non-interactive

# Test specific version
sh install.sh --version v0.0.1

# Clean up test installation
rm /tmp/cairn-test/cairn
```

### Key Decisions

1. **POSIX sh, not bash**: Maximum portability across Linux and macOS
2. **`~/.local/bin` default**: No sudo required, XDG-compliant
3. **Checksum verification**: Always verifies SHA-256 before installing
4. **Vicinae detection via `vici` CLI**: More reliable than path scanning
5. **Pre-built extension artifact**: Avoids requiring Node.js on the user's machine

### CI/Release Changes Needed

The existing release workflow needs to be extended to also build and publish the Vicinae extension as a release artifact (`cairn-vicinae-extension.tar.gz`). This enables the install script to download the pre-built extension without requiring Node.js.
