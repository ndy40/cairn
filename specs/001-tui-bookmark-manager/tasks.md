# Tasks: TUI Bookmark Manager

**Input**: Design documents from `/specs/001-tui-bookmark-manager/`
**Prerequisites**: plan.md ✓, spec.md ✓, research.md ✓, data-model.md ✓, contracts/ ✓, quickstart.md ✓

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story. No tests are included (none requested in specification).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1–US4)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and directory structure

- [ ] T001 Create directory structure: `cmd/bm/`, `internal/model/`, `internal/store/`, `internal/fetcher/`, `internal/search/`, `internal/clipboard/`
- [ ] T002 Initialize Go module with `go mod init github.com/<username>/bookmark-manager` in repository root
- [ ] T003 Add all dependencies to `go.mod`: `charmbracelet/bubbletea`, `charmbracelet/bubbles`, `charmbracelet/lipgloss`, `modernc.org/sqlite`, `PuerkitoBio/goquery`, `golang.org/x/net`, `sahilm/fuzzy`, `atotto/clipboard` via `go get`
- [ ] T004 Create empty stub files with package declarations for each internal package: `internal/store/store.go`, `internal/store/bookmark.go`, `internal/store/search.go`, `internal/fetcher/fetcher.go`, `internal/search/fuzzy.go`, `internal/clipboard/clipboard.go`, `internal/model/app.go`, `internal/model/browse.go`, `internal/model/search.go`, `internal/model/add.go`, `cmd/bm/main.go`
- [ ] T005 Verify `go build ./...` succeeds on empty stubs before proceeding

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that every user story depends on — database, clipboard, HTTP fetcher, root TUI shell

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T006 Implement `internal/store/store.go`: open SQLite database at OS-appropriate path (`$XDG_DATA_HOME/bookmark-manager/bookmarks.db` on Linux, `~/Library/Application Support/bookmark-manager/bookmarks.db` on macOS), enable WAL mode (`PRAGMA journal_mode=WAL`), run schema migrations by checking `schema_version` table, create directory if not exists
- [ ] T007 Implement SQLite schema in `internal/store/store.go`: create `schema_version` table, `bookmarks` table (id, url UNIQUE, domain, title, description, created_at), indexes on `domain` and `created_at DESC`, FTS5 virtual table `bookmarks_fts` (title, description, domain, content='bookmarks', content_rowid='id'), and after-insert/after-delete triggers to keep FTS in sync
- [ ] T008 [P] Implement `internal/store/bookmark.go`: `Insert(url, domain, title, description, createdAt)` returning error on duplicate, `List()` returning all bookmarks ordered by created_at DESC, `DeleteByID(id)`, `ExistsByURL(url) bool`, `GetByID(id)` — use `modernc.org/sqlite` driver with `database/sql`
- [ ] T009 [P] Implement `internal/clipboard/clipboard.go`: `Read() (string, error)` wrapping `atotto/clipboard.ReadAll()`, return descriptive error if clipboard is empty or read fails
- [ ] T010 [P] Implement `internal/fetcher/fetcher.go`: `Fetch(url string) (title, description string, err error)` — HTTP GET with 8-second timeout, 512 KB body limit via `io.LimitReader`, `BookmarkManager/1.0` User-Agent; wrap response body with `golang.org/x/net/html/charset.NewReader` for encoding detection; use `goquery` to extract og:title > title tag for title, og:description > meta[name=description] for description; return URL hostname as fallback title on any error
- [ ] T011 Implement `cmd/bm/main.go` skeleton: parse `os.Args` to dispatch to subcommands (`add`, `list`, `search`, `delete`, `version`, `help`); when no subcommand provided, open store and launch bubbletea TUI; accept `--db` flag and `BM_DB_PATH` env var to override database path; call `store.Open()` and pass store to TUI model or subcommand handler
- [ ] T012 Implement `internal/model/app.go` root bubbletea model: define `AppState` type with constants `StateBrowse`, `StateSearch`, `StateAdd`, `StateConfirmDelete`; define root `Model` struct containing current state, store reference, and sub-model fields for browse, search, add; implement `Init() tea.Cmd`, `Update(tea.Msg) (tea.Model, tea.Cmd)` dispatching to active sub-model, `View() string` rendering active sub-model view + persistent footer; handle `Ctrl+C` globally to quit

