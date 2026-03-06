# Research: Tags, Pinning, Archive & Startup Checks

**Feature**: 002-tags-pinning-archive
**Date**: 2026-03-06

---

## Decision 1: Tag Storage Format

**Decision**: Store tags as a JSON array in a single TEXT column (`tags TEXT NOT NULL DEFAULT '[]'`).

**Rationale**:
- Maximum 3 tags per bookmark makes a separate normalised table unnecessary overhead (extra JOIN, extra migration, extra store methods).
- `modernc.org/sqlite` ships with SQLite's built-in JSON functions (`json_each`, `json_extract`). If DB-side tag filtering is ever needed, it is trivially available.
- For ≤1,000 bookmarks, tag filtering is done in memory in Go after loading all active bookmarks — no DB query for filtering is needed.
- JSON array round-trips cleanly through Go's `encoding/json`; no custom parsing required.
- Comma-separated string was considered but rejected: corner cases around delimiter escaping and the need to write custom split/join logic outweigh any simplicity benefit.

**Alternatives considered**:
- Separate `tags` table with foreign key: Normalised, efficient at scale. Rejected — over-engineered for ≤3 tags, adds 3 new store methods, a JOIN in every List query, and a cascading delete concern.
- Comma-separated TEXT: Simple but rejected — requires custom parsing, LIKE query for filtering is fragile, no JSON tooling benefit.

---

## Decision 2: Display Environment Detection

**Decision**: Implement a new `internal/display` package with a single exported function `CheckPrerequisites() CheckResult`. Detection logic: `WAYLAND_DISPLAY` set → Wayland (check `wl-paste` via `exec.LookPath`); else `DISPLAY` set → X11 (check `xclip` or `xsel`); neither → unknown (non-blocking warning).

**Rationale**:
- The check is a TUI-only concern; isolating it in a dedicated package keeps `cmd/bm/main.go` clean and makes the logic unit-testable by injecting env vars.
- `exec.LookPath` is already used in `internal/clipboard/clipboard.go` for wl-paste; the same pattern applies here.
- When both `WAYLAND_DISPLAY` and `DISPLAY` are set (XWayland), Wayland takes precedence (per FR-001 and spec assumptions).
- The check must not block CLI subcommands (`bm add`, `bm list`, etc.) — it is called only in `runTUI()` before `tea.NewProgram()`.

**Alternatives considered**:
- Inline the check in `main.go`: Simple but mixes concerns and makes testing harder. Rejected.
- Reuse `internal/clipboard`: The clipboard package already does env detection, but its purpose is reading content, not reporting human-readable prerequisite errors. Reusing it would conflate two responsibilities. Rejected.

---

## Decision 3: Last-Visited Tracking

**Decision**: Record last-visited time in a new `last_visited_at TEXT` (nullable) column. Updated by a new `store.UpdateLastVisited(id int64) error` method called whenever a bookmark is opened via the TUI (inside the `openURL` tea.Cmd in `browse.go`).

**Rationale**:
- The `openURL` command already fires a `tea.Cmd` that opens the browser; adding a store call there is the minimal-touch approach.
- Storing as TEXT in ISO-8601 format (`2006-01-02T15:04:05Z`) is consistent with `created_at` already in the schema.
- Nullable (no default) means `NULL` = never visited — clean semantics for the "Never" display case.
- The `UpdateLastVisited` call happens after `cmd.Start()` succeeds, so failed opens do not update the timestamp.

**Alternatives considered**:
- Track visits in a separate `visit_log` table: Enables visit history but is far beyond what the spec requires. Rejected.
- Update last-visited via the existing `bookmark.go` methods: Would require adding a method to an already-cohesive file. Preferred to keep archive-related methods in the new `internal/store/archive.go` file for clarity.

---

## Decision 4: Archive Check Timing and Implementation

**Decision**: Run the archive check synchronously inside `runTUI()` immediately before `tea.NewProgram()`. Use a new `store.ArchiveStale() (int, error)` method that executes a single `UPDATE bookmarks SET is_archived=1, archived_at=datetime('now') WHERE ...` statement. The count of archived bookmarks is stored on the `App` model and shown in the footer status message at startup (cleared on first key press).

