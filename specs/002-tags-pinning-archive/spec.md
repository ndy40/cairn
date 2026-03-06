# Feature Specification: Tags, Pinning, Archive & Startup Checks

**Feature Branch**: `002-tags-pinning-archive`
**Created**: 2026-03-06
**Status**: Draft
**Input**: User description: "I would like the app to detect the display type wayland or X11. On load, it should make sure pre-requisites are installed and give appropriate instructions. I would also like the ability to add tags to bookmarks and to later filter them. Bookmarks can have multiple tags. Keep a max limit of three. If possible, track the last time a bookmark was visited or opened. We need a way to archive bookmarks that haven't been visited in the last 6 months. However, we also want certain bookmarks to be permanent as though they are not always in use, they come in handy when needed example flight booking site."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Startup Pre-requisite Check (Priority: P1)

When the application launches, it detects the display environment (Wayland or X11) and checks that the required clipboard tool is installed for that environment. If a required tool is missing, the application displays a clear message explaining what to install and how, then exits gracefully rather than silently failing later.

**Why this priority**: Without clipboard access, the core add-bookmark-from-clipboard feature (Ctrl+P) is broken for all users. Surfacing this at startup prevents confusing errors mid-session and makes the app self-documenting for new users.

**Independent Test**: Launch the app in a Wayland session without `wl-clipboard` installed and verify a clear install instruction is shown and the app exits. Then install the tool and verify normal launch proceeds without any warning.

**Acceptance Scenarios**:

1. **Given** a Wayland session where `wl-paste` is not installed, **When** the app starts, **Then** the app displays "Wayland detected. Please install wl-clipboard (e.g. sudo apt install wl-clipboard)" and exits with a non-zero code.
2. **Given** an X11 session where neither `xclip` nor `xsel` is installed, **When** the app starts, **Then** the app displays which tool to install and exits with a non-zero code.
3. **Given** all required tools are present for the detected display type, **When** the app starts, **Then** the app launches normally with no startup warning.
4. **Given** neither `WAYLAND_DISPLAY` nor `DISPLAY` is set, **When** the app starts, **Then** a non-blocking warning is shown that clipboard paste may not be available, and the app continues to launch.

---

### User Story 2 - Tag Bookmarks (Priority: P2)

A user categorises their bookmarks with short labels so they can find related links later. When adding or editing a bookmark, they can assign up to three tags. The system enforces the three-tag maximum and de-duplicates identical tags silently.

**Why this priority**: Tags are the prerequisite for filtering (US3). Without the ability to assign tags, the filter feature has nothing to work with. Tags also make the collection self-organising as it grows.

**Independent Test**: Add a bookmark and assign 1, 2, and 3 tags — verify all are saved and shown. Attempt to add a 4th tag and verify it is rejected with an explanatory message.

**Acceptance Scenarios**:

1. **Given** the add-bookmark modal is open, **When** the user enters one or more tags, **Then** the bookmark is saved with those tags attached and displayed in the list.
2. **Given** a bookmark already has 3 tags, **When** the user tries to enter a 4th tag, **Then** input is blocked and the message "Maximum 3 tags per bookmark" is shown.
3. **Given** a tag contains only whitespace or is empty, **When** the user submits, **Then** the empty tag is silently ignored and not saved.
4. **Given** the user enters duplicate tags (e.g. "work" and "work"), **When** the bookmark is saved, **Then** only one instance of the tag is stored.
5. **Given** an existing bookmark is edited, **When** the user changes its tags, **Then** the updated tags replace the previous ones and are immediately reflected in the list.

---

### User Story 3 - Filter Bookmarks by Tag (Priority: P3)

A user narrows the bookmark list to a specific topic by selecting one or more tags from a tag filter panel. Only bookmarks matching any selected tag are shown. The filter can be combined with the text search. Clearing the filter restores the full list.

**Why this priority**: Filtering is the primary payoff for tagging. Without it, tags are decorative labels with no functional benefit. This transforms tag data into real navigational value.

**Independent Test**: With bookmarks assigned to different tags, select a single tag filter and verify only matching bookmarks appear. Select a second tag and verify bookmarks matching either tag appear. Clear the filter and verify the full list is restored.

**Acceptance Scenarios**:

1. **Given** bookmarks with various tags exist, **When** the user selects a tag filter, **Then** only bookmarks tagged with that tag are shown.
2. **Given** a tag filter is active, **When** the user clears the filter, **Then** the full bookmark list is restored.
3. **Given** multiple tags are selected as filters, **When** the list updates, **Then** bookmarks matching any of the selected tags are shown (OR logic).
4. **Given** a tag filter is active and the user types a search query, **When** the list updates, **Then** results match the selected tag(s) AND the search query simultaneously.
5. **Given** a tag filter is active but no bookmarks match, **When** the list is empty, **Then** the message "No bookmarks found with the selected tag(s)" is shown.

---

### User Story 4 - Track Last Visited Date (Priority: P4)

The application records when each bookmark was last opened. This date is shown in the bookmark list so the user can see at a glance how recently they used each link. It also drives automatic archiving eligibility.

