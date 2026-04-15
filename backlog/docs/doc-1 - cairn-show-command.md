---
id: doc-1
title: cairn show command
type: other
created_date: '2026-04-15 15:11'
---
# cairn show

Displays a summary of the current cairn environment.

## Usage

```
cairn show
```

## Output

| Field       | Description                                      |
|-------------|--------------------------------------------------|
| Database    | Resolved path to the SQLite bookmark database    |
| Bookmarks   | Total number of bookmarks stored                 |
| Sync        | Sync status: `not configured` or `configured (<backend>)` |
| Last sync   | UTC timestamp of the last successful sync, or `never` |

## Example

```
$ cairn show
Database:    /home/user/.local/share/cairn/bookmarks.db
Bookmarks:   42
Sync:        configured (dropbox)
Last sync:   2026-04-14 08:30:00 UTC
```

## Notes

- `cairn show` does **not** trigger a background sync (unlike `add`, `list`, etc.).
- Sync fields are omitted when sync is not configured.
- Use `cairn sync status` for a more detailed sync report.
