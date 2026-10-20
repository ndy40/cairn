---
id: TASK-143
title: Cross-platform desktop quick-launcher (cairn-desktop)
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-20 11:10'
updated_date: '2026-10-20 12:23'
labels:
  - desktop
  - feature
dependencies: []
priority: medium
ordinal: 106000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Cairn's quick-access UI is the Linux-only Vicinae extension. Provide a Wails v3 tray app with a Spotlight-style popup to search, add and open bookmarks on macOS and Linux (Windows best-effort), sharing the CLI's database, search ranking and sync.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Popup can be toggled via a global hotkey (macOS/X11) or by running 'cairn-desktop --toggle' (Wayland)
- [x] #2 Typing in the popup fuzzy-searches bookmarks with the same ranking as the TUI; Enter opens the selected link in the default browser
- [x] #3 User can add a URL (pre-filled from clipboard) with tags; title/description are fetched and duplicates are rejected
- [ ] #4 Bookmarks added via CLI appear in the popup on next show, and vice versa
- [ ] #5 Writes trigger background sync push when sync is configured
- [x] #6 CLI still builds with CGO_ENABLED=0
- [ ] #7 macOS and Linux desktop artifacts are produced by the release workflow
- [x] #8 User guide docs (docs/content/docs) cover installing and using cairn-desktop: hotkey/--toggle setup per platform, search/add/open, desktop_hotkey config, tray caveats
- [ ] #9 Re-running install.sh --with-desktop updates the desktop app while it is running (Linux and macOS) and restarts it if it was running
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Full plan: ~/.claude/plans/i-would-like-to-noble-gizmo.md

1. Slice 1: desktop/ Wails v3 (v3.0.0-beta.33) app; popup window + Search/Open service over internal/store + internal/search
2. Slice 2: Add (fetcher.Fetch + store.Insert, ExistsByURL, clipboard prefill)
3. Slice 3: tray, SingleInstance --toggle, global shortcut (X11/Wayland portal/macOS)
4. Slice 4: AutoPull on show / AutoPush after writes
5. Slice 5: justfile targets, release.yml desktop jobs, install.sh --desktop
6. Slice 6: user guide docs: new docs/content/docs/desktop-app.md; update installation.md, configuration.md (desktop_hotkey), quickstart/_index links, README

Linux build uses the Wails default GTK4 + WebKitGTK 6.0 backend (no build tags).
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- Tooling installed: backlog 1.54.0, wails3 v3.0.0-beta.33 built with -tags gtk3; wails3 doctor passes

- Switched to GTK4: installed gtk4-devel + webkitgtk6.0-devel; rebuilt wails3 without -tags gtk3; doctor passes

