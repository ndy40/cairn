# Feature Specification: Async Autosync

**Feature Branch**: `007-async-autosync`
**Created**: 2026-03-11
**Status**: Draft
**Input**: User description: "The autosync call after inserting, editing or deleting a bookmark is blocking. It would be good to move it out of the main action to speed up user experience."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Non-blocking bookmark creation (Priority: P1)

A user adds a new bookmark via the CLI. The bookmark is saved locally and the command returns immediately. Sync with the cloud happens in the background without delaying the user's workflow.

**Why this priority**: Adding bookmarks is the most frequent operation. Eliminating the sync wait time on every add has the highest impact on perceived responsiveness.

**Independent Test**: Can be fully tested by adding a bookmark and measuring the time from command invocation to prompt return, confirming it completes without waiting for network activity.

**Acceptance Scenarios**:

1. **Given** sync is configured and the cloud is reachable, **When** a user adds a bookmark, **Then** the bookmark is saved locally and the command returns to the prompt before the sync operation completes.
2. **Given** sync is configured and the cloud is unreachable, **When** a user adds a bookmark, **Then** the bookmark is saved locally and the command returns to the prompt without any delay or error message related to sync.

---

### User Story 2 - Non-blocking bookmark deletion (Priority: P1)

A user deletes a bookmark via the CLI. The deletion is applied locally and the command returns immediately. Cloud sync occurs in the background.

**Why this priority**: Deletion is equally time-sensitive as creation; users expect instant feedback when removing items.

**Independent Test**: Can be fully tested by deleting a bookmark and confirming the command returns immediately, then verifying the sync eventually propagates.

**Acceptance Scenarios**:

1. **Given** sync is configured, **When** a user deletes a bookmark, **Then** the deletion is persisted locally and the command returns before the sync operation finishes.
2. **Given** sync is configured but the network is slow, **When** a user deletes a bookmark, **Then** the user is not blocked by the sync and receives immediate confirmation of the local deletion.

---

### User Story 3 - Non-blocking bookmark editing (Priority: P1)

A user edits a bookmark (e.g., updates tags) via the CLI. The edit is saved locally and the command returns immediately. Sync happens in the background.

**Why this priority**: Editing follows the same pattern as add/delete and users expect the same instant responsiveness.

**Independent Test**: Can be fully tested by editing bookmark tags and measuring prompt return time, confirming no sync-related delay.

**Acceptance Scenarios**:

1. **Given** sync is configured, **When** a user edits a bookmark, **Then** the edit is applied locally and the command returns before the sync completes.

---

### User Story 4 - Non-blocking startup auto-pull (Priority: P1)

A user launches any CLI command (e.g., `cairn list`, `cairn add`, or the TUI). The command starts immediately without waiting for a remote sync pull. The pull happens in a background process so remote changes are fetched without blocking the user.

**Why this priority**: The startup auto-pull is the single biggest source of latency — every CLI invocation pays a network round-trip cost before doing anything useful.

**Independent Test**: Run `time cairn list` with sync configured — should return in < 500ms regardless of network conditions.

**Acceptance Scenarios**:

1. **Given** sync is configured and the cloud is reachable, **When** a user runs any CLI command, **Then** the command begins executing immediately and the pull happens in the background.
2. **Given** sync is configured and the cloud is unreachable, **When** a user runs any CLI command, **Then** the command executes without any delay or error message related to sync.
3. **Given** sync is not configured, **When** a user runs any CLI command, **Then** no background pull process is spawned.

---

### User Story 5 - Background sync failure handling (Priority: P2)

When background sync fails (network error, auth expired, etc.), the failure is handled gracefully without disrupting the user or losing data. The pending changes remain queued for the next sync opportunity.

**Why this priority**: Without graceful failure handling, background sync could silently lose changes or produce confusing error output after the user has moved on.

**Independent Test**: Can be tested by simulating a network failure during background sync and verifying that pending changes are preserved and no error output disrupts the terminal.

**Acceptance Scenarios**:

1. **Given** sync is configured and the network fails during background sync, **When** the background sync encounters an error, **Then** the pending changes remain in the local queue for future sync and no error output is displayed to the user's active terminal session.
2. **Given** sync credentials have expired, **When** background sync attempts to push, **Then** the sync fails silently in the background and the next interactive sync command (e.g., `cairn sync push`) informs the user about the auth issue.

---

### Edge Cases

- What happens if the user runs multiple bookmark operations in rapid succession before the first background sync completes? The system must queue changes and not lose any pending sync records.
- What happens if the application exits immediately after a bookmark operation? The local change is already persisted; the background sync may not complete, but the pending change is preserved for the next run.
- What happens if background sync is still running when a new sync is triggered? The system must prevent concurrent sync operations to avoid conflicts.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST complete bookmark add, delete, and edit operations and return control to the user without waiting for the cloud sync to finish.
- **FR-002**: System MUST trigger the auto-push sync operation in the background after any bookmark-modifying operation (add, delete, edit tags).
- **FR-003**: System MUST preserve all pending sync changes locally regardless of whether the background sync succeeds or fails.
- **FR-004**: System MUST NOT display sync-related error output to the user's active terminal session when background sync fails.
- **FR-005**: System MUST prevent concurrent background sync operations from running simultaneously to avoid data conflicts.
- **FR-006**: System MUST ensure that rapid successive bookmark operations each record their pending changes before any background sync processes them.
- **FR-007**: System MUST NOT change the behavior of explicit sync commands (`cairn sync push`, `cairn sync pull`); these remain synchronous and interactive.
- **FR-008**: System MUST trigger the auto-pull sync operation in the background on CLI startup instead of blocking until the pull completes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All CLI commands (including startup) return to the user prompt in under 500 milliseconds, regardless of network conditions or sync configuration.
- **SC-002**: 100% of bookmark changes made while sync is configured are eventually synced to the cloud once connectivity is restored, with zero data loss.
- **SC-003**: No sync-related error messages appear in the user's terminal during normal bookmark operations (add, delete, edit).
- **SC-004**: Explicit sync commands (`cairn sync push/pull`) continue to provide synchronous feedback including error reporting.

## Assumptions

- The existing pending sync queue (pending_sync table) already reliably records changes before sync attempts. This feature only changes *when* the sync is triggered (background vs. blocking), not *what* is recorded.
- Auto-pull on CLI startup is also moved to a background process. Users do not need up-to-date remote data on every invocation; they can run `cairn sync pull` explicitly when freshness matters.
- The Vicinae extension (Raycast) delegates to the CLI, so it automatically benefits from this change without any extension-side modifications.
