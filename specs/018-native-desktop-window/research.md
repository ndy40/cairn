# Research: cairn-desktop as a native app window

All findings checked against Wails `v3.0.0-beta.33` (the version pinned in `desktop/go.mod`) and the current `desktop/` sources.

## R1. macOS window chrome

- **Decision**: Standard (framed) window with `Mac.TitleBar = application.MacTitleBarHiddenInsetUnified` (transparent title bar, hidden title, full-size content, unified toolbar so the traffic lights are vertically centred in a ~52px band). `Mac.Backdrop = MacBackdropTranslucent`; the result list and footer paint an opaque CSS background so only the title band shows the material. The band is `--wails-draggable: drag`; the search field and buttons are `no-drag`. `Title: "Cairn"` stays set for the Window menu and Mission Control.
- **Rationale**: A preset exists in Wails (`MacTitleBarHiddenInsetUnified`), so no native code is needed. Measure the actual band height on a build and set a CSS variable to match.
- **Alternatives**: `MacTitleBarHiddenInset` (smaller band, traffic lights sit higher than a 34px field); standard title bar plus a separate search row (wastes 28px and looks dated).

## R2. Linux window chrome

- **Decision**: Drop `Frameless`. GTK4 draws its own client-side title bar, which follows the user's theme and button layout on GNOME and KDE (Breeze GTK). No window menu bar: the application menu is only set on macOS, and Linux windows get `DisableMenu: true`. Add and the main menu (☰) are HTML buttons in Cairn's search bar below the title bar.
- **Rationale**: Wails cannot place widgets in GTK's header bar. Setting an app menu on Linux would add a GTK menu bar, which the GNOME HIG says not to use.
- **Alternatives**: Keep frameless and draw a header bar in HTML (rejected: fake buttons, wrong layout for some users; this is the current complaint).

## R3. Desktop detection and KDE variant

- **Decision**: Go reads `XDG_CURRENT_DESKTOP` and `XDG_SESSION_TYPE` at startup and exposes them through `WindowService.Platform()`. The frontend sets `data-os` (`darwin`/`linux`) and `data-desktop` (`gnome`/`kde`/`other`) on `<html>`; CSS handles radius, accent, labelled Add and footer visibility.
- **Rationale**: Pure CSS switches; one frontend build.

## R4. App presence: Dock, menus, reopen

- **Decision**: `Mac.ActivationPolicy = ActivationPolicyRegular` (Dock icon). Build the application menu with `app.NewMenu()` using Wails roles for the app, Edit and Window menus and custom File and Help items. Listen for `events.Mac.ApplicationShouldHandleReopen` to show the window when the Dock icon is clicked with no window visible. Closing the window hides it (the existing `WindowClosing` hook) on macOS.
- **Rationale**: All available in Wails; no native code.

## R5. Close and quit lifecycle on Linux; tray detection

- **Decision**: On Linux, the `WindowClosing` hook hides the window only when `runInBackground && trayAvailable`; otherwise it calls `app.Quit()`. Tray availability = the D-Bus name `org.kde.StatusNotifierWatcher` has an owner on the session bus, checked with `github.com/godbus/dbus/v5` (already in the module graph as an indirect dependency; becomes direct). Watch `NameOwnerChanged` to update availability live and notify the frontend.
- **Rationale**: Stock GNOME has no tray; hiding there would leave an unreachable app.
- **Alternatives**: Always hide (unreachable app on GNOME); always quit (loses shortcut and sync on KDE).

## R6. Quit while syncing

- **Decision**: Add `Syncer.Wait(timeout)` that blocks until in-flight push/pull goroutines finish or the timeout (5s) passes. Quit paths emit a `status` event "Finishing sync…" when a sync is running, then call `Wait` in `OnShutdown`.
- **Rationale**: Today shutdown can interrupt a push; the sync lock is reclaimed as stale, but finishing is cleaner.

## R7. Global shortcut behaviour

- **Decision**: Keep `app.GlobalShortcut.Register` (Wails already uses the GlobalShortcuts portal on Wayland and direct registration on X11 and macOS). Change the handler: if the window is visible and focused → macOS `app.Hide()`; Linux `window.Minimise()`, or hide when background mode is on. Otherwise → show, restore, focus, and emit `launcher:show`. Record the registration error and expose it via `WindowService.ShortcutStatus()`.
- **Rationale**: Matches the design's behaviour table; makes failures visible in Settings.

## R8. Raising the window from `--toggle` on Wayland

