# CLI Contract: `cairn pin`

## Command

```
cairn pin <id>
```

## Description

Toggles the `is_permanent` (pin) flag on the bookmark identified by the given numeric ID.
If the bookmark is currently unpinned, it becomes pinned. If it is pinned, it becomes unpinned.

## Arguments

| Argument | Type    | Required | Description                        |
|----------|---------|----------|------------------------------------|
| `id`     | integer | yes      | Numeric primary key of the bookmark |

## Exit Codes

| Code | Meaning                                |
|------|----------------------------------------|
| 0    | Success — pin state toggled            |
| 1    | Bookmark not found                     |
| 3    | Unexpected error (DB error, bad input) |

## Stdout

On success (exit 0): `Pinned: "<title>" (<domain>)` or `Unpinned: "<title>" (<domain>)`.

## Stderr

On error (exit 1 or 3): human-readable error message.

## Examples

```bash
$ cairn pin 42
Pinned: "Go spec" (go.dev)

$ cairn pin 42
Unpinned: "Go spec" (go.dev)

$ cairn pin 9999
Bookmark not found
# exit code 1
```
