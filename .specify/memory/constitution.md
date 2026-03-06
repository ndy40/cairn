# Project Constitution

**Project**: bookmark-manager
**Last Updated**: 2026-03-06

---

## Principles

These rules apply to ALL features in this project, enforced at every `/speckit.plan` Constitution Check.

---

## Task Management

**RULE: All tasks MUST be created and managed via the Backlog CLI.**

- After generating `tasks.md` with `/speckit.tasks`, ALL tasks in that file MUST be created in Backlog using `backlog task create`.
- Never manage tasks by directly editing files in `backlog/tasks/`.
- Use `backlog task edit` to update task status, not direct file edits.
- The `tasks.md` file serves as the design artifact (source of truth for task descriptions and ordering). Backlog is the operational tracker.
- When a task is completed during implementation, mark it done with `backlog task complete <id>`.

---

## Architecture

- **Zero CGO**: All Go dependencies must be pure Go to enable `CGO_ENABLED=0` static builds.
- **Single binary**: No external runtime, server, or database process.
- **No unnecessary abstractions**: New packages must have a single clear responsibility. Three similar lines of code are better than a premature abstraction.

---

## Technology Stack

- **Language**: Go 1.22+
- **TUI**: charmbracelet/bubbletea + bubbles + lipgloss
- **Storage**: modernc.org/sqlite (pure Go, WAL mode, FTS5)
- **Schema migrations**: Versioned, applied automatically at startup via `schema_version` table

---

## Constitution Check Gates

Every `/speckit.plan` MUST verify:

| Gate | Requirement |
|------|-------------|
| No CGO | All dependencies pure Go |
| Single binary | No external runtime |
| Task management | Tasks must go through Backlog CLI after `/speckit.tasks` |
| Backward-compatible migrations | ALTER TABLE with DEFAULT values only; no destructive schema changes |
