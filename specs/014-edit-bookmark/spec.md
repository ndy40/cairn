# Feature Specification: Edit Bookmark

**Feature Branch**: `014-edit-bookmark`
**Created**: 2026-04-01
**Status**: Draft
**Input**: User description: "We need to have the ability to edit bookmarks. We have a runEdit function in cmd/cairn/main.go. Verify that this meets the requirement for editing a bookmark. We also need to make this available in the CLI and TUI. We should be able to edit the bookmark url and tags. Also make sure we are able to do this from the vicinae-extension as well."

## User Scenarios & Testing

### User Story 1 - Edit Bookmark URL via CLI (Priority: P1)

A user realises a bookmarked URL has moved to a new address (e.g. a domain change or a typo on save). They want to correct the URL without deleting and re-adding the bookmark, preserving its tags, title, and creation date.

**Why this priority**: Correcting a bookmark's URL is the core missing capability. Tags and title are already editable via the CLI; URL editing is the gap that prompted this feature.

**Independent Test**: Can be fully tested by running `cairn edit <id> --url=<new-url>` and verifying the bookmark's URL and domain update while all other fields remain unchanged.

**Acceptance Scenarios**:

1. **Given** a bookmark exists with ID 5, **When** the user runs `cairn edit 5 --url=https://new-example.com`, **Then** the bookmark's URL is updated to `https://new-example.com`, its domain is recalculated, and the `updated_at` timestamp is refreshed.
2. **Given** a bookmark exists with ID 5, **When** the user runs `cairn edit 5 --url=https://new-example.com --tags=go,dev`, **Then** both the URL and tags are updated in a single operation.
3. **Given** a bookmark exists with ID 5 and another bookmark already has the URL `https://taken.com`, **When** the user runs `cairn edit 5 --url=https://taken.com`, **Then** the system rejects the edit with a duplicate URL error.
4. **Given** a bookmark exists with ID 5, **When** the user runs `cairn edit 5 --url=""`, **Then** the system rejects the edit with a validation error (URL cannot be empty).

---

### User Story 2 - Edit Bookmark URL and Tags via TUI (Priority: P2)

A user browsing their bookmarks in the TUI selects one and wants to update its URL, title, or tags without leaving the terminal interface. Currently the TUI edit panel only allows editing tags.

**Why this priority**: The TUI is the primary interactive interface. Extending the existing edit panel to support URL editing delivers a seamless experience for power users.

**Independent Test**: Can be fully tested by launching the TUI, pressing `e` on a bookmark, modifying the URL field, pressing Enter, and verifying the bookmark is updated in the database.

**Acceptance Scenarios**:

1. **Given** the user is viewing bookmarks in the TUI, **When** they press `e` on a selected bookmark, **Then** an edit panel opens showing editable fields for URL and tags (pre-filled with current values).
2. **Given** the edit panel is open with a bookmark's URL, **When** the user changes the URL and presses Enter, **Then** the bookmark's URL and domain are updated and the list view reflects the change.
3. **Given** the edit panel is open, **When** the user presses Escape, **Then** no changes are saved and the user returns to the list view.
4. **Given** the edit panel is open and the user changes the URL to one that already exists, **When** they press Enter, **Then** an error message is displayed and the edit is not saved.

---

### User Story 3 - Edit Bookmark from Vicinae Extension (Priority: P3)

A user working in their browser via the Vicinae extension wants to correct a bookmark's URL or update its tags without switching to the terminal.

**Why this priority**: The extension is a convenience layer. Adding edit support completes the CRUD feature set in the extension but is lower priority since users can fall back to the CLI.

**Independent Test**: Can be fully tested by opening the extension, selecting a bookmark, editing its URL or tags, submitting, and verifying the changes persist.

**Acceptance Scenarios**:

1. **Given** the user has the Vicinae extension open and a bookmark is displayed, **When** they trigger the edit action on that bookmark, **Then** an edit form is presented with the bookmark's current URL and tags pre-filled.
2. **Given** the edit form is open, **When** the user modifies the URL and/or tags and submits, **Then** the extension invokes the cairn CLI to persist the changes and the bookmark list refreshes with updated values.
3. **Given** the edit form is open, **When** the user cancels the edit, **Then** no changes are made and the view returns to the bookmark list.

---

### Edge Cases

- What happens when the user provides an invalid URL format (e.g. missing scheme)? The system should accept it as-is (same behaviour as `cairn add`).
- What happens when the user edits a bookmark that was deleted by another sync device? The system should return a "not found" error.
- What happens when the user edits only the URL to the same value it already has? The system should accept it (no-op on URL, still update `updated_at`).
- What happens when the user provides a URL that matches the same bookmark's own existing URL? The system should allow it (not a duplicate conflict with itself).

## Requirements

### Functional Requirements

- **FR-001**: System MUST allow editing a bookmark's URL via the CLI (`cairn edit <id> --url=<new-url>`).
- **FR-002**: System MUST recalculate and update the domain field when a bookmark's URL changes.
- **FR-003**: System MUST reject URL edits that would create a duplicate URL (another bookmark already has that URL), unless the URL belongs to the bookmark being edited.
- **FR-004**: System MUST reject empty URL values during edit.
- **FR-005**: System MUST support editing URL, title, and tags independently or in any combination in a single CLI command.
- **FR-006**: The TUI edit panel MUST allow editing the bookmark URL in addition to the existing tags field.
- **FR-007**: The Vicinae extension MUST provide an edit action for bookmarks that allows modifying URL and tags.
- **FR-008**: The Vicinae extension MUST invoke the cairn CLI to perform edit operations (consistent with existing extension patterns).
- **FR-009**: All edit operations MUST record a pending sync change so edits propagate to synced devices.
- **FR-010**: All edit operations MUST update the `updated_at` timestamp.

### Key Entities

- **Bookmark**: Existing entity. The `URL` and `Domain` fields become editable (previously write-once at insert time). All other fields remain unchanged.
- **BookmarkPatch**: Existing data transfer object. Gains a new optional `URL` field to support URL updates.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Users can update a bookmark's URL from the CLI and see the change reflected immediately in subsequent list/search commands.
- **SC-002**: Users can update a bookmark's URL and tags from the TUI edit panel in a single interaction.
- **SC-003**: Users can update a bookmark's URL and tags from the Vicinae extension without switching to the terminal.
- **SC-004**: Editing a bookmark to a duplicate URL is rejected with a clear error message across all interfaces (CLI, TUI, extension).
- **SC-005**: All bookmark edits (URL, title, tags) trigger sync propagation so changes appear on other synced devices.
