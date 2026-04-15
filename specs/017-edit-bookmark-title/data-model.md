# Data Model: Edit Bookmark Title (017)

## No schema changes required

The `bookmarks` table already has a `title TEXT` column. The `store.BookmarkPatch` struct already has a `Title *string` field. `store.UpdateFields()` already writes `title` to the database when the patch includes it.

---

## Existing Bookmark entity (reference)

| Field | Type | Notes |
|---|---|---|
| `id` | INTEGER PK | Auto-increment |
| `uuid` | TEXT | Stable external identifier |
| `url` | TEXT NOT NULL UNIQUE | Normalised URL |
| `domain` | TEXT | Extracted hostname |
| **`title`** | TEXT | **Target of this feature — already writable** |
| `description` | TEXT | Meta description |
| `tags` | TEXT | JSON array |
| `created_at` | TEXT | ISO-8601 |
| `updated_at` | TEXT | ISO-8601, updated on every `UpdateFields` call |
| `is_permanent` | INTEGER | Pinned flag |
| `is_archived` | INTEGER | Archive flag |
| `archived_at` | TEXT | Nullable |

## Validation rules (application layer)

| Rule | Enforcement point |
|---|---|
| Title must not be empty or whitespace-only | TUI `EditModel.Save()` + extension `handleSubmit()` |
| Title maximum 500 characters | TUI `textinput.CharLimit = 500` + extension `maxLength={500}` |
