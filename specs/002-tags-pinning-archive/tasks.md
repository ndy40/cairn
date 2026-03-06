# Tasks: Tags, Pinning, Archive & Startup Checks

**Input**: Design documents from `/specs/002-tags-pinning-archive/`
**Prerequisites**: plan.md ✓, spec.md ✓, research.md ✓, data-model.md ✓, contracts/ ✓

**Organization**: Tasks are grouped by user story. Each story is independently testable after its phase completes.

**Note on story ordering**: US6 (Permanent) is implemented before US5 (Archive) because `ArchiveStale()` depends on the `is_permanent` column being queryable — the schema is added in Phase 2 (Foundational) so SQL runs correctly from the start, but the store helper for toggling permanence is needed before the archive logic is wired to the UI.

## Format: `[ID] [P?] [Story] Description`

---

## Phase 1: Setup

**Purpose**: Confirm the existing project is ready for extension.

- [X] T001 Confirm `go build ./...` and `go vet ./...` succeed on the current codebase before making any changes in `cmd/bm/` and `internal/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Extend the schema and core data structures that every user story depends on.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T002 Add migration v2 to the `migrations` slice in `internal/store/store.go` using five `ALTER TABLE bookmarks ADD COLUMN` statements: `tags TEXT NOT NULL DEFAULT '[]'`, `last_visited_at TEXT`, `is_permanent INTEGER NOT NULL DEFAULT 0`, `is_archived INTEGER NOT NULL DEFAULT 0`, `archived_at TEXT`; also add indexes on `is_archived`, `archived_at DESC`, and `is_permanent`
- [X] T003 Extend the `Bookmark` struct in `internal/store/bookmark.go` with five new fields: `Tags []string`, `LastVisitedAt *time.Time`, `IsPermanent bool`, `IsArchived bool`, `ArchivedAt *time.Time`
- [X] T004 Update `List()` in `internal/store/bookmark.go` to add `WHERE is_archived = 0` filter and scan the five new columns; update `GetByID()` and `ListByIDs()` to scan new columns; update `Insert()` signature to `Insert(url, title, description string, tags []string) (int64, error)` and JSON-encode tags into the `tags` column; update all callers in `cmd/bm/main.go` (CLI `bm add` command) to pass `nil` tags
- [X] T005 Update `fetchAndSave` cmd in `internal/model/app.go` to accept `tags []string` parameter and pass it to `store.Insert()`; update the `fetchSaveResultMsg` if needed to carry tags; update the `updateAdd` handler call to pass empty tags for now (tags wiring comes in US2)

**Checkpoint**: `go build ./...` must succeed. Existing TUI and CLI behaviour must be unchanged.

---

## Phase 3: User Story 1 — Startup Prerequisite Check (Priority: P1) 🎯 MVP

**Goal**: Detect Wayland/X11 display environment at TUI launch; verify the required clipboard tool is installed; print a specific install instruction and exit(1) if missing; warn-and-continue when display is unknown.

**Independent Test**: Run `WAYLAND_DISPLAY="" DISPLAY="" ./bm` → warning shown, TUI launches. Rename `wl-paste` temporarily on a Wayland session → specific error shown, exit code 1. Run `./bm` with all tools present → no message, TUI launches normally.

- [X] T006 [US1] Create `internal/display/check.go` with `type DisplayType int` constants (`Wayland`, `X11`, `Unknown`), `type CheckResult struct { DisplayType, ToolFound bool, MissingTool, InstallHint string, ShouldBlock bool }`, and `func CheckPrerequisites() CheckResult` that reads `WAYLAND_DISPLAY` (Wayland, wins when both set) and `DISPLAY` (X11) env vars, then calls `exec.LookPath("wl-paste")` for Wayland or `exec.LookPath("xclip")`/`exec.LookPath("xsel")` for X11; sets `ShouldBlock = true` when display detected but tool absent; returns Unknown result with `ShouldBlock = false` when neither display var is set
- [X] T007 [US1] Modify `runTUI()` in `cmd/bm/main.go` to call `display.CheckPrerequisites()` immediately after parsing flags (before opening the store); if `result.ShouldBlock` is true, print `result.InstallHint` to `os.Stderr` and `os.Exit(1)`; if display is Unknown and `!result.ToolFound`, print the non-blocking warning to `os.Stderr` and continue

**Checkpoint**: US1 fully functional. Prerequisite check blocks or warns as specified. All existing CLI subcommands still work (check is TUI-only).

---

## Phase 4: User Story 2 — Tag Bookmarks (Priority: P2)

**Goal**: Users can assign up to 3 lowercase tags to a bookmark in the add modal; tags are stored, validated (dedup, truncate, max 3), and displayed in the browse list.

**Independent Test**: Add a bookmark via Ctrl+P and type `work, go, tools` in the Tags field — all 3 tags saved and visible. Type a 4th tag — rejected with message. Type `Work` — saved as `work`. Leave tags empty — bookmark saved with no tags. Edit tags on existing bookmark (open add modal pre-filled with existing tags).

- [X] T008 [US2] Add tag validation helper `func normaliseTags(input string) []string` in `internal/store/bookmark.go` (or a dedicated `internal/tags/tags.go`) that splits on commas, trims whitespace, lowercases, filters empty strings, deduplicates, truncates each tag to 32 runes, and returns the first 3 tags; export as needed by the model layer
- [X] T009 [US2] Add a second `bubbles/textinput` field (`tagsInput`) to `AddModel` in `internal/model/add.go`; label it "Tags (comma-separated, max 3)"; handle `Tab`/`Shift+Tab` to cycle focus between URL input and tags input; update `View()` to render both fields; add `Tags() string` accessor method returning `tagsInput.Value()`
- [X] T010 [US2] Update `updateAdd()` in `internal/model/app.go` to extract tags from `a.add.Tags()`, call `normaliseTags()`, and pass the result to `fetchAndSave(a.store, rawURL, tags)`; if `len(normalisedTags) < len(rawTagCount)` set `a.add.setStatus("Only first 3 tags saved")` before saving
- [X] T011 [US2] Update `BookmarkItem.Description()` in `internal/model/browse.go` to append tags as `#tag1 #tag2` after the domain and date, separated by ` · `; skip tag section when `len(b.Tags) == 0`
- [X] T012 [US2] Update `BookmarkItem.Title()` in `internal/model/browse.go` to prepend `[pin] ` when `b.IsPermanent == true` (added here for completeness, used actively in US6; costs nothing to add now)

