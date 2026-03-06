# Contract: CLI Interface

**Feature**: 001-tui-bookmark-manager
**Date**: 2026-03-06

This document defines the command-line interface contract for the bookmark manager binary.

---

## Binary Name

```
bm
```

---

## Usage

```
bm [command] [flags]
```

Running `bm` with no arguments launches the interactive TUI in browse mode.

---

## Commands

### `bm` (no arguments)

Launches the interactive TUI application.

**Behavior**:
- Reads all bookmarks from the local database and displays the list.
- If the database does not exist, creates it and displays an empty state.

---

### `bm add <url>`

Non-interactive bookmark creation. Fetches the URL, extracts metadata, and saves the bookmark without launching the TUI.

**Arguments**:
- `<url>` — Required. The URL to bookmark. Must be a valid `http://` or `https://` URL.

**Exit codes**:
- `0` — Bookmark saved successfully.
- `1` — URL already exists (duplicate).
- `2` — URL is invalid or unreachable (bookmark still saved with fallback title).
- `3` — Unexpected error (database, filesystem).

**Output** (stdout):
```
Saved: "Page Title" (example.com)
```

---

### `bm list`

Prints all bookmarks to stdout in plain text, one per line. Suitable for piping.

**Output format** (tab-separated):
```
<id>\t<title>\t<url>\t<domain>\t<created_at>
```

**Flags**:
- `--json` — Output as a JSON array instead of tab-separated text.

---

### `bm search <query>`

Non-interactive fuzzy search. Prints matching bookmarks to stdout.

**Arguments**:
- `<query>` — Required. The search term.

**Output format**: Same as `bm list`.

**Flags**:
- `--json` — Output as a JSON array.
- `--limit N` — Return at most N results (default: 10).

---

### `bm delete <id>`

Non-interactive deletion by bookmark ID. Does not prompt for confirmation.

**Arguments**:
- `<id>` — Required. The integer ID of the bookmark (visible in `bm list` output).

**Exit codes**:
- `0` — Deleted successfully.
- `1` — ID not found.

---

### `bm version`

Prints the application version and exits.

```
bm version 0.1.0
```

---

### `bm help`

Prints usage summary and exits.

---

## Global Flags

| Flag | Description |
|------|-------------|
| `--db <path>` | Override the default database file path |
| `--no-color` | Disable ANSI color output (for `list` and `search` commands) |

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `BM_DB_PATH` | Overrides the default database location (same as `--db`) |
| `NO_COLOR` | When set (any value), disables color output (respects the `no-color.org` standard) |

---

## Exit Code Summary

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Expected error (duplicate, not found) |
| `2` | Network or fetch error |
| `3` | Unexpected internal error |
