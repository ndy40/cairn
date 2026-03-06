# Feature Specification: Vicinae Extension for Bookmark Manager

**Feature Branch**: `005-vicinae-extension`
**Created**: 2026-03-06
**Status**: Draft
**Input**: User description: "Build a Vicinae extension for the bookmarking app. I want to be able to add bookmarks, search bookmarks and view bookmarks. It should work by calling the bm cli."

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Search Bookmarks (Priority: P1)

A user triggers the Vicinae launcher and types a search query. The extension instantly shows matching bookmarks filtered by the query — displaying each bookmark's title, domain, and tags. The user selects a result and the bookmark opens in their default browser.

**Why this priority**: Search is the primary reason to access bookmarks through a launcher. It replaces hunting through a terminal; the user should be able to find and open any saved URL in two to three keystrokes from the launcher.

**Independent Test**: Open the Vicinae launcher, invoke the search command, type a keyword that matches at least one saved bookmark, and verify matching results appear. Press Enter on a result and confirm the URL opens in the browser.

**Acceptance Scenarios**:

1. **Given** the Vicinae launcher is open and bookmarks exist, **When** the user invokes the "Search Bookmarks" command, **Then** a search input is shown and all bookmarks are listed below it.
2. **Given** the search view is active, **When** the user types a query, **Then** the result list filters in real time to show only matching bookmarks.
3. **Given** matching results are displayed, **When** the user selects a bookmark and confirms, **Then** the bookmark URL opens in the default browser and the launcher closes.
4. **Given** the user types a query that matches no bookmarks, **When** the result list updates, **Then** an empty-state message ("No bookmarks found") is shown.
5. **Given** the `bm` CLI is not installed or not accessible, **When** the extension loads, **Then** an error message is displayed explaining that the `bm` tool must be installed.

---

### User Story 2 - View All Bookmarks (Priority: P2)

A user opens the "List Bookmarks" command from the Vicinae launcher without typing a search query. All saved active bookmarks are displayed in a scrollable list ordered by newest first. The user can browse and open any bookmark directly.

**Why this priority**: Browse mode lets users discover bookmarks they saved without remembering a specific keyword. It is less frequently used than search but essential for reviewing saved content.

**Independent Test**: Open the Vicinae launcher, invoke the "List Bookmarks" command, and verify the full list of active bookmarks is displayed ordered by date saved. Select one and verify it opens in the browser.

**Acceptance Scenarios**:

1. **Given** the Vicinae launcher is open, **When** the user invokes the "List Bookmarks" command, **Then** all active (non-archived) bookmarks are displayed newest first.
2. **Given** the bookmark list is displayed, **When** the user selects a bookmark and confirms, **Then** the URL opens in the default browser.
3. **Given** no bookmarks have been saved, **When** the list view loads, **Then** an empty-state message ("No bookmarks saved yet") is shown.
4. **Given** the bookmark list is displayed, **When** the user types in the launcher search box, **Then** the list filters to show only matching results (inline filtering within the list view).

---

### User Story 3 - Add a Bookmark (Priority: P3)

A user copies a URL to the clipboard, opens the Vicinae launcher, and invokes the "Add Bookmark" command. A URL input form appears. The user can paste the copied URL, optionally add tags, and submit. The extension saves the bookmark using the `bm` CLI and shows a success or error message.

**Why this priority**: Adding a bookmark from the launcher is faster than switching to a terminal. It serves as a quick-capture flow. It is lower priority than search and view because bookmarks can still be added via the CLI directly.

**Independent Test**: Open the Vicinae launcher, invoke "Add Bookmark", paste a URL into the form, submit, then invoke "List Bookmarks" and verify the new bookmark appears.

**Acceptance Scenarios**:

