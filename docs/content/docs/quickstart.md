---
title: "Quickstart"
weight: 20
---

# Quickstart

## 1. Install Cairn

```sh
curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh
```

## 2. Launch the TUI

Run `cairn` with no arguments to open the interactive terminal UI:

```sh
cairn
```

### TUI Keybindings

| Key | Action |
|-----|--------|
| `↑` / `↓` or `j` / `k` | Navigate the bookmark list |
| `Enter` | Open selected bookmark in browser |
| `/` | Focus search bar |
| `Ctrl+P` | Add a bookmark from the clipboard |
| `e` | Edit tags on selected bookmark |
| `d` or `Delete` | Delete selected bookmark |
| `p` | Toggle pin on selected bookmark |
| `t` | Open the tag filter overlay |
| `a` | Open the archive view |
| `F1` or `?` | Toggle the help screen |
| `q` or `Ctrl+C` | Quit |

## 3. Add a Bookmark (non-interactive)

```sh
cairn add https://go.dev
```

Add with tags (up to 3):

```sh
cairn add https://pkg.go.dev --tags "go,reference,docs"
```

Cairn automatically fetches the page title and description.

## 4. Search

```sh
cairn search golang
```

Combine with JSON output for scripting:

```sh
cairn search "rust" --json | jq '.[].url'
```

## 5. List All Bookmarks

```sh
cairn list

# Newest first (default)
cairn list --order desc

# Oldest first
cairn list --order asc

# As JSON
cairn list --json
```

## 6. Delete a Bookmark

```sh
# Find the ID first
cairn list

# Then delete by ID
cairn delete 42
```

## 7. Pin a Bookmark

Pinned bookmarks are marked as permanent and won't be auto-archived.

```sh
cairn pin 42   # toggles pin on/off
```

## 8. Review and Restore Archived Bookmarks

Cairn keeps your active list tidy by **auto-archiving** stale bookmarks: on
TUI startup, any bookmark older than 30 days that isn't pinned is moved to the
archive. Archiving never deletes a bookmark — it just hides it from the main
list so you can review or bring it back later.

To browse the archive, launch the TUI and press `a`:

```sh
cairn
```

| Key | Action |
|-----|--------|
| `a` | Open the archive view (from the main list) |
| `↑` / `↓` or `j` / `k` | Navigate archived bookmarks |
| `r` | Restore the selected bookmark to the main list |
| `Esc` | Return to the main list |

Each archived entry shows its domain and the date it was archived. Pressing `r`
restores the selected bookmark: it clears the archive flag and date, moves the
bookmark back into your active list, and removes it from the archive view. The
restored bookmark appears in the main list immediately when you press `Esc`.

To keep a bookmark from being auto-archived in the first place, pin it (see the
previous step) — pinned bookmarks are treated as permanent and are never
archived. To turn auto-archiving off entirely, see
[Disable Auto-Archiving]({{< relref "/docs/configuration" >}}#disable-auto-archiving).

## 9. Check Configuration

```sh
cairn config
```

Prints the resolved database path and whether the Dropbox app key is set.

## 10. Update Cairn

Check whether a newer version is available:

```sh
cairn update --check
```

Apply the update (downloads, verifies checksum, and replaces the binary atomically):

```sh
cairn update
```

Update the Vicinae browser extension (if installed):

```sh
cairn update --extension
```

## Next Steps

- [Configuration]({{< relref "/docs/configuration" >}}) — set a custom DB path or Dropbox key
- [Dropbox Sync]({{< relref "/docs/sync" >}}) — sync bookmarks across devices
- [Vicinae Extension]({{< relref "/docs/browser-extension" >}}) — save bookmarks from the browser