**Checkpoint**: Foundation complete — `go build ./cmd/bm` succeeds; `./bm` launches a blank TUI that shows "No bookmarks" and quits on Ctrl+C

---

## Phase 3: User Story 1 — Add Bookmark via Ctrl+P (Priority: P1) 🎯 MVP

**Goal**: User presses Ctrl+P, clipboard URL is pasted into modal, page is fetched for title/description, bookmark is saved and visible.

**Independent Test**: Press Ctrl+P with a valid URL in clipboard → modal appears with URL pre-filled → press Enter → bookmark appears in list with correct title, domain, and creation date. Test with: valid URL, unreachable URL (saves with fallback title), empty clipboard (error message shown), duplicate URL (duplicate warning shown).

### Implementation for User Story 1

- [ ] T013 [US1] Implement `internal/model/add.go`: define `AddModel` struct with a `bubbles/textinput.Model` for URL entry and a `status string` field for error/success messages; `Init()` focuses the text input; `Update()` handles Enter (trigger fetch+save), Escape (cancel), and character input; expose `URL()` method returning current input value
- [ ] T014 [US1] Wire Ctrl+P shortcut in `internal/model/app.go` `Update()`: when `tea.KeyMsg` matches `"ctrl+p"`, read clipboard via `clipboard.Read()`, populate `AddModel` input with the clipboard value, transition to `StateAdd`; if clipboard is empty, show an inline error in the footer without entering add mode
- [ ] T015 [US1] Implement fetch-and-save flow in `internal/model/app.go` `Update()`: when `AddModel` emits a confirm message (Enter pressed), run `fetcher.Fetch(url)` and `store.Insert(...)` as a `tea.Cmd` (async command); on success, reload bookmark list and transition to `StateBrowse`; on duplicate error, display "Already bookmarked" in modal status; on fetch error, save bookmark with fallback title and show "Saved (title unavailable)" message
- [ ] T016 [US1] Implement `internal/model/app.go` `View()` for `StateAdd`: render the bookmark list dimmed in the background, then overlay a lipgloss-bordered modal box (min 60 chars wide) containing the label "Add Bookmark", the textinput field, the current `AddModel.status` message, and the footer `[Enter] Save  [Esc] Cancel`
- [ ] T017 [US1] Implement minimal bookmark list display in `internal/model/browse.go`: define `BrowseModel` with a `bubbles/list.Model`; implement a `BookmarkItem` struct satisfying `list.Item` interface with `Title()` returning page title and `Description()` returning domain + " · " + created_at date; `Load(bookmarks)` resets the list delegate items
- [ ] T018 [US1] Wire browse list into `internal/model/app.go` `View()` for `StateBrowse`: render `BrowseModel.View()` with `lipgloss` taking full terminal height minus footer row; render footer `[/] Search  [Ctrl+P] Add  [Enter] Open  [d] Delete  [?] Help  [Ctrl+C] Quit`
- [ ] T019 [US1] Load bookmarks from store on TUI startup: in `internal/model/app.go` `Init()`, return a `tea.Cmd` that calls `store.List()` and sends results as a message; `Update()` receives this message and calls `BrowseModel.Load(bookmarks)`

**Checkpoint**: `./bm` launches, Ctrl+P opens modal with clipboard URL, Enter saves and shows bookmark in list, error cases handled — US1 independently testable

---

## Phase 4: User Story 2 — Browse and Open Bookmarks (Priority: P2)

**Goal**: User navigates the bookmark list with arrow keys/j/k and opens a selected bookmark in the default browser.

**Independent Test**: With 5 pre-seeded bookmarks, launch app, scroll list with ↑/↓ and j/k, press Enter on a bookmark → correct URL opens in default browser. Test empty state message with no bookmarks.