1. **Given** the "Add Bookmark" command is invoked, **When** the form loads, **Then** a URL input field is shown (pre-filled with the current clipboard content if it is a valid URL).
2. **Given** the form is shown, **When** the user submits a valid URL, **Then** the extension saves the bookmark and shows a "Saved" confirmation message.
3. **Given** the form is shown, **When** the user submits a URL that is already bookmarked, **Then** an error message "Already bookmarked" is shown.
4. **Given** the form is shown, **When** the user submits an empty URL or no URL, **Then** validation prevents submission and shows "URL is required".
5. **Given** the form is shown, **When** the user optionally fills in a comma-separated tags field and submits, **Then** the bookmark is saved with those tags.
6. **Given** a network or CLI error occurs during save, **When** submission is attempted, **Then** an error message with the failure reason is shown.

---

### Edge Cases

- The `bm` CLI is not installed: all commands show a clear error with installation instructions rather than crashing silently.
- The bookmark database is missing or inaccessible: the error from `bm` is surfaced to the user in the extension UI.
- A URL pasted into the Add form is not a valid URL format: the form validates and rejects it before calling `bm`.
- The user has a large number of bookmarks (hundreds): the list view loads and scrolls without performance issues.
- The clipboard does not contain a URL when opening the Add form: the URL field is empty (no pre-fill) and the user types manually.
- Search with special characters or spaces: the query is passed safely to `bm search` without causing a shell error.

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The extension MUST provide three commands accessible from the Vicinae launcher: "Search Bookmarks", "List Bookmarks", and "Add Bookmark".
- **FR-002**: The "Search Bookmarks" command MUST accept a text query and display matching bookmarks from the `bm` CLI in real time as the user types.
- **FR-003**: The "List Bookmarks" command MUST display all active bookmarks ordered by date saved (newest first), with inline filtering as the user types in the launcher search box.
- **FR-004**: Each bookmark in any list or search result MUST display at minimum: title (or URL if no title), domain, and tags (if any).
- **FR-005**: Selecting a bookmark in any list view and confirming MUST open the bookmark URL in the default browser.
- **FR-006**: The "Add Bookmark" command MUST present a URL input field pre-filled with the clipboard content if it is a valid URL.
- **FR-007**: The "Add Bookmark" command MUST include an optional comma-separated tags input field.
- **FR-008**: On successful save, the "Add Bookmark" command MUST display a success confirmation.
- **FR-009**: On duplicate URL submission, the "Add Bookmark" command MUST display "Already bookmarked".
- **FR-010**: The extension MUST display a clear error message if the `bm` CLI is not found in the system path, with a hint to install it.
- **FR-011**: All communication with the bookmark store MUST go through the `bm` CLI (no direct database access from the extension).

### Key Entities

- **Bookmark** (read from CLI output): title, URL, domain, tags, date saved. Displayed in list and search results; opened in browser on selection.
- **Search Query**: free-text string entered by the user; passed to `bm search <query>`; results rendered as a list.
- **Add Form**: URL (required) and tags (optional, comma-separated); submitted via `bm add <url>`.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can find and open a saved bookmark in under 5 seconds from triggering the Vicinae launcher.
- **SC-002**: Search results appear within 1 second of the user finishing typing a query.
- **SC-003**: A user can save a new bookmark from the launcher in under 15 seconds (trigger launcher → invoke add → paste URL → submit).
- **SC-004**: All three commands (search, list, add) are accessible from the Vicinae launcher without any configuration beyond having `bm` installed.
- **SC-005**: 100% of CLI error messages (duplicate, not found, CLI missing) are surfaced in the extension UI rather than silently failing.

---

## Assumptions

- The `bm` CLI is installed and available in the system PATH on the machine running Vicinae. The extension does not bundle or install `bm` itself.
- The extension is a new standalone package (separate from the Go CLI codebase) that lives in its own directory within this repository (e.g., `vicinae-extension/`) or in a sibling repository.
- "View bookmarks" means the active (non-archived) list; archived bookmarks are not surfaced in the extension.
- The extension does not support delete, pin, archive, or tag-filter operations; those remain TUI-only for this feature.
- Pre-filling the Add form with clipboard content is a best-effort convenience; if clipboard access is unavailable, the field is simply empty.
- The `bm list --json` and `bm search <query> --json` output formats are used to parse bookmark data in the extension.
- Tags field in the Add form maps to the tags argument in `bm add`; tag normalisation (max 3, lowercase) is handled by the `bm` CLI.
