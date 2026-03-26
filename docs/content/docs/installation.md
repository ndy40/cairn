---
title: "Installation"
weight: 10
---

# Installation

## Quick Install (Linux / macOS)

```sh
curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh
```

This downloads the latest release binary for your platform and places it in `~/.local/bin`.

## Install Options

```sh
# Install to a custom directory
sh install.sh --install-dir /usr/local/bin

# Install a specific version
sh install.sh --version v0.0.1

# Non-interactive (CI/CD environments)
sh install.sh -y

# Non-interactive with Vicinae extension
sh install.sh -y --with-extension

# One-liner with Vicinae extension
curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh -s -- -y --with-extension
```

If [Vicinae](https://vicinae.com) is installed on your system, the installer will offer to set up the Cairn extension automatically.

## Build from Source

Requires Go 1.22+. The binary is built with `CGO_ENABLED=0` for a fully static, dependency-free executable.

```sh
git clone https://github.com/ndy40/cairn.git
cd cairn
go build -o cairn ./cmd/cairn/
```

To cross-compile for another platform:

```sh
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o cairn-linux-arm64 ./cmd/cairn/
```

## Pre-built Binaries

Pre-built binaries for Linux and macOS (amd64 and arm64) are available on the [GitHub Releases](https://github.com/ndy40/cairn/releases) page. Each release includes SHA256 checksums.

Supported platforms:

| OS      | Architectures          |
|---------|------------------------|
| Linux   | amd64, arm64           |
| macOS   | amd64, arm64 (Apple Silicon) |
| Windows | amd64                  |

## Verify Installation

```sh
cairn version
```

## Database Location

Cairn creates its SQLite database automatically on first run:

| OS      | Default path                                                        |
|---------|---------------------------------------------------------------------|
| Linux   | `$XDG_DATA_HOME/cairn/bookmarks.db` (default: `~/.local/share/cairn/bookmarks.db`) |
| macOS   | `~/Library/Application Support/cairn/bookmarks.db`                  |
| Windows | `%APPDATA%\cairn\bookmarks.db`                                       |

Override the path with `--db <path>`, the `CAIRN_DB_PATH` environment variable, or the `db_path` key in `cairn.json`.
