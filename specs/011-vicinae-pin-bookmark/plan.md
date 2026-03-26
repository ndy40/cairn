# Implementation Plan: Pin Bookmarks in Vicinae Extension

**Branch**: `011-vicinae-pin-bookmark` | **Date**: 2026-03-26 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/011-vicinae-pin-bookmark/spec.md`

---

## Summary

Add a "Toggle Pin" action to the Vicinae extension list view. The action calls a new `cairn pin <id>` CLI subcommand
that flips `is_permanent` in SQLite. The extension helper module gains a `bmPin()` function, and the list view
refreshes after the toggle. No schema change is required — `is_permanent` already exists from feature 002.

---

## Technical Context

**Language/Version**: TypeScript 5.9.2 (extension); Go 1.25.0 (CLI + store change)
**Primary Dependencies**: `@vicinae/api` (existing), React (existing); no new Go/TS dependencies
**Storage**: SQLite `is_permanent` column (existing, no migration needed)
**Testing**: `go test ./...` for store/CLI; manual `vici develop` for extension UI
**Target Platform**: Linux (Vicinae + cairn both run on Linux)
**Project Type**: Vicinae launcher extension (TypeScript) + CLI subcommand (Go)
**Performance Goals**: Pin toggle completes in < 500 ms; list refresh in < 1 s
**Constraints**: Extension must not access SQLite directly; all data through `cairn` CLI
**Scale/Scope**: Single user, local

---

## Constitution Check

| Gate | Status | Notes |
|------|--------|-------|
| No CGO | PASS | No new Go dependencies; pure-Go store unchanged |
| Single binary | PASS | Extension is a separate Vicinae package |
| Task management | PASS | Tasks via Backlog CLI after `/speckit.tasks` |
| Backward-compatible migrations | PASS | No schema changes; `is_permanent` column exists |
| Language: Go + TypeScript | JUSTIFIED (existing) | Same justification as feature 005 |

---

## Project Structure

### Documentation (this feature)

```text
specs/011-vicinae-pin-bookmark/
├── spec.md              # Feature spec
├── plan.md              # This file
├── contracts/
│   └── cli-pin-command.md   # cairn pin subcommand contract
└── tasks.md             # /speckit.tasks output (NOT created here)
```

### Source Code (new + modified)

```text
bookmark-manager/
├── internal/store/bookmark.go      # MODIFIED: add TogglePin(id int64) error
├── cmd/cairn/main.go                # MODIFIED: add `pin` subcommand, runPin()
│
└── vicinae-extension/src/
    ├── bm.ts                        # MODIFIED: add bmPin(id) export
    └── bm-list.tsx                  # MODIFIED: add "Toggle Pin" action to BookmarkListItem
```

**Structure Decision**: Minimal change set — two Go files, two TypeScript files. No new packages or files needed.

---

## Phase 0 Output: Research

### Decision: Toggle vs Set

**Decision**: Single `cairn pin <id>` command that toggles `is_permanent`.
**Rationale**: The TUI uses a toggle pattern (`[p] Pin`). A toggle keeps the CLI surface minimal — no need for
`pin --enable` / `pin --disable` flags. The extension can read current state from the list JSON and display the
correct action title ("Pin Bookmark" / "Unpin Bookmark").
**Alternatives considered**: Separate `pin` and `unpin` subcommands — rejected as unnecessary complexity for a
single boolean flag.

### Decision: List Refresh Strategy

**Decision**: After a successful pin toggle, invalidate the list cache (`listCache = null`) and call `bmList()`
to refresh state, then call `setBookmarks()`.
**Rationale**: This is the same pattern used by `bmDelete` in the existing code. Keeps state consistent without
a full page reload.
**Alternatives considered**: Optimistic UI update (flip `IsPermanent` in local state) — simpler but risks
showing stale state if the CLI call fails silently.

---

## Phase 1 Output: Design

### Data Model

No new entities or schema changes. The `Bookmark` interface in `bm.ts` already has `IsPermanent: boolean`. The
`is_permanent` column exists in SQLite. The store `Bookmark` struct already has `IsPermanent bool`.

The only data-layer addition is:

```go
// internal/store/bookmark.go
func (s *Store) TogglePin(id int64) error {
    _, err := s.db.Exec(
        `UPDATE bookmarks SET is_permanent = NOT is_permanent, updated_at = ? WHERE id = ?`,
        time.Now().UTC(), id,
    )
    return err
}
```

If no row is affected (id not found), a sentinel error `ErrNotFound` should be returned to allow the CLI to exit 1.

### CLI Contract

See [contracts/cli-pin-command.md](./contracts/cli-pin-command.md).

```
cairn pin <id>

Exit codes:
  0  Success — pin state toggled
  1  Bookmark not found
  3  Unexpected error
```

### Extension Changes

1. **`bm.ts`**: Add `bmPin(id: number)` that calls `runCairn(["pin", String(id)])`, invalidates `listCache`, and
   returns `{ exitCode, stderr }`.

2. **`bm-list.tsx`**: In `BookmarkListItem`, add an `onPin` prop and a third `<Action>` in the `<ActionPanel>`:
   - Title: `bookmark.IsPermanent ? "Unpin Bookmark" : "Pin Bookmark"`
   - On action: call `bmPin(bookmark.ID)`, show success/failure toast, call `onPin()` to trigger list refresh.

3. **`ListBookmarks`**: Add `handlePin` callback (same pattern as `handleDelete`): sets loading, re-fetches list,
   updates state. Pass it as `onPin` to each `BookmarkListItem`.
