# Contract: CLI Delete Command

**Feature**: 006-vicinae-delete-bookmark
**Date**: 2026-03-11

## Existing Contract (no changes)

The `cairn delete <id>` CLI command is already implemented and used by the extension.

### bmDelete() TypeScript Wrapper

```typescript
function bmDelete(id: number): { exitCode: number; stderr: string }
```

| Parameter | Type   | Description               |
|-----------|--------|---------------------------|
| `id`      | number | Bookmark ID from list/search JSON output |

| Exit Code | Meaning                  | Extension Behavior          |
|-----------|--------------------------|-----------------------------|
| 0         | Successfully deleted     | Show success toast, refresh list |
| 1         | Bookmark not found       | Show error toast: "Bookmark not found" |
| 3         | Error                    | Show error toast with stderr message |

### Confirmation Dialog

```typescript
confirmAlert({
  title: "Delete Bookmark?",
  message: `"${bookmark.Title || bookmark.URL}" will be permanently deleted.`,
  primaryAction: {
    title: "Delete",
    style: Alert.ActionStyle.Destructive,
  },
})
```

Returns `Promise<boolean>` — `true` if user confirms, `false` if cancelled.