### Implementation for User Story 2

- [ ] T020 [US2] Add full keyboard navigation to `internal/model/browse.go`: configure `bubbles/list` with default key bindings (↑/↓ arrows); add additional bindings for `j`/`k` (vim keys), `g` (jump to top), `G` (jump to bottom) via the list's `AdditionalShortHelpKeys` / `AdditionalFullHelpKeys`
- [ ] T021 [US2] Implement open-URL action in `internal/model/browse.go` `Update()`: on Enter keypress with a selected bookmark, run a `tea.Cmd` that executes the OS-appropriate open command (`xdg-open` on Linux, `open` on macOS, `start` on Windows) with the bookmark URL using `os/exec`; surface any exec error in the footer status
- [ ] T022 [US2] Implement empty state display in `internal/model/browse.go`: when the bookmark list is empty, render a centered message "No bookmarks yet. Press Ctrl+P to add your first bookmark." using lipgloss centering within the available terminal height
- [ ] T023 [US2] Configure `bubbles/list` delegate in `internal/model/browse.go`: set list title to "Bookmarks", disable built-in filtering (filtering is handled by the search mode), set item height to 2 lines (title on line 1, domain + date on line 2), configure selected item highlight style via lipgloss

**Checkpoint**: `./bm` shows full bookmark list with navigation, Enter opens URL in browser, empty state message shown when list is empty — US2 independently testable

---

## Phase 5: User Story 3 — Fuzzy Search Bookmarks (Priority: P3)

**Goal**: User types `/` to enter search mode, types a partial/misspelled term, list filters in real time across title, domain, and description.

**Independent Test**: With 20+ bookmarks, press `/`, type a partial word → list narrows instantly; type a misspelled term → near matches surface; press Esc → full list restores; press Enter on a result → URL opens in browser.

### Implementation for User Story 3

- [ ] T024 [US3] Implement `internal/store/search.go`: `FTSSearch(db, term string) ([]int64, error)` queries `bookmarks_fts` using FTS5 `MATCH` syntax to return matching bookmark IDs; used as a pre-filter when search term is 3+ characters; return all IDs (no filter) for terms under 3 characters
- [ ] T025 [US3] Implement `internal/search/fuzzy.go`: `Search(query string, bookmarks []Bookmark) []Bookmark` — runs `sahilm/fuzzy.FindFrom` separately against title (weight 3), domain (weight 2), description (weight 1) fields; merges results by bookmark ID taking the highest weighted score; returns bookmarks sorted by composite score descending; handles empty query by returning full slice unchanged
- [ ] T026 [US3] Implement `internal/model/search.go`: define `SearchModel` with a `bubbles/textinput.Model` (placeholder "Search bookmarks…") and the current filtered `[]Bookmark` slice; `Update()` on each keystroke calls `search.Search(term, allBookmarks)` and updates the filtered list; on Escape clear the input and signal parent to return to `StateBrowse`; on Ctrl+A clear only the search term; expose `Results() []Bookmark`
- [ ] T027 [US3] Wire search mode into `internal/model/app.go`: on `/` keypress in `StateBrowse`, transition to `StateSearch` and focus `SearchModel` input; on each search update message, call `BrowseModel.Load(searchResults)` so the list updates in real time; on Escape from search, reload full bookmark list from store and return to `StateBrowse`
- [ ] T028 [US3] Implement two-stage search in `internal/model/app.go`: when search term length ≥ 3, first call `store.FTSSearch(term)` to get candidate IDs, filter `allBookmarks` to candidates, then pass candidates to `search.Search()`; when term length < 3, skip FTS and pass full `allBookmarks` to `search.Search()` directly
- [ ] T029 [US3] Implement `View()` for `StateSearch` in `internal/model/app.go`: render `SearchModel` text input at the top of the terminal in a lipgloss-styled input bar (above the list); render filtered list below; render footer `[Esc] Clear  [Enter] Open  [Ctrl+P] Add  [Ctrl+C] Quit`
- [ ] T030 [US3] Handle "no results" state in `internal/model/search.go`: when `search.Search()` returns an empty slice, render a centered message "No results for «{term}»" in the list area instead of an empty list

