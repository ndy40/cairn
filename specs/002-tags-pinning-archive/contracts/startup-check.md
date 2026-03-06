# Contract: Startup Prerequisite Check

**Feature**: 002-tags-pinning-archive
**Date**: 2026-03-06

---

## Overview

When the TUI is launched (`bm` with no subcommand, or `bm tui`), the application performs a prerequisite check before opening the TUI. The check detects the display environment and verifies the required clipboard tool is available. CLI subcommands (`bm add`, `bm list`, `bm search`, `bm delete`) are **not affected** by this check.

---

## Detection Logic

| Condition | Detected Environment | Required Tool |
|-----------|---------------------|---------------|
| `WAYLAND_DISPLAY` is set (regardless of `DISPLAY`) | Wayland | `wl-paste` (from `wl-clipboard`) |
| `WAYLAND_DISPLAY` unset, `DISPLAY` is set | X11 | `xclip` **or** `xsel` (either is sufficient) |
| Neither `WAYLAND_DISPLAY` nor `DISPLAY` is set | Unknown | (no check performed) |

Note: When both `WAYLAND_DISPLAY` and `DISPLAY` are set (XWayland), Wayland takes precedence.

---

## Outcomes

### Outcome A: Tool present (any detected environment)

- No output; TUI launches normally.

### Outcome B: Wayland detected, `wl-paste` missing

- Prints to stderr:
  ```
  Wayland detected. Please install wl-clipboard:
    sudo apt install wl-clipboard    # Debian/Ubuntu
    sudo pacman -S wl-clipboard      # Arch
    sudo dnf install wl-clipboard    # Fedora
  ```
- Exits with **non-zero status code** (code 1).
- TUI does **not** open.

### Outcome C: X11 detected, neither `xclip` nor `xsel` found

- Prints to stderr:
  ```
  X11 detected. Please install xclip or xsel:
    sudo apt install xclip           # Debian/Ubuntu (recommended)
    sudo pacman -S xclip             # Arch
    sudo dnf install xclip           # Fedora
  ```
- Exits with **non-zero status code** (code 1).
- TUI does **not** open.

### Outcome D: Unknown display environment

- Prints to stderr:
  ```
  Warning: display environment not detected. Clipboard paste (Ctrl+P) may not be available.
  ```
- TUI **continues to launch** (non-blocking).
- This message is shown once and does not appear again during the session.

---

## Function Signature

```
display.CheckPrerequisites() CheckResult
```

`CheckResult` fields:
- `DisplayType` — `Wayland`, `X11`, or `Unknown`
- `ToolFound` — `true` if the required tool was found (or environment is Unknown)
- `MissingTool` — name of the missing tool (`"wl-paste"`, `"xclip/xsel"`, or `""`)
- `InstallHint` — multi-line human-readable install instructions (`""` if not needed)
- `ShouldBlock` — `true` when the check result requires the app to exit before launching the TUI

---

## Exit Codes (TUI mode only)

| Scenario | Exit Code |
|----------|-----------|
| Successful launch | 0 (on quit) |
| Clipboard tool missing (Wayland or X11) | 1 |

---

## Interaction with Archive Check

The startup sequence for TUI mode is:

1. Parse CLI flags
2. Open database (run migrations)
3. **Run prerequisite check** → if `ShouldBlock`, print hint and exit(1)
4. **Run archive check** (`store.ArchiveStale()`) → record count
5. Launch TUI (`tea.NewProgram`) with archive count passed to `App` model
6. Display archive count in footer status on first render if count > 0

---

## Not Affected

- `bm add <url>` — no prerequisite check
- `bm list` — no prerequisite check
- `bm search <query>` — no prerequisite check
- `bm delete <id>` — no prerequisite check
- `bm version` — no prerequisite check
- `bm help` — no prerequisite check