**Why this priority**: Last-visited data is the prerequisite for meaningful archiving (US5). Without it, the archive feature has no signal to act on. It also helps users identify stale bookmarks themselves.

**Independent Test**: Open a bookmark and verify its last-visited date updates to the current date and time. Open it again later and verify the date reflects the most recent visit.

**Acceptance Scenarios**:

1. **Given** a bookmark has never been opened via the app, **When** it is displayed, **Then** the last-visited field shows "Never".
2. **Given** the user opens a bookmark by pressing Enter, **When** the browser launches, **Then** the bookmark's last-visited date is updated to the current date and time.
3. **Given** a bookmark was opened previously, **When** the user opens it again, **Then** the last-visited date updates to the new visit time, replacing the previous value.

---

### User Story 5 - Auto-Archive Stale Bookmarks (Priority: P5)

On every startup, the application checks for bookmarks not visited in the past 6 months (and not marked permanent). Those bookmarks are silently moved to an archive. Archived bookmarks are hidden from the main list but remain visible in a dedicated archive view. The user can restore any archived bookmark at any time.

**Why this priority**: As the bookmark collection grows, stale links accumulate and reduce the signal-to-noise ratio for search and browsing. Auto-archiving keeps the active list focused without destroying data.

**Independent Test**: Seed bookmarks with last-visited dates older than 6 months (non-permanent). Launch the app and verify those bookmarks move to the archive. Verify a permanent bookmark with the same old date remains in the active list.

**Acceptance Scenarios**:

1. **Given** non-permanent bookmarks with no visit in 6+ months exist, **When** the app starts, **Then** those bookmarks are moved to the archive and disappear from the active list.
2. **Given** a permanent bookmark has not been visited in 6+ months, **When** the app starts, **Then** it remains in the active list and is not archived.
3. **Given** one or more bookmarks are archived on startup, **When** the app finishes loading, **Then** a brief status message shows how many bookmarks were archived (e.g. "2 bookmarks archived").
4. **Given** bookmarks have been archived, **When** the user opens the archive view, **Then** archived bookmarks are visible with their title, URL, tags, creation date, and last-visited date.
5. **Given** a bookmark is in the archive, **When** the user chooses to restore it, **Then** it reappears in the active list with its archived status cleared.
6. **Given** a bookmark was added but never visited and is older than 6 months, **When** the archive check runs, **Then** it is treated as stale and moved to the archive (creation date used as fallback threshold).

---

### User Story 6 - Mark Bookmark as Permanent (Priority: P6)

A user marks certain bookmarks as permanent to protect them from ever being automatically archived. Permanent bookmarks remain in the active list indefinitely, regardless of how long they have been unvisited. The permanent flag can be toggled on or off at any time.

**Why this priority**: Without a protection mechanism, the archive feature would remove genuinely important but rarely-used bookmarks. Permanent marking is the explicit user override that makes archiving safe to enable.

**Independent Test**: Mark a bookmark as permanent. Manually trigger the archive check (or wait for startup). Verify the bookmark remains in the active list. Remove the permanent flag and re-run the archive check; verify it is now eligible for archiving.

**Acceptance Scenarios**:

1. **Given** a bookmark is selected, **When** the user marks it as permanent, **Then** a permanent indicator is displayed next to the bookmark in the list.
2. **Given** a bookmark is marked as permanent and has not been visited in 6+ months, **When** the archive check runs, **Then** the bookmark is not archived.
3. **Given** a bookmark is marked as permanent, **When** the user removes the permanent flag, **Then** the bookmark becomes eligible for archiving on the next check.
4. **Given** the bookmark list is displayed, **When** permanent bookmarks are present, **Then** they are visually distinguished from non-permanent bookmarks.

---

### Edge Cases

