# Contract: Vicinae Extension Commands

**Feature**: 005-vicinae-extension
**Date**: 2026-03-06

---

## Extension Manifest (package.json)

The extension registers three commands with the Vicinae platform:

```
name:        "bm-bookmarks"
title:       "Bookmark Manager"
description: "Search, browse, and save bookmarks via the bm CLI"

commands:
  - name:        "search-bookmarks"
    title:       "Search Bookmarks"
    description: "Search your saved bookmarks and open in browser"
    mode:        view

  - name:        "list-bookmarks"
    title:       "List Bookmarks"
    description: "Browse all saved bookmarks"
    mode:        view

  - name:        "add-bookmark"
    title:       "Add Bookmark"
    description: "Save a new bookmark from a URL"
    mode:        view
```

---

## Command 1: Search Bookmarks

**Entry point**: `src/search.tsx`

### UI Contract

```
State: loading  →  Show spinner/loading indicator
State: results  →  Show List of BookmarkItems
State: empty    →  Show "No bookmarks found" empty state
State: error    →  Show error with message (e.g., "bm not found")
```

### List Item Contract

Each result row displays:
```
Title:     bookmark.title || bookmark.URL
Subtitle:  bookmark.Domain
Accessories:
  - Created date (YYYY-MM-DD)
  - Tags joined as "#tag1 #tag2" (omitted if empty)
  - "📌" if IsPermanent
```

### Actions on each item

```
Primary action (Enter):  Open bookmark.URL in default browser
Secondary action:        Copy URL to clipboard
```

### Search Behaviour

```
Query empty   →  call: bm list --json
               display all active bookmarks

Query ≥ 1 char →  call: bm search <query> --json --limit 20
               display ranked results
```

---

## Command 2: List Bookmarks

**Entry point**: `src/list.tsx`

### UI Contract

```
State: loading  →  Show spinner
State: results  →  Show scrollable List of all active BookmarkItems
State: empty    →  Show "No bookmarks saved yet"
State: error    →  Show error with message
```

### List Item Contract

Same as Search Bookmarks (title, subtitle, accessories).

### Actions on each item

```
Primary action (Enter):  Open bookmark.URL in default browser
Secondary action:        Copy URL to clipboard
```

### Search Behaviour

The Vicinae launcher's built-in search box filters the rendered list client-side (no additional CLI call). The list is loaded once on mount via `bm list --json`.

### CLI Call

```
bm list --json
```

---

## Command 3: Add Bookmark

**Entry point**: `src/add.tsx`

### UI Contract

```
State: form     →  Show Form with URL and Tags fields
State: loading  →  Show spinner after submission
State: success  →  Show "Saved" toast and close
State: error    →  Show error message inline
```

### Form Fields

```
Field:       URL
Type:        text input
Required:    yes
Placeholder: "https://example.com"
Pre-fill:    clipboard content if it starts with http:// or https://

Field:       Tags
Type:        text input
Required:    no
Placeholder: "work, go, tools  (comma-separated, max 3)"
```

### Submit Action

```
Validate:  URL must be non-empty and start with http:// or https://
Call:      bm add <url> [--tags <tags>]
On exit 0: show "Saved" toast, close form
On exit 1: show inline error "Already bookmarked"
On exit 2: show "Saved (title unavailable)" toast, close form
On exit 3: show inline error with stderr message
```

---

## CLI Dependency Contract

All three commands depend on `bm` being available in `$PATH`.

**Detection check** (run once on extension load):

```
which bm  →  if exit != 0: show error panel
             "bm is not installed. Install from: https://github.com/..."
             with a link to installation instructions
```

**Minimum bm version required**: any version that supports `--json` flag on `list` and `search`, and `--tags` on `add` (added in feature 005).
