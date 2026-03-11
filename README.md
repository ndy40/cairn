# Cairn

A CLI bookmark manager with a terminal UI and a [Vicinae](https://vicinae.com) browser extension.

Cairn stores bookmarks locally in SQLite, supports full-text search, tags, pinning, archiving, and optional Dropbox sync across machines.

## Features

- Interactive TUI for browsing and managing bookmarks
- Add, search, list, and delete bookmarks from the command line
- Full-text search (SQLite FTS5)
- Tags, pinning, and archiving
- Dropbox sync with automatic push/pull
- JSON output mode for scripting
- Vicinae browser extension for searching and saving bookmarks

## Installation

### Quick install (Linux / macOS)

```sh
curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh
```

This downloads the latest release binary for your platform and places it in `~/.local/bin`.

### Install options

```sh
# Install to a custom directory
sh install.sh --install-dir /usr/local/bin

# Install a specific version
sh install.sh --version v0.0.1

# Non-interactive (CI/CD)
sh install.sh -y

# Non-interactive with Vicinae extension
sh install.sh -y --with-extension
```

If Vicinae is installed on your system, the installer will offer to set up the Cairn extension automatically.

### Build from source

Requires Go 1.22+.

```sh
git clone https://github.com/ndy40/cairn.git
cd cairn
go build -o cairn ./cmd/cairn/
```

## Usage

Launch the interactive TUI:

```sh
cairn
```

### Commands

```
cairn add <url> [title]   Save a bookmark
cairn list                List all bookmarks
cairn search <query>      Search bookmarks
cairn delete <id>         Delete a bookmark by ID
cairn sync <subcommand>   Manage Dropbox sync
cairn version             Print version
cairn help                Show help
```

### Sync

Cairn can sync bookmarks to Dropbox. Set the `CAIRN_DROPBOX_APP_KEY` environment variable and run:

```sh
cairn sync setup
```

Other sync commands: `push`, `pull`, `status`, `auth`, `unlink`.

## Vicinae Extension

The `vicinae-extension/` directory contains a Vicinae extension that provides:

- Search bookmarks from the launcher
- Browse and list all bookmarks
- Save the current page as a bookmark

See [vicinae-extension/README.md](vicinae-extension/README.md) for details.

## License

MIT
