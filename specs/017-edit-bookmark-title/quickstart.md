# Quickstart: Edit Bookmark Title (017)

Manual smoke-test scenarios for both surfaces. Run after implementation.

---

## Prerequisites

```bash
# Build and install the binary
go build -o cairn ./cmd/cairn
export PATH="$PWD:$PATH"

# Add a test bookmark with a fallback title to work with
cairn add https://pixabay.com/vectors/
# Expected: saved with title "pixabay.com" (Cloudflare blocks fetch)
```

---

## Surface 1: TUI

```bash
cairn          # open TUI
```

**Test 1 — happy path**: Edit title successfully

1. Navigate to the "pixabay.com" bookmark (arrow keys).
2. Press `e` (edit key).
3. The edit panel opens. **Title field should be focused and pre-populated with "pixabay.com"**.
4. Clear the field and type `Free Vector Graphics`.
5. Press Enter.
6. **Expected**: panel closes; list shows `Free Vector Graphics` for that bookmark.

**Test 2 — empty title rejected**

1. Open edit panel for any bookmark.
2. Clear the title field entirely.
3. Press Enter.
4. **Expected**: panel stays open; inline error "Title cannot be empty" appears below the title field.

**Test 3 — tab navigation includes title**

1. Open edit panel.
2. Press Tab repeatedly.
3. **Expected**: focus cycles Title → URL → Tags → Title.
4. Press Shift+Tab.
5. **Expected**: focus cycles in reverse.

**Test 4 — Esc cancels without saving**

1. Open edit panel, change the title to something different.
2. Press Esc.
3. **Expected**: panel closes; original title unchanged in list.

---

## Surface 2: Vicinae Extension

```bash
cd vicinae-extension
vici dev          # or however you run the extension locally
```

**Test 1 — happy path**: Edit title successfully

1. Open the extension list.
2. Find the "pixabay.com" bookmark, press the edit action.
3. **Expected**: form opens with Title field pre-populated with "pixabay.com".
4. Change title to `Free Vector Graphics`. Leave URL and tags unchanged.
5. Submit.
6. **Expected**: "Bookmark updated" toast; list reflects new title.

**Test 2 — empty title rejected**

1. Open edit form for any bookmark.
2. Clear the title field.
3. Submit.
4. **Expected**: form stays open; field shows "Title cannot be empty" validation error.

**Test 3 — no-op when nothing changed**

1. Open edit form. Make no changes.
2. Submit.
3. **Expected**: "No changes" toast; no CLI call made.

**Test 4 — title-only change passes correct flag**

1. Open edit form. Change only the title. Leave URL and tags untouched.
2. Submit.
3. **Expected**: "Bookmark updated" toast. Verify via `cairn list` that only the title changed.

---

## CLI verification (both surfaces)

After any edit, verify the change persisted:

```bash
cairn list | grep -i "vector"
# or
cairn search "vector"
```