**Checkpoint**: `/` enters search, real-time fuzzy filtering works across title/domain/description, Esc restores full list, Enter opens result — US3 independently testable

---

## Phase 6: User Story 4 — Delete Bookmark (Priority: P4)

**Goal**: User presses `d` or Delete on a selected bookmark, a confirmation dialog appears, `y`/Enter confirms deletion, `n`/Esc cancels.

**Independent Test**: Select a bookmark, press `d` → confirmation dialog shown with bookmark title; press `y` → bookmark removed from list and database; repeat but press `n` → bookmark still present.

### Implementation for User Story 4

- [ ] T031 [US4] Add `StateConfirmDelete` handling to `internal/model/app.go`: store the selected bookmark's ID and title when entering confirm-delete state; `View()` renders the existing browse list dimmed behind a centered lipgloss modal showing "Delete «{title}»?" and footer `[y/Enter] Confirm Delete  [n/Esc] Cancel`
- [ ] T032 [US4] Wire `d` and `Delete` keypresses in `internal/model/browse.go` `Update()`: when a bookmark is selected and `d` or `Delete` is pressed, emit a message to parent `app.go` containing the selected bookmark's ID and title to trigger `StateConfirmDelete`
- [ ] T033 [US4] Implement confirm/cancel in `internal/model/app.go` `Update()` for `StateConfirmDelete`: on `y` or Enter, run `store.DeleteByID(id)` as a `tea.Cmd`; on success, reload bookmark list from store, return to `StateBrowse`; on `n` or Escape, return to `StateBrowse` without any change; surface store errors in footer

**Checkpoint**: `d` on any bookmark shows confirmation dialog, `y` deletes and list refreshes, `n` cancels — US4 independently testable

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: CLI subcommands, edge-case handling, help overlay, and final validation

- [ ] T034 [P] Implement `bm add <url>` subcommand in `cmd/bm/main.go`: non-interactive mode; call `fetcher.Fetch(url)`, `store.Insert(...)`, print "Saved: «title» (domain)" on success, "Already bookmarked" on duplicate, "Saved (title unavailable): (domain)" on fetch error; use exit codes per CLI contract
- [ ] T035 [P] Implement `bm list` subcommand in `cmd/bm/main.go`: call `store.List()`, print tab-separated rows (id, title, url, domain, created_at) to stdout; support `--json` flag to output as JSON array using `encoding/json`
- [ ] T036 [P] Implement `bm search <query>` subcommand in `cmd/bm/main.go`: call `store.FTSSearch` + `search.Search`, print matching bookmarks in same format as `bm list`; support `--json` and `--limit N` flags
- [ ] T037 [P] Implement `bm delete <id>` subcommand in `cmd/bm/main.go`: call `store.DeleteByID(id)`, print "Deleted" on success, "Not found" on missing ID; use exit codes per CLI contract
- [ ] T038 Implement `--db` flag and `BM_DB_PATH` environment variable handling in `cmd/bm/main.go`: override default database path when either is provided; `--db` flag takes precedence over env var
- [ ] T039 Implement `NO_COLOR` env var and `--no-color` flag: when set, pass a `lipgloss.NewRenderer` with `termenv.Ascii` to disable ANSI color output in list/search subcommands
- [ ] T040 Implement help overlay in TUI (`?` key in any mode): `internal/model/app.go` renders a full-screen lipgloss overlay listing all keyboard shortcuts from `contracts/keyboard-shortcuts.md`; any keypress closes the overlay
- [ ] T041 Implement `bm version` subcommand: print version string (embed via `go build -ldflags="-X main.version=0.1.0"`)
- [ ] T042 Add `BM_DB_PATH` auto-creation: on store open, create parent directories with `os.MkdirAll` before opening SQLite file; ensure this is tested with a non-existent path on a fresh run
- [ ] T043 Validate quickstart.md by performing a clean build (`CGO_ENABLED=0 go build -ldflags="-s -w" -o bm ./cmd/bm`), running `./bm`, adding a bookmark non-interactively, listing it, and searching for it; update quickstart.md if any step is incorrect

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 (T001–T005) — **BLOCKS all user stories**
- **US1 (Phase 3)**: Depends on Phase 2 complete (T006–T012)
- **US2 (Phase 4)**: Depends on Phase 2 complete; integrates with browse model from US1 (T017–T019)
- **US3 (Phase 5)**: Depends on Phase 2 complete; reads from store and browse model established in US1/US2
- **US4 (Phase 6)**: Depends on Phase 2 complete; uses browse model from US1/US2 and store from foundation
- **Polish (Phase 7)**: Depends on all user stories complete

