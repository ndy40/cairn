# Vicinae Extension Contract: Edit Bookmark

**Feature**: 014-edit-bookmark

## API Function

```typescript
export async function bmEdit(
  id: number,
  url?: string,
  tags?: string,
): Promise<{ exitCode: number; stderr: string }>
```

### Parameters

| Parameter | Type     | Required | Description                     |
| --------- | -------- | -------- | ------------------------------- |
| `id`      | number   | Yes      | Bookmark ID to edit             |
| `url`     | string   | No       | New URL (omit to keep current)  |
| `tags`    | string   | No       | Comma-separated tags (omit to keep current) |

### Return Value

| Field      | Type   | Description                      |
| ---------- | ------ | -------------------------------- |
| `exitCode` | number | 0 = success, 1 = not found/dup, 3 = validation error |
| `stderr`   | string | Error message if non-zero exit   |

### Behaviour

- Calls: `cairn edit <id> [--url=<url>] [--tags=<tags>]`
- Invalidates list cache on success (exitCode 0)
- At least one of `url` or `tags` must be provided

## UI: Edit Action

Added to bookmark list and search result items as an Action button.

### User Flow

1. User selects a bookmark in list/search view
2. User triggers "Edit Bookmark" action
3. A form opens with pre-filled URL and tags fields
4. User modifies fields and submits
5. Extension calls `bmEdit()` with changed values
6. On success: toast notification, list refreshes
7. On error: form field error displayed (e.g. "Duplicate URL", "Not found")

### Form Fields

| Field | Type   | Pre-filled | Validation       |
| ----- | ------ | ---------- | ---------------- |
| URL   | string | Current URL| Required, non-empty |
| Tags  | string | Current tags (comma-joined) | Optional |