**Checkpoint**: US2 fully functional. Tags saved, normalised, and displayed. Existing bookmarks show no tags gracefully.

---

## Phase 5: User Story 3 — Filter Bookmarks by Tag (Priority: P3)

**Goal**: Users press `t` in browse mode to open a tag filter overlay; selecting tags narrows the list using OR logic; tag filter composes with the existing text search.

**Independent Test**: With bookmarks tagged `work` and `go`, select `work` filter — only `work`-tagged bookmarks appear. Also type a search term — only results matching tag AND search show. Clear filter (`c`) — full list restored.

- [X] T013 [US3] Add `activeTagFilter []string` field to `App` struct in `internal/model/app.go`; add `func tagFilter(bookmarks []*store.Bookmark, tags []string) []*store.Bookmark` (OR logic: include if bookmark has any selected tag) in `internal/model/app.go`; update the browse load path (`bookmarksLoadedMsg` handler and `twoStageSearch` call site) to apply `tagFilter(a.allBookmarks, a.activeTagFilter)` as the first stage before FTS5/fuzzy
- [X] T014 [US3] Create `internal/model/tagfilter.go` with `TagFilterModel` struct holding a `bubbles/list` of tag items, a `selected map[string]bool`, and `allTags []string` derived from the active bookmark list; implement `newTagFilterModel(bookmarks []*store.Bookmark) TagFilterModel`, `Update(msg tea.Msg) (TagFilterModel, tea.Cmd)` (handle `j`/`k` navigate, `Space`/`Enter` toggle, `c` clear-all, `Esc`/`t` close), and `View() string`; add `SelectedTags() []string` accessor
- [X] T015 [US3] Add `StateTagFilter` to the `AppState` enum in `internal/model/app.go`; add `tagFilter TagFilterModel` field to `App`; handle `t` key in `updateBrowse` to enter `StateTagFilter` by calling `newTagFilterModel(a.allBookmarks)`; handle `StateTagFilter` in `Update()` to delegate to `a.tagFilter.Update()`; on close (Esc/t/c), read `a.tagFilter.SelectedTags()`, store in `a.activeTagFilter`, reload browse list
- [X] T016 [US3] Add `StateTagFilter` case to `View()` in `internal/model/app.go` to render `a.tagFilter.View()`; add filter status indicator to `browseView()` footer (e.g., `[Filter: work, go]`) when `activeTagFilter` is non-empty
- [X] T017 [US3] Add `t` key to browse mode footer hint and help screen in `internal/model/app.go`; add `[t] Tag filter` to `browseView()` footer and to `helpView()` Browse Mode section

