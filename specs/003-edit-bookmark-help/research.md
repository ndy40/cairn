# Research: Edit Bookmark Tags, Last-Visited Visibility & CLI Help

**Feature**: 003-edit-bookmark-help
**Date**: 2026-03-06

---

## Decision 1: Edit Panel Trigger Key

**Decision**: Use `e` key in browse mode to open the edit panel.

**Rationale**: `e` is the universal convention for "edit" in terminal UIs (used by vim, ranger, nnn, lf, etc.). It is currently unbound in the browse mode keyset. No conflict with existing bindings (`/`, `g`, `G`, `d`, `p`, `t`, `a`, `enter`, `ctrl+p`, `ctrl+c`, `?`).

**Alternatives considered**:
- `E` (shift-e): Would work but requires shift; `e` is more ergonomic and conventional.
- `Tab`: Already used for field cycling in the add modal context.

---

## Decision 2: Edit Panel Scope (Tags Only)

**Decision**: The edit panel shows the bookmark title as a read-only label and a single editable tags field. Title and URL editing is out of scope.

**Rationale**: Editing a URL would require re-fetching the page title and meta description (an async network operation with error handling), and potentially checking for duplicate URLs. This is a significantly larger feature. Tags-only editing is a clean, bounded increment that delivers the missing capability immediately. The spec explicitly excludes title/URL editing.

**Alternatives considered**:
- Full edit modal (title + URL + tags): Larger scope, requires network fetch, duplicate URL check, and re-normalisation of domain. Deferred to a future feature.
- Inline editing (edit tags directly in the list row): Would require a custom list delegate and complex focus management. Rejected in favour of the simpler modal pattern already established in the codebase.

---

## Decision 3: Store Method for Tag Update

**Decision**: Add `UpdateTags(id int64, tags []string) error` to `internal/store/bookmark.go`. Executes `UPDATE bookmarks SET tags = ? WHERE id = ?` with JSON-encoded tags.

**Rationale**: A dedicated `UpdateTags` method is the minimal-touch approach: it updates exactly one column, reuses the existing `NormaliseTags` and JSON encoding logic already in `bookmark.go`, and leaves all other fields untouched (satisfying FR-006). No new store file is needed.

**Alternatives considered**:
- Generic `Update(b *Bookmark) error` method: Would update all fields in one call but risks accidentally overwriting last_visited_at or other fields if the caller has a stale copy. Too broad for this use case. Rejected.
- Putting it in `archive.go`: That file is for lifecycle operations (archive, restore, permanent flag). Tag editing is a content operation, better placed in `bookmark.go`.

---

## Decision 4: Last-Visited Update — Confirming Current Implementation

**Decision**: No code change needed for US2. The current `openBookmarkCmd` in `internal/model/app.go` already: (1) opens the browser via `openURLRaw`, (2) calls `s.UpdateLastVisited(b.ID)` on success, and (3) returns `loadBookmarks(s)()` which triggers a list reload showing the new date. This matches the spec requirements exactly.

**Rationale**: Review of the implementation confirms the correct behavior is already in place. The user's question ("when is the last visited link updated?") was a documentation/understanding gap, not a bug. US2 therefore requires no implementation work — only this documentation and the spec confirming the expected behavior. The `loadBookmarks` reload after `UpdateLastVisited` ensures the UI reflects the new date without a restart.

**What was verified**:
- `openBookmarkCmd` is called from both `updateBrowse` (Enter key) and `updateSearch` (Enter key) — both modes are covered.
- `openURLErrMsg` is returned on `openURLRaw` failure, preventing `UpdateLastVisited` from being called — matching FR-008.

---

## Decision 5: CLI Help Implementation

**Decision**: For root `--help`/`-h`: intercept before `flag.Parse()` by checking `os.Args` for `-h` or `--help`, call `printHelp()`, and exit 0. For subcommands: each subcommand's `flag.FlagSet` already calls `fs.Parse(args)` — add `-h`/`--help` detection before `fs.Parse()` or use `fs.Usage` to set a custom usage function that calls a per-subcommand help printer and exits 0.

**Rationale**: The Go standard `flag` package handles `-h`/`--help` by calling `Usage()` and exiting with code 2 by default. We want exit 0 (per FR-009). The cleanest approach: before calling `flag.Parse()` on the root, check for `-h`/`--help` manually and call `printHelp()` + `os.Exit(0)`. For subcommands, set `fs.Usage` to a function that prints subcommand-specific help and exits 0; call `fs.Bool("help", false, "show help")` as a flag on each FlagSet and check it after parse.

**Alternatives considered**:
- Using a third-party CLI framework (cobra, urfave/cli): Would add a dependency and require significant refactoring of the existing `main.go`. Rejected — the existing hand-rolled flag approach is sufficient and the constitution requires minimal abstractions.
- Overriding `flag.Usage` globally: Affects all flag sets; harder to provide per-subcommand help. Rejected in favour of per-FlagSet `Usage` override.
