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
| `a` | Add a new bookmark |
| `d` | Delete selected bookmark |
| `p` | Toggle pin on selected bookmark |
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

## 8. Check Configuration

```sh
cairn config
```

Prints the resolved database path and whether the Dropbox app key is set.

## 9. Update Cairn

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
