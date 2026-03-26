# Feature Specification: Pin Bookmarks in Vicinae Extension

**Feature Branch**: `011-vicinae-pin-bookmark`
**Created**: 2026-03-26
**Status**: Draft
**Input**: User description: "In the vicinae extension, we need the ability to pin bookmarks in the list view."

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Toggle Pin on a Bookmark (Priority: P1)

A user opens the "List Bookmarks" command in the Vicinae launcher. They find a bookmark they want to pin for quick
access. They invoke the "Toggle Pin" action on that bookmark. The bookmark's pin state is toggled and a confirmation
message is shown. Pinned bookmarks display the 📌 indicator.

**Why this priority**: Pinning is a core organisational feature that lets users mark important bookmarks for easy
identification. It is the sole focus of this feature.

**Independent Test**: Open the Vicinae launcher, invoke "List Bookmarks", find an unpinned bookmark, invoke "Toggle
Pin", and verify the 📌 indicator appears. Invoke "Toggle Pin" again and verify the indicator disappears.

**Acceptance Scenarios**:

1. **Given** the list view is open, **When** the user invokes "Toggle Pin" on an unpinned bookmark, **Then** the
   bookmark's `is_permanent` field is set to true and the 📌 indicator appears on the item.
2. **Given** the list view is open, **When** the user invokes "Toggle Pin" on a pinned bookmark (📌 visible), **Then**
   the bookmark's `is_permanent` field is set to false and the 📌 indicator is removed.
3. **Given** the pin action succeeds, **When** the list refreshes, **Then** the updated pin state persists across
   re-opens of the list view.
4. **Given** the `cairn pin` command fails, **When** the action completes, **Then** a failure toast is shown with the
   error message and the list state is unchanged.

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The "List Bookmarks" view MUST expose a "Toggle Pin" action on each bookmark item.
- **FR-002**: The "Toggle Pin" action MUST call `cairn pin <id>` to toggle the bookmark's `is_permanent` state.
- **FR-003**: On success, the list MUST refresh to reflect the updated pin state.
- **FR-004**: On failure, a toast with the error message MUST be shown; the list state MUST remain unchanged.
- **FR-005**: The `cairn` CLI MUST gain a `pin` subcommand: `cairn pin <id>` that toggles `is_permanent`.
- **FR-006**: The store MUST gain a `TogglePin(id int64) error` method that flips `is_permanent` in SQLite.
- **FR-007**: The `bm.ts` helper module MUST export a `bmPin(id: number)` function that invokes `cairn pin <id>`.

### Key Entities

- **Bookmark** (existing): gains toggled `IsPermanent` field via the new `cairn pin` command.
- **`cairn pin <id>`** (new CLI subcommand): takes a numeric bookmark ID; exits 0 on success, 1 if not found, 3 on
  error.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can pin or unpin a bookmark from the Vicinae extension list view without opening the TUI.
- **SC-002**: The pin state change persists: a subsequent `cairn list --json` reflects the updated `IsPermanent` value.
- **SC-003**: All CLI errors (ID not found, DB error) are surfaced as toast messages in the extension UI.

---

## Assumptions

- The `is_permanent` column already exists in the SQLite schema (added in feature 002); no migration is needed.
- "Toggle Pin" means flip the current `is_permanent` value (true → false, false → true) with a single command.
- The search view (`bm-search.tsx`) is out of scope for this feature; only the list view gains the pin action.
- The `cairn pin` subcommand uses the numeric `id` field (integer primary key), consistent with `cairn delete`.
