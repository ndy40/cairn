# Contract: TUI Edit Panel (017)

**Surface**: `internal/model/edit.go` — `EditModel`

---

## Field layout (tab order)

```
┌─ Edit Bookmark ──────────────────────────────────────┐
│                                                       │
│  Title                          ← field 0 (focused)  │
│  ┌─────────────────────────────────────────────────┐ │
│  │ Never Gonna Give You Up                         │ │
│  └─────────────────────────────────────────────────┘ │
│  [error: Title cannot be empty]  ← shown on bad save │
│                                                       │
│  URL                                                  │
│  ┌─────────────────────────────────────────────────┐ │
│  │ https://youtube.com/watch?v=dQw4w9WgXcQ         │ │
│  └─────────────────────────────────────────────────┘ │
│                                                       │
│  Tags                                                 │
│  ┌─────────────────────────────────────────────────┐ │
│  │ music, classic                                  │ │
│  └─────────────────────────────────────────────────┘ │
│                                                       │
│  Tab: switch fields  Enter: save  Esc: cancel         │
└───────────────────────────────────────────────────────┘
```

---

## Keyboard contract

| Key | Active field | Behaviour |
|---|---|---|
| Tab | any | Advance: Title → URL → Tags → Title (cycle) |
| Shift+Tab | any | Reverse: Title → Tags → URL → Title (cycle) |
| Enter | any | Validate title (non-empty); if valid, emit save message with title, URL, tags |
| Esc | any | Cancel — discard changes, return to browse view |
| Any printable | active field | Insert character into the active textinput |

---

## Validation contract

| Condition | Behaviour |
|---|---|
| Title empty / whitespace-only on save | Show inline error below title field; do not persist |
| Title > 500 characters | Input blocked at field level (CharLimit) |
| URL or Tags unchanged | Patch omits those fields (nil pointer in BookmarkPatch) |

---

## State fields added to EditModel

```go
// New fields alongside existing urlInput / tagsInput:
titleInput  textinput.Model  // CharLimit = 500, pre-populated with b.Title
titleErr    string           // inline validation message, cleared on each save attempt

// New constant:
editFieldTitle = 0  // Title is first in tab order
editFieldURL   = 1  // was 0
editFieldTags  = 2  // was 1
```
