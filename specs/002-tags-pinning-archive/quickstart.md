# Quickstart: Tags, Pinning, Archive & Startup Checks

**Feature**: 002-tags-pinning-archive
**Date**: 2026-03-06
**Extends**: `specs/001-tui-bookmark-manager/quickstart.md`

---

## Prerequisites (updated)

In addition to the requirements from feature 001:

**On Wayland** (check `echo $WAYLAND_DISPLAY`):
```bash
sudo apt install wl-clipboard    # Debian/Ubuntu
sudo pacman -S wl-clipboard      # Arch
sudo dnf install wl-clipboard    # Fedora
```

**On X11** (check `echo $DISPLAY`):
```bash
sudo apt install xclip           # Debian/Ubuntu (or xsel)
```

The app now verifies these tools on startup and provides a specific error message if missing.

---

## Build

No changes to the build process:

```bash
CGO_ENABLED=0 go build -o bm ./cmd/bm
```

---

## Running

```bash
# Launch TUI (performs prerequisite check + archive check)
./bm

# CLI subcommands (prerequisite check NOT performed)
./bm add https://example.com
./bm list
./bm search "go tutorial"
./bm delete 3
```

---

## New Features in This Release

### Tags

When adding a bookmark (`Ctrl+P`), a **Tags** field appears below the URL:
- Type up to 3 tags, comma-separated: `work, go, tools`
- Tags are saved lowercase; duplicates are removed
- Tags appear in the bookmark list next to the domain

### Tag Filtering

Press `t` from the main browse view to open the tag filter overlay:
- Navigate with `j`/`k`, toggle tags with `Space` or `Enter`
- Press `c` to clear all filters, `Esc` to close

### Permanent Bookmarks

Press `p` on any bookmark to mark it permanent (exempt from auto-archiving). Permanent bookmarks show a `[pin]` prefix. Press `p` again to remove the flag.

### Archive View

Press `a` from browse mode to view archived bookmarks. Press `r` to restore any archived bookmark back to the active list.

### Startup Archive Check

On every TUI launch, bookmarks not visited in 6+ months (183 days) are automatically moved to the archive. A footer message shows the count (e.g., `"2 bookmarks archived"`).

---

## Database Migration

Schema migration v2 runs automatically on first launch after upgrading. Existing bookmarks receive safe defaults:
- `tags = '[]'` (no tags)
- `last_visited_at = NULL` (never visited)
- `is_permanent = 0` (not permanent)
- `is_archived = 0` (active)
- `archived_at = NULL`

No manual steps required.

---

## Testing New Behaviour

### Prerequisite Check

```bash
# Simulate missing clipboard tool (unset WAYLAND_DISPLAY/DISPLAY temporarily)
WAYLAND_DISPLAY="" DISPLAY="" ./bm
# Expected: warning message, TUI still launches

# On Wayland without wl-paste (rename or remove temporarily for testing)
# Expected: error message, exit code 1
```

### Archive Check

To test archiving without waiting 183 days, temporarily insert a bookmark with an old `created_at`:

```bash
sqlite3 ~/.local/share/bookmark-manager/bookmarks.db \
  "INSERT INTO bookmarks (url, domain, title, description, created_at)
   VALUES ('https://old.example.com', 'old.example.com', 'Old Site', '', datetime('now', '-200 days'))"

# Then launch the TUI — the inserted bookmark should be archived
./bm
```

### Tag Filtering

```bash
# Add a bookmark with tags via CLI isn't available (tags are TUI-only in the add modal)
# Use sqlite3 to seed test data
sqlite3 ~/.local/share/bookmark-manager/bookmarks.db \
  "UPDATE bookmarks SET tags = '[\"work\",\"go\"]' WHERE id = 1"
```

---

## Running Tests

```bash
go test ./...
```
