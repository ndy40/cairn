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

# Non-interactive with Vicinae extension (curl installer)
curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh -s -- -y --with-extension
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

#### Enable Dropbox integration

1. Create or obtain a Dropbox app key.
2. Provide the app key to Cairn:
   - Preferred: add `"dropbox_app_key"` to your `cairn.json`
   - Alternative (quick test): export `CAIRN_DROPBOX_APP_KEY`
3. Run `cairn sync setup` to authenticate, then use `cairn sync push|pull|status`.

## Configuration

Cairn can be configured via an optional JSON config file, environment variables, and CLI flags.

### Precedence (highest to lowest)

1. **Environment variables** — always win
2. **CLI flags** (e.g., `--db`)
3. **Config file** (`cairn.json`)
4. **Defaults** (OS-appropriate paths)

### Config file

Create a `cairn.json` file in your OS config directory:

| OS      | Path                                                      |
|---------|-----------------------------------------------------------|
| Linux   | `$XDG_CONFIG_HOME/cairn/cairn.json` (default: `~/.config/cairn/cairn.json`) |
| macOS   | `~/Library/Application Support/cairn/cairn.json`          |
| Windows | `%APPDATA%/cairn/cairn.json`                              |

All keys are optional. Example:

```json
{
  "db_path": "/path/to/bookmarks.db",
  "dropbox_app_key": "your-app-key"
}
```

### Supported settings

| JSON key          | Environment variable      | CLI flag | Description                          |
|-------------------|---------------------------|----------|--------------------------------------|
| `db_path`         | `CAIRN_DB_PATH`           | `--db`   | Path to the SQLite bookmark database |
| `dropbox_app_key` | `CAIRN_DROPBOX_APP_KEY`   | —        | Dropbox app key for sync             |

Run `cairn config` to see the resolved configuration values.

## Security

### Credentials Storage

**Dropbox App Key** (`dropbox_app_key`):
- Store in `cairn.json` (config directory), NOT as an environment variable
- **Do not** use `CAIRN_DROPBOX_APP_KEY` environment variable in production — it's visible in process listings (`ps aux`), logs, and CI/CD output
- The app key is a shared secret for the entire application; keep it secure

**OAuth2 Tokens** (`sync.json`):
- Sync credentials are stored in `~/.config/cairn/sync.json` (Linux), `~/Library/Application Support/cairn/sync.json` (macOS), or `%APPDATA%/cairn/sync.json` (Windows)
- **Tokens are stored in plaintext** — ensure the config directory is protected with restrictive file permissions (owner-only)
- File is created with `0600` permissions (owner-only) automatically
- **If `sync.json` is compromised, anyone with access can:**
  - Read and write your Dropbox bookmarks
  - Potentially access other files in your Dropbox account through the app's granted permissions
  - Impersonate your Dropbox session

### Best Practices

1. **Restrict config file access** — Never share or commit `cairn.json` or `sync.json` to version control
2. **Don't use env vars for secrets in production** — Environment variables are visible in process listings and logs
3. **Protect your machine** — Since tokens are stored in plaintext, ensure your system is secured with disk encryption and strong authentication
4. **Unlink sync if you lose access** — Run `cairn sync unlink` to revoke access; local bookmarks are preserved

For responsible security disclosure, see [SECURITY.md](SECURITY.md).

## Vicinae Extension

The `vicinae-extension/` directory contains a Vicinae extension that provides:

- Search bookmarks from the launcher
- Browse and list all bookmarks
- Save the current page as a bookmark

See [vicinae-extension/README.md](vicinae-extension/README.md) for details.

## License

MIT
