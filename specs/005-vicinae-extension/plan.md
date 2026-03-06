# Implementation Plan: Vicinae Extension for Bookmark Manager

**Branch**: `005-vicinae-extension` | **Date**: 2026-03-06 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/005-vicinae-extension/spec.md`

---

## Summary

A new TypeScript/React Vicinae extension providing three launcher commands — Search Bookmarks, List Bookmarks, and Add Bookmark — all delegating to the `bm` CLI for data access. Additionally, `bm add` gains a `--tags` flag so the extension can pass tags when saving a bookmark. The extension lives in `vicinae-extension/` at the repository root.

---

## Technical Context

**Language/Version**: TypeScript 5.x (extension); Go 1.22+ (one-line CLI change only)
**Primary Dependencies**:
- Extension: `@vicinae/api` (Vicinae React SDK), React
- CLI change: no new dependencies
**Storage**: None in the extension (all via `bm` CLI); SQLite unchanged
**Testing**: `vici develop` for manual hot-reload; `go test ./...` for CLI change
**Target Platform**: Linux (Vicinae runs on Linux)
**Project Type**: Vicinae launcher extension (TypeScript) + minor CLI patch (Go)
**Performance Goals**: Search results appear within 1 second of typing; list loads within 500 ms
**Constraints**: Extension must not access the SQLite database directly; all data flows through `bm` CLI
**Scale/Scope**: Single user, local (same as main app)

---

## Constitution Check

| Gate | Status | Notes |
|------|--------|-------|
| No CGO dependencies | PASS | Extension is TypeScript; Go CLI change adds no new Go dependencies |
| Single binary | PASS | The extension is a separate Vicinae package, not bundled into `bm` |
| Task management | PASS | Tasks through Backlog CLI after `/speckit.tasks` |
| Backward-compatible migrations | PASS | No schema changes |
| **Language: Go only** | JUSTIFIED EXCEPTION | Vicinae extensions MUST be TypeScript — there is no Go SDK. The main `bm` binary remains Go. The extension is a companion launcher package, not part of the core tool. |

---

## Project Structure

### Documentation (this feature)

```text
specs/005-vicinae-extension/
├── spec.md
├── plan.md              # This file
├── research.md
├── data-model.md
├── contracts/
│   ├── extension-commands.md
│   └── cli-interface.md
├── checklists/
│   └── requirements.md
└── tasks.md             # /speckit.tasks output (NOT created here)
```

### Source Code (new + modified)

```text
bookmark-manager/
├── cmd/bm/
│   └── main.go               # MODIFIED: add --tags flag to bm add subcommand
│
└── vicinae-extension/         # NEW: standalone TypeScript package
    ├── package.json           # Extension manifest + @vicinae/api dependency
    ├── tsconfig.json          # TypeScript configuration
    └── src/
        ├── search.tsx         # "Search Bookmarks" command
        ├── list.tsx           # "List Bookmarks" command
        └── add.tsx            # "Add Bookmark" command
```

**Structure Decision**: Extension is a self-contained package at repository root. It is not nested inside the Go module and has its own `node_modules`. The Go CLI is modified in exactly one place.

---

## Phase 0 Output

- [x] research.md

## Phase 1 Output

- [x] data-model.md
- [x] contracts/extension-commands.md
- [x] contracts/cli-interface.md
- [ ] tasks.md — `/speckit.tasks` (next step)
