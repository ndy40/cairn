# Feature Specification: TUI Bookmark Manager

**Feature Branch**: `001-tui-bookmark-manager`
**Created**: 2026-03-06
**Status**: Draft
**Input**: User description: "I want to build a simple but useful cli based TUI application that I can use to hold bookmarks to important links. I want to be able to add new links and bookmark them. Once app is running, I want to be able to use shortcut keys like Ctrl+P to paste links. The program then fetches the page and scanning meta tags, save the title of page + link and the date it was created. As the number of links increases, I want to be able to run quick fuzzy searches. I need to be able to search by title, content or domain."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Add a Bookmark via Keyboard Shortcut (Priority: P1)

A user is running the TUI application and wants to save a URL they have copied to their clipboard. They press Ctrl+P to trigger the bookmark entry flow. The application pastes the URL from their clipboard into an input field, fetches the page in the background, extracts its title and any relevant metadata from the page's meta tags, then saves the bookmark with the title, URL, and timestamp automatically.

**Why this priority**: This is the core interaction loop. Without the ability to add bookmarks, the application delivers no value. All other features depend on bookmarks existing.

**Independent Test**: Can be fully tested by pressing Ctrl+P with a URL in the clipboard, observing the auto-populated title, and verifying the bookmark appears in the list with correct title, URL, and creation date.

**Acceptance Scenarios**:

1. **Given** the app is open and a valid URL is in the clipboard, **When** the user presses Ctrl+P, **Then** the URL is pasted into an input field and the app fetches the page title and metadata automatically.
2. **Given** the page has been fetched successfully, **When** the bookmark is confirmed, **Then** it is saved with the page title, original URL, and current date/time.
3. **Given** the page fetch fails (e.g., network error or invalid URL), **When** the user confirms, **Then** the bookmark is saved with the URL and date, with the title marked as unavailable, and the user sees a clear error message.
4. **Given** a URL is already bookmarked, **When** the user tries to add the same URL again, **Then** the app notifies the user that the bookmark already exists and does not create a duplicate.

---

### User Story 2 - Browse Saved Bookmarks (Priority: P2)

A user opens the TUI application and wants to browse through all their saved bookmarks in a readable list. They can scroll through entries, each showing the page title, URL, and date added. They can select a bookmark and open it in their default browser.

**Why this priority**: Viewing the bookmark list is foundational — users need to see what they have saved and act on it. This creates a usable MVP when combined with P1.

**Independent Test**: Can be fully tested with a pre-populated set of bookmarks, verifying they display correctly and that selecting one opens the correct URL.

**Acceptance Scenarios**:

1. **Given** the app has saved bookmarks, **When** the user opens the app, **Then** all bookmarks are displayed in a scrollable list showing title, URL (or domain), and date added.
2. **Given** the bookmark list is displayed, **When** the user selects a bookmark and confirms, **Then** the URL opens in the system's default browser.
3. **Given** the app has no bookmarks saved, **When** the user opens the app, **Then** a helpful empty state message is shown guiding them to add their first bookmark.

---

### User Story 3 - Fuzzy Search Bookmarks (Priority: P3)

A user has accumulated many bookmarks and needs to find a specific one quickly. They start typing in a search bar and the list narrows down in real time, matching against bookmark titles, page content descriptions (meta description), and domain names. The search is tolerant of typos and partial matches.

**Why this priority**: As the collection grows, manual scrolling becomes impractical. Fuzzy search is essential for long-term usability and is the primary differentiator over simple list management.

**Independent Test**: Can be fully tested with a collection of 20+ bookmarks by searching partial terms and verifying relevant results appear instantly, including matches on title, domain, and description fields.

**Acceptance Scenarios**:

1. **Given** bookmarks are saved, **When** the user types a partial word in the search bar, **Then** the bookmark list filters in real time to show only matching results.
2. **Given** a search query is entered, **When** the search runs, **Then** matches are found across title, meta description/content, and domain name independently.
3. **Given** the user types a slightly misspelled term, **When** the search runs, **Then** fuzzy matching surfaces near-matching bookmarks rather than returning zero results.
4. **Given** a search returns no matches, **When** the list is empty, **Then** the user sees a "no results found" message with their search term shown.
5. **Given** a search is active, **When** the user clears the search field, **Then** the full bookmark list is restored immediately.

---

### User Story 4 - Delete a Bookmark (Priority: P4)

A user wants to remove a bookmark that is no longer relevant. They select it from the list and delete it using a keyboard shortcut. The app asks for confirmation before permanently removing the entry.

