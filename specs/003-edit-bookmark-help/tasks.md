# Tasks: Edit Bookmark Tags, Last-Visited Visibility & CLI Help

**Feature**: 003-edit-bookmark-help
**Branch**: `003-edit-bookmark-help`
**Spec**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md)
**Generated**: 2026-03-06

---

## Summary

| Phase | Story | Tasks | Parallelisable |
|-------|-------|-------|----------------|
| 1 – Setup | — | 1 | 0 |
| 2 – Foundational | — | 1 | 0 |
| 3 – US1: Edit Tags | US1 | 4 | 2 |
| 4 – US2: Last Visited | US2 | 1 | 0 |
| 5 – US3: CLI Help | US3 | 1 | 0 |
| 6 – Polish | — | 2 | 0 |
| **Total** | | **10** | |

---

## Phase 1: Setup

**Goal**: Verify project builds cleanly on current branch before changes.

- [X] T001 Verify clean build: run `go build ./...` and `go vet ./...` from repo root

---

## Phase 2: Foundational

**Goal**: Add the store method that US1 depends on. Must complete before Phase 3.

- [X] T002 Add `UpdateTags(id int64, tags []string) error` to `internal/store/bookmark.go` — normalises via `NormaliseTags()`, JSON-encodes, executes `UPDATE bookmarks SET tags = ? WHERE id = ?`

---

## Phase 3: US1 — Edit Tags on Existing Bookmarks

**Story goal**: User presses `e` in browse mode, an edit panel opens pre-filled with current tags; user edits and saves with Enter (or cancels with Esc); tags are saved and list reloads.

**Independent test criteria**: After implementation, pressing `e` on a bookmark in browse mode opens the edit panel, editing tags and pressing Enter persists the change, and pressing Esc discards it.

- [X] T003 [P] [US1] Create `internal/model/edit.go`: define `EditModel` struct with a `textinput.Model` for tags and a read-only title string; implement `New(b Bookmark) EditModel`, `Update(msg tea.Msg) (EditModel, tea.Cmd)`, `View() string`, and `Tags() []string` (split on comma, return slice)
- [X] T004 [P] [US1] Add `editKey` binding (key `e`) to `internal/model/browse.go`
- [X] T005 [US1] Add `StateEdit` to the `AppState` const block in `internal/model/app.go`; add `editModel EditModel` field to `Model` struct
- [X] T006 [US1] Wire edit flow in `internal/model/app.go`: handle `editKey` press in `updateBrowse()` to enter `StateEdit`; add `updateEdit()` routing `Enter` (call `s.UpdateTags`, then `loadBookmarks`) and `Esc` (return to `StateBrowse`); add `editView()` returning `m.editModel.View()`; route `StateEdit` in `Update()` and `View()`

---

## Phase 4: US2 — Last-Visited Visibility (Documentation Only)

**Story goal**: Confirm and document that last-visited dates update immediately after opening a bookmark. No code changes required per research decision 4.

**Independent test criteria**: Manual smoke test — open a bookmark with Enter, return to browse list, verify the "Last:" date on that row has updated without restarting the app.

- [X] T007 [US2] Add a comment block in `internal/model/app.go` above `openBookmarkCmd` documenting the last-visited update flow (matches the flow described in `data-model.md` section "Confirmed: Last-Visited Update Flow")

---

## Phase 5: US3 — Per-Subcommand CLI Help

**Story goal**: `bm --help`, `bm -h`, and `bm <subcommand> --help` all print usage text and exit 0.

**Independent test criteria**: `bm --help`, `bm add --help`, `bm list --help`, `bm search --help`, `bm delete --help`, `bm version --help`, and `bm help --help` each print the expected usage text and exit with code 0.

- [X] T008 [US3] Update `cmd/bm/main.go`: (1) check `os.Args` for `-h`/`--help` before `flag.Parse()` on the root command, call `printHelp()`, exit 0; (2) for each subcommand FlagSet (`add`, `list`, `search`, `delete`, `version`, `help`) set `fs.Usage` to a function that prints the per-subcommand help text (per `contracts/cli-interface.md`) and exits 0; add a `--help` bool flag to each FlagSet and check it after `fs.Parse()`

---

## Phase 6: Polish & Cross-Cutting Concerns

**Goal**: Update the TUI help screen, do a final build+vet pass.

- [X] T009 Update the help screen text in `internal/model/app.go` `helpView()` to include `e  Edit tags on selected bookmark` in the Browse Mode section (per `contracts/keyboard-shortcuts.md`)
- [X] T010 Final build and vet: run `go build ./...` and `go vet ./...`; confirm zero errors

---

## Dependencies

```
T001 → T002 → T003, T004 (parallel)
              T003, T004 → T005 → T006
T006 → T007, T008, T009 (parallel)
T007, T008, T009 → T010
```

## Parallel Execution Examples

**Phase 3** — T003 and T004 touch different files and can be done simultaneously:
- Agent A: `internal/model/edit.go` (T003)
- Agent B: `internal/model/browse.go` (T004)

**Phase 6** — T007, T008, T009 touch different files:
- Agent A: `internal/model/app.go` comment (T007)
- Agent B: `cmd/bm/main.go` (T008)
- Agent C: `internal/model/app.go` helpView (T009) — coordinate with Agent A on same file

## Implementation Strategy

**MVP scope (US1 only)**: Complete T001 → T002 → T003 → T004 → T005 → T006.
This delivers the primary user-facing feature (edit tags). US2 requires no code. US3 is a polish task that can follow independently.

**Suggested order for single-agent execution**: T001, T002, T003, T004, T005, T006, T007, T008, T009, T010.
