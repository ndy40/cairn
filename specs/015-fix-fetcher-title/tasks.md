# Tasks: Fix Fetcher Title Extraction (015)

**Input**: `specs/015-fix-fetcher-title/`
**Prerequisites**: plan.md ✓, spec.md ✓

**Organization**: Tasks are grouped by user story. US1 (priority swap) is independent of US2 (header improvements) and can land in isolation.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2)

---

## Phase 1: Setup

**Purpose**: Confirm test harness and baseline before making changes.

- [ ] T001 Confirm `go test ./internal/fetcher/...` runs without errors (baseline green) in `internal/fetcher/`

---

## Phase 2: User Story 1 — Priority Swap (`<title>` over `og:title`) (Priority: P1) 🎯 MVP

**Goal**: Make `extractTitle` return the HTML `<title>` tag value when present, falling back to `og:title`, then hostname.

**Independent Test**: A page that has both `<title>My Page</title>` and `<meta property="og:title" content="OG Title">` must resolve to `"My Page"`.

### Implementation for User Story 1

- [ ] T002 [US1] Reorder `extractTitle` in `internal/fetcher/fetcher.go` — check `<title>` first, then `og:title`, then fallback
- [ ] T003 [P] [US1] Write unit tests for `extractTitle` covering all three branches (title-only, og-only, neither) in `internal/fetcher/fetcher_test.go`

**Checkpoint**: `go test ./internal/fetcher/...` passes; `extractTitle` prefers `<title>` over `og:title`.

---

## Phase 3: User Story 2 — Browser-Like Headers for Bot-Gated Pages (Priority: P2)

**Goal**: The fetcher sends realistic browser headers so that YouTube and similar sites return HTML with a proper `<title>` tag instead of a minimal bot-response.

**Independent Test**: `Fetch("https://www.youtube.com/watch?v=dQw4w9WgXcQ")` returns a non-empty title that is not `"youtube.com"`.

### Implementation for User Story 2

- [ ] T004 [US2] Replace the `userAgent` constant in `internal/fetcher/fetcher.go` with a realistic browser UA string (Chrome/Linux)
- [ ] T005 [US2] Add `Accept` and `Accept-Language` request headers in the `Fetch` function in `internal/fetcher/fetcher.go`
- [ ] T006 [P] [US2] Add or extend unit/integration tests in `internal/fetcher/fetcher_test.go` to assert non-empty title for a YouTube-style response fixture

**Checkpoint**: Manual smoke test `bm add https://www.youtube.com/watch?v=<id>` stores the video title, not `youtube.com`.

---

## Phase 4: Polish & Cross-Cutting Concerns

- [ ] T007 [P] Run `go build ./... && go vet ./...` and confirm clean build
- [ ] T008 [P] Run `go test ./...` to confirm no regressions across the whole project

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — run first.
- **Phase 2 (US1)**: Depends on Phase 1 only. Independent of US2.
- **Phase 3 (US2)**: Depends on Phase 1. Can start after Phase 1; integrates cleanly with the reordered `extractTitle` from US1.
- **Phase 4 (Polish)**: Depends on Phases 2 and 3.

### User Story Dependencies

- **US1 (T002–T003)**: No dependency on US2 — can be merged separately as a minimal fix.
- **US2 (T004–T006)**: No dependency on US1, but benefits from it (better headers + correct priority together solve the YouTube problem end-to-end).

### Parallel Opportunities

- T003 (write tests) and T002 (reorder code) can be written in parallel since they touch the same file but are independent edits.
- T004 and T005 can be implemented in a single commit — they are both small edits to `fetcher.go`.
- T007 and T008 (build/test checks) can run in parallel.

---

## Parallel Example: User Story 1

```bash
# Both can be drafted in parallel:
Task T002: Reorder extractTitle in internal/fetcher/fetcher.go
Task T003: Write unit tests for all three branches in internal/fetcher/fetcher_test.go
```

---

## Implementation Strategy

### MVP First (US1 Only)

1. Complete Phase 1: baseline check.
2. Complete Phase 2: swap `<title>` / `og:title` priority (T002 + T003).
3. **STOP and VALIDATE**: all tests pass, `extractTitle` behaves correctly.
4. This alone fixes pages where `<title>` contains the real title but `og:title` is absent or misleading.

### Full Fix (US1 + US2)

1. After US1 lands, add US2 header improvements (T004–T006).
2. Together, these solve the YouTube case end-to-end: browser headers ensure the video title appears in the `<title>` tag, and the priority swap ensures it is picked over any `og:title`.

---

## Notes

- Total tasks: **8** (T001–T008)
- US1 tasks: **2** (T002–T003)
- US2 tasks: **3** (T004–T006)
- Parallel opportunities: T002‖T003, T004‖T005, T007‖T008
- No new dependencies or schema changes required
- All changes are confined to `internal/fetcher/fetcher.go` and `internal/fetcher/fetcher_test.go`
