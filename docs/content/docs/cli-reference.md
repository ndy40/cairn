---
title: "CLI Reference"
weight: 30
---

# CLI Reference

## Global Usage

```
cairn [--db <path>] [command] [flags]
```

Run `cairn` with no arguments to launch the interactive TUI.

## Global Flags

| Flag | Description |
|------|-------------|
| `--db <path>` | Override the default database path |

## Commands

### `cairn add`

Save a bookmark by URL. The page title and description are fetched automatically.

```
cairn add <url> [--tags <comma-separated>]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `<url>` | The URL to bookmark (required) |

**Flags:**

| Flag | Description |
|------|-------------|
| `--tags` | Comma-separated tags, e.g. `"work,go,tools"` — maximum 3 tags |

**Exit codes:**

| Code | Meaning |
|------|---------|
| `0` | Saved successfully |
| `1` | Already bookmarked (duplicate URL) |
| `2` | Saved but page title could not be fetched |
| `3` | Error (invalid arguments, database error) |

**Examples:**

```sh
cairn add https://go.dev
cairn add https://pkg.go.dev/net/http --tags "go,stdlib"
```

---

### `cairn list`

List all bookmarks.

```
cairn list [--json] [--order asc|desc]
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--json` | false | Output as JSON array instead of tab-separated text |
| `--order` | `desc` | Sort order: `asc` (oldest first) or `desc` (newest first) |

**Tab-separated output columns:** `ID`, `Title`, `URL`, `Domain`, `CreatedAt`

**Examples:**

```sh
cairn list
cairn list --json
cairn list --order asc
cairn list --json | jq '.[0].title'
```

---

### `cairn search`

Search bookmarks by title, domain, and description using full-text search (FTS5) combined with fuzzy matching.

```
cairn search <query> [--json] [--limit N]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `<query>` | Search query (required) |

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--json` | false | Output as JSON array |
| `--limit` | `10` | Maximum number of results to return |

**Examples:**

```sh
cairn search golang
cairn search "rust async" --limit 5
cairn search api --json | jq '.[].url'
```

---

### `cairn delete`

Delete a bookmark by its numeric ID.

```
cairn delete <id>
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `<id>` | Bookmark ID (use `cairn list` to find IDs) |

**Exit codes:**

| Code | Meaning |
|------|---------|
| `0` | Deleted successfully |
| `1` | Bookmark not found |
| `3` | Error |

**Example:**

```sh
cairn delete 42
```

---

### `cairn pin`

Toggle the pin (permanent) flag on a bookmark. Pinned bookmarks are never auto-archived.

```
cairn pin <id>
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `<id>` | Bookmark ID to pin or unpin |

**Example:**

```sh
cairn pin 42   # pins if unpinned, unpins if pinned
```

---

### `cairn sync`

Manage bookmark synchronisation via Dropbox. See [Dropbox Sync]({{< relref "/docs/sync" >}}) for setup instructions.

```
cairn sync <subcommand>
```

**Subcommands:**

| Subcommand | Description |
|------------|-------------|
| `setup` | Connect to Dropbox and perform initial sync |
| `push` | Push local changes to Dropbox |
| `pull` | Pull remote changes from Dropbox |
| `status` | Show sync configuration and pending changes |
| `auth` | Re-authenticate with Dropbox (refresh OAuth2 token) |
| `unlink` | Disconnect sync; local bookmarks are preserved |

#### `cairn sync setup`

Connects cairn to Dropbox and performs the initial sync. If `CAIRN_DROPBOX_APP_KEY` is not set via environment variable or `cairn.json`, the command prompts interactively:

```
Enter your Dropbox App Key: <key>
Enter database path (press Enter for default: /home/user/.local/share/cairn/bookmarks.db):
```

Values entered during the prompt are saved to `cairn.json` for future runs. The App Key is required for `setup` and `auth`; the database path defaults to the platform default if left blank.

You can skip the interactive prompt by supplying the App Key in advance:

```sh
export CAIRN_DROPBOX_APP_KEY=<your-app-key>
cairn sync setup
```

#### `cairn sync auth`

Re-authenticates with Dropbox. Requires `CAIRN_DROPBOX_APP_KEY` to be set via environment variable or `cairn.json` (no interactive prompt).

---

#### Auto-sync (background)

When sync is configured, most commands trigger a background sync automatically — no manual `cairn sync pull` or `cairn sync push` is needed for routine use.

| Trigger | Direction |
|---------|-----------|
| Before `add`, `list`, `search`, `delete`, `pin` | Pull (background) |
| After `add`, `edit`, `delete`, `pin` | Push (background) |

Background syncs run as a detached subprocess with no output to the terminal. A lock file prevents concurrent syncs. If a sync is already running, the new request is silently skipped.

#### First-run prompt

On the very first run (before any sync configuration exists), cairn asks whether you want to set up Dropbox sync:

```
No sync configured — connect to Dropbox? (y/N)
```

- Answering `y` prints instructions to run `cairn sync setup`.
- Any other response records a declined state so the prompt does not appear again.

---

### `cairn update`

Check for and apply updates to the cairn binary, the [desktop app]({{< relref "/docs/desktop-app" >}}#updating) or the Vicinae extension. When the desktop app is installed, `cairn update` updates it after the CLI so both stay on the same release. No network requests are ever made by other subcommands — updates are strictly opt-in.

```
cairn update [--check] [--extension | --desktop]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--check` | Report version status only; do not download or modify any files |
| `--extension` | Target the Vicinae browser extension instead of the CLI binary |
| `--desktop` | Target only the desktop app |

**Flag combinations:**

| Invocation | Behaviour |
|------------|-----------|
| `cairn update` | Check and apply CLI binary update, then desktop app update if it is installed |
| `cairn update --check` | Report CLI and desktop app version status; no changes made |
| `cairn update --desktop` | Check and apply desktop app update only |
| `cairn update --desktop --check` | Report desktop app version status; no changes made |
| `cairn update --extension` | Check and apply Vicinae extension update |
| `cairn update --extension --check` | Report extension version status; no changes made |

**Exit codes:**

| Code | Meaning |
|------|---------|
| `0` | Success (updated, already up to date, or `--check` completed) |
| `1` | Error (network failure, API error, or unknown platform) |
| `3` | Checksum verification failed |
| `4` | Permission denied (cannot write to install directory) |

**Example output:**

```sh
# Update available
$ cairn update
cairn: current version v0.1.0, latest v0.2.0
cairn: downloading cairn-linux-amd64...
cairn: verifying checksum...
cairn: updated to v0.2.0

