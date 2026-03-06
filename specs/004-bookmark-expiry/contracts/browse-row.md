# Contract: Browse Row Display (Updated for Feature 004)

**Feature**: 004-bookmark-expiry
**Date**: 2026-03-06
**Extends**: `specs/003-edit-bookmark-help/contracts/keyboard-shortcuts.md`

---

## Changes

The bookmark description line in the browse list no longer shows last-visited information.

---

## Browse Row Format

Each active bookmark in the list renders two lines:

### Title line
```
[pin] <title or URL>
```
- `[pin]` prefix is shown only when `is_permanent = true`.
- If title is empty, the URL is shown as the title.

### Description line

**Previous format** (features 001–003):
```
<domain> · <created_date> · Last: <visited_date>     (when visited)
<domain> · <created_date> · Never visited             (when never visited)
<domain> · <created_date> · Last: <visited_date> · #tag1 #tag2
```

**New format** (feature 004 onwards):
```
<domain> · <created_date>
<domain> · <created_date> · #tag1 #tag2
```

Where:
- `<domain>` — stripped hostname (no `www.`)
- `<created_date>` — `YYYY-MM-DD` format
- Tags — space-separated, each prefixed with `#`, omitted if no tags

---

## Expiry Notification (Startup)

Unchanged from feature 002. When bookmarks are archived at startup, the footer shows:
```
N bookmark(s) archived
```
Where N ≥ 1. No message is shown when 0 bookmarks are archived.

The threshold that triggers this message changes from 183 days to **30 days** from creation date.

---

## Keyboard Shortcuts

All keyboard shortcuts remain unchanged from feature 003. No new bindings are added or removed.
