# Feature Specification: Edit Bookmark Title

**Feature Branch**: `017-edit-bookmark-title`
**Created**: 2026-04-15
**Status**: Draft
**Input**: User description: "Implement the ability to edit the title of a bookmark. Sometimes behind cloudflare or sites like youtube block bots and so we are unable to fetch accurate title. Give the user the ability to update the title. Update the vicinae-extension as well and allow editing the title."

## Background

When a bookmark is saved, the title is auto-fetched from the web page. Some sites (Cloudflare-protected pages, YouTube under certain conditions) block automated requests and return a challenge page, resulting in titles like "Just a moment..." or the hostname fallback being stored. Users currently have no way to correct these titles after the fact through the interactive surfaces.

The `cairn edit` CLI command already accepts `--title` to correct titles. This feature brings that capability to the two interactive surfaces that don't yet expose it: the TUI and the Vicinae extension.

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Correct a Title in the TUI (Priority: P1)

A user saved a Pixabay or YouTube bookmark whose title was auto-set to a fallback value. They open the cairn TUI, navigate to that bookmark, open the edit panel, and are able to type a correct title before saving.

**Why this priority**: The TUI is the primary interactive interface for cairn. This is the most direct path for users who live in the terminal.

**Independent Test**: Open the TUI, select any bookmark, press the edit key, edit the title field, save, and confirm the updated title appears in the bookmark list.

**Acceptance Scenarios**:

1. **Given** a bookmark with an incorrect title (e.g. "youtube.com"), **When** the user opens the edit panel and changes the title field to "Never Gonna Give You Up", **Then** the bookmark is saved with the new title and it appears immediately in the list view.
2. **Given** the edit panel is open, **When** the user clears the title field and saves, **Then** the save is rejected and an inline error message is shown ("Title cannot be empty").
3. **Given** the edit panel is open, **When** the user navigates between fields using Tab/Shift+Tab, **Then** the title field is reachable in the tab order alongside URL and Tags.
4. **Given** the edit panel is open, **When** the user presses Esc without saving, **Then** the original title is unchanged.

---

### User Story 2 — Correct a Title in the Vicinae Extension (Priority: P2)

A user saved a bookmark via the Vicinae extension whose title was auto-set to a fallback. They open the extension's edit form for that bookmark and see a pre-populated title field they can correct before submitting.

**Why this priority**: The Vicinae extension is the secondary editing surface. Users who bookmark via the extension should have full parity with the CLI and TUI for post-save corrections.

**Independent Test**: Open the Vicinae extension, find a bookmark, open the edit form, change the title, submit, and confirm the updated title is reflected when listing bookmarks.

**Acceptance Scenarios**:

1. **Given** a bookmark with title "youtube.com", **When** the user opens its edit form and types "My Favourite Video" in the title field, **Then** submitting saves the new title.
2. **Given** the edit form is open, **When** the user submits with the title field empty or whitespace-only, **Then** the form shows a validation error and does not call the CLI.
3. **Given** a bookmark is loaded into the edit form, **When** the title is unchanged and other fields are also unchanged, **Then** submitting shows "No changes" without calling the CLI.
4. **Given** a bookmark is loaded into the edit form, **When** only the title has changed, **Then** submitting calls the CLI with only the `--title` flag (URL and tags are left untouched).

---

### Edge Cases

- What if the title the user types is only whitespace? It should be treated as empty and rejected with a clear error.
- What if the title is extremely long (e.g. 1000+ characters)? The field should enforce a maximum of 500 characters.
- What if the bookmark no longer exists by the time the user submits the edit? The system should show an error and not crash.
- What if the title is unchanged but other fields are changed? Save should proceed using only the changed fields.

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The TUI edit panel MUST expose a text input field for the bookmark title, pre-populated with the current title.
- **FR-002**: The TUI edit panel MUST include the title field in the Tab/Shift+Tab navigation order alongside URL and Tags.
- **FR-003**: The TUI MUST reject a save attempt when the title field is empty or whitespace-only, showing an inline error message.
- **FR-004**: The TUI MUST persist the updated title when the user saves and refresh the list view to reflect the change immediately.
- **FR-005**: The Vicinae extension edit form MUST include a title field, pre-populated with the bookmark's current title.
- **FR-006**: The Vicinae extension MUST pass the updated title to the CLI only when the title value has changed from its original.
- **FR-007**: The Vicinae extension MUST reject submission when the title field is empty or whitespace-only, showing a field-level validation error.
- **FR-008**: Both surfaces MUST enforce a maximum title length of 500 characters.

### Key Entities

- **Bookmark**: Existing entity. The `Title` field is already stored and writable via `cairn edit --title`. No schema changes are required.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can correct an incorrect bookmark title in the TUI in under 30 seconds from opening the edit panel.
- **SC-002**: A user can correct an incorrect bookmark title in the Vicinae extension in under 30 seconds from opening the edit form.
- **SC-003**: 100% of title edits submitted through either surface are reflected immediately in the bookmark list without requiring a restart or manual refresh.
- **SC-004**: Empty or whitespace-only titles are rejected 100% of the time on both surfaces, with a clear inline error message shown to the user.

---

## Assumptions

- The CLI `cairn edit <id> --title "<value>"` command already works correctly and requires no changes.
- The backing store's `UpdateFields` method already handles title updates — no schema migration needed.
- The Vicinae extension delegates to the installed `cairn` binary; no new CLI flags or API changes are required beyond passing `--title`.
- Maximum title length of 500 characters matches the existing URL field limit as a reasonable default; no existing bookmarks exceed this.