- **Decision**: Pass the second instance's `XDG_ACTIVATION_TOKEN` (and `DESKTOP_STARTUP_ID` for X11) through `SingleInstanceOptions.AdditionalData`, which each process builds from its own environment at startup. In the first instance, a small Linux-only cgo helper (`activation_linux.go`) calls `gtk_window_set_startup_id(window, token)` before `Show()/Focus()`, so GTK's `gtk_window_present` uses the token.
- **Rationale**: Wails exposes no API for this. The desktop module already requires cgo and GTK4 headers, so a ~20-line helper adds no new build requirement.
- **Risk**: Getting the `GtkWindow*` from a Wails window may need `NativeWindow()`; if it isn't reachable, fall back to the current behaviour (GNOME shows "Cairn is ready") and document it. Spike this first.

## R9. Open at login

- **Decision**: Use Wails' built-in `app.Autostart` manager (`Enable(AutostartOptions{Identifier, Arguments: ["--background"]})`, `Disable()`, `Status()`). Identifier `io.github.ndy40.cairn.desktop` on macOS (same as the installer's LaunchAgent label) and `cairn-desktop` on Linux (same file name as the installer's autostart entry). Settings reads `Status()` live, so existing installer-created entries show as on with no migration code. `install.sh` stops writing login items; `--no-autostart` is kept as an accepted no-op for scripts. The app keeps accepting `--hidden` as an alias for `--background`.
- **Rationale**: Wails already implements SMAppService (macOS 13+ bundles), LaunchAgent fallback and XDG autostart, with atomic writes and stale-entry cleanup. Matching identifiers delivers FR-042's "existing users keep it" for free.
- **Alternatives**: Hand-written plist and `.desktop` writers (duplicates Wails).
- **Check**: On macOS 13+ Wails may prefer SMAppService and not see the old plist. In that case, `Status()` falling back to the plist path must be verified; if it doesn't, add a one-time check for `~/Library/LaunchAgents/io.github.ndy40.cairn.desktop.plist`.

## R10. Window size and position persistence

- **Decision**: Extend `windowState` with `x`, `y` and `screen` (screen ID/name). Save on resize and move (debounced, as today) and on shutdown. Restore position only when not on Wayland and when the saved screen still exists; otherwise centre on first launch, or let the compositor place the window on Wayland. Defaults 720×520, minimum 520×360; clamp to the current screen's work area.
- **Rationale**: Keeps the existing file and its fallback logic; older files without `x/y` still load.

## R11. Settings storage

- **Decision**: New `desktop-settings.json` next to `desktop-state.json` (both next to `cairn.json`): `showInMenuBar`, `runInBackground` (nullable: unset means the per-desktop default, on for KDE), `hideAfterOpen`. Open-at-login is not stored; it is read from `app.Autostart.Status()`.
- **Rationale**: FR-043 (never rewrite `cairn.json`); separates user choices from window geometry.

## R12. System accent colour

- **Decision**: macOS: CSS system colour `AccentColor` inside `@supports (color: AccentColor)`, with the fallback tokens otherwise. Linux: read `org.freedesktop.appearance accent-color` from the Settings portal via godbus (pure Go) and pass it to the frontend as a CSS variable; ignore it when absent. Darken the colour in the frontend if white text would fail 4.5:1.
- **Rationale**: No native code; graceful fallback (spec assumption).

## R13. Add Bookmark modal

- **Decision**: An HTML `<dialog>` element shown with `showModal()`. On macOS it is styled as a sheet that slides down from the title band; on Linux it is a centred libadwaita-style dialog over a dimmed backdrop. Same form logic as today (`BookmarkService.Add`, `ClipboardURL`).
- **Rationale**: `<dialog>` gives focus trapping, Esc handling and inert background for free and stays inside the window (no Wayland placement problem).

## R14. Pin and Copy Link

- **Decision**: Add `BookmarkService.SetPinned(id, pinned)` (wraps `store.SetPermanent`, same as `cairn pin`) and `BookmarkService.CopyLink(id)` (writes the URL with `app.Clipboard.SetText`). Row context menu = HTML menu (keeps platform styling consistent and works on Linux, where Wails' context menus are GTK-styled anyway). The menu also appears through `contextmenu` and Shift+F10.
- **Note**: `cairn pin` does not queue a sync change today; desktop pin keeps parity, and syncing pin state is out of scope.

## R15. Removing the focus guard

- **Decision**: Delete `focus.go`, `focus_test.go`, `WindowService.PointerDown` and the frontend pointer-down call; drop the `WindowLostFocus` hide logic.
- **Rationale**: FR-009. Removes the GNOME/Wayland workaround.

## R16. Testing approach

- **Decision**: Unit tests in Go for pure logic: window-state load/save/clamp, settings load/save/defaults, close-decision function (`closeAction(os, runInBackground, trayAvailable)`), shortcut toggle decision, platform detection from env, `AdditionalData` token extraction, service `SetPinned`/`CopyLink`. The frontend is type-checked (`tsc`), and the manual QA matrix in `quickstart.md` covers macOS, GNOME Wayland, GNOME X11 and KDE.
- **Rationale**: Matches existing `desktop/*_test.go` style; the window manager behaviour can only be checked by hand.
