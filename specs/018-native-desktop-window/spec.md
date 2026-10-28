# Feature Specification: cairn-desktop as a native app window

**Feature Branch**: `018-native-desktop-window`
**Created**: 2026-10-27
**Status**: Draft
**Input**: User description: "cairn-desktop native app window. Replace the frameless Spotlight-style popup with a regular desktop app that uses native window chrome on macOS 12+ and Linux (GNOME on Wayland/X11, KDE Plasma), per the design canvas. Users complain the current popup has no close/minimise/zoom buttons, no visible resize/move affordance, no obvious quit, doesn't look like a GTK/KDE window, and hides whenever focus moves away."

**Design reference**: [cairn-desktop native window design guide](https://claude.ai/artifact/Qvs5oe46kZEPWXYcJmhrd5) (artboards: Overview, Behaviour, Do and don't, macOS search/dark/add/menus, GNOME search/add, KDE, Shared parts, Settings).

## Overview

cairn-desktop today is a frameless, always-on-top popup that hides whenever another app takes focus. Users expect it to behave like any other app on their platform: a real window frame with the platform's own buttons, a visible way to move, resize, close and quit, and a window that stays where they left it. This feature turns cairn-desktop into a regular app window on macOS and Linux while keeping the core flow (summon, type a few letters, Enter) as fast as it is today.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - A window that looks and behaves like a native app (Priority: P1)

A user opens Cairn and sees a normal window for their platform: on macOS, the red/yellow/green window buttons and the system shadow and rounded corners; on GNOME and KDE, the desktop's own title bar with its close (and, where configured, minimise and maximise) buttons. They can drag the window by its title area, resize it from any edge, maximise or zoom it, and it stays open when they click into another app.

**Why this priority**: This is the complaint that triggered the work. Without it, nothing else in the feature addresses users' expectations.

**Independent Test**: Launch the app on each platform; confirm the native window controls are present and work, the window can be moved and resized, and it remains visible after focusing another app.

**Acceptance Scenarios**:

1. **Given** Cairn is not running, **When** the user launches it on macOS, **Then** a window opens with native close, minimise and zoom buttons inset in a title band that also holds the search field, and the search field has focus.
2. **Given** Cairn is not running, **When** the user launches it on GNOME or KDE, **Then** a window opens with the desktop's own title bar titled "Cairn", following the user's button layout and theme, with Cairn's search bar directly below it and focused.
3. **Given** the window is open, **When** the user drags the title area (anywhere that is not a control), **Then** the window moves; **When** the user drags the search field, **Then** text is selected and the window does not move.
4. **Given** the window is open, **When** the user clicks into another application, **Then** the Cairn window stays open with its query and selection unchanged.
5. **Given** the window is open, **When** the user presses Esc with text in the search field, **Then** the text is cleared; **When** they press Esc again on an empty field, **Then** nothing happens and the window stays open.
6. **Given** the window is open, **When** the user resizes or moves it and later reopens it (including after quitting and relaunching), **Then** it reappears at the same size and, where the platform allows it, the same position.

---

### User Story 2 - Find and open a bookmark in seconds (Priority: P1)

The user is working in another app, presses the global shortcut (or clicks the Dock, taskbar or tray icon), types two or three letters, and presses Enter. The link opens in the default browser.

**Why this priority**: This is the main daily journey. Moving to a regular window must not slow it down.

**Independent Test**: With Cairn running in the background, press the shortcut, type a query and press Enter; time the flow and confirm focus lands in the search field with no extra click.

**Acceptance Scenarios**:

1. **Given** Cairn is running and not in front, **When** the user presses the global shortcut, **Then** the window comes forward (reopening it if it was closed) with focus in the search field and any existing query selected.
2. **Given** the Cairn window is in front, **When** the user presses the global shortcut, **Then** Cairn gets out of the way and focus returns to the previous app.
3. **Given** results are listed, **When** the user presses the arrow keys, **Then** the selection moves while keyboard focus stays in the search field.
4. **Given** a result is selected, **When** the user presses Enter, **Then** the link opens in the default browser and the Cairn window stays where it is (unless the user has turned on "Hide Cairn after opening a link").
5. **Given** Cairn is already running, **When** the user launches it again (from the app menu, Dock, or `cairn-desktop`), **Then** the existing window is raised and no second copy starts.
6. **Given** the user is on GNOME Wayland and has bound a custom shortcut to `cairn-desktop --toggle`, **When** they press it, **Then** the window comes forward and is focused rather than only showing a "Cairn is ready" notification.

---

### User Story 3 - Close, quit and run in the background predictably (Priority: P2)

The user wants to put Cairn away without losing the global shortcut and sync, and also needs an obvious way to quit it completely.

**Why this priority**: The lack of a visible close and quit is a named complaint, and the close behaviour decides whether the shortcut keeps working.

**Independent Test**: Close the window via its close button and via the keyboard on each platform, then use the shortcut and the quit command; confirm the app's state matches the rules below.

**Acceptance Scenarios**:

1. **Given** macOS, **When** the user closes the window, **Then** the window closes, Cairn keeps running (Dock icon and menu bar icon remain) and the shortcut still brings the window back.
2. **Given** macOS and no Cairn window is open, **When** the user clicks the Dock icon, **Then** the window reopens.
3. **Given** Linux with a system tray available and "Run in background" on, **When** the user closes the window, **Then** Cairn keeps running in the tray.
4. **Given** Linux without a system tray, or with "Run in background" off, **When** the user closes the window, **Then** Cairn quits.
5. **Given** any platform, **When** the user chooses Quit (Cmd+Q / Ctrl+Q, the app menu, main menu or tray menu), **Then** Cairn quits without asking for confirmation; if a sync upload is in progress, the window shows "Finishing sync…" and quits once it completes or after a short timeout.
6. **Given** Cairn was started at login, **Then** no window is shown until the user summons it.

---

### User Story 4 - Add a bookmark from the window (Priority: P2)

The user copies a URL in the browser, brings Cairn forward, presses Cmd/Ctrl+N (or clicks Add), optionally types tags, and presses Enter.

**Why this priority**: Saving is the second most common journey; it must move from a view swap to a platform-appropriate modal.

**Independent Test**: With a URL on the clipboard, open the add form, save, and confirm the new bookmark is selected in the list; repeat with a duplicate URL.

**Acceptance Scenarios**:

1. **Given** the clipboard holds a link, **When** the user opens Add Bookmark, **Then** a modal attached to the window appears (a sheet under the title band on macOS, a dialog over the dimmed window on Linux) with the URL pre-filled and selected and focus in the URL field.
2. **Given** the user types `+go.dev` in the search field, **Then** the same modal opens with `go.dev` in the URL field and the cursor at the end.
3. **Given** the modal is open, **When** the user presses Enter (or Add), **Then** the bookmark is saved, the modal closes, the new row is selected, and a status message "Added. Fetching title…" appears briefly.
4. **Given** the URL is already bookmarked, **When** the user saves, **Then** the modal closes, the existing row is selected and an "Already bookmarked" confirmation appears; no error dialog is shown.
5. **Given** the modal is open, **When** the user presses Esc or Cancel, **Then** the modal closes without saving and focus returns to the search field.

---

### User Story 5 - Discoverable menus, settings and shortcuts (Priority: P2)

The user looks for Settings, Quit, keyboard shortcuts and per-bookmark actions where their platform normally puts them.

**Why this priority**: Makes quitting, start-at-login and shortcuts discoverable, which today are hidden in the tray menu or docs.

**Independent Test**: On each platform, find Settings, Keyboard Shortcuts, About and Quit from the standard places, change each setting and confirm it takes effect.

**Acceptance Scenarios**:

1. **Given** macOS, **Then** the menu bar shows Cairn, File, Edit, Window and Help menus; the Cairn menu holds About, Settings… (Cmd+,), Hide and Quit (Cmd+Q); the File menu holds New Bookmark…, Open Link, Copy Link, Pin/Unpin, Sync Now and Close Window.
2. **Given** GNOME or KDE, **Then** the window offers an Add button and a main menu button containing Add Bookmark…, Sync Now, Preferences, Keyboard Shortcuts, About Cairn and Quit, in that order.
3. **Given** a result row, **When** the user right-clicks it, **Then** a context menu offers Open Link, Copy Link and Pin (Unpin when already pinned).
4. **Given** Settings is open, **Then** it offers: Open Cairn at login; Show in menu bar (macOS only); Run in background (Linux only, disabled with an explanation when no tray is available); Hide Cairn after opening a link; and the global shortcut row.
5. **Given** the global shortcut cannot be registered (taken, refused, or unsupported by the desktop), **Then** Settings says so and shows how to bind `cairn-desktop --toggle` in the system keyboard settings with a button that copies the command.
6. **Given** the user presses Cmd+/ (macOS) or Ctrl+? (Linux), **Then** a list of keyboard shortcuts is shown.

---

### User Story 6 - Refreshed, readable contents (Priority: P3)

The list, messages and empty states follow the design: clear rows, a proper pin indicator, sectioned Pinned and Recent lists, and messages that match their urgency.

**Why this priority**: Improves look and feel but builds on the window work above.

**Independent Test**: Compare the app in light and dark mode on each platform against the design artboards for rows, selection, empty states and messages.

**Acceptance Scenarios**:

1. **Given** an empty query, **Then** pinned bookmarks are listed under a "Pinned" heading, followed by recent bookmarks under "Recent".
2. **Given** a query, **Then** each row shows a monogram tile, the title with the matched text in bold, the domain, up to three tags followed by "+N" for the rest, and a pin glyph for pinned bookmarks; the selected row also shows an "Open" hint.
3. **Given** no bookmarks exist, **Then** the list says "No bookmarks yet" and how to add one; **Given** no results match, **Then** it says "No bookmarks match "<query>"" and suggests typing `+` before a URL.
4. **Given** sync fails with a persistent problem such as an expired Dropbox sign-in, **Then** a banner appears under the title area and stays until the problem is fixed or dismissed; on macOS the menu bar icon also shows a warning badge.
5. **Given** a passing confirmation ("Already bookmarked", "Link copied"), **Then** it appears as a toast for about 4 seconds.
6. **Given** the system is in dark mode, **Then** the window, list and messages use the dark palette, and all text meets a 4.5:1 contrast ratio.

---

### Edge Cases

- The remembered window position is on a display that is no longer connected: the window opens centred on the main display at the remembered size.
- The remembered size is below the minimum or larger than the current display: it is clamped to fit.
- On Wayland, where apps cannot place their own windows: only the size is remembered and the desktop chooses the position.
- The global shortcut is pressed while Add Bookmark is open: the window comes forward (or goes away) with the modal state preserved.
- The global shortcut is already taken by another app: the app keeps running, Settings reports the conflict, and the tray/Dock/`--toggle` still work.
- A tray host appears or disappears while Cairn is running on Linux: "Run in background" availability updates; if the tray disappears while the window is closed, Cairn stays running until the next quit.
- The user closes the window on Linux during a sync upload with background mode off: quit behaves as in Story 3, scenario 5.
- An existing user upgrading from the popup version, who had start-at-login set up by the installer: start-at-login stays on after upgrade.
- "Show in menu bar" is turned off on macOS while the window is closed: the Dock icon and shortcut still bring Cairn back.
- A very long title or many tags: the title is truncated with an ellipsis on one line and tags collapse to "+N".
- The system accent colour is too light for white text: the selection colour is darkened until it meets contrast.

## Requirements *(mandatory)*

### Functional Requirements

**Window chrome and look**

- **FR-001**: The app MUST present a standard window with the platform's own window controls: on macOS the native close/minimise/zoom buttons inset in a title band; on Linux the desktop's own title bar with the user's button layout and theme.
- **FR-002**: On macOS, the search field MUST sit in the title band; the rest of the band MUST be draggable and double-clicking it MUST zoom the window. Only the title band MAY be translucent; the result list MUST stay opaque.
- **FR-003**: On Linux, Cairn's search bar MUST sit directly below the title bar and hold the Add button and the main menu button.
- **FR-004**: On KDE Plasma, the app MUST apply the KDE variant from the design (smaller corner radius, Breeze accent, labelled Add button, footer kept); on GNOME the footer MUST be omitted.
- **FR-005**: The app MUST NOT draw its own window border, shadow or window buttons.
- **FR-006**: The app MUST follow the system light/dark setting and use the system UI font.

**Window behaviour**

- **FR-007**: The window MUST open at 720×520 on first launch, centred, and MUST NOT be resizable below 520×360.
- **FR-008**: The app MUST remember the window size on all platforms and the window position on macOS and X11, restoring both on reopen and relaunch; on Wayland it MUST remember size only.
- **FR-009**: The window MUST NOT be always-on-top and MUST NOT hide when focus moves to another app.
- **FR-010**: Every way of bringing the window forward (launch, relaunch, shortcut, Dock, taskbar, tray, `--toggle`) MUST focus the search field and select its text.
- **FR-011**: Esc MUST clear a non-empty query or cancel an open modal, and MUST NOT close or hide the window.
- **FR-012**: Opening a link MUST leave the window open unless the "Hide Cairn after opening a link" setting is on (default off).
- **FR-013**: Only one instance MUST run; a second launch MUST raise the existing window, and `--add` MUST open the Add Bookmark modal.
- **FR-014**: The window MUST be maximisable/zoomable and usable full screen, with rows not stretched beyond a readable width.

**App presence and lifecycle**

- **FR-015**: On macOS the app MUST show a Dock icon and a menu bar; clicking the Dock icon with no window open MUST reopen it.
- **FR-016**: On macOS, closing the window MUST keep the app running; Quit MUST end it.
- **FR-017**: On Linux, closing the window MUST quit the app unless "Run in background" is on and a system tray is available, in which case it MUST keep running in the tray.
- **FR-018**: On Linux the app MUST appear in the taskbar/dash as "Cairn" with its icon, grouped with its launcher entry.
- **FR-019**: A tray (Linux) or menu bar (macOS) icon MUST offer Search Bookmarks…, Add Bookmark…, sync status, Settings… and Quit Cairn; on macOS it MUST be shown by default and be hideable from Settings; on Linux it MUST be shown only when a tray is available.
- **FR-020**: Quit MUST NOT ask for confirmation; if a sync upload is in progress it MUST show "Finishing sync…" and wait at most 5 seconds before quitting.
- **FR-021**: When started at login, the app MUST start without showing a window.

**Global shortcut**

- **FR-022**: The configured global shortcut (`desktop_hotkey`, default Cmd/Ctrl+Shift+Space, `off` to disable) MUST bring Cairn forward when it is not in front and put it away when it is (hide the app on macOS; minimise, or close in background mode, on Linux).
- **FR-023**: On Wayland desktops the app MUST request the shortcut through the desktop's shortcut approval mechanism where available, and `cairn-desktop --toggle` MUST be able to raise and focus the window.
- **FR-024**: If the shortcut cannot be registered, the app MUST keep running and Settings MUST explain the failure and offer to copy the `cairn-desktop --toggle` command for a custom system shortcut.

**Menus and actions**

- **FR-025**: On macOS the app MUST provide Cairn, File, Edit, Window and Help menus as listed in User Story 5, with system-standard Edit and Window behaviour.
- **FR-026**: On Linux the main menu MUST contain Add Bookmark…, Sync Now, Preferences, Keyboard Shortcuts, About Cairn and Quit.
- **FR-027**: Each result row MUST have a context menu with Open Link, Copy Link and Pin/Unpin; File-menu actions MUST act on the selected row.
- **FR-028**: The app MUST support the keyboard map: Cmd/Ctrl+F focus search; ↑/↓ move selection; Enter open; Shift+Cmd/Ctrl+C copy link; Cmd/Ctrl+N add; Esc clear/cancel; Cmd/Ctrl+W close window; Cmd/Ctrl+Q quit; Cmd/Ctrl+, settings; Cmd+/ (macOS) or Ctrl+? (Linux) show shortcuts.
- **FR-029**: Sync Now MUST trigger a pull and push and report the result in the status area.

**Add bookmark**

- **FR-030**: Add Bookmark MUST open as a modal attached to the main window (sheet on macOS, dialog on Linux), never as a separate top-level window or a full view swap.
- **FR-031**: The modal MUST pre-fill the URL from the clipboard when it holds an http(s) link, select it, and accept comma-separated tags; `+<text>` typed in search MUST open it with that text.
- **FR-032**: The modal MUST order its buttons Cancel then Add, with Add as the default action; Enter adds, Esc cancels, Tab moves between fields.
- **FR-033**: After a successful add the modal MUST close, the new row MUST be selected, and "Added. Fetching title…" MUST show until the title is fetched; on a duplicate the existing row MUST be selected and "Already bookmarked" shown as a toast.
- **FR-034**: Adding MUST keep today's behaviour: scheme added if missing, only http/https accepted, title and description fetched, sync push afterwards.

**Contents**

- **FR-035**: With an empty query, the list MUST show a "Pinned" section followed by a "Recent" section; with a query, results MUST keep the current ranking.
- **FR-036**: Each row MUST show a monogram tile coloured consistently per domain, the title on one line with matched text in bold, the domain, up to 3 tags then "+N", a pin glyph when pinned (no emoji), and an Open hint on the selected row.
- **FR-037**: Selection MUST use the system accent colour where available, falling back to the design's blue, and MUST meet 4.5:1 text contrast (using a tint with outline where white text would not).
- **FR-038**: Empty states MUST give one short line and the next step to take.
- **FR-039**: Passing confirmations MUST appear as toasts (~4 s); persistent problems MUST appear as a dismissible banner under the title area that returns on the next failure; background progress MUST appear in the status area (macOS and KDE footer, or a toast on GNOME).
- **FR-040**: Controls MUST be accessible: a labelled search field, a list announced as a listbox with the selected item exposed, real buttons with accessible names for icon-only buttons, and toasts/banners announced as status.

**Settings and login**

- **FR-041**: Settings (Cmd/Ctrl+,) MUST offer Open Cairn at login, Show in menu bar (macOS), Run in background (Linux; on by default on KDE, off elsewhere; disabled with an explanation when no tray is available), Hide Cairn after opening a link, and the global shortcut row; changes MUST take effect without restarting.
- **FR-042**: Open at login MUST be off for new installs; users upgrading with an existing installer-created login item MUST keep it on, and the app MUST become the only place that creates or removes the login item (the installer no longer does).
- **FR-043**: Settings MUST be stored apart from `cairn.json` so the app never rewrites the user's config file.

**Documentation**

- **FR-044**: The Desktop App user guide MUST be updated to describe the window, close/quit behaviour, menus, settings, start-at-login and keyboard map.

### Key Entities

- **Window state**: the last window size and, where supported, position and the display it was on; remembered between runs.
- **Desktop settings**: open at login, show in menu bar, run in background, hide after opening a link; stored per user alongside the window state, separate from `cairn.json`.
- **Login item**: the platform's start-at-login entry for Cairn, owned and managed by the app.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: From pressing the global shortcut to the link opening in the browser takes no more keystrokes than today (shortcut, 2–3 letters, Enter) and the window is ready for typing within 300 ms of the shortcut on a warm app.
- **SC-002**: On macOS, GNOME and KDE, 100% of the window operations a user expects (close, minimise where the platform shows it, maximise/zoom, move, resize, snap) work through the platform's own controls with no Cairn-specific instructions.
- **SC-003**: A first-time user can find how to quit Cairn and how to turn on start-at-login within 30 seconds without reading documentation.
- **SC-004**: The window never disappears unless the user closes it, hides it, presses the shortcut, or has turned on "Hide Cairn after opening a link".
- **SC-005**: After relaunch, the window reopens at its previous size (and position, where supported) in 100% of cases where the previous display is still connected.
- **SC-006**: All text in the window meets 4.5:1 contrast in light and dark mode, and every control is reachable and operable by keyboard.
- **SC-007**: Users upgrading from the popup version who had start-at-login keep it with no action on their part.
- **SC-008**: Support requests and issues about missing window controls, inability to resize/move, or not knowing how to quit drop to zero after release.

## Assumptions

- The design canvas linked above is the source for visual details (sizes, colours, spacing, copy); where this spec and the canvas differ, this spec wins.
- Monogram tiles stand in for favicons; stored favicons are a follow-up that needs a database change.
- Reading the live system accent colour is included only where it needs no custom native code; otherwise the fallback colours apply and live accent support is a follow-up.
- On GNOME without a tray extension, Cairn quits when its window closes by default; users wanting it always available bind `cairn-desktop --toggle` to a system shortcut, which relaunches it.
- The existing search ranking, add/fetch logic, duplicate check and background sync are reused unchanged.
- The app continues to read `cairn.json` (including `db_path` and `desktop_hotkey`) as today.

## Out of Scope

- Stored favicons.
- Archive and Delete from the desktop app, and undo.
- Editing tags or titles in the desktop app.
- Setting up Dropbox sync from the desktop app.
- Windows support.
