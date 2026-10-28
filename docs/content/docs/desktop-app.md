---
title: "Desktop App"
weight: 25
---

# Desktop App

`cairn-desktop` is a desktop app for your bookmarks. It opens as a normal app
window, and a global shortcut brings it forward from anywhere. Type a few
letters to search and press `Enter` to open the link. You can also add a
bookmark without opening a terminal.

It works on Linux and macOS and uses the same database, search ranking and
Dropbox sync as the `cairn` CLI. A bookmark you add in one appears in the other.

## Installation

### Via Installer (Recommended)

```sh
curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh -s -- -y --with-desktop
```

Without `-y`, the installer asks whether to install the desktop app. The
installer adds:

| | Linux | macOS |
|---|---|---|
| App | `cairn-desktop` in the install directory (default `/usr/local/bin`) | `~/Applications/cairn-desktop.app` |
| Menu entry | `~/.local/share/applications/cairn-desktop.desktop` | — |

The installer does not set Cairn to start at login. Turn that on in
[Settings](#settings) if you want it. `--no-autostart` is still accepted so
existing scripts keep working, but it does nothing.

{{< hint info >}}
**Upgrading from the popup version?** If an earlier installer set Cairn to
start at login, that login item stays in place and keeps working. Settings
shows **Open Cairn at login** as on.
{{< /hint >}}

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

## Opening Cairn

Open **Cairn** from your applications menu (Linux) or the Applications folder
(macOS), or run `cairn-desktop`. The window opens with the cursor in the search
field.

Once Cairn is running, you can bring the window forward in any of these ways:

| Method | How |
|---|---|
| Global shortcut | `Ctrl+Shift+Space` (Linux), `Cmd+Shift+Space` (macOS) |
| Dock or taskbar | Click the Cairn icon |
| Menu bar or tray icon | Click the Cairn icon, or choose **Search Bookmarks…** from its menu |
| Command | `cairn-desktop --toggle` |

Launching Cairn again while it is running does not start a second copy. It
brings the existing window forward. Each of these puts the cursor in the search
field and selects what you typed last, so you can start typing straight away.

## The Window

Cairn uses your platform's own window frame:

- **macOS**: the close, minimise and zoom buttons sit in a title band next to
  the search field.
- **GNOME and KDE**: the desktop draws its usual title bar, with your theme and
  button layout. Cairn's search bar sits just below it.

Move the window by dragging its title area. Dragging inside the search field
selects text instead. Resize it from any edge, or maximise or zoom it. The
smallest size is 520×360.

The window stays open when you click into another app. It only goes away when
you close it, hide it, press the global shortcut while it is in front, or turn
on **Hide Cairn after opening a link**.

Cairn remembers the window size between runs. On macOS and on Linux with X11 it
also remembers the position. On Wayland, apps can't place their own windows, so
only the size is remembered and the desktop picks the position. If the display
the window was on is no longer connected, the window opens in the centre of the
main display.

## Closing and Quitting

Closing the window and quitting Cairn are different on each platform.

| | Close the window | Quit |
|---|---|---|
| macOS | Cairn keeps running. The Dock icon, menu bar icon and global shortcut bring the window back. | `Cmd+Q`, **Cairn → Quit Cairn**, or **Quit Cairn** in the menu bar icon menu |
| Linux, **Run in background** on and a tray available | Cairn keeps running in the tray. | `Ctrl+Q`, **Quit** in the main menu, or **Quit Cairn** in the tray menu |
| Linux, otherwise | Cairn quits. | `Ctrl+Q`, or **Quit** in the main menu |

**Run in background** is on by default on KDE Plasma and off on other desktops.
Change it in [Settings](#settings).

Quitting never asks for confirmation. If a sync upload is in progress, the
window shows **Finishing sync…** and Cairn quits once it is done, or after 5
seconds at most.

{{< hint info >}}
**No tray icon on GNOME?** Stock GNOME doesn't show tray icons, so closing the
window quits Cairn. To see the tray icon, install the AppIndicator extension
(`gnome-shell-extension-appindicator` on Fedora). Or bind a
[custom shortcut](#wayland-and-custom-shortcuts) to `cairn-desktop --toggle`:
it starts Cairn when it isn't running.
{{< /hint >}}

## Global Shortcut

The global shortcut works from any app while Cairn is running:

- If Cairn is not in front, it brings the window forward. A closed window opens
  again.
- If Cairn is in front, it gets out of the way and focus returns to the app you
  were using. On macOS the app is hidden. On Linux the window is minimised, or
  hidden to the tray when **Run in background** is on.

The default is `CmdOrCtrl+Shift+Space`. To change it, set `desktop_hotkey` in
[`cairn.json`]({{< relref "/docs/configuration" >}}), then restart the app:

| JSON key | Environment variable | Description |
|---|---|---|
| `desktop_hotkey` | `CAIRN_DESKTOP_HOTKEY` | Global shortcut for the Cairn window. Default `CmdOrCtrl+Shift+Space`. Set to `off` to disable it. |

Shortcuts are written as modifiers and a key joined with `+`, for example
`Ctrl+Alt+B` or `CmdOrCtrl+Shift+K`. `CmdOrCtrl` means `Cmd` on macOS and
`Ctrl` on Linux.

```json
{
  "desktop_hotkey": "Ctrl+Alt+B"
}
```

If the shortcut can't be registered, for example because another app already
uses it, Cairn keeps running. Settings shows the problem and the command to use
instead.

### Wayland and Custom Shortcuts

On Wayland, apps can't claim a global shortcut by themselves; the desktop
decides. On GNOME, KDE and other desktops that support it, the first launch
asks you to approve the shortcut. The desktop may also let you choose a
different key.

If the shortcut does nothing, or you'd rather pick the key yourself, add a
custom keyboard shortcut in your system settings that runs the toggle command.
On GNOME:

1. Open **Settings → Keyboard → View and Customise Shortcuts → Custom Shortcuts**.
2. Add a shortcut with the command `cairn-desktop --toggle`.
3. Assign the keys you want.

Settings has a **Copy Command** button for this. A second shortcut with
`cairn-desktop --add` opens the window straight on **Add Bookmark**.

{{< hint info >}}
If GNOME shows a "Cairn is ready" notification instead of bringing the window
forward, click the notification to switch to Cairn.
{{< /hint >}}

## Searching and Opening

Type to search. Results use the same ranking as the TUI: title, domain, tags and
description, with exact matches first. With an empty search field, pinned
bookmarks are listed under **Pinned**, followed by recent ones under
**Recent**.

| Key | Action |
|---|---|
| `↑` / `↓` | Move the selection; the cursor stays in the search field |
| `Enter` or click | Open the bookmark in your default browser |
| `Esc` | Clear the search |

The window stays open after you open a link, unless **Hide Cairn after opening
a link** is on. `Esc` clears the search field; on an empty field it does
nothing. It never closes the window.

## Adding a Bookmark

Press `Ctrl+N` (`Cmd+N` on macOS), click the **+** button, or type `+` in the
search field to open **Add Bookmark**. On macOS it opens as a sheet under the
title band. On Linux it opens as a dialog over the window.

If the clipboard holds a link, it is filled in and selected for you. Typing
`+go.dev` in the search field opens the form with `go.dev` already entered.

| Key | Action |
|---|---|
| `Tab` | Move between the URL and Tags fields |
| `Enter` | Add the bookmark |
| `Esc` | Cancel and return to search |

- **URL**: if you leave out `https://`, it is added for you. Only `http` and
  `https` links are accepted.
- **Tags**: optional, separated by commas.

Cairn fetches the page title and description, just like `cairn add`. If the page
can't be fetched, the bookmark is still saved, using the site name as its title.
After it is saved, the form closes and the new bookmark is selected in the list.

If you have already saved the link, the form closes, the existing bookmark is
selected, and **Already bookmarked** appears briefly.

## Menus

### macOS Menu Bar

| Menu | Items |
|---|---|
| Cairn | About Cairn, Settings… (`Cmd+,`), Services, Hide Cairn, Hide Others, Show All, Quit Cairn (`Cmd+Q`) |
| File | New Bookmark… (`Cmd+N`), Open Link, Copy Link (`Shift+Cmd+C`), Pin / Unpin, Sync Now, Close Window (`Cmd+W`) |
| Edit | The standard editing items, plus Find Bookmarks (`Cmd+F`) |
| Window | The standard window items |
| Help | Cairn Help, Keyboard Shortcuts (`Cmd+/`) |

File menu items act on the selected bookmark.

### Linux Main Menu

The search bar has an **Add** button and a main menu button (☰). The main menu
holds Add Bookmark…, Sync Now, Preferences, Keyboard Shortcuts, About Cairn and
Quit.

### Bookmark Menu

Right-click a bookmark, or press `Shift+F10` or the Menu key, for Open Link,
Copy Link and Pin (Unpin if it is already pinned). Pinning works like
`cairn pin`.

### Menu Bar and Tray Icon

On macOS, Cairn shows an icon in the menu bar. You can hide it in Settings. On
Linux, it shows a tray icon when your desktop has a system tray. The icon's
menu holds Search Bookmarks…, Add Bookmark…, the sync status, Settings… and
Quit Cairn.

## Settings

Open Settings with `Cmd+,` / `Ctrl+,`, from **Cairn → Settings…** on macOS, or
from **Preferences** in the main menu on Linux. Changes take effect straight
away.

| Setting | Platform | Default | What it does |
|---|---|---|---|
| Open Cairn at login | Both | Off | Starts Cairn in the background when you log in. No window is shown until you summon it. |
| Show in menu bar | macOS | On | Shows the Cairn icon in the menu bar. With it off, the Dock icon and global shortcut still bring Cairn back. |
| Run in background | Linux | On for KDE, off elsewhere | Keeps Cairn running in the tray when you close the window. Unavailable when your desktop has no system tray. |
| Hide Cairn after opening a link | Both | Off | Puts Cairn away after you open a bookmark. |

Settings also shows the global shortcut. If it couldn't be registered, it says
why and offers a **Copy Command** button for `cairn-desktop --toggle`.

### Open at Login

Cairn creates and removes its own login item. Turning the setting off removes
it.

| Platform | Login item |
|---|---|
| macOS | `~/Library/LaunchAgents/io.github.ndy40.cairn.desktop.plist`. System Settings lists it under **General → Login Items → Allow in the Background**, not under **Open at Login**. |
| Linux | `~/.config/autostart/cairn-desktop.desktop` |

## Keyboard Shortcuts

Press `Cmd+/` (macOS) or `Ctrl+?` (Linux) to see this list in the app.

| Action | macOS | Linux |
|---|---|---|
| Focus search | `Cmd+F` | `Ctrl+F` |
| Move the selection | `↑` / `↓` | `↑` / `↓` |
| Open the selected bookmark | `Enter` | `Enter` |
| Copy link | `Shift+Cmd+C` | `Shift+Ctrl+C` |
| Add a bookmark | `Cmd+N` | `Ctrl+N` |
| Clear the search, or cancel a dialog | `Esc` | `Esc` |
| Close the window | `Cmd+W` | `Ctrl+W` |
| Quit | `Cmd+Q` | `Ctrl+Q` |
| Settings | `Cmd+,` | `Ctrl+,` |
| Show keyboard shortcuts | `Cmd+/` | `Ctrl+?` |
| Bring Cairn forward or put it away (global) | `Cmd+Shift+Space` | `Ctrl+Shift+Space` |

## Sync

If you have set up [Dropbox Sync]({{< relref "/docs/sync" >}}), the desktop app
syncs automatically:

- **Pulls** when it starts and when you bring the window forward (at most once
  every two minutes)
- **Pushes** after you add a bookmark
- **Sync Now** in the menus pulls and pushes straight away

Messages appear in the window by how urgent they are:

- **Banner**: problems that need you, such as an expired Dropbox sign-in, show
  as a banner under the title area. It stays until the next successful sync or
  until you click **Dismiss**, and comes back if sync fails again. On macOS the
  menu bar icon also shows a warning.
- **Toast**: short confirmations, such as **Already bookmarked** or **Link
  copied**, show for about 4 seconds.
- **Status**: background progress, such as **Finishing sync…**, shows in the
  footer on macOS and KDE, and as a toast on GNOME.

Sync is set up from a terminal with `cairn sync setup`. The desktop app picks
up the change the next time it syncs, without a restart.

## Command-Line Options

| Option | First launch | When Cairn is already running |
|---|---|---|
| *(none)* | Open the window | Bring the window forward |
| `--toggle` | Open the window | Bring the window forward, or put it away if it is in front (like the global shortcut) |
| `--add` | Open the window on **Add Bookmark** | Bring the window forward on **Add Bookmark** |
| `--background` | Start without a window (used at login) | Ignored |
| `--hidden` | Deprecated. Same as `--background`; kept for login items made by earlier installers | Ignored |

## Files

The desktop app reads the same [`cairn.json`]({{< relref "/docs/configuration" >}})
as the CLI, including `db_path`. It never writes to `cairn.json`. Its own state
lives in two files next to it:

| File | Contents |
|---|---|
| `desktop-state.json` | Window size and, where the platform allows it, position and display |
| `desktop-settings.json` | The Settings choices: `showInMenuBar`, `runInBackground`, `hideAfterOpen` |

Delete `desktop-state.json` to go back to the default size and position.
Delete `desktop-settings.json` to reset Settings to the defaults. Quit Cairn
first, as it rewrites the files when it runs. **Open Cairn at login** is not
stored here; it is the login item itself.

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

To remove only the desktop app and keep the CLI:

1. In Settings, turn off **Open Cairn at login**. This removes the login item.
2. Quit Cairn (`Cmd+Q` / `Ctrl+Q`).
3. Remove its files.

**Linux** (replace `/usr/local/bin` if you installed elsewhere):

```sh
rm /usr/local/bin/cairn-desktop \
   ~/.local/share/applications/cairn-desktop.desktop \
   ~/.local/share/cairn/cairn-desktop.png \
   ~/.local/share/cairn/desktop-version.txt
rm -f ~/.config/autostart/cairn-desktop.desktop
```

**macOS:**

```sh
rm -rf ~/Applications/cairn-desktop.app \
   ~/Library/Application\ Support/cairn/desktop-version.txt
rm -f ~/Library/LaunchAgents/io.github.ndy40.cairn.desktop.plist
```

The `rm -f` lines remove the login item if you skipped step 1. To also remove
the window state and settings, delete `desktop-state.json` and
`desktop-settings.json` from the folder that holds `cairn.json`.

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
