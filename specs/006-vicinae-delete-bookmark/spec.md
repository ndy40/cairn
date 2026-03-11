# Feature Specification: Delete Bookmarks from Vicinae Extension

**Feature Branch**: `006-vicinae-delete-bookmark`
**Created**: 2026-03-11
**Status**: Draft
**Input**: User description: "in the vicinae-extension, I want to be able to delete bookmarks from the list view."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Delete a Bookmark from the List (Priority: P1)

A user browsing their bookmarks in the Vicinae extension's list view wants to remove a bookmark they no longer need. They select a bookmark, trigger a delete action from the action panel, confirm the deletion, and the bookmark is removed from both the list and the underlying store.

**Why this priority**: This is the core feature — without it, users must switch to the CLI or TUI to delete bookmarks, breaking workflow continuity.

**Independent Test**: Can be fully tested by launching the List Bookmarks command, selecting any bookmark, triggering delete, confirming, and verifying the bookmark no longer appears in the list.

**Acceptance Scenarios**:

1. **Given** the user has bookmarks in the list view, **When** they select a bookmark and choose "Delete Bookmark" from the action panel and confirm, **Then** the bookmark is permanently deleted and the list refreshes without the deleted item.
2. **Given** the user has bookmarks in the list view, **When** they select a bookmark and choose "Delete Bookmark" but cancel the confirmation, **Then** the bookmark remains in the list and nothing is changed.
3. **Given** a deletion succeeds, **When** the list refreshes, **Then** the user sees a brief success notification confirming the deletion.

---

### User Story 2 - Delete a Bookmark from Search Results (Priority: P2)

A user searching for bookmarks in the Vicinae extension finds an outdated result and wants to delete it directly from the search results view without navigating back to the full list.

**Why this priority**: Extends the delete action to the search view for workflow consistency, but is secondary since the list view is the primary browsing surface.

**Independent Test**: Can be fully tested by launching the Search Bookmarks command, searching for a bookmark, triggering delete from the action panel, confirming, and verifying the result disappears.

**Acceptance Scenarios**:

1. **Given** the user has search results showing, **When** they select a result and choose "Delete Bookmark" and confirm, **Then** the bookmark is deleted and the search results refresh without the deleted item.

---

### Edge Cases

- What happens when the user tries to delete a bookmark that was already deleted (e.g., deleted from CLI while extension is open)? The extension should show an appropriate error message and refresh the list.
- What happens when the `cairn` CLI is not available? The delete action should not be shown if the CLI is unavailable (consistent with existing behavior where the entire list shows an error).
- What happens when the deletion fails due to a database or CLI error? The user should see an error notification with a human-readable message.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The extension MUST display a "Delete Bookmark" action in the action panel for each bookmark item in the list view.
- **FR-002**: The extension MUST display a "Delete Bookmark" action in the action panel for each bookmark item in the search results view.
- **FR-003**: The extension MUST show a confirmation prompt before deleting a bookmark to prevent accidental deletions.
- **FR-004**: The extension MUST call the `cairn delete <id>` CLI command to perform the deletion.
- **FR-005**: The extension MUST refresh the bookmark list after a successful deletion so the deleted item is no longer visible.
- **FR-006**: The extension MUST display a success notification after a bookmark is deleted.
- **FR-007**: The extension MUST display an error notification if the deletion fails, including the reason (e.g., "not found", CLI error).
- **FR-008**: The "Delete Bookmark" action MUST be a secondary action (not the primary Enter action), preserving "Open in Browser" as the default action.

### Key Entities

- **Bookmark**: Existing entity with an `ID` field used to identify the bookmark for deletion via the CLI.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can delete a bookmark from the list view in under 5 seconds (select, action, confirm, done).
- **SC-002**: Deleted bookmarks disappear from the list immediately after confirmation without requiring a manual refresh.
- **SC-003**: Accidental deletions are prevented — 100% of delete operations require explicit user confirmation.
- **SC-004**: Failed deletions show a clear error message — users are never left wondering if the action succeeded or failed.

## Assumptions

- The `cairn delete <id>` CLI command is the only interface for deletion — the extension does not access the database directly.
- The Vicinae API provides a confirmation dialog or alert mechanism (consistent with the existing `showHUD` and `Clipboard` API usage).
- The bookmark `ID` field is available in the JSON output from `cairn list --json` and `cairn search --json`.

## Dependencies

- Existing `cairn delete <id>` CLI command (already implemented).
- Existing `bmList()` and `bmSearch()` wrapper functions in the extension.