- Slice 1 (search + open) implemented in desktop/ as a nested Go module (github.com/ndy40/cairn/desktop, replace ../). Reasons: Wails beta.33 needs Go 1.26 (root/CI on 1.25) and CGO; nesting keeps root go test/lint/CGO_ENABLED=0 builds untouched. Imports internal/* (allowed by path prefix, same as gopls/x-tools).
- BookmarkService.Search reuses store.List + search.Search (pinned-first on empty query, cap 50); Open uses app.Browser.OpenURL via injected func for testability.
- Frontend: vanilla TS + Vite, runtime/vite/ts pinned; all page-derived text rendered via textContent.
- Root .gitignore ignores build/ and dist/: desktop/.gitignore re-includes them; dist/.gitkeep is kept alive via public/.gitkeep (Vite empties dist).
- Build assets regenerated from build/config.yml; removed unused Windows MSIX templates.
- Verified: desktop go test/vet pass; wails3 build OK; frontend tsc clean; root go test ./... and CGO_ENABLED=0 CLI build still pass; binary runs 8s on GNOME Wayland without errors. UI not yet visually verified.

- Slice 2 (add) implemented: BookmarkService.Add(url, tags) mirrors `cairn add` (fetcher.Fetch then store.Insert; fetch failure still saves with hostname title, reported via AddResult.TitleFetched). Duplicates rejected via ExistsByURL before any network call (plus ErrDuplicate fallback for races). Light URL validation: trim, default https://, require http(s)+host; URL stored as typed (not re-encoded) so dedupe matches CLI-stored URLs.
- ClipboardURL() pre-fills the add form only when clipboard holds an http(s) URL.
- UI: Ctrl/Cmd+N or typing "+<url>" in search opens the add form; Enter saves, Esc cancels; status line shows Saved / Saved (title unavailable) / errors. Go errors surface as RuntimeError.message.
- Verified: 11 service tests pass (add/dedupe/invalid/fetch-failure/clipboard), tsc clean, wails3 build OK, CLI CGO_ENABLED=0 build OK, app smoke-runs. Awaiting user UI check before ticking AC #3.

- Slice 3 (summoning) implemented in desktop/launcher.go: frameless, always-on-top, hidden-at-start popup (HideOnFocusLost; close = hide). SingleInstance (io.github.ndy40.cairn.desktop) forwards later launches: --toggle (default) toggles, --add opens the add form; --hidden starts in the tray. Tray menu: Search / Add / Quit, click toggles. macOS runs as accessory (menu-bar) app.
- Global hotkey via Wails GlobalShortcut: X11 XGrabKey, Wayland xdg-desktop-portal GlobalShortcuts, macOS native. Default CmdOrCtrl+Shift+Space; configurable with new `desktop_hotkey` key in cairn.json / CAIRN_DESKTOP_HOTKEY ("off" disables). Registration failure is logged, not fatal.
- Show event buffered until WindowRuntimeReady so a first launch with --add is not lost. Frontend re-queries on every show (covers AC #4 freshness), Esc on empty search hides, opening a link hides.
- Verified: config tests (+3) and desktop tests pass; root tests + static CLI build pass; second-instance --toggle hands off and exits 0 in ~70ms while first instance keeps running. Not yet verified by user: popup toggling visually, portal hotkey prompt on GNOME Wayland, tray icon (GNOME needs AppIndicator extension).

- Added AC #8 + slice 6 (user guide docs) at user request.
- Slice 4 (sync) implemented: desktop/sync.go Syncer runs engine.AutoPull/AutoPush in-process (no dependency on `cairn` being on PATH, e.g. inside a macOS .app). Pull on startup/--hidden and on every popup show, throttled to once per 2 min; push after each successful Add. Sync config re-read each time so `cairn sync setup/unlink` apply without restart. Successful pull emits bookmarks:changed (frontend re-runs current search); warnings (e.g. expired auth) emit sync:warning to the status line.
- Refactor: moved the CLI sync lockfile helpers (acquireSyncLock/cleanupStaleLock/readLockPID/processAlive) from cmd/cairn into internal/sync (AcquireLock) unchanged, so CLI and desktop share one lock implementation and never overlap syncs. Added internal/sync/lock_test.go (exclusive, independent kinds, dead-PID + TTL reclaim).
- Verified: root go test + golangci-lint (0 issues) + gofmt + static CLI build; cross-compiles for linux/darwin/windows; desktop tests (+2) and lint pass; app smoke-runs. Not verified end-to-end: sync is not configured on this machine.

- Slice 5 (build & distribution) implemented:
  - release.yml: build-desktop-linux (native ubuntu-24.04 + ubuntu-24.04-arm runners, GTK4/WebKitGTK 6.0) -> cairn-desktop-linux-{amd64,arm64}.tar.gz (binary + icon); build-desktop-macos (macos-15, wails3 task darwin:package:universal, ad-hoc signed, not notarised) -> cairn-desktop-darwin-universal.zip. Release tag written into build/config.yml + update:build-assets so Info.plist carries the version. wails3 CLI version read from desktop/go.mod. Release job now needs both; artifacts flow into existing checksums.txt.
  - build.yml: new `desktop` PR job (go test, golangci-lint, wails3 build, tsc), since root jobs do not reach the nested module.
  - justfile: build-desktop / dev-desktop / test-desktop / package-desktop; Makefile: desktop / desktop-test.
  - install.sh: --with-desktop (or interactive prompt), --no-autostart. Linux: binary to INSTALL_DIR, .desktop entry + icon under XDG_DATA_HOME, autostart entry (--hidden), warns if WebKitGTK 6.0 runtime missing. macOS: ~/Applications/cairn-desktop.app + LaunchAgent (--hidden). Checksum mismatch aborts (exit 3) like the CLI path; other desktop failures skip only the desktop app.
- Verified: workflows parse as YAML; CI version lookup + sed tested locally; installer Linux path tested offline against a fake release (happy path, tampered archive -> exit 3, missing checksum -> skip, --no-autostart); generated .desktop passes desktop-file-validate. Not verified: workflows not run on GitHub; macOS installer path untested (no Mac).
- Pre-existing bug spotted (not fixed, out of scope): install.sh downloads `cairn-vicinae-extension.tar.gz` but release.yml publishes `vicinae-extension-<tag>.tar.gz`.

- Created TASK-144 (connect Dropbox from the desktop app) as a follow-up, per user.
- Slice 6 (user guide) done: new docs/content/docs/desktop-app.md (install, requirements, opening the popup incl. GNOME/Wayland shortcut setup + tray caveat, search/add keys, sync, desktop_hotkey, flags, updating, uninstall, build from source). Updated configuration.md (desktop_hotkey), installation.md (--with-desktop/--no-autostart), quickstart next steps, sync.md (desktop auto-sync + shared lock), cli-reference.md (`cairn update` does not cover the desktop app), architecture.md (desktop module + in-process sync), docs home, README (also added missing disable_auto_archive row).
- Verified: Hugo extended v0.168 build with the pinned hugo-book theme (97d7ee4) succeeds with no errors; all relrefs resolve and linked anchors (#configuration, #sync, #updating) exist in rendered output.

- User feedback fixes: (1) popup now resizable (DisableResize removed, MinWidth 480 / MinHeight 300); on Linux frameless windows the Wails runtime provides 5px edge resize handles, macOS resizes natively; size persists for the session. (2) Footer layout bug: a long "Saved: <title>" status squeezed the key hints so labels wrapped under their keys. Hints are now non-shrinking/nowrap and the status flexes, truncates with an ellipsis and shows full text as a tooltip. Reproduced and verified by rendering the popup markup + stylesheet in headless Firefox at 680x440 and 480x300. Docs mention resizing.

- User request: remember size + allow dragging. Size saved to desktop-state.json beside cairn.json (not cairn.json, so the app never rewrites user-edited config): debounced 500ms after WindowDidResize (mapped on Linux + macOS; Linux has no hide event) and on shutdown; atomic write; sizes below the 480x300 minimum (e.g. 0x0 while hidden) ignored on save and load. Dragging via --wails-draggable on empty results area, footer and add form; result rows and form fields opt out.
- Verified: 4 new windowstate tests; lint/vet/tsc clean; with an isolated config dir the app writes 680x440 on first show, and reopens at a seeded 820x520 (re-saved unchanged). Not verified: actual mouse drag/resize on GNOME Wayland (possible risk: HideOnFocusLost firing during a compositor move/resize grab).

- User confirmed dragging the popup works on their desktop (GNOME Wayland).

- Bug (user, GNOME 51 Wayland): any click / press-and-hold / edge resize hid the popup. Cause: Wails HideOnFocusLost hides immediately on GTK focus-leave, and pressing inside the window (compositor grab for move/resize, and clicks) produces focus-leave. Fix: dropped HideOnFocusLost; desktop/focus.go focusGuard + WindowService.PointerDown (frontend reports every pointerdown, capture phase). On focus loss, wait 250ms then hide only if visible, not gtk_window_is_active, and no press inside the popup in the last 600ms; still skipped on focus-follows-mouse WMs (Hyprland/Sway/i3), matching Wails. CAIRN_DESKTOP_DEBUG=1 logs focus events and decisions.
- Verified: 4 focus tests (fake clock), build/tsc/lint clean. Not yet verified on the user's desktop. Known trade-off: if GNOME does not hand focus back after a drag/resize, the popup stays open until toggled or Esc.

- User confirmed resizing works after the focus fix.
- Committed on branch feat/desktop-app (not pushed): dd7bd83 refactor(sync), 6e1386d feat(config), 53aeb45 feat(desktop), ef3021a ci, 21cbad2 docs.

- Installer update fix (AC #9): install_binary now copies to a temp file in the target dir and renames it over the old binary, so updating works while the binary is running (in-place cp failed with "Text file busy" and was misreported as "Permission denied", exit 4). Exit 4 now only when the dir is not writable; other failures show the real error (exit 1). Applies to the CLI too.
- Desktop update: if cairn-desktop is running (pgrep -x -u $UID), the installer stops it after installing and starts the new version with --hidden (Linux: only if a graphical session is present). macOS: new .app unpacked to a temp dir before the old one is removed; app restarted with `open -g ... --args --hidden`. Under sudo the user's app isn't visible, so no restart (documented).
- Verified offline with pgrep/pkill shims that only see a fake old process (never the user's real app) and a stub new binary: replaced while running, old stopped, new started with --hidden; fresh install does not start it; read-only dir -> exit 4; running CLI binary replaced. macOS path untested (no Mac).

- AC #9 left unchecked until the macOS update path is tested on a Mac; Linux path verified.
<!-- SECTION:NOTES:END -->
