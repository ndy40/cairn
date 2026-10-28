# Contract: Go ↔ frontend bindings and events

The frontend talks to Go only through Wails-generated bindings (`frontend/bindings/...`) and events. Methods marked **new** or **changed** are added by this feature.

## BookmarkService

| Method | Status | Behaviour |
|---|---|---|
| `Search(query string) ([]Bookmark, error)` | unchanged | Empty query: pinned first, then most recent. Max 50. |
| `Open(id int64) error` | unchanged | Opens in the default browser. |
| `Add(rawURL, tags string) (*AddResult, error)` | unchanged | Errors: `ErrInvalidURL`, `ErrAlreadyBookmarked`. |
| `ClipboardURL() string` | unchanged | http(s) link on the clipboard, or "". |
| `FindByURL(rawURL string) (*Bookmark, error)` | **new** | Normalises like `Add` and returns the existing bookmark, so a duplicate add can select it (FR-033). nil if not found. |
| `SetPinned(id int64, pinned bool) error` | **new** | Same effect as `cairn pin`. Emits `bookmarks:changed`. |
| `CopyLink(id int64) error` | **new** | Puts the bookmark URL on the clipboard. |
| `SyncNow()` | **new** | Starts a pull then a push, ignoring the pull throttle. Result reported through `sync:*` events. |

## WindowService

| Method | Status | Behaviour |
|---|---|---|
| `PointerDown()` | **removed** | Focus guard deleted. |
| `Platform() PlatformInfo` | **new** | See data-model.md. |
| `Settings() DesktopSettings` | **new** | Effective settings, with per-desktop defaults resolved. |
| `SetSetting(key string, value bool) error` | **new** | `showInMenuBar`, `runInBackground`, `hideAfterOpen`. Applies immediately (e.g. shows/hides the menu bar icon). Error if the key is unknown or the setting is unavailable (background without a tray). |
| `OpenAtLogin() (bool, error)` | **new** | `app.Autostart.Status().Enabled`. |
| `SetOpenAtLogin(enabled bool) error` | **new** | Enable with `--background`, or disable. |
| `ShortcutStatus() ShortcutStatus` | **new** | See data-model.md. |
| `AfterOpen()` | **new** | Called after a link opens; hides the app/window if `hideAfterOpen` is on. |
| `CloseWindow()` | **new** | Same as the window's close button (Cmd/Ctrl+W). |
| `Quit()` | **new** | Quit path with the "Finishing sync…" wait. |
| `ShowSettings()` / `ShowShortcuts()` / `ShowAbout()` | **new** | Emit `ui:open` with the panel name, so menu items and HTML buttons share one path. |

## Events (Go → frontend)

| Event | Payload | Status | Meaning |
|---|---|---|---|
| `launcher:show` | `"search"` \| `"add"` | changed meaning | Window came forward: focus search and select its text, or open the add modal. |
| `bookmarks:changed` | `"pull"` \| `"pin"` | extended | Refresh the list, keeping the selection by id. |
| `sync:warning` | string | changed display | Show the persistent banner (until dismissed or a later sync succeeds). |
| `sync:ok` | — | **new** | Clears the banner; tray/menu bar badge cleared. |
| `status` | string | **new** | Background progress ("Finishing sync…", "Syncing…"); "" clears it. |
| `platform:changed` | `PlatformInfo` | **new** | Tray availability or accent colour changed. |
| `ui:open` | `"settings"` \| `"shortcuts"` \| `"about"` \| `"add"` | **new** | Open a panel or modal; sent by menus, the Dock menu and the tray. |
| `menu:action` | `"open"` \| `"copy"` \| `"pin"` \| `"find"` | **new** | macOS File/Edit menu items acting on the selected row. |

## Frontend obligations

- On every `launcher:show`, focus `#query` and select its text (FR-010).
- Esc: clear a non-empty query, or close the open `<dialog>`; otherwise do nothing (FR-011).
- After `Open`, call `AfterOpen()`.
- Keyboard map per FR-028 is handled in the web view. Where the macOS menu also binds the same accelerator, the menu item is the single handler and emits `menu:action`/`ui:open`.
- All bookmark text is set with `textContent` (titles come from fetched pages).
