---
title: "Desktop App"
weight: 25
---

# Desktop App

`cairn-desktop` is a small quick-launch popup for your bookmarks. Press a
hotkey from anywhere, start typing to search, and press `Enter` to open the
link. You can also add a bookmark without opening a terminal.

It works on Linux and macOS and uses the same database, search ranking and
Dropbox sync as the `cairn` CLI. A bookmark you add in one appears in the other.

## Installation

### Via Installer (Recommended)

```sh
curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh -s -- -y --with-desktop
```

Without `-y`, the installer asks whether to install the desktop app. The
installer:

| | Linux | macOS |
|---|---|---|
| App | `cairn-desktop` in the install directory (default `/usr/local/bin`) | `~/Applications/cairn-desktop.app` |
| Menu entry | `~/.local/share/applications/cairn-desktop.desktop` | — |
| Start at login | `~/.config/autostart/cairn-desktop.desktop` | `~/Library/LaunchAgents/io.github.ndy40.cairn.desktop.plist` |

At login the app starts hidden in the background. Pass `--no-autostart` to
the installer to skip this.

### Requirements

**Linux** needs GTK 4 and WebKitGTK 6.0. These ship with most current desktops;
the installer warns you if WebKitGTK 6.0 is missing.

| Distribution | Command |
|---|---|
| Fedora | `sudo dnf install webkitgtk6.0` |
| Debian / Ubuntu (24.04+) | `sudo apt install libwebkitgtk-6.0-4` |
| Arch | `sudo pacman -S webkitgtk-6.0` |

The Linux builds are produced on Ubuntu 24.04, so they need a distribution of
a similar age or newer.

**macOS** 12 (Monterey) or later. A single download runs on both Apple Silicon
and Intel Macs.

{{< hint warning >}}
**macOS security warning.** The app is not yet notarised by Apple. If you
download it from the releases page in a browser, macOS may refuse to open it
the first time. Right-click `cairn-desktop.app`, choose **Open**, then confirm.
Installing with `install.sh` avoids this in most cases.
{{< /hint >}}

### Manual Download