**Checkpoint**: US3 fully functional. Tag filter overlay opens, toggles work, filter composes with text search, clear resets to full list.

---

## Phase 6: User Story 4 — Track Last Visited Date (Priority: P4)

**Goal**: Record the timestamp each time a bookmark is opened via the TUI; display "Never" or the date in the bookmark list.

**Independent Test**: Open a bookmark — its description line updates to show today's date as last visited. Open it again — date updates to the new time. A bookmark never opened shows "Never".

- [X] T018 [US4] Create `internal/store/archive.go` with `func (s *Store) UpdateLastVisited(id int64) error` that executes `UPDATE bookmarks SET last_visited_at = datetime('now') WHERE id = ?`
- [X] T019 [US4] Add `func openBookmarkCmd(s *store.Store, b *store.Bookmark) tea.Cmd` in `internal/model/app.go` that runs the OS browser open command (reusing the `openURL` logic from `browse.go`) and, on success (`cmd.Start()` returns nil), calls `s.UpdateLastVisited(b.ID)`; returns `openURLErrMsg` on failure
- [X] T020 [US4] Replace all `openURL(sel.URL)` calls in `internal/model/app.go` (in `updateBrowse` and `updateSearch`) and in `internal/model/browse.go` (in the `enter` key handler) with `openBookmarkCmd(a.store, sel)` — pass the store reference from `App` to the relevant handlers; update `BookmarkItem.Description()` in `internal/model/browse.go` to append ` · Last: <date>` (formatted `2006-01-02`) or ` · Never visited` based on `b.LastVisitedAt`

**Checkpoint**: US4 fully functional. Last-visited date updates on every open. "Never" shown for unvisited bookmarks.

---

## Phase 7: User Story 5 — Auto-Archive Stale Bookmarks (Priority: P5)

**Goal**: On every TUI startup, bookmarks unvisited for 183+ days (and not permanent) are silently archived; a count is shown in the footer; archived bookmarks are visible in a dedicated archive view accessible via `a`; bookmarks can be restored via `r`.

**Independent Test**: Seed a non-permanent bookmark with `created_at = datetime('now', '-200 days')`; launch the TUI — footer shows "1 bookmark archived"; press `a` — archive view lists the bookmark; press `r` — bookmark returns to active list; a permanent bookmark with same old date stays in active list.