# Already up to date
$ cairn update
cairn: already up to date (v0.2.0)

# Check only
$ cairn update --check
cairn: current version v0.1.0, latest v0.2.0 (update available)

# Extension not installed
$ cairn update --extension
cairn: extension not installed; run the install script with --with-extension to install it
```

**Notes:**

- The existing binary is backed up before replacement and restored if the update fails — your installation is never left in a broken state.
- The desktop app archive must match `checksums.txt`; on a mismatch the installed app is left untouched. A running desktop app is restarted on the new version.
- On Windows, in-process update is not supported. Re-run the install script to update.
- Requires write permission to the directory where cairn is installed (typically `~/.local/bin`). If permission is denied, exit code `4` is returned with a descriptive error.

---

### `cairn uninstall`

Remove the cairn binary and, if installed, the Vicinae extension and the [desktop app]({{< relref "/docs/desktop-app" >}}). A running desktop app is stopped first, and its launcher, login item, icon and version file are removed with it. Your bookmarks database, `cairn.json` and `sync.json` are kept unless you pass `--purge`.

```
cairn uninstall [--purge] [-y]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--purge` | Also delete the database (with its `-wal`/`-shm` files), `cairn.json` and `sync.json` |
| `-y`, `--yes` | Skip the confirmation prompt |

**Exit codes:**

| Code | Meaning |
|------|---------|
| `0` | Success (uninstalled, or cancelled at the prompt) |
| `1` | Error removing a file |
| `3` | Invalid flags |
| `4` | Permission denied (cannot remove a binary) |

**Example output:**

```sh
$ cairn uninstall
This will remove:
  /home/me/.local/bin/cairn
  /home/me/.local/share/vicinae/extensions/cairn (Vicinae extension)
  /home/me/.local/bin/cairn-desktop (desktop app; stopped first if running)
  /home/me/.local/share/cairn/desktop-version.txt
  /home/me/.local/share/cairn/cairn-desktop.png
  /home/me/.local/share/applications/cairn-desktop.desktop
  /home/me/.config/autostart/cairn-desktop.desktop

Your bookmarks and settings will be kept. Use --purge to delete them too.

Continue? [y/N] y
cairn has been uninstalled.

Your bookmarks and settings were kept:
  /home/me/.local/share/cairn/bookmarks.db
  /home/me/.config/cairn/sync.json (contains your Dropbox access token)

To delete them, run:
  rm '/home/me/.local/share/cairn/bookmarks.db' '/home/me/.config/cairn/sync.json'

The copy of your bookmarks in Dropbox has not been deleted.
```

**Notes:**

- The database path respects `--db`, `CAIRN_DB_PATH` and `db_path` in `cairn.json`, so a custom location is found too. With `--purge`, only cairn's own directories are removed once empty; the directory of a custom database path is left alone.
- On Linux, only a desktop app installed beside the `cairn` binary is removed, as the installer places it. One found elsewhere on your `PATH` is left in place.
- If you can't write to the install directory (for example `/usr/local/bin`), cairn changes nothing, prints the commands to run instead (`sudo rm` for the binaries, plain `rm` for the rest) and exits with code `4`. Avoid `sudo cairn uninstall`: under sudo, cairn may look for your data in root's home directory.
- On Windows a running program can't delete itself, so cairn prints the paths to remove manually.
- Removing cairn never deletes the copy of your bookmarks in Dropbox.

---

### `cairn version`

Print the application version.

```sh
cairn version
```

---

### `cairn config`

Print the resolved configuration (database path and whether the Dropbox app key is set).

```sh
cairn config
```

**Output example (Dropbox key configured):**

```
CAIRN_DB_PATH=/home/user/.local/share/cairn/bookmarks.db
CAIRN_DROPBOX_APP_KEY=(set)
```

**Output example (no Dropbox key):**

```
CAIRN_DB_PATH=/home/user/.local/share/cairn/bookmarks.db
```

The `CAIRN_DROPBOX_APP_KEY` line is only printed when the key is set.

---

### `cairn help`

Print the top-level help text.

```sh
cairn help
cairn --help
cairn -h
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `CAIRN_DB_PATH` | Override the default database path |
| `CAIRN_DROPBOX_APP_KEY` | Dropbox app key for sync setup and auth |
