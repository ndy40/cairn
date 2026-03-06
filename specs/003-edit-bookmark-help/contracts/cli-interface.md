# Contract: CLI Interface (Updated for Feature 003)

**Feature**: 003-edit-bookmark-help
**Date**: 2026-03-06
**Extends**: `specs/001-tui-bookmark-manager/contracts/cli-interface.md`

---

## Changes

All subcommands now respond to `-h` / `--help` flags, printing usage text and exiting with code 0.

---

## Root Command

```
bm [--db <path>] [--help|-h] [subcommand]
```

`bm --help` or `bm -h` prints:

```
bm - terminal bookmark manager

Usage:
  bm                    Launch interactive TUI
  bm add <url>          Save a bookmark non-interactively
  bm list [--json]      List all bookmarks
  bm search <query> [--json] [--limit N]  Search bookmarks
  bm delete <id>        Delete a bookmark by ID
  bm version            Print version
  bm help               Show this help

Flags:
  --db <path>           Override default database path

Environment:
  BM_DB_PATH            Override default database path

Exit code: 0
```

---

## Subcommand Help

### `bm add --help`

```
Usage: bm add <url>

Save a bookmark by URL. The page title and description are fetched automatically.

Arguments:
  <url>    The URL to bookmark (required)

Exit codes:
  0  Saved successfully
  1  Already bookmarked (duplicate URL)
  2  Saved but title could not be fetched
  3  Error (invalid arguments, database error)
```

### `bm list --help`

```
Usage: bm list [--json]

List all bookmarks ordered by date added (newest first).

Flags:
  --json    Output as JSON array instead of tab-separated text
```

### `bm search --help`

```
Usage: bm search <query> [--json] [--limit N]

Search bookmarks by title, domain, and description.

Arguments:
  <query>  Search query (required)

Flags:
  --json       Output as JSON array
  --limit N    Maximum number of results to return (default: 10)
```

### `bm delete --help`

```
Usage: bm delete <id>

Delete a bookmark by its numeric ID.

Arguments:
  <id>    Bookmark ID (required, use bm list to find IDs)

Exit codes:
  0  Deleted successfully
  1  Bookmark not found
  3  Error
```

### `bm version --help`

```
Usage: bm version

Print the application version and exit.
```

### `bm help --help`

```
Usage: bm help

Print the full usage guide and exit.
```

---

## Exit Codes (all help flags)

All `-h`/`--help` responses exit with code **0**.