**Rationale**:
- Synchronous pre-TUI execution guarantees the archive check completes before the bookmark list loads, so the first render already shows the post-archive active list.
- A single SQL UPDATE for all eligible bookmarks is efficient (one statement, one transaction) and well within the 2-second performance target for ≤1,000 bookmarks.
- The 183-day threshold matches the spec (approximately 6 months). Stored as an integer constant `archiveThresholdDays = 183`.
- Eligibility SQL: `is_permanent = 0 AND is_archived = 0 AND (last_visited_at IS NULL AND created_at <= datetime('now', '-183 days') OR last_visited_at <= datetime('now', '-183 days'))`.

**Alternatives considered**:
- Run archive check as a `tea.Cmd` on `Init()`: Async, would complicate startup sequencing — the bookmark list would load before the check completes, causing a flicker or second reload. Rejected.
- Archive to a separate table: Preserves the main table cleanly. Rejected — the `is_archived` flag approach keeps the schema simple and restoring is a one-column UPDATE; no data movement needed.

---

## Decision 5: Permanent Flag Toggle

**Decision**: `p` key in browse mode toggles the `is_permanent` flag on the selected bookmark via a new `store.SetPermanent(id int64, permanent bool) error` method. The browse list immediately reflects the change via a reload (`loadBookmarks` cmd). Permanent bookmarks display a `[pin]` prefix in the list item title.

**Rationale**:
- `p` is currently unused in browse mode and is mnemonic for "pin/permanent".
- Toggling (same key for on and off) is simpler than two separate keys and consistent with the spec (FR-022, FR-023).
- `[pin]` prefix in `BookmarkItem.Title()` is the simplest visual distinction without requiring a custom `list.Delegate`.

**Alternatives considered**:
- Separate `P` (shift-P) for unpin: Two keys for one concept adds cognitive overhead. Rejected.
- Custom list delegate with icon column: More visual but requires a full custom `list.ItemDelegate` implementation. Over-engineered for this feature. Rejected.

---

## Decision 6: Tag Filtering in TUI

**Decision**: Tag filter state is held on the `App` model as `activeTagFilter []string`. When non-empty, `allBookmarks` is pre-filtered by tag match (OR logic) before being passed to the two-stage text search. A new `StateTagFilter` AppState drives a tag selection overlay, accessible via the `t` key from browse mode.

**Rationale**:
- In-memory filtering over ≤1,000 bookmarks is instantaneous; no DB query needed.
- OR logic for multiple selected tags: a bookmark is shown if it has any of the selected tags. This matches FR-012 and the spec's OR filter assumption.
- Composability with text search: `filtered = tagFilter(allBookmarks, activeTagFilter)` then `twoStageSearch(term, filtered)`. Clean pipeline with no special cases.
- The tag filter overlay shows all unique tags present in the active bookmark list and lets the user toggle them with Space/Enter.

**Alternatives considered**:
- DB query for tag filtering: Adds a round-trip per keystroke. Unnecessary for the scale. Rejected.
- Single-tag filter only: Spec explicitly requires multi-tag OR selection. Rejected.

---

## Decision 7: Archive View

**Decision**: A new `StateArchive` AppState drives a new `ArchiveModel` (mirrors `BrowseModel` pattern). Accessible via `a` key from browse mode. Shows archived bookmarks in reverse `archived_at` order (per spec edge cases). Restore is triggered by `r` key.

**Rationale**:
- Mirroring `BrowseModel` with a new `ArchiveModel` struct keeps the pattern consistent and the code readable.
- `a` is currently unused in browse mode and is mnemonic for "archive".
- `r` for restore is mnemonic and unused in the archive state.
- Archive list is loaded fresh from the DB when `StateArchive` is entered (separate `loadArchivedBookmarks` cmd).

**Alternatives considered**:
- Reuse `BrowseModel` with an archive flag: Would require conditional logic throughout browse rendering. Rejected — separation of concerns is cleaner.
- Dedicated archive DB table: Adds complexity for no benefit; `is_archived` flag with a filtered query is sufficient. Rejected.
