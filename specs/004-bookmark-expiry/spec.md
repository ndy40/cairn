# Feature Specification: Bookmark Expiry & Last-Visited Removal

**Feature Branch**: `004-bookmark-expiry`
**Created**: 2026-03-06
**Status**: Draft
**Input**: User description: "We not able to detect when a bookmark has been visited or not. Remove this feature. Make an update to how bookmarks work. Every bookmark by default should expire after 30 days unless it has been pinned."

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Bookmarks Expire After 30 Days (Priority: P1)

A user saves bookmarks throughout the day. After 30 days, any bookmark that has not been pinned is automatically moved to the archive at next startup. This keeps the active list fresh and removes stale saves the user likely no longer needs, without requiring them to manually curate.

**Why this priority**: This is the core behavioural change requested. It replaces the previous archive threshold (which was based on last-visited detection) with a simpler, time-since-creation rule that works regardless of visit tracking.

**Independent Test**: Add a bookmark, simulate 31 days passing, restart the application, and verify the bookmark no longer appears in the active browse list and is instead visible in the archive view.

**Acceptance Scenarios**:

1. **Given** an active bookmark was created 30 or more days ago and is not pinned, **When** the application starts, **Then** the bookmark is automatically archived and removed from the active browse list.
2. **Given** an active bookmark was created 30 or more days ago and IS pinned, **When** the application starts, **Then** the bookmark remains in the active browse list and is NOT archived.
3. **Given** an active bookmark was created fewer than 30 days ago (not pinned), **When** the application starts, **Then** the bookmark remains in the active browse list unchanged.
4. **Given** multiple bookmarks exceed the 30-day threshold at startup, **When** the application starts, **Then** all qualifying bookmarks are archived in one pass and a count is displayed (e.g., "3 bookmark(s) archived").
5. **Given** no bookmarks exceed the 30-day threshold, **When** the application starts, **Then** no archiving occurs and no count message is shown.

---

### User Story 2 - Remove Last-Visited Tracking (Priority: P2)

Last-visited date is currently displayed on each bookmark row. Since visit detection does not work reliably, this data is misleading. The last-visited label and any update logic should be removed entirely from the visible interface.

**Why this priority**: The last-visited feature is broken and shows incorrect (stale) information. Removing it eliminates misleading UI content and simplifies the codebase. It is also a prerequisite for the expiry model, since expiry is now based on creation date only.

**Independent Test**: Launch the application and confirm no "Last:" or "Never visited" text appears on any bookmark row in the browse list or search results. Open a bookmark and return — the row still shows no visit-related label.

**Acceptance Scenarios**:

1. **Given** the application is running, **When** the user views the active bookmark list, **Then** no last-visited date or "Never visited" label is shown on any bookmark row.
2. **Given** the user opens a bookmark by pressing Enter, **When** the browser launches successfully, **Then** no visit timestamp is recorded and the list reloads showing the same row without last-visited data.
3. **Given** existing bookmarks that previously stored a last-visited value, **When** the application displays them, **Then** no last-visited information is surfaced in the UI regardless of what is stored.

---

### Edge Cases

- A bookmark already in the archive before this change: it remains archived. The 30-day rule only applies to active bookmarks.
- A bookmark pinned after it was archived: the user can restore it from the archive; once active and pinned, it is exempt from expiry.
- A bookmark created exactly 30 days ago: it is considered expired (threshold is ≥ 30 days).
- The device clock is wrong: expiry uses the stored UTC creation date compared to the current UTC time; no retroactive changes occur to previously stored creation dates.
- Bookmarks already in the archive at upgrade time: they remain archived unchanged; no re-archiving or purging of existing archived bookmarks occurs.

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST automatically archive all active non-pinned bookmarks whose creation date is 30 or more days before today, at every application startup.
- **FR-002**: Pinned bookmarks MUST be exempt from automatic expiry regardless of how old they are.
- **FR-003**: When at least one bookmark is archived at startup, the system MUST display a count message (e.g., "3 bookmark(s) archived"). When zero are archived, no message is shown.
- **FR-004**: The system MUST remove the last-visited date label ("Last: YYYY-MM-DD") from all bookmark rows in the browse list and search results.
- **FR-005**: The system MUST remove the "Never visited" label from all bookmark rows.
- **FR-006**: Opening a bookmark MUST NOT trigger any last-visited update; it only opens the URL and reloads the list.
- **FR-007**: The expiry check MUST evaluate only active (non-archived) bookmarks; bookmarks already in the archive are not affected.
- **FR-008**: The permanent-flag toggle, archive view, and restore flow remain fully functional and unchanged.

### Key Entities

- **Bookmark**: A saved web page with URL, title, domain, description, creation date, tags, permanent flag, archived flag, and archived date. The last-visited date field is no longer displayed or written to.
- **Expiry Rule**: A bookmark is subject to expiry if `is_archived = false` AND `is_permanent = false` AND `created_at` is 30 or more days before the current date.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All active non-pinned bookmarks older than 30 days are archived before the first frame of the TUI renders at startup.
- **SC-002**: Zero occurrences of "Last:", "Never visited", or any visit-date text appear in the browse list or search results after the change.
- **SC-003**: Pinned bookmarks of any age remain in the active list across 100% of application restarts.
- **SC-004**: The startup archive count message appears when ≥ 1 bookmark was expired, and is absent when 0 were expired.
- **SC-005**: The browse list correctly shows the post-expiry state on the very first render — no expired bookmarks are visible even briefly.

---

## Assumptions

- The expiry check runs at startup only (not as a background timer), consistent with the existing archive-check pattern.
- The 30-day threshold is: creation date + 30 calendar days ≤ today (UTC). A bookmark created on day 0 expires starting on day 30.
- The `last_visited_at` database column is left in place (no destructive schema migration to drop it); it is simply no longer written or read.
- The browse row description line retains domain and creation date but drops the last-visited segment.
- The startup count message format ("N bookmark(s) archived") remains the same as the existing implementation; only the threshold changes.
- No user-configurable expiry period is added; 30 days is the fixed default for this feature.