- [X] T021 [US5] Add `func (s *Store) ArchiveStale() (int, error)` to `internal/store/archive.go` executing the eligibility UPDATE (see data-model.md) and returning `db.RowsAffected()` as the count
- [X] T022 [US5] [P] Add `func (s *Store) ListArchived() ([]*Bookmark, error)` to `internal/store/archive.go` querying `WHERE is_archived = 1 ORDER BY archived_at DESC` scanning all bookmark fields
- [X] T023 [US5] [P] Add `func (s *Store) RestoreByID(id int64) error` to `internal/store/archive.go` executing `UPDATE bookmarks SET is_archived = 0, archived_at = NULL WHERE id = ?`
- [X] T024 [US5] Create `internal/model/archive.go` with `ArchiveModel` struct wrapping a `bubbles/list`; implement `newArchiveModel() ArchiveModel`, `load(bookmarks []*store.Bookmark)`, `setSize(w, h int)`, `selected() *store.Bookmark`, `Update(msg tea.Msg) (ArchiveModel, tea.Cmd)` (handle `j`/`k` navigate, `r` sets `restoreRequested = true`, `Esc` returns), and `View() string` (empty-state message: "No archived bookmarks"); add `ArchiveBookmarkItem` implementing `list.Item` showing title, domain, archived date, last-visited date
- [X] T025 [US5] Add `StateArchive` to `AppState` enum in `internal/model/app.go`; add `archive ArchiveModel` and `startupArchiveCount int` fields to `App`; update `New(s *store.Store, archiveCount int) App` to accept the count and store it; add `loadArchivedBookmarks(s *store.Store) tea.Cmd` function
- [X] T026 [US5] Modify `runTUI()` in `cmd/bm/main.go` to call `store.ArchiveStale()` after opening the store (and after the prerequisite check); pass the returned count to `model.New(s, count)`
- [X] T027 [US5] Handle `StateArchive` in `App.Update()` and `App.View()` in `internal/model/app.go`: delegate key events to `a.archive.Update()`; handle `a` key in `updateBrowse` to enter `StateArchive` and fire `loadArchivedBookmarks`; handle `Esc` in archive state to return to browse; handle `archive.restoreRequested` by calling `store.RestoreByID()` and reloading both archived and active lists; show `startupArchiveCount` in footer status on first render (store in `footerMsg` and clear on first key press as normal)

**Checkpoint**: US5 fully functional. Archive check runs on startup. Archive view lists archived bookmarks. Restore works. Permanent bookmarks never appear in archive.

---

## Phase 8: User Story 6 — Mark Bookmark as Permanent (Priority: P6)

**Goal**: Users press `p` on any active bookmark to toggle the permanent flag; permanent bookmarks show a `[pin]` prefix; they are never archived.

**Independent Test**: Press `p` on a bookmark — `[pin]` prefix appears. Launch app again — the bookmark is not archived despite having a stale date. Press `p` again — prefix removed. Run the startup archive check — bookmark is now eligible.

- [X] T028 [US6] Add `func (s *Store) SetPermanent(id int64, permanent bool) error` to `internal/store/archive.go` executing `UPDATE bookmarks SET is_permanent = ? WHERE id = ?` with `1` or `0`
- [X] T029 [US6] Handle `p` key in `updateBrowse` in `internal/model/app.go`: call `store.SetPermanent(sel.ID, !sel.IsPermanent)` as a `tea.Cmd` returning a `bookmarkUpdatedMsg`; on `bookmarkUpdatedMsg`, reload bookmarks via `loadBookmarks()`; add `bookmarkUpdatedMsg` struct to `internal/model/app.go`
- [X] T030 [US6] Add `p` key binding to `BrowseModel.AdditionalShortHelpKeys` and `AdditionalFullHelpKeys` in `internal/model/browse.go` alongside the existing delete key; add `permanentKey` binding constant with help text `"p", "toggle permanent"`
- [X] T031 [US6] Add `[p] Pin` to browse mode footer hint in `browseView()` in `internal/model/app.go` and add the permanent toggle to `helpView()` Browse Mode section

**Checkpoint**: US6 fully functional. Permanent flag persists across restarts. `[pin]` prefix visible. Archive check skips permanent bookmarks.

---

## Phase 9: Polish & Cross-Cutting Concerns

