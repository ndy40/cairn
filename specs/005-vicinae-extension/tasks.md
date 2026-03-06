# Tasks: Vicinae Extension for Bookmark Manager

**Feature**: 005-vicinae-extension
**Branch**: `005-vicinae-extension`
**Spec**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md)
**Generated**: 2026-03-06

---

## Summary

| Phase | Story | Tasks | Parallelisable |
|-------|-------|-------|----------------|
| 1 – Setup | — | 2 | 0 |
| 2 – Foundational | — | 2 | 1 |
| 3 – US1: Search Bookmarks | US1 | 2 | 0 |
| 4 – US2: List Bookmarks | US2 | 1 | 0 |
| 5 – US3: Add Bookmark | US3 | 2 | 0 |
| 6 – Polish | — | 2 | 1 |
| **Total** | | **11** | |

---

## Phase 1: Setup

**Goal**: Bootstrap the extension package and verify the Go CLI builds cleanly.

- [X] T001 Verify Go CLI builds cleanly: run `go build ./...` and `go vet ./...` from repo root
- [X] T002 Create `vicinae-extension/` directory at repository root with `package.json`, `tsconfig.json`, and `src/` directory; `package.json` must declare: `name: "bm-bookmarks"`, `title: "Bookmark Manager"`, three commands (`search-bookmarks → src/search.tsx`, `list-bookmarks → src/list.tsx`, `add-bookmark → src/add.tsx`), dependency `@vicinae/api`, devDependency `typescript`; `tsconfig.json` targets ESNext with JSX support; per `specs/005-vicinae-extension/contracts/extension-commands.md`

---

## Phase 2: Foundational

**Goal**: Shared utilities that all three commands depend on, plus the CLI prerequisite.

- [X] T003 Add `--tags` flag to `bm add` subcommand in `cmd/bm/main.go`: create a FlagSet for add, add `--tags string` flag, parse remaining args after the URL positional argument, split on comma and pass to `s.Insert()`; update `printAddHelp()` to document the flag; per `specs/005-vicinae-extension/contracts/cli-interface.md`
- [X] T004 [P] Create `vicinae-extension/src/bm.ts`: a shared utility module that wraps `child_process.spawnSync` to call the `bm` CLI; export three functions — `bmList(): Bookmark[]`, `bmSearch(query: string): Bookmark[]`, `bmAdd(url: string, tags?: string): { exitCode: number; stderr: string }`; define and export a `Bookmark` TypeScript interface matching the JSON schema in `specs/005-vicinae-extension/data-model.md`; include a `bmAvailable(): boolean` check using `which bm`

---

## Phase 3: US1 — Search Bookmarks

**Story goal**: User invokes "Search Bookmarks" from the launcher, types a query, results filter in real time, pressing Enter opens the selected bookmark in the browser.

**Independent test criteria**: Open the Vicinae launcher, invoke "Search Bookmarks", type a keyword, verify matching bookmark rows appear; press Enter on a result and verify the browser opens to that URL. With an unknown query, verify "No bookmarks found" empty state.

- [X] T005 [US1] Create `vicinae-extension/src/search.tsx`: implement the Search Bookmarks command using a `List` component from `@vicinae/api`; on empty query call `bmList()`, on non-empty query call `bmSearch(query)`; display each result as a `List.Item` with title (fallback to URL), subtitle (domain), and accessories (date, tags, pin badge); primary action opens URL in browser, secondary action copies URL to clipboard; show `bmAvailable()` error state if CLI not found; per `specs/005-vicinae-extension/contracts/extension-commands.md`
- [ ] T006 [US1] Install extension dependencies and verify `vici develop` starts without errors for the search command: run `npm install` in `vicinae-extension/` and confirm no TypeScript compilation errors — NOTE: requires Vicinae SDK (`@vicinae/api`) to be installed locally via `vici develop` toolchain

---

## Phase 4: US2 — List Bookmarks

**Story goal**: User invokes "List Bookmarks", sees all active bookmarks newest-first, and can open any by pressing Enter. Typing in the launcher filters the list client-side.

**Independent test criteria**: Open the launcher, invoke "List Bookmarks", verify the full bookmark list loads. Type a word in the launcher — verify the list filters without an additional CLI call. Press Enter on an item and verify the browser opens.

- [X] T007 [US2] Create `vicinae-extension/src/list.tsx`: implement the List Bookmarks command using a `List` component; call `bmList()` once on mount; render all results with the same `List.Item` structure as search.tsx (title, subtitle, date, tags, pin); primary action opens URL, secondary copies URL; show CLI error state if `bmAvailable()` returns false; Vicinae's built-in search box handles client-side filtering with no additional CLI call; per `specs/005-vicinae-extension/contracts/extension-commands.md`

---

## Phase 5: US3 — Add Bookmark

**Story goal**: User invokes "Add Bookmark", sees a form pre-filled with clipboard URL (if valid), optionally adds tags, submits, and receives a success or error message.

**Independent test criteria**: Copy a URL to clipboard, invoke "Add Bookmark", verify the URL field is pre-filled. Submit → verify "Saved" confirmation and that the bookmark appears in "List Bookmarks". Submit a duplicate URL → verify "Already bookmarked" error. Submit empty URL → verify validation blocks submission.

- [X] T008 [US3] Create `vicinae-extension/src/add.tsx`: implement the Add Bookmark command using a `Form` component from `@vicinae/api`; URL field (required, pre-filled from clipboard if valid URL) and Tags field (optional, placeholder "work, go, tools  (comma-separated, max 3)"); on submit: validate URL is non-empty and starts with `http://` or `https://`; call `bmAdd(url, tags)`; on exit 0 or 2 show success toast and close; on exit 1 show inline "Already bookmarked" error; on exit 3 show inline error with stderr; per `specs/005-vicinae-extension/contracts/extension-commands.md`
- [X] T009 [US3] Verify end-to-end add flow: run `go build ./...` to confirm `bm add --tags` compiles; manually test `bm add https://example.com --tags "test"` in terminal and confirm the bookmark is saved with the tag

---

## Phase 6: Polish & Cross-Cutting Concerns

**Goal**: CLI help update, final build verification across both Go and TypeScript.

- [X] T010 [P] Run `go build ./...` and `go vet ./...` from repo root to confirm the `--tags` CLI change is clean
- [ ] T011 [P] Run `npm run build` in `vicinae-extension/` to produce the production bundle; confirm zero TypeScript compilation errors — NOTE: requires Vicinae SDK to be installed via `vici` toolchain

---

## Dependencies

```
T001 → T002 → T003, T004 (parallel)
T003, T004 → T005 → T006
T006 → T007
T007 → T008 → T009
T009 → T010, T011 (parallel)
```

Note: T004 (bm.ts utility) is a shared dependency for T005, T007, and T008.

## Parallel Execution Examples

**Phase 2** — T003 (Go CLI) and T004 (TS utility module) touch different files:
- Agent A: `cmd/bm/main.go` (T003)
- Agent B: `vicinae-extension/src/bm.ts` (T004)

**Phase 6** — T010 (Go build) and T011 (TS build) are independent:
- Agent A: Go repo root (T010)
- Agent B: `vicinae-extension/` (T011)

## Implementation Strategy

**MVP scope (US1 only)**: T001 → T002 → T003 → T004 → T005 → T006.
Delivers a working "Search Bookmarks" command — the highest-value launcher interaction. US2 (list) and US3 (add) can be added incrementally.

**Recommended single-agent order**: T001, T002, T003, T004, T005, T006, T007, T008, T009, T010, T011.