Download the archive for your platform from
[GitHub Releases](https://github.com/ndy40/cairn/releases) and check it against
`checksums.txt`:

| Platform | File |
|---|---|
| Linux x86-64 | `cairn-desktop-linux-amd64.tar.gz` |
| Linux ARM64 | `cairn-desktop-linux-arm64.tar.gz` |
| macOS (Apple Silicon and Intel) | `cairn-desktop-darwin-universal.zip` |

## Opening the Popup

Start the app once, either from your applications menu or with
`cairn-desktop &`. It then waits in the background, and you can open the
popup in any of these ways:

| Method | How |
|---|---|
| Global hotkey | `Ctrl+Shift+Space` (Linux), `Cmd+Shift+Space` (macOS) |
| Command | `cairn-desktop --toggle` |
| Tray icon | Click the Cairn icon in the menu bar or system tray |

Running `cairn-desktop` again while it is already running does not start a
second copy. It shows or hides the existing popup instead.

### Linux: GNOME and Other Wayland Desktops

On Wayland, apps can't claim a global hotkey by themselves; the desktop decides.
On GNOME, KDE and other desktops that support it, the first launch asks you to
approve the shortcut. The desktop may also let you choose a different key.

If the hotkey does nothing, or you'd rather pick the key yourself, add a custom
keyboard shortcut that runs the toggle command:

1. Open **Settings → Keyboard → View and Customise Shortcuts → Custom Shortcuts**.
2. Add a shortcut with the command `cairn-desktop --toggle`.
3. Assign the keys you want.

A second shortcut with `cairn-desktop --add` opens the popup straight on the
**Add bookmark** form.

{{< hint info >}}
**No tray icon on GNOME?** Stock GNOME doesn't show tray icons. Install the
AppIndicator extension (`gnome-shell-extension-appindicator` on Fedora) to see
it. The hotkey and `--toggle` work without it.
{{< /hint >}}

### macOS

The app runs in the menu bar and has no Dock icon. If `Cmd+Shift+Space` is
already used by another app, set a different key (see
[Configuration](#configuration)).

## Searching and Opening

Type to search. Results use the same ranking as the TUI: title, domain, tags and
description, with exact matches first. With an empty search box, pinned
bookmarks are listed first, then the most recently added.

| Key | Action |
|---|---|
| `↑` / `↓` | Move the selection |
| `Enter` or click | Open the bookmark in your default browser and close the popup |
| `Ctrl+N` (`Cmd+N` on macOS) | Add a bookmark |
| `Esc` | Clear the search; press again on an empty box to close the popup |

The popup also closes when you click outside it.

Drag any edge of the popup to resize it. The size is remembered, including
after a restart, in `desktop-state.json` next to `cairn.json`; delete that file
to go back to the default size. To move the popup, drag an empty part of it,
such as the space below the results or the footer.

## Adding a Bookmark

Press `Ctrl+N`, or type `+` in the search box, to open the **Add bookmark**
form. If the clipboard holds a link, it is filled in for you. Typing
`+go.dev` opens the form with `go.dev` already entered.

| Key | Action |
|---|---|
| `Tab` | Move between the URL and Tags fields |
| `Enter` | Save |
| `Esc` | Cancel and return to search |

- **URL**: if you leave out `https://`, it is added for you. Only `http` and
  `https` links are accepted.
- **Tags**: optional, separated by commas.

Cairn fetches the page title and description, just like `cairn add`. If the page
can't be fetched, the bookmark is still saved, using the site name as its title.
Adding a link you have already saved shows **Already bookmarked**.

## Sync

If you have set up [Dropbox Sync]({{< relref "/docs/sync" >}}), the desktop app
syncs automatically:

- **Pulls** when it starts and when you open the popup (at most once every two minutes)
- **Pushes** after you add a bookmark

Problems such as an expired Dropbox login appear at the bottom of the popup.

Sync is set up from a terminal with `cairn sync setup`. The desktop app picks
up the change the next time it syncs, without a restart.

## Configuration

The desktop app reads the same [`cairn.json`]({{< relref "/docs/configuration" >}})
as the CLI, including `db_path`. One setting applies only to the desktop app:

| JSON key | Environment variable | Description |
|---|---|---|
| `desktop_hotkey` | `CAIRN_DESKTOP_HOTKEY` | Global shortcut that opens the popup. Default `CmdOrCtrl+Shift+Space`. Set to `off` to disable it. |

Shortcuts are written as modifiers and a key joined with `+`, for example
`Ctrl+Alt+B` or `CmdOrCtrl+Shift+K`. `CmdOrCtrl` means `Cmd` on macOS and
`Ctrl` on Linux. Restart the app after changing the setting.

```json
{
  "desktop_hotkey": "Ctrl+Alt+B"
}
```

## Command-Line Options

| Option | Description |
|---|---|
| *(none)* | Start the app and open the popup. If it is already running, show or hide the popup. |
| `--toggle` | Show or hide the popup. |
| `--add` | Open the popup on the **Add bookmark** form. |
| `--hidden` | Start in the background without opening the popup (used at login). |

## Updating

`cairn update` updates the desktop app along with the CLI, so both stay on the
same release. To update only the desktop app:

```sh
cairn update --desktop
```

You don't need to quit the app first. If it's running, `cairn update` replaces
it, quits the old version and starts the new one in the background. If you run
it with `sudo`, it can't see your running app; quit and restart the app
yourself after updating.

Re-running the installer with `--with-desktop` also updates the app:

```sh
curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh -s -- -y --with-desktop
```

## Uninstalling

`cairn uninstall` removes the desktop app along with the CLI: it stops the app
if it is running, then deletes it with its launcher, login item, icon and
version file. Your bookmarks are kept unless you add `--purge`. See
[`cairn uninstall`]({{< relref "/docs/cli-reference" >}}#cairn-uninstall).

To remove only the desktop app and keep the CLI, quit the app from its tray
menu (**Quit Cairn**), then remove its files.

**Linux** (replace `/usr/local/bin` if you installed elsewhere):

```sh
rm /usr/local/bin/cairn-desktop \
   ~/.local/share/applications/cairn-desktop.desktop \
   ~/.config/autostart/cairn-desktop.desktop \
   ~/.local/share/cairn/cairn-desktop.png \
   ~/.local/share/cairn/desktop-version.txt
```

**macOS:**

```sh
rm -rf ~/Applications/cairn-desktop.app
rm ~/Library/LaunchAgents/io.github.ndy40.cairn.desktop.plist \
   ~/Library/Application\ Support/cairn/desktop-version.txt
```

Your bookmarks are untouched; they belong to Cairn, not the desktop app.

## Build from Source

The desktop app lives in `desktop/` as its own Go module. It needs Go 1.26+,
Node.js 20.19+, the [Wails v3](https://v3.wails.io) CLI, and on Linux the GTK 4
and WebKitGTK 6.0 development packages (`gtk4-devel webkitgtk6.0-devel` on
Fedora, `libgtk-4-dev libwebkitgtk-6.0-dev` on Debian/Ubuntu).

```sh
cd desktop
go install github.com/wailsapp/wails/v3/cmd/wails3@$(go list -m -f '{{.Version}}' github.com/wailsapp/wails/v3)
wails3 build          # binary at desktop/bin/cairn-desktop
```

Use `wails3 dev` to run with live reload, or `just build-desktop` /
`just test-desktop` from the repository root.
