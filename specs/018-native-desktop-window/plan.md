# Implementation Plan: cairn-desktop as a native app window

**Branch**: `018-native-desktop-window` | **Date**: 2026-10-27 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/018-native-desktop-window/spec.md`

## Summary

Turn cairn-desktop from a frameless, always-on-top popup into a regular app window. On macOS: inset traffic lights in a unified title band that holds the search field, a Dock icon, a full menu bar and Dock-click reopen. On Linux: GTK4's own title bar with Cairn's search bar below it, Add and main-menu buttons, and a KDE variant. Behaviour changes: no hide on focus loss, Esc never closes, size and position remembered, platform close/quit rules, global shortcut that brings the window forward or puts it away, an Add Bookmark sheet/dialog, a refreshed list, toasts and banners, a Settings panel (open at login, menu bar icon, run in background, hide after open) and a row context menu (Open, Copy Link, Pin).

Almost everything uses existing Wails v3 APIs (`MacTitleBarHiddenInsetUnified`, `ActivationPolicyRegular`, `ApplicationShouldHandleReopen`, app menus with roles, `app.Autostart`, `GlobalShortcut`, `SystemTray.Show/Hide`). The only native code is a small Linux cgo helper to raise the window with a forwarded Wayland activation token; it is spiked first.

## Technical Context

**Language/Version**: Go 1.26 (desktop module, as today), TypeScript 5.x (frontend)
**Primary Dependencies**: Wails v3.0.0-beta.33 (unchanged), `@wailsio/runtime`, Vite; `github.com/godbus/dbus/v5` promoted from indirect to direct (pure Go) for tray detection and the Settings portal accent colour
**Storage**: SQLite unchanged (no schema change). New `desktop-settings.json`; `desktop-state.json` gains `x`, `y`, `screen`
**Testing**: `go test ./...` in `desktop/`, `tsc --noEmit` for the frontend, manual QA matrix in [quickstart.md](quickstart.md)
**Target Platform**: macOS 12+ (universal), Linux GTK4/WebKitGTK 6.0 on GNOME (Wayland, X11) and KDE Plasma
**Project Type**: Desktop app (separate Go module `desktop/` with an embedded web frontend)
**Performance Goals**: Window ready for typing ≤ 300 ms after the global shortcut on a warm app (SC-001)
**Constraints**: No heavy UI framework; system fonts; text contrast ≥ 4.5:1; no change to the CLI binary
**Scale/Scope**: ~10 Go files and 3 frontend files in `desktop/`, `install.sh`, one docs page

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Gate | Status | Notes |
|---|---|---|
| No CGO | Pass (existing exception) | The CLI stays `CGO_ENABLED=0`. `desktop/` is a separate module that already needs cgo for Wails/GTK (accepted in TASK-143). This feature adds one small Linux cgo helper inside that module and one pure-Go dependency (godbus). See Complexity Tracking. |
| Single binary | Pass | `cairn-desktop` stays one binary; no external runtime. |
| Task management | Pass (pending) | After `/speckit.tasks`, every task in `tasks.md` is created in Backlog with `backlog task create`. |
| Backward-compatible migrations | Pass | No SQLite schema change. JSON state files stay backward compatible (new optional fields; missing file → defaults). |

**Post-design re-check**: still passing. The design adds no packages outside `desktop/` and no new abstractions beyond small files with one job each (settings, platform detection, tray watch, menus, activation).

## Project Structure

### Documentation (this feature)

```text
specs/018-native-desktop-window/
├── plan.md              # This file
├── research.md          # Phase 0: decisions R1–R16
├── data-model.md        # Phase 1: state files, runtime info, lifecycle
├── quickstart.md        # Phase 1: build, spike, QA matrix
├── contracts/
│   ├── frontend-bindings.md   # Go ↔ frontend methods and events
│   └── cli-and-install.md     # flags, files, login items, install.sh
├── checklists/requirements.md
└── tasks.md             # Phase 2 (/speckit.tasks, not created here)
```

### Source Code

```text
desktop/
├── main.go              # app options: ActivationPolicyRegular, menus, reopen, --background, AdditionalData token
├── launcher.go          # window options (framed, Mac title bar, min size), show/hide/toggle, close hook, tray
├── lifecycle.go         # NEW: closeAction(), shortcutAction() decision functions + quit/wait
├── lifecycle_test.go    # NEW
├── menus.go             # NEW: macOS application menu, Dock menu, tray menu
├── platform.go          # NEW: PlatformInfo from env; tray watch (godbus); portal accent colour
├── platform_test.go     # NEW
├── activation_linux.go  # NEW: cgo helper, gtk_window_set_startup_id before present
├── activation_other.go  # NEW: no-op for darwin
├── settings.go          # NEW: DesktopSettings load/save/defaults
├── settings_test.go     # NEW
├── windowstate.go       # extended: x, y, screen; clamp; restore rules
├── windowstate_test.go  # extended
├── window_service.go    # NEW: WindowService bindings (Platform, Settings, OpenAtLogin, ShortcutStatus, …)
├── service.go           # + FindByURL, SetPinned, CopyLink, SyncNow
├── service_test.go      # extended
├── sync.go              # + Wait(timeout), sync:ok event, Running()
├── sync_test.go         # extended
├── focus.go, focus_test.go   # DELETED
└── frontend/
    ├── index.html       # title band/search bar, list, footer, <dialog> add, settings, shortcuts, about, toast, banner, context menu
    └── src/
        ├── main.ts      # rewritten around views/modals; keyboard map; events
        └── style.css    # tokens, data-os/data-desktop variants, light/dark

install.sh                       # stop creating login items; --no-autostart no-op; restart with --background
docs/content/docs/desktop-app.md # rewritten for the window app
```

**Structure Decision**: Keep the existing flat `package main` layout of `desktop/`, adding one file per new responsibility. Decision logic (close action, shortcut action, restore position, settings defaults) is written as pure functions so it can be unit tested without a window.

## Delivery Order

1. **Spike**: activation token on GNOME Wayland (R8). Decides whether `--toggle` raising is in or documented as a limitation.
2. **P1 window**: framed window per platform, remove the focus guard, Esc rules, size/position persistence, single-instance raise, focus-on-show. (Stories 1–2)
3. **P2 lifecycle**: Dock/reopen, close/quit rules, tray detection, Syncer wait, shortcut toggle and status. (Story 3)
4. **P2 add modal**: `<dialog>` sheet/dialog, `FindByURL` duplicate selection. (Story 4)
5. **P2 menus and settings**: macOS menu bar, Linux main menu, context menu, Settings panel, autostart, shortcuts list, install.sh change. (Story 5)
6. **P3 visual refresh**: rows, sections, empty states, toasts and banner, tokens, accent colour, KDE variant. (Story 6)
7. **Docs** and the QA matrix.

## Risks

| Risk | Mitigation |
|---|---|
| Wayland activation token can't reach GTK through Wails | Spike first; fall back to current behaviour and document it |
| `app.Autostart` on macOS 13+ uses SMAppService and doesn't see the installer's plist | Verify in the spike; if needed, a one-time check of the plist path (R9) |
| Title band height differs between macOS versions | Measure on a build; drive it from one CSS variable |
| Wails beta API changes | Version stays pinned; no upgrade in this feature |

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|---|---|---|
| cgo helper `activation_linux.go` (inside the already-cgo `desktop/` module) | Raising the window from `--toggle` on Wayland needs `gtk_window_set_startup_id`, which Wails doesn't expose | Without it, GNOME shows "Cairn is ready" instead of raising the window, a known complaint path for keyboard users |
