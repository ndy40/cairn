# Quickstart: build and verify the native window

## Build and run

```sh
cd desktop
wails3 dev                      # live reload
# or
wails3 build && ./bin/cairn-desktop
go test ./...                   # unit tests
(cd frontend && npx tsc --noEmit)
```

From the repo root: `just build-desktop`, `just test-desktop`.

## Spike first

Before building the rest, confirm on GNOME Wayland that `cairn-desktop --toggle` from a custom shortcut raises and focuses the window once the activation token is forwarded (research R8). If it can't, keep today's behaviour and note it in the docs.

## Manual QA matrix

Run each check on macOS 12+, GNOME Wayland, GNOME X11 and KDE Plasma (Wayland). Use light and dark mode.

| # | Check | Expected |
|---|---|---|
| 1 | First launch | 720×520, centred, native buttons, search focused |
| 2 | Drag the title band / search field | Window moves / text is selected |
| 3 | Resize below the minimum | Stops at 520×360 |
| 4 | Click another app | Cairn stays open, query kept |
| 5 | Esc with text, then Esc again | Clears; then nothing happens |
| 6 | Move and resize, quit, relaunch | Same size; same position except on Wayland |
| 7 | Unplug the remembered display, relaunch (macOS/X11) | Centred on the main display |
| 8 | Global shortcut from another app | Window comes forward, query selected |
| 9 | Global shortcut while Cairn is in front | macOS: app hidden; Linux: minimised (or hidden in background mode) |
| 10 | `cairn-desktop --toggle` from a custom shortcut (GNOME Wayland) | Window raised and focused |
| 11 | Close the window (macOS) | App keeps running; Dock click reopens |
| 12 | Close the window (GNOME, no tray) | App quits |
| 13 | Close the window (KDE, defaults) | App stays in the tray |
| 14 | Cmd/Ctrl+Q during a sync push | "Finishing sync…", quits within 5 s |
| 15 | Add with a URL on the clipboard | Sheet/dialog, URL selected, Enter saves, new row selected |
| 16 | Add a duplicate | Existing row selected, "Already bookmarked" toast |
| 17 | `+go.dev` in search | Add modal with `go.dev` |
| 18 | Right-click a row | Open Link, Copy Link, Pin/Unpin |
| 19 | Settings: Open at login on, log out and in | Cairn running, no window |
| 20 | Upgrade from v0.7.x with installer login item | Settings shows Open at login on |
| 21 | Take the global shortcut in another app, relaunch | Settings shows the conflict and the copy-command fallback |
| 22 | Expire the Dropbox token | Banner under the title area; badge on the menu bar icon (macOS) |
| 23 | Keyboard only: Tab through all controls | Every control reachable; focus visible |
| 24 | Contrast check (light/dark, each desktop) | All text ≥ 4.5:1 |
