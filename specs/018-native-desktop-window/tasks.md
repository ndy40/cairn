---
description: "Task list for 018-native-desktop-window"
---

# Tasks: cairn-desktop as a native app window

**Input**: Design documents from `specs/018-native-desktop-window/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Unit tests are included for the pure decision logic only (research R16), following the existing `desktop/*_test.go` style. Window-manager behaviour is checked with the manual QA matrix in `quickstart.md`.

**Organization**: Tasks are grouped by user story. All Go paths are in the `desktop/` module (`package main`); run Go commands from `desktop/`.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: User story from spec.md (US1–US6)

---

## Phase 1: Setup

**Purpose**: Dependencies and the two spikes that decide risky parts of the design.

- [ ] T001 Promote `github.com/godbus/dbus/v5` from indirect to a direct dependency in `desktop/go.mod` (`go get github.com/godbus/dbus/v5@v5.2.2 && go mod tidy`), keeping the CLI module untouched
- [ ] T002 Spike the Wayland activation token (research R8): add `desktop/activation_linux.go` (cgo, `//go:build linux`, `#cgo pkg-config: gtk4`) with `applyActivationToken(w *application.WebviewWindow, token string)` that calls `gtk_window_set_startup_id` on the window's native `GtkWindow*` (via `w.NativeWindow()`), and `desktop/activation_other.go` (`//go:build !linux`) as a no-op; on GNOME Wayland, bind a custom shortcut to `cairn-desktop --toggle` and confirm the window is raised and focused; record the outcome (works / falls back) under R8 in `specs/018-native-desktop-window/research.md`
- [ ] T003 Spike open-at-login detection (research R9): on macOS 13+ with the installer's `~/Library/LaunchAgents/io.github.ndy40.cairn.desktop.plist` present, and on Linux with `~/.config/autostart/cairn-desktop.desktop` present, confirm `app.Autostart.Status()` with identifiers `io.github.ndy40.cairn.desktop` / `cairn-desktop` reports `Enabled: true`; record the outcome under R9 in `specs/018-native-desktop-window/research.md` and, if macOS does not detect it, note that T037 must add the one-time plist check

---

## Phase 2: Foundational (blocking prerequisites)

**Purpose**: State, settings, platform detection and the WindowService binding surface that every story uses.

- [ ] T004 [P] Create `desktop/platform.go` with `type PlatformInfo struct{ OS, Desktop string; Wayland, TrayAvailable bool; AccentColor string }` (JSON tags per `data-model.md`) and a pure `detectPlatform(goos string, getenv func(string) string) PlatformInfo` mapping `XDG_CURRENT_DESKTOP` (case-insensitive, colon-separated) to `gnome`/`kde`/`other` and `XDG_SESSION_TYPE=wayland` to `Wayland`; macOS sets `TrayAvailable=true`; add table tests in `desktop/platform_test.go`
- [ ] T005 [P] Create `desktop/settings.go` with `DesktopSettings{ShowInMenuBar bool; RunInBackground *bool; HideAfterOpen bool}`, `defaultSettingsPath()` (`desktop-settings.json` next to `cairn.json`, mirroring `defaultStatePath()`), `loadSettings(path)` (missing/invalid file → defaults, invalid JSON logged), atomic `save(path)` (temp file + rename), and `resolvedRunInBackground(desktop string) bool` (nil → true on `kde`, false elsewhere); tests in `desktop/settings_test.go`
- [ ] T006 [P] Extend `desktop/windowstate.go`: defaults 720×520, minimum 520×360; add optional `X, Y *int` and `Screen string` fields; add pure `restorePlan(state windowState, wayland bool, screens []screenRect) (w, h int, x, y int, hasPos bool)` that clamps size to the target screen work area, restores position only when not Wayland and the saved screen exists, and otherwise reports no position; keep old `{width,height}` files loading; extend `desktop/windowstate_test.go` (old file, missing screen, Wayland, clamp)
- [ ] T007 Create `desktop/window_service.go`: move `WindowService` out of `focus.go`, give it fields for the app, launcher, settings path/settings and platform info, and add `Platform() PlatformInfo`, `Settings() DesktopSettings` (with `runInBackground` resolved) and `SetSetting(key string, value bool) error` (persist only; side effects added in US3/US5); register it in `desktop/main.go`; run `wails3 generate bindings` so `frontend/bindings/` includes the new methods
- [ ] T008 At startup in `desktop/frontend/src/main.ts`, call `WindowService.Platform()` and set `data-os` and `data-desktop` on `document.documentElement`; in `desktop/frontend/src/style.css` replace the `:root` palette with the token set from the design (`--accent`, `--text`, `--text-2`, `--list`, `--hover`, `--warn-bg`, `--warn-fg`, light and dark, fallback accent `#1766d1` / `#2f6fd6`)

**Checkpoint**: Settings and state load, platform info reaches the frontend; story work can start.

---

## Phase 3: User Story 1 – Native window (Priority: P1) 🎯 MVP

**Goal**: A framed window with the platform's own controls that can be moved, resized, stays open on focus loss and remembers its size/position.

**Independent Test**: QA matrix rows 1–7 in `quickstart.md` on macOS, GNOME and KDE.

- [ ] T009 [US1] In `desktop/launcher.go` `NewLauncher`, remove `Frameless`, `AlwaysOnTop` and `Windows.HiddenOnTaskbar`; set `MinWidth: 520`, `MinHeight: 360`, `Title: "Cairn"`; add `Mac: application.MacWindow{TitleBar: application.MacTitleBarHiddenInsetUnified, Backdrop: application.MacBackdropTranslucent}` and leave `Linux.Menu` unset so GTK draws its own title bar with no menu bar
- [ ] T010 [US1] Delete `desktop/focus.go` and `desktop/focus_test.go`; remove the `WindowLostFocus`/`WindowFocus` handlers, `hideIfUnfocused`, the `focus` field and `newFocusGuard()` from `desktop/launcher.go` and `desktop/main.go`; move `debugf` into `desktop/launcher.go` if still used; remove the `WindowService.PointerDown()` call from `desktop/frontend/src/main.ts`
- [ ] T011 [US1] In `desktop/launcher.go`, apply `restorePlan` (T006) when creating the window using `app.Screen` data and `PlatformInfo.Wayland`; centre only when there is no saved state; save size, position and screen (debounced, existing `saveSizeDelay`) on `WindowDidResize` and `WindowDidMove`, skipping position on Wayland; stop calling `window.Center()` in `Show`
- [ ] T012 [US1] Restructure the top of `desktop/frontend/index.html`: a `<header class="titleband">` containing a `<label>`-wrapped `<input type="search" id="query" aria-label="Search bookmarks">`, an Add `<button id="add-btn" aria-label="Add bookmark">` with an inline SVG plus icon, and a main-menu `<button id="menu-btn" aria-label="Main menu" aria-haspopup="menu">` (shown on Linux only); keep `#results` as `role="listbox"`; remove the old placeholder text "(type + to add)"
- [ ] T013 [US1] In `desktop/frontend/src/style.css`: remove the `.launcher` border; add `--titleband-h` (52px macOS, measured on a build) and make `.titleband` `--wails-draggable: drag` with the input and buttons `no-drag`; on `[data-os=darwin]` pad the band left for the traffic lights and give only the band a translucent background while `.results` and the footer are opaque; on `[data-os=linux]` style it as a search bar under the GTK title bar (title-bar fill, hairline below); remove the old drag rules on `.results`, `.hints` and `.add`
- [ ] T014 [US1] In `desktop/frontend/src/main.ts`, change Esc handling: non-empty query → clear and re-search; empty query → do nothing (no hide call); remove any code path that hides the window on Esc

**Checkpoint**: The app is a normal window on every platform. MVP is shippable here with the old list styling.

---

## Phase 4: User Story 2 – Find and open in seconds (Priority: P1)

**Goal**: Shortcut, relaunch, `--toggle`, Dock and tray all bring the window forward with search focused; the shortcut also puts it away.

**Independent Test**: QA matrix rows 8–10 plus "shortcut → 2–3 letters → Enter" timing.

- [ ] T015 [P] [US2] Create `desktop/lifecycle.go` with pure `shortcutAction(os string, visible, focused, background bool) action` returning `actionShow`, `actionHideApp` (macOS), `actionMinimise` or `actionHideWindow` (Linux background mode); table tests in `desktop/lifecycle_test.go`
- [ ] T016 [US2] In `desktop/launcher.go`, replace `Show`/`Toggle` with `BringForward(view string, token string)` (apply activation token via T002's helper, `Restore()` if minimised, `Show()`, `Focus()`, emit `launcher:show` with the pending-view logic kept) and `Toggle(token string)` that uses `shortcutAction` (`app.Hide()` on macOS, `window.Minimise()` or `window.Hide()` on Linux)
- [ ] T017 [US2] In `desktop/main.go`: set `SingleInstanceOptions.AdditionalData` to `{"activation_token": os.Getenv("XDG_ACTIVATION_TOKEN"), "startup_id": os.Getenv("DESKTOP_STARTUP_ID")}`; add pure `activationToken(data map[string]string) string` (prefers activation_token) with a test in `desktop/lifecycle_test.go`; in `Launcher.HandleArgs` route no args → `BringForward("search")`, `--toggle` → `Toggle`, `--add` → `BringForward("add")`, all passing the token
- [ ] T018 [US2] In `desktop/launcher.go` `RegisterHotkey`, store a `ShortcutStatus{Hotkey, Registered, Error, ToggleCommand: "cairn-desktop --toggle"}` (per `data-model.md`) instead of only logging; the hotkey calls `Toggle("")`; expose `ShortcutStatus()` on `WindowService` in `desktop/window_service.go`
- [ ] T019 [US2] In `desktop/frontend/src/main.ts`, on every `launcher:show` focus `#query` and call `select()`; keep arrow keys moving the selection while focus stays in the field; after `BookmarkService.Open` succeeds call `WindowService.AfterOpen()`; add `AfterOpen()` to `desktop/window_service.go` that hides the app (macOS) or window (Linux) only when `hideAfterOpen` is on

**Checkpoint**: Journey 1 works end to end in the new window.

---

## Phase 5: User Story 3 – Close, quit and background (Priority: P2)

**Goal**: Predictable close/quit per platform, Dock reopen, tray-aware background mode, clean quit during sync.

**Independent Test**: QA matrix rows 11–14.

- [ ] T020 [P] [US3] In `desktop/platform.go`, add `watchTray(ctx, onChange func(bool))` using godbus: `NameHasOwner("org.kde.StatusNotifierWatcher")` on the session bus at start, then subscribe to `NameOwnerChanged` for that name; on change update `PlatformInfo.TrayAvailable` and emit `platform:changed`; no-op on macOS; a missing session bus means `false`
- [ ] T021 [P] [US3] In `desktop/lifecycle.go`, add pure `closeAction(os string, runInBackground, trayAvailable bool) action` (macOS → hide window; Linux → hide only when both true, else quit) with table tests in `desktop/lifecycle_test.go`
- [ ] T022 [US3] In `desktop/sync.go`, add `Running() bool` and `Wait(timeout time.Duration) bool` (track in-flight runs with a `sync.WaitGroup` or counter) and emit a new `sync:ok` event after a successful run; extend `desktop/sync_test.go`
- [ ] T023 [US3] In `desktop/launcher.go`, change the `WindowClosing` hook to use `closeAction` (hide + cancel, or cancel and call the quit path); in `desktop/window_service.go` add `CloseWindow()` (same path) and `Quit()` which, if `Syncer.Running()`, emits `status` "Finishing sync…", then calls `app.Quit()`; in `desktop/main.go` `OnShutdown` call `syncer.Wait(5*time.Second)` before saving window state
- [ ] T024 [US3] In `desktop/main.go`, set `Mac.ActivationPolicy` to `application.ActivationPolicyRegular` and register `events.Mac.ApplicationShouldHandleReopen` to call `launcher.BringForward("search", "")` when the window is not visible
- [ ] T025 [US3] In `desktop/main.go` and `desktop/launcher.go`, add the `--background` flag (keep `--hidden` as an alias) that starts without a window and pulls sync on start; make the Linux tray only appear when `TrayAvailable`; make `SetSetting("runInBackground", true)` return an error when no tray is available

**Checkpoint**: Window and app lifecycle match the design's behaviour table.

---

## Phase 6: User Story 4 – Add bookmark sheet/dialog (Priority: P2)

**Goal**: Add Bookmark as a modal attached to the window, with clipboard prefill and duplicate handling.

**Independent Test**: QA matrix rows 15–17.

- [ ] T026 [P] [US4] In `desktop/service.go`, add `FindByURL(rawURL string) (*Bookmark, error)` that normalises like `Add` and returns the matching bookmark or nil; test in `desktop/service_test.go`; regenerate bindings
- [ ] T027 [US4] In `desktop/frontend/index.html`, replace the `#add-view` form with `<dialog id="add-dialog">` containing a heading "Add Bookmark", a labelled URL input, a labelled Tags input ("comma separated"), helper text "Title and description are fetched automatically.", and Cancel/Add buttons (Add is `type="submit"`, default)
- [ ] T028 [US4] In `desktop/frontend/src/style.css`, style `#add-dialog` as a macOS sheet under the title band (`[data-os=darwin]`, slide-down, attached to the band, list visible behind) and as a centred libadwaita-style dialog with a dimmed `::backdrop` and Cancel/Add in its header on Linux
- [ ] T029 [US4] In `desktop/frontend/src/main.ts`, replace the add view swap with `openAdd(prefill?)` using `showModal()`: prefill from `ClipboardURL()` and select it, or from `+text` typed in search with the cursor at the end; Enter submits, Esc/Cancel closes and refocuses search; on success close, refresh, select the new row and set status "Added. Fetching title…"; on `ErrAlreadyBookmarked` close, `FindByURL`, select that row and show an "Already bookmarked" toast (add a minimal `showToast(text)` helper, styled in US6); wire the Add button, `Ctrl/Cmd+N` and `launcher:show "add"`

**Checkpoint**: Journey 3 and 4 work in the new window.

---

## Phase 7: User Story 5 – Menus, settings, shortcuts (Priority: P2)

**Goal**: Settings, Quit, shortcuts and row actions are where each platform puts them.

**Independent Test**: QA matrix rows 18–21 and Story 5 acceptance scenarios.

- [ ] T030 [P] [US5] In `desktop/service.go`, add `SetPinned(id int64, pinned bool) error` (wraps `store.SetPermanent`, emits `bookmarks:changed` "pin" via an injected callback), `CopyLink(id int64) error` (injected clipboard writer using `app.Clipboard.SetText`) and `SyncNow()` (calls a new `Syncer.PullNow()` that skips the throttle, then `Push`, in `desktop/sync.go`); tests in `desktop/service_test.go`; regenerate bindings
- [ ] T031 [US5] Create `desktop/menus.go` with `buildAppMenu(app, launcher, svc)` (darwin only): Cairn menu (About Cairn, Settings… `CmdOrCtrl+,`, Services role, Hide/Hide Others/Show All roles, Quit `CmdOrCtrl+Q`), File (New Bookmark… `CmdOrCtrl+N`, Open Link, Copy Link `Shift+CmdOrCtrl+C`, Pin, Sync Now, Close Window `CmdOrCtrl+W`), Edit role plus Find Bookmarks `CmdOrCtrl+F`, Window role, Help (Cairn Help → docs URL, Keyboard Shortcuts `CmdOrCtrl+/`); items emit `menu:action` / `ui:open` or call `WindowService` methods; set it with `app.Menu.Set` from `desktop/main.go`
- [ ] T032 [US5] Move tray setup from `desktop/launcher.go` into `desktop/menus.go` `buildTray`: items Search Bookmarks…, Add Bookmark…, a disabled sync status line, Settings…, Quit Cairn (via `WindowService.Quit`); macOS uses the template icon and respects `showInMenuBar` (`tray.Show()/Hide()`); add a Dock menu with New Bookmark if Wails exposes one, otherwise note it in `research.md`
- [ ] T033 [US5] In `desktop/window_service.go`, add `OpenAtLogin() (bool, error)` and `SetOpenAtLogin(enabled bool) error` using `app.Autostart` with identifier `io.github.ndy40.cairn.desktop` (macOS) / `cairn-desktop` (Linux) and arguments `--background` (plus the one-time plist check if T003 found it necessary); make `SetSetting("showInMenuBar", …)` show/hide the tray immediately
- [ ] T034 [US5] In `desktop/frontend/index.html` and `desktop/frontend/src/main.ts`, add the Linux main menu popover for `#menu-btn` (`role="menu"`, items Add Bookmark…, Sync Now, Preferences, Keyboard Shortcuts, About Cairn, Quit) with arrow-key navigation and Esc to close
- [ ] T035 [US5] In `desktop/frontend/index.html` and `desktop/frontend/src/main.ts`, add the row context menu (`contextmenu` event and Shift+F10 / Menu key) with Open Link, Copy Link and Pin/Unpin (label follows `pinned`), and handle `menu:action` (`open`, `copy`, `pin`, `find`) against the selected row
- [ ] T036 [US5] Add a Settings `<dialog id="settings-dialog">` in `desktop/frontend/index.html` and its logic in `desktop/frontend/src/main.ts`: switches (`role="switch"`) for Open Cairn at login, Show in menu bar (macOS only), Run in background (Linux only; disabled with "Requires a system tray…" when `!trayAvailable`), Hide Cairn after opening a link; a shortcut row that shows the accelerator when registered, or the error and the `cairn-desktop --toggle` command with a Copy Command button when not; the closing note about Quit on macOS; open it on `ui:open "settings"` and `Cmd/Ctrl+,`; style per the Settings artboard in `desktop/frontend/src/style.css`
- [ ] T037 [US5] Add Keyboard Shortcuts and About `<dialog>`s in `desktop/frontend/index.html` and `desktop/frontend/src/main.ts` (shortcuts list per FR-028 with platform-correct key names; About shows name and version from a new `WindowService.Version()` reading the build version); open on `ui:open` and `Cmd+/` / `Ctrl+?`
- [ ] T038 [US5] Implement the FR-028 keyboard map in `desktop/frontend/src/main.ts`: `Cmd/Ctrl+F` focus+select search, `Shift+Cmd/Ctrl+C` copy link, `Cmd/Ctrl+N` add, `Cmd/Ctrl+W` `CloseWindow()`, `Cmd/Ctrl+Q` `Quit()`, `Cmd/Ctrl+,` settings, `Cmd+/` / `Ctrl+?` shortcuts; on macOS skip keys the app menu already handles so actions do not fire twice
- [ ] T039 [P] [US5] In `install.sh`: stop writing `~/.config/autostart/cairn-desktop.desktop` and `~/Library/LaunchAgents/io.github.ndy40.cairn.desktop.plist` for new installs (leave existing ones untouched); keep `--no-autostart` accepted as a no-op that prints "Start at login is now set in Cairn's Settings"; restart a running app with `--background` instead of `--hidden`; update the usage text

**Checkpoint**: Quit, settings and start-at-login are discoverable on every platform.

---

## Phase 8: User Story 6 – Refreshed contents (Priority: P3)

**Goal**: Rows, sections, messages and colours per the Shared parts artboard.

**Independent Test**: Compare light/dark on each platform against the design; QA matrix rows 22–24.

- [ ] T040 [US6] Rewrite `render()` in `desktop/frontend/src/main.ts`: empty query → "Pinned" and "Recent" section headings (non-selectable `role="presentation"` items) with selection skipping them; rows with a 28px monogram tile (first letter of the domain, tint from a stable hash of the domain over a fixed accessible palette), title with case-insensitive query-term matches wrapped in `<b>` built with `textContent` only, domain, first 3 tags then "+N", an SVG pin glyph with `aria-label="Pinned"`, and an "↩ Open" hint on the selected row; empty states "No bookmarks yet" + add hint, and `No bookmarks match "<query>"` + "Type + before a URL to add it."; keep the selection by id across refreshes
- [ ] T041 [US6] Style rows, sections, states (default, hover, selected, window inactive via `blur`/`focus` on `window`) and empty states in `desktop/frontend/src/style.css` per the Components artboard: 52px rows, 7px radius, 8px inset; selection = accent fill with white text on macOS, accent tint (16% GNOME, 18% + 1px outline KDE) with dark text on Linux
- [ ] T042 [US6] Add toast and banner elements to `desktop/frontend/index.html` (`role="status"`) and helpers in `desktop/frontend/src/main.ts`: `showToast(text, action?)` for ~4s; `showBanner(text)` with Dismiss that persists until dismissed or `sync:ok`; route `sync:warning` → banner, `sync:ok` → clear, `status` → footer status area on macOS/KDE and a toast on GNOME; style them in `desktop/frontend/src/style.css` with the warning icon and `--warn-*` tokens
- [ ] T043 [US6] In `desktop/menus.go`, reflect sync state on the tray/menu bar icon: the status line shows "Synced N minutes ago" or the warning text, and the macOS menu bar icon gets a warning label/badge (`SetLabel` or an alternate template icon) on `sync:warning`, cleared on `sync:ok`
- [ ] T044 [US6] Accent colour: in `desktop/frontend/src/style.css` use `AccentColor` under `@supports (color: AccentColor)` on macOS; in `desktop/platform.go` read `org.freedesktop.appearance accent-color` from the Settings portal via godbus (RGB 0–1 tuple → hex, ignore when absent or out of range) into `PlatformInfo.AccentColor`; in `desktop/frontend/src/main.ts` apply it to `--accent` and darken it until white text reaches 4.5:1 (small contrast helper)
- [ ] T045 [US6] Desktop variants in `desktop/frontend/src/style.css`: `[data-desktop=kde]` 4px radius, Breeze accent fallback `#3daee9`, labelled Add button, footer shown; `[data-desktop=gnome]` footer hidden; footer with key hints and result count kept on macOS; system-ui font stack and the sizes from the Components artboard (search 16/15/14px by platform)

**Checkpoint**: Visual refresh complete.

---

## Phase 9: Polish & cross-cutting

- [ ] T046 [P] Rewrite `docs/content/docs/desktop-app.md` for the window app: opening, window behaviour, close/quit per platform, menus, Add Bookmark, Settings (start at login, menu bar, background, hide after open), global shortcut and the Wayland fallback, keyboard map, `--background`/`--hidden`, `desktop-settings.json`, updated uninstall steps
- [ ] T047 [P] Update the desktop app description in `README.md` ("a hotkey popup" → a native desktop app with a global shortcut) and the `--with-desktop`/`--no-autostart` notes
- [ ] T048 Remove dead code and constants left from the popup in `desktop/launcher.go`, `desktop/main.go` and `desktop/frontend/src/main.ts` (old `viewAdd` view swap, `#hints-add`, focus-guard leftovers); regenerate bindings
- [ ] T049 Run `go test ./...`, the desktop lint job, `npx tsc --noEmit` in `desktop/frontend` and `wails3 build` (or `just test-desktop` / `just build-desktop`); fix failures
- [ ] T050 Run the manual QA matrix in `specs/018-native-desktop-window/quickstart.md` on macOS, GNOME Wayland, GNOME X11 and KDE; record results and any platform limitations in the task notes

---

## Dependencies & execution order

- **Setup (T001–T003)**: no dependencies. T002 and T003 are spikes whose outcome adjusts T016–T017 and T033.
- **Foundational (T004–T008)**: after T001. T004–T006 in parallel; T007 after T004–T005; T008 after T007.
- **US1 (T009–T014)**: after Foundational. T009 → T010 → T011 (same file); T012 → T013; T014 any time.
- **US2 (T015–T019)**: after US1 (uses the new launcher). T015 in parallel with US1 work.
- **US3 (T020–T025)**: after US2 (close/quit uses `BringForward`). T020–T022 in parallel.
- **US4 (T026–T029)**: after US1 (needs the new header and `launcher:show`); independent of US2/US3. T026 in parallel.
- **US5 (T030–T039)**: after US3 (Quit path, tray) and US4 (add entry points). T030 and T039 in parallel.
- **US6 (T040–T045)**: after US4 (toast helper) and US5 (tray menu for T043). T044 in parallel with T040–T042.
- **Polish (T046–T050)**: after all stories; T046 and T047 in parallel.

### Story completion order

US1 (MVP) → US2 → US3 → US4 (can run alongside US2/US3) → US5 → US6 → Polish.

## Parallel examples

- Foundational: T004 `platform.go`, T005 `settings.go`, T006 `windowstate.go` together.
- US1 with US2 prep: T012/T013 (frontend) alongside T009–T011 (Go) and T015 (`lifecycle.go`).
- US3: T020 (tray watch), T021 (`closeAction`), T022 (`sync.go`) together.
- US5: T030 (`service.go`) and T039 (`install.sh`) alongside T031 (`menus.go`).

## Implementation strategy

1. **Spikes first** (T002, T003) so the riskiest decisions are settled before building on them.
2. **MVP = Setup + Foundational + US1**: a native window with the old list. This alone answers the main complaint and can ship.
3. Add US2 next to restore the keyboard journey fully, then US3 and US4, then US5 menus and settings, then the US6 visual refresh.
4. Run T049 at the end of each phase; run the full QA matrix (T050) before release.
