# Data Model: cairn-desktop as a native app window

No SQLite schema changes. All new state lives in small JSON files next to `cairn.json` (never inside it), or is read live from the OS.

## WindowState (`desktop-state.json`, extended)

| Field | Type | Rules |
|---|---|---|
| `width` | int | ≥ 520; default 720 |
| `height` | int | ≥ 360; default 520 |
| `x` | int, optional | Saved on macOS and X11 only |
| `y` | int, optional | Saved on macOS and X11 only |
| `screen` | string, optional | ID or name of the screen the window was on |

- Older files with only `width`/`height` still load. Sizes below the minimum fall back to defaults (existing behaviour).
- On load: if `x/y` are present, not on Wayland, and `screen` still exists, restore the position clamped to that screen's work area; otherwise centre on first launch (no file) or let the compositor choose (Wayland).
- Size is clamped to the current screen's work area.
- Written debounced (500 ms) after resize/move and on shutdown; zero sizes reported while hidden are ignored (existing rule).

## DesktopSettings (`desktop-settings.json`, new)

| Field | Type | Default | Platform |
|---|---|---|---|
| `showInMenuBar` | bool | `true` | macOS |
| `runInBackground` | bool or null | null → `true` on KDE, `false` elsewhere | Linux |
| `hideAfterOpen` | bool | `false` | all |

- Missing file or unknown keys → defaults; invalid JSON is logged and treated as defaults (the file is not overwritten until the user changes a setting).
- Written atomically (temp file + rename).

## Open at login (OS-owned, not stored)

- Source of truth: `app.Autostart.Status()`.
- Identifier: `io.github.ndy40.cairn.desktop` (macOS), `cairn-desktop` (Linux); launch arguments `--background`.
- State: `enabled` ⇄ `disabled`, toggled only from Settings. The installer no longer creates it; entries it created earlier are detected as enabled.

## PlatformInfo (runtime, sent to the frontend)

| Field | Type | Source |
|---|---|---|
| `os` | `"darwin"` \| `"linux"` | `runtime.GOOS` |
| `desktop` | `"gnome"` \| `"kde"` \| `"other"` \| `""` | `XDG_CURRENT_DESKTOP` |
| `wayland` | bool | `XDG_SESSION_TYPE == "wayland"` |
| `trayAvailable` | bool | macOS: always true; Linux: StatusNotifierWatcher has an owner; updates live |
| `accentColor` | string, optional | Settings portal on Linux; empty on macOS (CSS `AccentColor` used) |

## ShortcutStatus (runtime)

| Field | Type | Meaning |
|---|---|---|
| `hotkey` | string | Effective accelerator (configured or default), or `off` |
| `registered` | bool | Registration succeeded |
| `error` | string | Reason when not registered (taken, refused, unsupported) |
| `toggleCommand` | string | `cairn-desktop --toggle`, for the copy button |

## Bookmark DTO (unchanged)

`id`, `title`, `url`, `domain`, `tags`, `pinned`. The frontend derives: the section (`Pinned`/`Recent` for an empty query), the monogram letter and tint (hash of `domain`), match highlighting (case-insensitive match of query terms in `title`), and the visible tags (first 3 + `+N`).

## Window lifecycle (state transitions)

```
            launch / --toggle / shortcut / Dock / tray
 [Hidden] ─────────────────────────────────────────────▶ [Shown, focused]
    ▲  ▲                                                     │   │
    │  └──── close (macOS; or Linux with background+tray) ◀──┘   │
    │                                                            │ shortcut while in front
    └──── macOS: app hidden / Linux: minimised or hidden ◀───────┘
 [Shown] ── close on Linux without background or tray ──▶ [Quit]
 any ── Quit (menu, Cmd/Ctrl+Q, tray) ──▶ [Finishing sync ≤5s] ──▶ [Quit]
 launch --background ──▶ [Hidden] (no window shown)
```
