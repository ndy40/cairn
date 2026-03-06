# Implementation Plan: TUI Bookmark Manager

**Branch**: `001-tui-bookmark-manager` | **Date**: 2026-03-06 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-tui-bookmark-manager/spec.md`

---

## Summary

A single-user, keyboard-driven TUI application for saving and finding web bookmarks. The user presses Ctrl+P to paste a URL from the clipboard; the app fetches the page's title and meta description, saves the bookmark to a local SQLite database, and makes it instantly searchable via fuzzy search across title, domain, and description fields. Built in Go as a single statically linked binary with zero CGO dependencies.

---

## Technical Context

**Language/Version**: Go 1.22+
**Primary Dependencies**:
- `charmbracelet/bubbletea` + `charmbracelet/bubbles` + `charmbracelet/lipgloss` — TUI framework (MVU architecture)
- `modernc.org/sqlite` — Pure Go embedded SQLite with FTS5
- `PuerkitoBio/goquery` + `golang.org/x/net/html/charset` — HTML meta tag extraction
- `sahilm/fuzzy` — Fuzzy search scoring (Sublime Text-style algorithm)
- `atotto/clipboard` — Clipboard read (zero CGO)

**Storage**: SQLite via `modernc.org/sqlite` (pure Go, no CGO). Database at `$XDG_DATA_HOME/bookmark-manager/bookmarks.db` (Linux), `~/Library/Application Support/bookmark-manager/bookmarks.db` (macOS).

**Testing**: `go test ./...` (stdlib testing package). No external test framework required.

**Target Platform**: Linux, macOS, Windows (single binary per platform). Primary development target: Linux.

**Project Type**: CLI/TUI application

**Performance Goals**:
- Fuzzy search results update within 100ms of each keystroke for ≤1,000 bookmarks (per SC-002)
- Application launches and displays list in under 2 seconds for ≤1,000 bookmarks (per SC-006)
- Bookmark save (clipboard-to-confirmed) completes in under 5 seconds on responsive network (per SC-001)

**Constraints**:
- Zero CGO — all dependencies must be pure Go to enable cross-compilation and simple deployment
- Single binary distribution — no external runtime, server, or database process
- Linux clipboard requires `xclip` or `xsel` installed on host (documented prerequisite)

**Scale/Scope**: Single user, local machine, up to ~1,000 bookmarks initial target. No sync, no multi-user, no cloud.

---

## Constitution Check

No `.specify/memory/constitution.md` exists for this project. This is a greenfield repository with no prior architectural constraints. The following lightweight gates apply by default:

| Gate | Status | Notes |
|------|--------|-------|
| No CGO dependencies | PASS | All 8 selected dependencies are pure Go |
| Single binary deployment | PASS | Achieved via `CGO_ENABLED=0 go build` |
| No unnecessary abstractions | PASS | Internal packages map 1:1 to responsibility areas |
| No external runtime dependencies | PASS | SQLite embedded via modernc; no database server |

No violations. No Complexity Tracking table required.

---

## Project Structure

### Documentation (this feature)

```text
specs/001-tui-bookmark-manager/
├── spec.md              # Feature specification (/speckit.specify output)
├── plan.md              # This file (/speckit.plan output)
├── research.md          # Phase 0 research findings
├── data-model.md        # Entity definitions and storage schema
├── quickstart.md        # Developer setup guide
├── contracts/
│   ├── keyboard-shortcuts.md   # TUI keyboard interface contract
│   └── cli-interface.md        # CLI command/flag contract
├── checklists/
│   └── requirements.md  # Spec quality checklist
└── tasks.md             # Phase 2 output (/speckit.tasks — NOT created here)
```

### Source Code (repository root)

```text
bookmark-manager/
├── cmd/
│   └── bm/
│       └── main.go              # Entry point: CLI arg parsing, TUI launch or subcommand dispatch
├── internal/
│   ├── model/
│   │   ├── app.go               # Root bubbletea Model: state machine (browse/search/add/confirm-delete)
│   │   ├── browse.go            # Browse mode: list navigation, shortcut handling
│   │   ├── search.go            # Search mode: text input wired to fuzzy filter
│   │   └── add.go               # Add modal: URL input, fetch trigger, save confirmation
│   ├── store/
│   │   ├── store.go             # SQLite open, WAL mode, schema migration runner
│   │   ├── bookmark.go          # CRUD: insert, delete, list, get-by-id, check-duplicate
│   │   └── search.go            # FTS5 pre-filter query returning candidate bookmark IDs
│   ├── fetcher/
│   │   └── fetcher.go           # HTTP GET + goquery: extract title, og:title, description, og:description
│   ├── search/
│   │   └── fuzzy.go             # Multi-field fuzzy wrapper: runs sahilm/fuzzy per field, merges with weights
│   └── clipboard/
│       └── clipboard.go         # Thin wrapper: reads text from system clipboard via atotto/clipboard
├── go.mod
└── go.sum
```

**Structure Decision**: Single-project layout (Option 1). One `cmd/bm` entry point, feature logic in `internal/` sub-packages, each with a single clear responsibility. No public API surface — all packages are `internal/`.

---

## Phase 0 Output

- [x] research.md — All technical decisions resolved, all NEEDS CLARIFICATION cleared

## Phase 1 Output

- [x] data-model.md — Bookmark entity, SQLite schema, FTS5 setup, two-stage search architecture
- [x] contracts/keyboard-shortcuts.md — Full keyboard interface, mode transitions, status bar content
- [x] contracts/cli-interface.md — Binary name, commands, flags, exit codes, env vars
- [x] quickstart.md — Prerequisites, project init, build, run, test, database location
- [ ] tasks.md — Generated by `/speckit.tasks` (next step)
