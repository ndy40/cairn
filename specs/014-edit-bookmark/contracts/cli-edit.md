# CLI Contract: `cairn edit`

**Feature**: 014-edit-bookmark

## Command Syntax

```
cairn edit <id> [--url=<url>] [--title=<title>] [--tags=<tags>]
```

### Arguments

| Argument | Required | Description              |
| -------- | -------- | ------------------------ |
| `<id>`   | Yes      | Numeric bookmark ID      |

### Flags

| Flag      | Type   | Default | Description                        |
| --------- | ------ | ------- | ---------------------------------- |
| `--url`   | string | ""      | New URL for the bookmark           |
| `--title` | string | ""      | New title for the bookmark         |
| `--tags`  | string | ""      | Comma-separated tags (max 3)       |

At least one of `--url`, `--title`, or `--tags` must be provided.

## Exit Codes

| Code | Meaning                          | Stderr                    |
| ---- | -------------------------------- | ------------------------- |
| 0    | Success                          | —                         |
| 1    | Bookmark not found / duplicate   | "Not found" / "Duplicate URL" |
| 3    | Validation error / usage error   | Error message             |

## Stdout Messages

| Condition           | Output                     |
| ------------------- | -------------------------- |
| URL only            | `Updated URL`              |
| Title only          | `Updated title`            |
| Tags only           | `Updated tags`             |
| Multiple fields     | `Updated URL, title and tags` (or subset) |

## Examples

```bash
# Edit URL only
cairn edit 5 --url=https://new-example.com

# Edit tags only (existing behaviour)
cairn edit 5 --tags=go,dev

# Edit URL and tags together
cairn edit 5 --url=https://new-example.com --tags=go,dev

# Edit all three
cairn edit 5 --url=https://new-example.com --title="New Title" --tags=go,dev
```

## Error Cases

```bash
# No flags provided
$ cairn edit 5
# stderr: "edit: specify at least one of --url, --title, or --tags"
# exit: 3

# Empty URL
$ cairn edit 5 --url=""
# stderr: "edit: --url cannot be empty"
# exit: 3

# Duplicate URL
$ cairn edit 5 --url=https://already-exists.com
# stderr: "Duplicate URL"
# exit: 1

# Not found
$ cairn edit 999 --title="New"
# stderr: "Not found"
# exit: 1
```
