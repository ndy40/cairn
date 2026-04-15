# Contract: Vicinae Extension Edit Form (017)

**Surface**: `vicinae-extension/src/bm-edit.tsx` + `vicinae-extension/src/bm.ts`

---

## Form layout

```
┌─ Edit Bookmark ───────────────────────────────────────┐
│                                                        │
│  Title *                                               │
│  ┌──────────────────────────────────────────────────┐ │
│  │ Never Gonna Give You Up                          │ │
│  └──────────────────────────────────────────────────┘ │
│  [Title cannot be empty]  ← shown on submit if blank  │
│                                                        │
│  URL                                                   │
│  ┌──────────────────────────────────────────────────┐ │
│  │ https://youtube.com/watch?v=dQw4w9WgXcQ          │ │
│  └──────────────────────────────────────────────────┘ │
│                                                        │
│  Tags                                                  │
│  ┌──────────────────────────────────────────────────┐ │
│  │ music, classic                                   │ │
│  └──────────────────────────────────────────────────┘ │
│                                                        │
│  [Save Changes]                                        │
└────────────────────────────────────────────────────────┘
```

---

## `bmEdit()` updated signature — `bm.ts`

```typescript
export async function bmEdit(
  id: number,
  url?: string,
  tags?: string,
  title?: string,           // ← new optional parameter
): Promise<{ exitCode: number; stderr: string }>
```

**CLI args constructed**:

```
cairn edit <id>
  [--url <url>]      if url !== undefined && url.trim() !== ""
  [--tags <tags>]    if tags !== undefined
  [--title <title>]  if title !== undefined && title.trim() !== ""
```

---

## Change detection contract — `bm-edit.tsx`

| Condition | `bmEdit` called? | `--title` passed? |
|---|---|---|
| Title changed, URL and tags unchanged | Yes | Yes |
| Title unchanged, URL or tags changed | Yes | No |
| Title, URL, and tags all unchanged | No (show "No changes" toast) | — |
| Title empty or whitespace on submit | No (show field error) | — |

---

## Validation contract

| Condition | Behaviour |
|---|---|
| Title empty / whitespace-only on submit | Set `titleError` state; do not call `bmEdit` |
| Title > 500 characters | Blocked at field level via `maxLength={500}` |
| Bookmark not found (exitCode 1, stderr contains "not found") | Show toast error "Bookmark not found" |
