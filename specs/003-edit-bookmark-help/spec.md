# Feature Specification: Edit Bookmark Tags, Last-Visited Visibility & CLI Help

**Feature Branch**: `003-edit-bookmark-help`
**Created**: 2026-03-06
**Status**: Draft
**Input**: User description: "Add a cli help for the program. Also in the TUI, add a way to add tags to already existing bookmarks. Also when is the last visited link in the bookmark list updated?"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Edit Tags on Existing Bookmarks (Priority: P1)

A user looks at their bookmark list and realises a bookmark has no tags (or the wrong tags). They select the bookmark, press a key to open an edit panel, modify the tags, and save. The updated tags are immediately reflected in the list without losing any other data (title, URL, visit dates).

**Why this priority**: Tags were introduced in feature 002 but there is no way to add or change tags on bookmarks that already existed before tags were available, or on bookmarks added before the user thought of a tag. Without edit support, the tagging feature is incomplete for real-world use.

**Independent Test**: Select a bookmark with no tags in the TUI; press `e`; type a tag in the tags field; press Enter to save. Verify the tag appears next to the bookmark in the list. Open the edit panel again and verify the tag is pre-filled. Clear all tags and save; verify no tags are shown.

**Acceptance Scenarios**:

1. **Given** a bookmark is selected in the browse list, **When** the user presses `e`, **Then** an edit panel opens showing the bookmark's title (read-only) and its current tags (editable, pre-filled).
2. **Given** the edit panel is open, **When** the user types new tags and presses Enter, **Then** the tags are saved and the browse list immediately shows the updated tags next to the bookmark.
3. **Given** the edit panel is open with existing tags, **When** the user clears the tags field and saves, **Then** the bookmark is saved with no tags and no tags are shown in the list.
4. **Given** the edit panel is open, **When** the user presses Esc, **Then** the panel closes with no changes saved.
5. **Given** the user enters more than 3 tags in the edit panel, **When** they save, **Then** only the first 3 tags are saved and a note "Only first 3 tags saved" is shown briefly.
6. **Given** a bookmark has the permanent flag or any other fields set, **When** the user edits its tags, **Then** only the tags field changes; all other fields (title, URL, dates, permanent flag) are preserved.

---

### User Story 2 - Last-Visited Date: Visibility and Update Timing (Priority: P2)

A user opens a bookmark by pressing Enter in the TUI and wants to confirm that the "last visited" date in the list immediately reflects that action — without restarting the app. They also want to understand what event triggers the update: it is the successful launch of the browser, not the moment the webpage loads or the browser closes.

**Why this priority**: The last-visited date drives the auto-archive feature. If users do not see it update immediately after opening a link, they lose trust in the archiving logic and may believe the feature is broken.

**Independent Test**: Open the TUI. Find a bookmark showing "Never visited". Press Enter to open it. Without quitting or restarting, observe the bookmark description line. Verify it now shows "Last: <today's date>" rather than "Never visited". Repeat from search mode to confirm both modes update correctly.

**Acceptance Scenarios**:

1. **Given** a bookmark has never been opened via the TUI, **When** it is displayed in the list, **Then** its description line shows "Never visited".
2. **Given** the user presses Enter on a bookmark in browse mode and the browser launches successfully, **When** the TUI list refreshes, **Then** the bookmark's description line shows "Last: <today's date>" without the user doing anything further.
3. **Given** the user presses Enter on a bookmark in search mode and the browser launches successfully, **When** the TUI list refreshes, **Then** the last-visited date updates identically to browse mode.
4. **Given** the browser fails to launch (e.g., no browser is installed), **When** an error message is shown, **Then** the last-visited date is NOT updated.
5. **Given** a bookmark was opened before and shows a previous date, **When** the user opens it again, **Then** the date updates to today's date, replacing the previous value.

---

### User Story 3 - Per-Subcommand CLI Help (Priority: P3)

A user types `bm --help` or `bm add --help` expecting standard help output. The application should respond to `-h` and `--help` flags on both the root command and every subcommand, printing actionable usage text and exiting cleanly.

**Why this priority**: Standard `--help` support is the expected convention for any CLI tool. Without it, users who follow standard CLI habits get unhelpful output or errors, reducing discoverability.

**Independent Test**: Run `bm --help` → full usage guide printed, exit 0. Run `bm add --help` → add subcommand usage printed, exit 0. Run `bm search --help` → search subcommand usage with flags printed, exit 0. Repeat for all subcommands.

**Acceptance Scenarios**:

