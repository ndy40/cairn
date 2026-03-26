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

**Required environment variable for `setup` and `auth`:**

```sh
CAIRN_DROPBOX_APP_KEY=<your-app-key>
```

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

**Output example:**

```
CAIRN_DB_PATH=/home/user/.local/share/cairn/bookmarks.db
CAIRN_DROPBOX_APP_KEY=(set)
```

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
