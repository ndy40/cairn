# Contract: command line, installer and files

## `cairn-desktop` flags

| Flag | Status | First launch | Forwarded to a running instance |
|---|---|---|---|
| *(none)* | changed | Open the window, focus search | Bring the window forward (no longer toggles) |
| `--toggle` | unchanged name | Open the window | Bring forward, or put away if already in front (same as the global shortcut) |
| `--add` | unchanged | Open the window with the add modal | Bring forward with the add modal |
| `--background` | **new** | Start without a window (used at login) | Ignored |
| `--hidden` | deprecated alias | Same as `--background` | Ignored |

A second launch also forwards `XDG_ACTIVATION_TOKEN` / `DESKTOP_STARTUP_ID` from its environment so the running instance can raise its window on Wayland.

## `cairn.json`

Unchanged. `desktop_hotkey` (env `CAIRN_DESKTOP_HOTKEY`): accelerator, empty for the default `CmdOrCtrl+Shift+Space`, `off` to disable.

## Files next to `cairn.json`

| File | Status | Contents |
|---|---|---|
| `desktop-state.json` | extended | `width`, `height`, optional `x`, `y`, `screen` |
| `desktop-settings.json` | **new** | `showInMenuBar`, `runInBackground`, `hideAfterOpen` |

## Login items

| Platform | Identifier | Created by |
|---|---|---|
| macOS | `io.github.ndy40.cairn.desktop` (SMAppService on 13+ for the bundled app, otherwise `~/Library/LaunchAgents/io.github.ndy40.cairn.desktop.plist`) | The app (Settings) |
| Linux | `~/.config/autostart/cairn-desktop.desktop` | The app (Settings) |

## `install.sh`

- Stops creating login items for new installs. Existing login items are left in place (they keep working because the app still accepts `--hidden`).
- `--no-autostart` stays accepted, does nothing, and prints a notice that start-at-login is now set in Cairn's Settings.
- Still installs the app and the Linux menu entry (`StartupWMClass=cairn-desktop`), and still restarts a running app after an update (now with `--background`).