1. **Given** the user runs `bm --help` or `bm -h`, **When** the output is shown, **Then** it lists all subcommands with brief descriptions, all global flags, and the `BM_DB_PATH` environment variable, then exits with code 0.
2. **Given** the user runs `bm add --help`, **When** the output is shown, **Then** it shows the required URL argument and a description of what the command does, then exits with code 0.
3. **Given** the user runs `bm search --help`, **When** the output is shown, **Then** it shows the query argument, `--json` flag, `--limit` flag with its default, and a brief description, then exits with code 0.
4. **Given** the user runs `bm list --help`, **When** the output is shown, **Then** it shows the `--json` flag and description, then exits with code 0.
5. **Given** the user runs `bm delete --help`, **When** the output is shown, **Then** it shows the required numeric ID argument, then exits with code 0.
6. **Given** the user runs `bm version --help` or `bm help --help`, **When** the output is shown, **Then** it shows a one-line description of the subcommand, then exits with code 0.

---

### Edge Cases

- If the user presses `e` when the bookmark list is empty, nothing happens (no edit panel opens).
- Editing tags on an archived bookmark is out of scope; `e` only works in the active browse list, not in the archive view.
- If the user opens the edit panel and saves without changing anything, the tags are re-saved with identical values — no error, no visible change.
- Tag editing applies the same normalisation rules as tag creation: whitespace trimming, lowercase conversion, deduplication, 32-character truncation per tag, maximum 3 tags.
- `bm --help` and `bm help` (as a subcommand) both print the same content.
- The last-visited update triggers on a successful `cmd.Start()` — meaning the browser process started. If the process starts but immediately fails to show a page, the date is still updated.

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: In the TUI browse mode, pressing `e` on a selected active bookmark MUST open an edit panel.
- **FR-002**: The edit panel MUST display the bookmark's title as a read-only label and the bookmark's current tags in an editable text field, pre-filled and comma-separated.
- **FR-003**: When the user confirms the edit panel (Enter), the updated tags MUST be persisted and the browse list MUST refresh immediately to show the new tags.
- **FR-004**: Pressing Esc in the edit panel MUST discard all changes and return to browse mode.
- **FR-005**: Tag editing MUST apply the same normalisation rules as tag creation: comma-separated, lowercase, deduplicated, truncated to 32 characters per tag, maximum 3 tags.
- **FR-006**: Editing tags MUST NOT modify any other bookmark field (title, URL, creation date, last-visited date, permanent flag, archive status).
- **FR-007**: The last-visited date MUST update and be visible in the browse list immediately after the user successfully opens a bookmark via Enter in both browse mode and search mode, without requiring an app restart.
- **FR-008**: If the browser launch command fails, the last-visited date MUST NOT be updated and an error message MUST be shown to the user.
- **FR-009**: The `bm` root command MUST respond to `--help` and `-h` flags with the full usage guide (all subcommands, global flags, environment variables) and exit with code 0.
- **FR-010**: Each subcommand (`add`, `list`, `search`, `delete`, `version`, `help`) MUST respond to `--help` and `-h` flags with subcommand-specific usage text and exit with code 0.

### Key Entities

- **Bookmark** (no new fields): The edit operation updates only the `tags` field on an existing bookmark record. All other fields are unchanged.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can update tags on an existing bookmark in under 15 seconds from pressing `e` to seeing the updated tags in the browse list.
- **SC-002**: 100% of bookmarks opened via the TUI (both browse and search modes) show an updated last-visited date in the list within 2 seconds of the browser launching, without any manual refresh.
- **SC-003**: `bm --help` and `bm <subcommand> --help` always exit with code 0 and print usage text within 1 second.
- **SC-004**: A new user can identify all subcommands, required arguments, and available flags from `bm --help` output alone, without consulting external documentation.

---

## Assumptions

- The edit panel uses the same tag input UX as the add modal: a single comma-separated text field. No per-tag chip UI is needed.
- Only tags are editable in this iteration; title and URL editing is out of scope.
- The `e` key is used to open the edit panel in browse mode; it is currently unbound.
- Last-visited updates immediately on successful `cmd.Start()` (browser process started), not after the page loads or browser closes.
- The same `e` key and edit panel are not available in archive view, search mode, or tag filter overlay — only browse mode.
- Subcommand help is triggered by adding a `-h`/`--help` flag to each subcommand's flag set; the flag prints usage and exits 0 before any operation is performed.
- `bm help` (subcommand) and `bm --help` (flag) are aliases and print identical content.
