---
title: "Installation"
weight: 10
---

# Installation

## Quick Install (Linux / macOS)

```sh
curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh
```

This downloads the latest release binary for your platform and places it in `/usr/local/bin`. Since `/usr/local/bin` is typically in your system PATH, `cairn` should be accessible immediately after installation.

### Permission Requirements

The installer may prompt for `sudo` to write to `/usr/local/bin`. If you don't have sudo access or prefer user-only installation, use:

```sh
sh install.sh --install-dir ~/.local/bin
```

Then ensure `~/.local/bin` is in your shell's `PATH`. Add this line to your shell profile (`~/.zshrc`, `~/.bashrc`, or `~/.profile`):

```sh
export PATH="$HOME/.local/bin:$PATH"
```

Then reload your shell:

```sh
source ~/.zshrc  # or ~/.bashrc, or ~/.profile
```

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

# Also install the desktop app
sh install.sh -y --with-desktop

# Desktop app without starting it at login
sh install.sh -y --with-desktop --no-autostart
```

If [Vicinae](https://vicinae.com) is installed on your system, the installer will offer to set up the Cairn extension automatically. In interactive mode it also offers to install the [desktop app]({{< relref "/docs/desktop-app" >}}), a quick-launch popup for searching, adding and opening bookmarks.

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

## Uninstall

```sh
cairn uninstall          # Remove cairn, the Vicinae extension and desktop app; keep your bookmarks
cairn uninstall --purge  # Also delete the database, cairn.json and sync.json
```

### Remove cairn and keep your bookmarks

Run `cairn uninstall`. It lists what it will remove and asks before going ahead. It removes:

- the `cairn` binary
- the Vicinae extension, if installed
- the [desktop app]({{< relref "/docs/desktop-app" >}}), if installed. A running app is stopped first.
  - **Linux:** `cairn-desktop`, plus its launcher, autostart entry, icon and version file
  - **macOS:** `~/Applications/cairn-desktop.app`, plus its login item and version file

Your bookmarks database, `cairn.json` and `sync.json` stay where they are. When it finishes, cairn lists them with an `rm` command you can run later. If you reinstall, cairn finds them again and your bookmarks are back.

### Remove everything

Run `cairn uninstall --purge` to also delete:

- the database, with its `-wal` and `-shm` files
- `cairn.json`
- `sync.json`

This can't be undone. To keep a copy of your bookmarks, export them first:

```sh
cairn list --json > bookmarks.json
```

Add `-y` to either command to skip the confirmation prompt, for example in scripts.

### If cairn is in `/usr/local/bin`

The quick installer puts cairn in `/usr/local/bin`, which usually needs `sudo` to change. In that case `cairn uninstall` changes nothing. It prints the commands to run instead and exits with code `4`:

```
cairn uninstall: permission denied: cannot remove /usr/local/bin/cairn, /usr/local/bin/cairn-desktop
Run these commands instead:
  pkill -x cairn-desktop
  sudo rm '/usr/local/bin/cairn' '/usr/local/bin/cairn-desktop'
  rm '/home/me/.local/share/cairn/desktop-version.txt' '/home/me/.local/share/applications/cairn-desktop.desktop' ...
```

Don't run `sudo cairn uninstall`. Under `sudo`, cairn may look for your files in root's home directory and miss them.

### What uninstalling doesn't touch

- **Dropbox:** the copy of your bookmarks in Dropbox (`cairn/sync.json`) is never deleted. Remove it from Dropbox yourself if you no longer need it.
- **Keyboard shortcuts:** if you added a custom shortcut for `cairn-desktop --toggle` in your desktop's keyboard settings, for example on GNOME, remove it there. cairn doesn't change desktop settings.
- **A desktop app installed elsewhere:** on Linux, only a `cairn-desktop` installed next to the `cairn` binary is removed. One elsewhere on your `PATH` is left in place.

### Uninstall manually

Use this if your cairn is too old to have `cairn uninstall`, or if you're on Windows, where a running program can't delete itself. On Windows, `cairn uninstall` prints the paths for you. Remove the binary, then whichever of these files you no longer need:

| What | Linux | macOS |
|------|-------|-------|
| Binary | `/usr/local/bin/cairn` (or your `--install-dir`) | `/usr/local/bin/cairn` (or your `--install-dir`) |
| Bookmarks | `~/.local/share/cairn/` | `~/Library/Application Support/cairn/` |
| Config and Dropbox token | `~/.config/cairn/cairn.json`, `~/.config/cairn/sync.json` | `~/Library/Application Support/cairn/cairn.json`, `~/Library/Application Support/cairn/sync.json` |
| Vicinae extension | `~/.local/share/vicinae/extensions/cairn/` | `~/Library/Application Support/vicinae/extensions/cairn/` |

To remove the desktop app by hand, see [Desktop App → Uninstalling]({{< relref "/docs/desktop-app" >}}#uninstalling).

See [`cairn uninstall`]({{< relref "/docs/cli-reference" >}}#cairn-uninstall) for all flags and exit codes.
