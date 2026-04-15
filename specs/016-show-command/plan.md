# Implementation Plan: `cairn show` Command (016)

## Tech Stack

- **Language**: Go 1.25.0, `CGO_ENABLED=0`
- **Libraries**: all existing — `internal/store`, `internal/sync`, `internal/config`; no new dependencies
- **CLI pattern**: mirrors existing commands — handler in `cmd/cairn/commands.go`, runner in `cmd/cairn/run_show.go`

---

## Affected Files

| File | Change |
|---|---|
| `internal/store/bookmark.go` | Add `Count() (int64, error)` method |
| `cmd/cairn/run_show.go` | New file — `runShow(db string, cfgManager)` prints the four fields |
| `cmd/cairn/commands.go` | Add `cmdShow(ctx cmdContext)` handler |
| `cmd/cairn/main.go` | Register `"show"` in the `commands` map |

---

## Implementation Notes

### `Store.Count()`
Simple `SELECT COUNT(1) FROM bookmarks` query on the existing `*sql.DB`. Matches the pattern of `ExistsByURL` in `internal/store/bookmark.go`.

### `runShow`
1. Open store (`openStore(db)`), defer close.
2. Call `s.Count()` for bookmark count.
3. Load sync config via `csync.LoadConfig(csync.DefaultConfigPath())`.
4. Determine sync status string using `csync.IsConfigured(cfg)`.
5. Print four labelled lines (aligned with `%-12s` format like `cmdConfig`):

```
Database:    /home/user/.local/share/bookmark-manager/bookmarks.db
Bookmarks:   42
Sync:        configured (dropbox)
Last sync:   2026-04-10 08:23:11 UTC
```

### `cmdShow`
Thin dispatcher in `commands.go` — checks `-h`/`--help`, then calls `runShow(ctx.db, ctx.cfgManager)`. No extra flags needed.

### Registration
Add `"show": {run: cmdShow, autoSync: false}` to the `commands` map in `main.go`. No auto-sync; `show` is a read-only introspection command.
