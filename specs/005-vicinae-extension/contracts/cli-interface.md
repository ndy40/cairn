# Contract: CLI Interface (Updated for Feature 005)

**Feature**: 005-vicinae-extension
**Date**: 2026-03-06
**Extends**: `specs/003-edit-bookmark-help/contracts/cli-interface.md`

---

## Changes

`bm add` gains an optional `--tags` flag.

---

## Updated `bm add` Signature

```
bm add <url> [--tags <tags>] [--help|-h]
```

### `bm add --help`

```
Usage: bm add <url> [--tags <comma-separated>]

Save a bookmark by URL. The page title and description are fetched automatically.

Arguments:
  <url>    The URL to bookmark (required)

Flags:
  --tags   Comma-separated tags (e.g. "work, go, tools") — max 3 tags

Exit codes:
  0  Saved successfully
  1  Already bookmarked (duplicate URL)
  2  Saved but title could not be fetched
  3  Error (invalid arguments, database error)
```

### Behaviour

- `--tags` value is split on commas, normalised (lowercase, dedup, truncate to 32 chars, max 3), and stored with the bookmark.
- If `--tags` is omitted, tags default to empty (unchanged from previous behaviour).
- All other `bm add` exit codes and output messages remain unchanged.

---

## All Other Subcommands

Unchanged from feature 003. See `specs/003-edit-bookmark-help/contracts/cli-interface.md`.