- A bookmark added for the first time has a null last-visited date. The archive check should not archive it based solely on never having been visited; creation date is the fallback — only bookmarks older than 6 months qualify.
- A bookmark whose permanent flag is removed becomes immediately eligible for archiving on the next startup check (no grace period).
- Tag names are stored in lowercase; input "Work" and "work" are treated as the same tag.
- A tag name exceeding 32 characters is truncated to 32 characters at the input level without an error.
- If the clipboard environment check fails for an ambiguous reason (display variable set but display unreachable), the app warns and continues rather than blocking the launch.
- Restoring an archived bookmark should not reset its last-visited date; it retains the original value.
- The archive view shows bookmarks in reverse archived-at order (most recently archived first).
- If both `WAYLAND_DISPLAY` and `DISPLAY` are set (XWayland session), Wayland takes precedence and `wl-paste` is checked first.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: On TUI launch, the application MUST detect the display environment by checking `WAYLAND_DISPLAY` (Wayland) and `DISPLAY` (X11) environment variables; Wayland takes precedence when both are set.
- **FR-002**: On TUI launch, the application MUST verify the required clipboard tool is available: `wl-paste` for Wayland; `xclip` or `xsel` for X11.
- **FR-003**: If the required clipboard tool is absent, the application MUST display a specific install instruction for the detected display type and exit with a non-zero status code before opening the TUI.
- **FR-004**: If no display environment can be detected, the application MUST display a non-blocking warning and continue launching without exiting.
- **FR-005**: The prerequisite check MUST only apply to TUI mode; CLI subcommands (`bm add`, `bm list`, etc.) are not affected.
- **FR-006**: Each bookmark MUST support between 0 and 3 tags; the application MUST reject any attempt to add a 4th tag with a message "Maximum 3 tags per bookmark".
- **FR-007**: Tags MUST be stored in lowercase; input casing is normalised on save.
- **FR-008**: Duplicate tags on the same bookmark MUST be silently de-duplicated before saving.
- **FR-009**: Tag input MUST be truncated to a maximum of 32 characters per tag.
- **FR-010**: Tags MUST be visible in the bookmark list view alongside the title and domain.
- **FR-011**: Users MUST be able to add, edit, and remove tags from a bookmark through the add/edit modal.
- **FR-012**: Users MUST be able to filter the active bookmark list by one or more tags; multiple selected tags use OR logic.
- **FR-013**: Tag filtering and text search MUST be composable; both filters apply simultaneously when both are active.
- **FR-014**: The application MUST record the date and time a bookmark is opened each time the user opens it via the TUI.
- **FR-015**: Bookmarks that have never been opened MUST display "Never" in the last-visited field.
- **FR-016**: On every TUI startup, the application MUST run an archive check that identifies bookmarks eligible for archiving.
- **FR-017**: A bookmark is eligible for archiving if it is not permanent AND (last-visited is more than 183 days ago OR has never been visited AND was created more than 183 days ago).
- **FR-018**: Eligible bookmarks MUST be automatically moved to the archive without user confirmation; a count of archived bookmarks MUST be displayed in a status message at startup.
- **FR-019**: Archived bookmarks MUST be hidden from the main active list and all main-list searches.
- **FR-020**: A dedicated archive view MUST be accessible from the browse mode and MUST display archived bookmarks with their title, URL, tags, creation date, last-visited date, and archived date.
- **FR-021**: Users MUST be able to restore any archived bookmark to the active list from the archive view.
- **FR-022**: Users MUST be able to mark any active bookmark as permanent via a keyboard shortcut.
- **FR-023**: Users MUST be able to remove the permanent flag from a bookmark via the same mechanism.
- **FR-024**: Permanent bookmarks MUST be visually distinguished in the list view (e.g. a pin icon or label).
- **FR-025**: Permanent bookmarks MUST never be moved to the archive by the automated check, regardless of last-visited date or creation date.

### Key Entities

- **Bookmark** (extended from feature 001): Gains new attributes — `tags` (list of up to 3 lowercase strings), `last_visited_at` (datetime, nullable — null means never visited), `is_permanent` (boolean, default false), `is_archived` (boolean, default false), `archived_at` (datetime, nullable).
- **Tag**: A short text label (1–32 lowercase characters) stored as part of the bookmark record. Not a standalone entity.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: On a system missing the required clipboard tool, the app displays a specific install instruction within 1 second of launch and exits cleanly with a non-zero status.
- **SC-002**: Users can assign all 3 tags to a bookmark in under 20 seconds from the add/edit modal.
- **SC-003**: Tag filtering narrows the bookmark list within 200 milliseconds of selecting a tag for collections up to 1,000 bookmarks.
- **SC-004**: 100% of bookmarks opened via the TUI have their last-visited date updated to the correct timestamp.
- **SC-005**: The startup archive check completes and any newly stale bookmarks are archived within 2 seconds for collections up to 1,000 bookmarks.
- **SC-006**: Zero permanent bookmarks are moved to the archive in any automated archive run.
- **SC-007**: Restored archived bookmarks retain 100% of their original data (title, URL, tags, dates).

## Assumptions

- Wayland detection: `WAYLAND_DISPLAY` set → Wayland session; otherwise `DISPLAY` set → X11; neither set → unknown.
- When both `WAYLAND_DISPLAY` and `DISPLAY` are set (XWayland), Wayland is assumed and `wl-paste` is checked first.
- Required tools: Wayland → `wl-paste` (from `wl-clipboard`); X11 → `xclip` or `xsel` (either is sufficient).
- The prerequisite check only blocks TUI launch; CLI subcommands do not perform this check.
- Archive threshold is 183 days (approximately 6 months) measured against UTC timestamps.
- Never-visited bookmarks use creation date as the fallback age for archive eligibility.
- Archive check runs silently on startup; no interactive confirmation is required from the user.
- Tag filtering uses OR logic across selected tags; combined with text search using AND logic.
- Tags are stored in lowercase; display shows them as stored (lowercase).
- The archive view is a separate list accessible via a keyboard shortcut from browse mode; archived bookmarks are not included in main-mode fuzzy search.
- This feature extends the data model introduced in feature 001-tui-bookmark-manager; migration of existing bookmark records to include the new fields is handled automatically on startup.