- [X] T032 Update `helpView()` in `internal/model/app.go` to document all new shortcuts added across US1–US6: `p` (pin), `t` (tag filter), `a` (archive view), `r` (restore in archive view), `Tab`/`Shift+Tab` (add modal field cycling)
- [X] T033 Update the browse mode footer string in `browseView()` in `internal/model/app.go` to include the most important new shortcuts without exceeding terminal width; suggested: `[/] Search  [Ctrl+P] Add  [Enter] Open  [d] Delete  [p] Pin  [t] Tags  [a] Archive  [?] Help  [Ctrl+C] Quit`
- [X] T034 Verify `go build ./...` and `go vet ./...` produce no errors or warnings
- [ ] T035 Manually run the quickstart.md test scenarios from `specs/002-tags-pinning-archive/quickstart.md` to confirm all six user stories work end-to-end

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 — **BLOCKS all user stories**
- **US1 (Phase 3)**: Depends on Phase 2 only — no dependency on other user stories
- **US2 (Phase 4)**: Depends on Phase 2 only — no dependency on US1
- **US3 (Phase 5)**: Depends on US2 (tags must exist to filter by)
- **US4 (Phase 6)**: Depends on Phase 2 only — no dependency on US1–US3
- **US5 (Phase 7)**: Depends on US4 (UpdateLastVisited in archive.go) — T021 extends archive.go created in T018
- **US6 (Phase 8)**: Depends on Phase 2 only; the `[pin]` prefix in T012 was added early as a no-op
- **Polish (Phase 9)**: Depends on all user story phases complete

### User Story Dependencies

```
Phase 2 (Foundational)
  ├── Phase 3 (US1) — independent
  ├── Phase 4 (US2) — independent
  │     └── Phase 5 (US3) — depends on US2
  ├── Phase 6 (US4) — independent
  │     └── Phase 7 (US5) — depends on US4
  └── Phase 8 (US6) — independent
```

### Within Each User Story

- Store methods before model changes
- Model data structures before event wiring
- Event wiring before UI rendering
- UI rendering before footer/help text updates

### Parallel Opportunities

- After Phase 2: US1, US2, US4, and US6 can all begin in parallel
- Within US5: T022 `ListArchived` and T023 `RestoreByID` are marked [P] (different methods, same file, no mutual dependency)

---

## Parallel Example: Foundational Phase

```
T002 (migration v2 SQL in store.go)
T003 (Bookmark struct extension)     ← can be drafted in parallel with T002
T004 (List/Insert/Get scan updates)  ← depends on T002 and T003
T005 (fetchAndSave wiring in app.go) ← depends on T004
```

## Parallel Example: After Foundational Complete

```
Developer A: US1 (T006–T007) — display check package + main.go integration
Developer B: US2 (T008–T012) — tags store helper + add modal + browse display
Developer C: US4 (T018–T020) — last visited store + openBookmarkCmd + description
Developer D: US6 (T028–T031) — SetPermanent store + 'p' key handler
```

---

## Implementation Strategy

### MVP (US1 only — startup safety check)

1. Phase 1: Verify build
2. Phase 2: Migration v2 + struct extension
3. Phase 3: Prerequisite check
4. **Validate**: Launch on Wayland without wl-paste → error. Launch normally → no error.

### Incremental Delivery

1. **After Phase 3**: Prerequisite check working (US1)
2. **After Phase 5**: Tags saved, displayed, and filterable (US2+US3)
3. **After Phase 6**: Last-visited timestamps tracking (US4)
4. **After Phase 8**: Permanent flag + full archive lifecycle (US5+US6)
5. **After Phase 9**: All shortcuts documented, build clean

---

## Notes

- `[P]` = different files or independent methods, no ordering constraint between them
- `[USN]` maps each task to the user story that drives it
- The `internal/store/archive.go` file is created in T018 (US4) and extended in T021–T023 (US5) and T028 (US6) — a single file collects all archive/lifecycle store methods
- The `openURL()` function in `browse.go` becomes `openBookmarkCmd()` in `app.go` — the browse.go `enter` handler is updated to delegate to the app-level command so the store reference is accessible
- Tag filter state lives on `App` (not `BrowseModel`) so it survives mode transitions
- After each task group, run `go build ./...` to catch regressions early
