# Tasks: Bookmark Expiry & Last-Visited Removal

**Feature**: 004-bookmark-expiry
**Branch**: `004-bookmark-expiry`
**Spec**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md)
**Generated**: 2026-03-06

---

## Summary

| Phase | Story | Tasks | Parallelisable |
|-------|-------|-------|----------------|
| 1 – Setup | — | 1 | 0 |
| 2 – US1: 30-Day Expiry | US1 | 2 | 1 |
| 3 – US2: Remove Last-Visited | US2 | 2 | 2 |
| 4 – Polish | — | 1 | 0 |
| **Total** | | **6** | |

---

## Phase 1: Setup

**Goal**: Verify project builds cleanly before changes.

- [X] T001 Verify clean build: run `go build ./...` and `go vet ./...` from repo root

---

## Phase 2: US1 — Bookmarks Expire After 30 Days

**Story goal**: At startup, archive all active non-pinned bookmarks whose `created_at` is 30 or more days ago; pinned bookmarks are exempt.

**Independent test criteria**: Add a bookmark, set its `created_at` to 31 days ago via SQL, restart the app, and verify it appears in the archive view and not the browse list.

- [X] T002 [US1] Simplify `ArchiveStale()` in `internal/store/archive.go`: replace the compound `last_visited_at`/`created_at` WHERE clause with `created_at <= datetime('now', '-30 days')` (keep `is_permanent = 0` and `is_archived = 0` conditions); remove the `last_visited_at` branch entirely
- [X] T003 [US1] Remove `UpdateLastVisited(id int64) error` method from `internal/store/archive.go` (the entire function body and comment)

---

## Phase 3: US2 — Remove Last-Visited Tracking

**Story goal**: No "Last: …" or "Never visited" text appears anywhere in the browse list or search results; opening a bookmark no longer records a visit timestamp.

**Independent test criteria**: Launch the app and verify no last-visited labels appear on any bookmark row. Open a bookmark; verify the row still shows no visit date after the list reloads.

- [X] T004 [P] [US2] Remove last-visited display from `BookmarkItem.Description()` in `internal/model/browse.go`: delete the `if i.b.LastVisitedAt != nil { ... } else { desc += " · Never visited" }` block entirely; the description line becomes `domain · YYYY-MM-DD [· #tags]`
- [X] T005 [P] [US2] Remove the `_ = s.UpdateLastVisited(b.ID)` call from `openBookmarkCmd` in `internal/model/app.go`; update the function comment to reflect that the command now only opens the URL and reloads the list (no visit recording)

---

## Phase 4: Polish & Cross-Cutting Concerns

**Goal**: Final build and vet pass to confirm zero errors across all changes.

- [X] T006 Final build and vet: run `go build ./...` and `go vet ./...`; confirm zero errors

---

## Dependencies

```
T001 → T002 → T003
T001 → T004, T005 (parallel, independent of T002/T003)
T003, T004, T005 → T006
```

Note: US1 (T002, T003) and US2 (T004, T005) touch different files and can be developed in parallel after T001.

## Parallel Execution Examples

**US1 and US2 are fully independent** — they touch different files:
- Agent A: `internal/store/archive.go` (T002, T003)
- Agent B: `internal/model/browse.go` + `internal/model/app.go` (T004, T005)

**Within US2** — T004 and T005 touch different files and can run simultaneously:
- Agent B1: `internal/model/browse.go` (T004)
- Agent B2: `internal/model/app.go` (T005)

## Implementation Strategy

**MVP scope (US1 only)**: T001 → T002 → T003. Delivers correct expiry threshold immediately; last-visited labels remain until US2 is applied (acceptable interim state — labels may show stale data but don't affect correctness).

**Recommended order for single-agent execution**: T001, T002, T003, T004, T005, T006.
