# Data Model: Vicinae Extension for Bookmark Manager

**Feature**: 005-vicinae-extension
**Date**: 2026-03-06

---

## No Database Changes

The extension reads and writes data exclusively through the `bm` CLI. No changes to the Go data model or SQLite schema are required.

---

## CLI Output Schema (bm list --json / bm search --json)

The extension consumes the JSON array produced by `bm list --json` and `bm search <query> --json`. Each element in the array has the following fields (from the Go `Bookmark` struct):

```
Bookmark
  id            integer      Numeric database ID
  URL           string       Full URL (e.g. https://example.com/page)
  Domain        string       Stripped hostname (e.g. example.com)
  Title         string       Page title (may be empty if fetch failed)
  Description   string       Meta description (may be empty)
  CreatedAt     string       RFC3339 UTC timestamp
  Tags          []string     Normalised tags (0–3 items)
  LastVisitedAt null|string  Always null (feature removed in 004)
  IsPermanent   boolean      Pinned flag
  IsArchived    boolean      Always false (list/search returns active only)
  ArchivedAt    null|string  Always null for active bookmarks
```

**Display mapping in the extension:**

| CLI field | Extension display |
|-----------|------------------|
| `Title` (fallback to `URL`) | List item title |
| `Domain` | List item subtitle / secondary text |
| `CreatedAt` | List item accessory (formatted as YYYY-MM-DD) |
| `Tags` | List item accessory (comma-joined, prefixed with #) |
| `URL` | Used for Open in Browser action |
| `IsPermanent` | Shown as "📌" accessory badge when true |

---

## CLI Input Schema (bm add)

After the `--tags` flag is added to the CLI in this feature:

```
bm add <url> [--tags <comma-separated>]
```

| Argument | Required | Description |
|----------|----------|-------------|
| `<url>` | Yes | The URL to bookmark |
| `--tags` | No | Comma-separated tags (e.g. "go, tools") |

**Exit codes consumed by the extension:**

| Code | Meaning | Extension response |
|------|---------|--------------------|
| 0 | Saved successfully | Show "Saved" toast, close form |
| 1 | Already bookmarked | Show "Already bookmarked" error |
| 2 | Saved but title unavailable | Show "Saved (title unavailable)" toast |
| 3 | Error (invalid args / DB error) | Show error message from stderr |

---

## CLI Prerequisite Change: bm add --tags

A small change to `cmd/bm/main.go` is required as part of this feature:

**Current `runAdd` signature**: `runAdd(dbPath, rawURL string)` — tags hard-coded to `nil`.

**New flow**:
```
bm add <url> [--tags <comma-separated-tags>]
```

Implementation in `main.go`:
- In the `"add"` case, check for `--help`/`-h` as before
- Create a FlagSet for add: `fs := flag.NewFlagSet("add", flag.ContinueOnError)`
- Add `tagsFlag := fs.String("tags", "", "comma-separated tags")`
- Parse remaining args after the URL
- Pass split tags to `s.Insert()`

This is the only Go source change in this feature.