**Why this priority**: Data hygiene is important for long-term usability. Users need to be able to remove outdated links to keep their collection useful.

**Independent Test**: Can be fully tested by deleting an existing bookmark and verifying it no longer appears in the list or search results.

**Acceptance Scenarios**:

1. **Given** a bookmark is selected, **When** the user presses the delete key, **Then** a confirmation prompt appears before deletion.
2. **Given** the confirmation prompt is shown, **When** the user confirms, **Then** the bookmark is permanently removed and the list updates.
3. **Given** the confirmation prompt is shown, **When** the user cancels, **Then** the bookmark is not deleted and the list remains unchanged.

---

### Edge Cases

- What happens when the clipboard is empty when Ctrl+P is pressed? The app shows a message indicating no URL was found in the clipboard.
- What happens when a URL points to a page that requires authentication or returns a non-200 response? The bookmark is saved with the URL and date; the title is left blank or marked as "Untitled".
- What happens when the page has no `<title>` or meta tags? The bookmark is saved using the URL as the display title.
- What happens when the app cannot reach the internet at all? The bookmark is saved with the URL and date; a warning is shown that metadata could not be fetched.
- What happens when the bookmark list is very large (hundreds of entries)? The list must remain scrollable and search must remain responsive.
- What happens when the user types a search query of a single character? Fuzzy search still runs and surfaces matching results.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The application MUST launch as an interactive TUI from the command line without additional arguments.
- **FR-002**: The application MUST respond to the Ctrl+P keyboard shortcut to trigger the add-bookmark flow using the current clipboard content.
- **FR-003**: The application MUST attempt to fetch the web page at the provided URL and extract the page title and meta description upon bookmark creation.
- **FR-004**: The application MUST save each bookmark with at minimum: the URL, the extracted page title (or fallback), the domain, the meta description (if available), and the date/time the bookmark was created.
- **FR-005**: The application MUST prevent saving duplicate bookmarks for the same URL and notify the user when a duplicate is attempted.
- **FR-006**: The application MUST display all saved bookmarks in a scrollable list, each showing the title, domain, and date added.
- **FR-007**: Users MUST be able to open a selected bookmark in the system default browser using a keyboard shortcut.
- **FR-008**: Users MUST be able to delete a selected bookmark via a keyboard shortcut, with a confirmation step before permanent removal.
- **FR-009**: The application MUST provide an always-visible search input field that performs fuzzy search in real time as the user types.
- **FR-010**: Fuzzy search MUST match against bookmark title, meta description/content, and domain name as independent searchable fields.
- **FR-011**: The application MUST store bookmarks persistently on the local filesystem so they survive application restarts.
- **FR-012**: The application MUST display clear error or status messages when a page fetch fails, a URL is invalid, or the clipboard is empty.
- **FR-013**: A help overlay or footer MUST be visible showing the available keyboard shortcuts at all times.

### Key Entities

- **Bookmark**: Represents a saved web page. Attributes: unique URL, page title (string, may be "Untitled"), domain (extracted from URL), meta description (string, optional), date/time created, date/time last accessed (optional).
- **Search Query**: A user-entered string used to filter the bookmark list. Matched against title, domain, and description fields using fuzzy logic.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can save a new bookmark (from clipboard to confirmed save) in under 5 seconds when the network is responsive.
- **SC-002**: Fuzzy search results update within 100 milliseconds of each keystroke for collections of up to 1,000 bookmarks.
- **SC-003**: 100% of bookmarks saved during a session are still present after the application is closed and reopened.
- **SC-004**: Users can find a specific bookmark using a partial or misspelled search term within 3 keystrokes for any collection up to 500 bookmarks.
- **SC-005**: All primary actions (add, search, delete, open) are achievable without a mouse, using only keyboard shortcuts.
- **SC-006**: The application starts and displays the bookmark list in under 2 seconds regardless of collection size (up to 1,000 bookmarks).

## Assumptions

- The application targets a single user on a local machine; no multi-user or sync capability is in scope.
- Bookmarks are stored locally in a file-based format (e.g., JSON or similar); no external database is required.
- The user's system has a default browser configured for opening URLs.
- Clipboard access is available on the host operating system; the app reads from the system clipboard on demand.
- Page fetching is done at bookmark creation time only; the app does not re-fetch or refresh metadata after saving.
- The app runs in a standard terminal emulator that supports TUI rendering.
- No authentication, user accounts, or cloud sync features are in scope for this version.
- Meta description is stored as the "content" field for search purposes; full page body text indexing is out of scope.
