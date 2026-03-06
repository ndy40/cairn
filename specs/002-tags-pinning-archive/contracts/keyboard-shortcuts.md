# Contract: Keyboard Shortcuts (Updated for Feature 002)

**Feature**: 002-tags-pinning-archive
**Date**: 2026-03-06
**Extends**: `specs/001-tui-bookmark-manager/contracts/keyboard-shortcuts.md`

---

## Summary of Changes

This document lists only the **new and modified** shortcuts introduced by feature 002. All shortcuts from the feature 001 contract remain unchanged.

---

## Browse Mode (additions)

| Key | Action | Notes |
|-----|--------|-------|
| `p` | Toggle permanent flag | Flips `is_permanent` on the selected bookmark; reloads list |
| `t` | Open tag filter overlay | Enters `StateTagFilter`; shows all unique tags in active list |
| `a` | Open archive view | Enters `StateArchive`; loads archived bookmarks |

### Permanent Indicator

When `is_permanent = true`, the bookmark title is prefixed with `[pin] ` in the list view. The prefix is rendered by `BookmarkItem.Title()`.

### Tag Display

Tags are shown in `BookmarkItem.Description()` after the domain and date, separated by `·`:
```
example.com · 2025-01-15 · #work #go #tools
```
Tags are prefixed with `#` for readability. If no tags are assigned, the description renders as before (domain and date only).

---

## Tag Filter Overlay (new state: StateTagFilter)

| Key | Action |
|-----|--------|
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `Space` / `Enter` | Toggle the highlighted tag on/off |
| `Esc` / `t` | Close overlay; return to browse with active filter applied |
| `c` | Clear all selected tags; return to browse |

### Behaviour

- The overlay lists all unique tags present in the current active bookmark list.
- Selected tags are marked with `[x]`; unselected with `[ ]`.
- The browse list updates in real time as tags are toggled (no need to close the overlay to see results).
- If the text search is also active, both filters apply simultaneously (tag filter AND text search).
- When no tags are selected, the full active list is shown.

---

## Archive View (new state: StateArchive)

| Key | Action |
|-----|--------|
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `r` | Restore selected bookmark to active list |
| `Esc` | Return to browse mode |
| `Ctrl+C` | Quit application |

### Behaviour

- Archived bookmarks are displayed in reverse `archived_at` order (most recently archived first).
- Each item shows: title, URL domain, creation date, last-visited date, archived date.
- Restoring a bookmark clears `is_archived` and `archived_at`; `last_visited_at` is preserved.
- After restore, the archive list reloads automatically.
- Permanent bookmarks never appear in this view (they are never archived).

---

## Add Modal (modified)

| Key | Action |
|-----|--------|
| `Tab` | Move to next field (URL → Tags → back to URL) |
| `Shift+Tab` | Move to previous field |
| All other keys | Unchanged from feature 001 contract |

### Tag Input Field

- A single text input field labelled "Tags (comma-separated, max 3)".
- User may type tags separated by commas: `work, go, tools`.
- On save (`Enter`): tags are split on commas, trimmed of whitespace, lowercased, deduplicated, truncated to 32 chars each, and only the first 3 are kept.
- Empty or whitespace-only tags are silently discarded.
- If the user enters more than 3 valid tags, the first 3 are saved; a note "Only first 3 tags saved" is shown in the status line.
- Tag field is optional; leaving it empty saves the bookmark with no tags.

---

## Updated Help Screen

The help screen (`?` key) is updated to include the new shortcuts:

```
Browse Mode
  /          Enter search
  Ctrl+P     Add bookmark from clipboard
  Enter      Open selected in browser
  d, Delete  Delete selected bookmark
  p          Toggle permanent (pin) flag
  t          Open tag filter
  a          Open archive view
  j / ↓      Move down
  k / ↑      Move up
  g          Jump to top
  G          Jump to bottom

Tag Filter Overlay
  j / ↓      Move down
  k / ↑      Move up
  Space/Enter Toggle selected tag
  c          Clear all tag filters
  Esc / t    Close overlay

Archive View
  r          Restore selected bookmark
  Esc        Return to browse

Add Mode
  Tab        Next field
  Enter      Save bookmark
  Esc        Cancel
```
