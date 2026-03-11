# Data Model: Installation Script

**Feature**: 009-install-script
**Date**: 2026-03-11

## Overview

This feature has no persistent data model. The installation script operates on files and directories without maintaining any state. This document captures the key entities involved in the installation process.

## Entities

### Release Artifact

Represents a downloadable file from a GitHub Release.

- **Name**: Binary filename (e.g., `cairn-linux-amd64`, `cairn-darwin-arm64`)
- **Version**: Semantic version tag (e.g., `v0.0.1`)
- **Download URL**: GitHub Release download URL
- **Checksum**: SHA-256 hash from `checksums.txt`

### Installation Target

Represents where the binary is placed on the user's system.

- **Directory**: Default `~/.local/bin`, overridable via `--install-dir`
- **Binary name**: `cairn` (renamed from the architecture-specific filename)
- **Permissions**: Executable (`chmod +x`)

### Vicinae Extension

Represents the browser extension artifact.

- **Source**: Pre-built extension archive from GitHub Release (future: `cairn-vicinae-extension.tar.gz`)
- **Target directory**: Vicinae extensions directory (OS-dependent, discovered via `vici` CLI)

## State Transitions

```
Not Installed → Installing → Installed
                    ↓
                  Failed (with error message)

Installed → Upgrading → Installed (newer version)
                ↓
              Failed (rollback: previous binary preserved)
```