### User Story Dependencies

- **US1 (P1)**: Depends only on Foundational — no other story dependency
- **US2 (P2)**: Depends on Foundational; browse model (`BrowseModel`) established in US1 (T017–T019) is extended — must complete US1 first
- **US3 (P3)**: Depends on Foundational + US2 browse list being available (uses `BrowseModel.Load()`)
- **US4 (P4)**: Depends on Foundational + US2 browse list (needs selected item from browse model)

### Within Each User Story

- Model/data tasks before UI wiring tasks
- Store queries before model consumption
- Sub-model implementation before root app.go wiring
- Core flow before error-case handling

### Parallel Opportunities

- T003 (dependencies), T006–T010 within Phase 2 can run in parallel (separate files)
- T034–T037 in Polish phase can all run in parallel (separate subcommands, separate code paths)
- T038, T039, T040, T041, T042 in Polish phase are all independent and can run in parallel

---

## Parallel Example: Phase 2 Foundational

```
# These tasks operate on completely separate files — run in parallel:
Task T008: internal/store/bookmark.go   (CRUD operations)
Task T009: internal/clipboard/clipboard.go  (clipboard wrapper)
Task T010: internal/fetcher/fetcher.go  (HTTP + HTML parsing)
```

## Parallel Example: Phase 7 Polish

```
# CLI subcommands are independent — run in parallel:
Task T034: bm add subcommand
Task T035: bm list subcommand
Task T036: bm search subcommand
Task T037: bm delete subcommand
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001–T005)
2. Complete Phase 2: Foundational (T006–T012) — critical blocker
3. Complete Phase 3: User Story 1 (T013–T019)
4. **STOP and VALIDATE**: Ctrl+P adds a bookmark, bookmark appears in list, error cases work
5. You now have a usable bookmark saver with minimal display

### Incremental Delivery

1. Setup + Foundational → blank TUI shell launches
2. + US1 → bookmark saving works (MVP)
3. + US2 → full list navigation + open in browser
4. + US3 → fuzzy search (the key usability feature at scale)
5. + US4 → delete support
6. + Polish → CLI subcommands, help overlay, release build

### Single Developer Sequential Order

```
T001 → T002 → T003 → T004 → T005
T006 → T007 → T008 → T009 → T010 → T011 → T012
T013 → T014 → T015 → T016 → T017 → T018 → T019
T020 → T021 → T022 → T023
T024 → T025 → T026 → T027 → T028 → T029 → T030
T031 → T032 → T033
T034 → T035 → T036 → T037 → T038 → T039 → T040 → T041 → T042 → T043
```

---

## Notes

- `[P]` tasks operate on different files with no inter-dependencies — safe to implement in parallel
- `[Story]` label maps each task to a specific user story for traceability
- Each user story phase ends with a **Checkpoint** — validate independently before moving to the next phase
- The `bubbles/list` built-in filter is disabled in favor of the custom fuzzy search in US3
- All store operations that modify data should be run as `tea.Cmd` (async) to keep the TUI responsive
- On Linux, remind users to install `xclip` if clipboard paste fails (surface as a readable error, not a panic)
- Use `os.UserCacheDir()` / `os.UserConfigDir()` from stdlib to resolve the platform-appropriate database path
