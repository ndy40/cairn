# Contract: Keyboard Shortcuts (Updated for Feature 003)

**Feature**: 003-edit-bookmark-help
**Date**: 2026-03-06
**Extends**: `specs/002-tags-pinning-archive/contracts/keyboard-shortcuts.md`

---

## Browse Mode (additions)

| Key | Action | Notes |
|-----|--------|-------|
| `e` | Open edit panel for selected bookmark | Enters `StateEdit`; pre-fills tags field |

---

## Edit Panel (new state: StateEdit)

| Key | Action |
|-----|--------|
| Any character | Types into the tags input field |
| `Enter` | Save updated tags and return to browse |
| `Esc` | Cancel; return to browse with no changes |

### Behaviour

- The panel shows the bookmark title as a read-only heading (not an input field).
- The tags field is pre-filled with the bookmark's current tags, comma-separated (e.g., `work, go, tools`).
- On save (`Enter`): tags are split on commas, normalised (lowercase, dedup, truncate, max 3), saved, and the browse list reloads.
- If more than 3 tags are entered, only the first 3 are saved and a status note "Only first 3 tags saved" is displayed briefly before returning to browse.
- Tab cycling between fields is not needed (only one editable field).
- Pressing `Esc` discards all edits; the original tags remain unchanged.
- `e` on an empty list (no bookmarks) does nothing.
- `e` is only active in browse mode; it has no effect in search mode, archive view, or tag filter overlay.

---

## Updated Help Screen

The help screen (`?` key) is updated to include the `e` key:

```
Browse Mode
  /          Enter search
  Ctrl+P     Add bookmark from clipboard
  Enter      Open selected in browser
  e          Edit tags on selected bookmark
  d, Delete  Delete selected bookmark
  p          Toggle permanent (pin) flag
  t          Open tag filter
  a          Open archive view
  j / ↓      Move down
  k / ↑      Move up
  g          Jump to top
  G          Jump to bottom
  ...
```
