# Feature Specification: Bookmark Cloud Sync

**Feature Branch**: `001-bookmark-sync`
**Created**: 2026-03-11
**Status**: Draft
**Input**: User description: "Using Option 3, we will like to build data sync capability with the option to support different backend cloud storage for files like Dropbox or S3. Using Dropbox as the first implementation. We want to be able to use the cli program across different devices and keep bookmarks in sync."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - First-Time Sync Setup (Priority: P1)

A user installs cairn on a second device and wants to connect it to their existing bookmarks stored in Dropbox. On their first run of any cairn command, the CLI detects no sync is configured and prompts: "No sync configured — connect to Dropbox? (y/N)". The user accepts, authenticates with Dropbox, and their full bookmark collection is downloaded to the new device. Alternatively, the user can run `cairn sync setup` directly at any time.

**Why this priority**: This is the entry point for the entire sync feature. Without the ability to connect a device to a cloud backend, no other sync story is possible. It delivers immediate value by giving users access to their existing bookmarks on a new device.

**Independent Test**: Can be fully tested by installing cairn on a fresh device, running any cairn command, accepting the sync setup prompt, authenticating with Dropbox, and verifying all existing bookmarks appear in the local database.

**Acceptance Scenarios**:

1. **Given** a new device with no sync config, **When** the user runs any cairn command for the first time, **Then** the CLI prompts once to set up sync; accepting begins the Dropbox authentication flow.
2. **Given** the first-run sync prompt, **When** the user declines (N), **Then** the CLI proceeds normally without sync and does not prompt again unless the user runs `cairn sync setup`.
3. **Given** a new device with no cairn bookmarks, **When** the user completes sync setup and Dropbox authentication, **Then** all bookmarks from the user's cloud storage are downloaded and available locally.
4. **Given** the user has already set up sync, **When** they run `cairn sync setup` again, **Then** the system informs them sync is already configured and offers to reconfigure or cancel.
5. **Given** Dropbox authentication fails or is denied, **When** setup is attempted, **Then** the system displays a clear error, local bookmarks remain untouched, and the user can retry.

---

### User Story 2 - Push Local Bookmarks to Cloud (Priority: P2)

A user adds new bookmarks on their primary device and wants to make them available on their other devices. They run a sync command and their new bookmarks are uploaded to Dropbox.

**Why this priority**: Uploading changes is the core outbound half of the sync loop. Once setup is done, users need to be able to push their local changes to make sync useful across devices.

**Independent Test**: Can be fully tested by adding a bookmark locally, running the sync push command, then verifying the change appears in Dropbox and can be fetched on another device.

**Acceptance Scenarios**:

1. **Given** a user has added new bookmarks since the last sync, **When** they run the sync command, **Then** all new bookmarks are uploaded and the last-sync timestamp is updated.
2. **Given** a user has deleted a bookmark since the last sync, **When** they run the sync command, **Then** the deletion is recorded and propagated so other devices remove it on next pull.
3. **Given** no changes were made since the last sync, **When** the user runs the sync command, **Then** the system reports "already up to date" and exits cleanly.

---

### User Story 3 - Pull Updates from Cloud (Priority: P3)

A user who has been adding bookmarks on another device wants to get those additions on their current device. They run a sync pull command and the new bookmarks appear locally.

**Why this priority**: Pulling updates from the cloud completes the bidirectional sync loop. Combined with P1 and P2, it delivers a fully functional multi-device sync experience.

**Independent Test**: Can be fully tested by syncing from Device A, adding bookmarks there, then running a pull on Device B and verifying the new bookmarks appear.

**Acceptance Scenarios**:

1. **Given** another device has uploaded new bookmarks, **When** the user runs the sync pull command, **Then** all new bookmarks from the cloud appear in the local database without duplicates.
2. **Given** a bookmark deleted on another device was synced to the cloud, **When** the user pulls, **Then** that bookmark is removed from the local database.
3. **Given** the cloud storage is unreachable, **When** a pull is attempted, **Then** the local database is unchanged and a clear error message is displayed.

---

### User Story 4 - Automatic Sync During Normal Use (Priority: P4)

A user adds a bookmark on their laptop. Without running any extra command, the bookmark is automatically pushed to Dropbox. Later, they open cairn on their desktop and their new bookmark is already there — pulled automatically at startup.

**Why this priority**: Automatic sync is the feature that makes multi-device use seamless. Manual sync is the fallback; automatic sync is the day-to-day experience that eliminates friction.

**Independent Test**: Can be fully tested by adding a bookmark on Device A (no manual sync), then launching cairn on Device B and verifying the bookmark is present without the user running any sync command.

**Acceptance Scenarios**:

1. **Given** sync is configured, **When** the user adds a bookmark via any CLI operation, **Then** the new bookmark is automatically pushed to cloud storage before the CLI exits.
2. **Given** sync is configured and another device has uploaded new bookmarks, **When** the user starts the CLI on any device, **Then** those bookmarks are automatically pulled and available locally.
3. **Given** sync is configured but the network is unavailable, **When** the CLI starts or a modifying operation completes, **Then** the auto-sync fails gracefully and the primary CLI operation (add/delete/list) still succeeds.

---

### User Story 5 - Switch or Add Cloud Backend (Priority: P5)

A power user wants to switch from Dropbox to S3, or add S3 as an additional sync target. They update their configuration and the system uses the new backend going forward.

**Why this priority**: The pluggable backend design future-proofs the feature. Dropbox is first, but the architecture must accommodate S3 and others. This story validates the extensibility of the design without requiring S3 to be implemented now.

**Independent Test**: Can be validated by verifying the configuration structure supports multiple backend types and that adding a new backend requires only a single new adapter, not changes to core sync logic.

**Acceptance Scenarios**:

1. **Given** a user with Dropbox configured, **When** they update their configuration to point to a different backend type, **Then** subsequent syncs use the new backend without requiring code changes.
2. **Given** an unsupported backend type is specified in configuration, **When** any sync operation is attempted, **Then** the system reports a clear error identifying the unsupported backend.

---

### Edge Cases

- What happens when two devices add a bookmark with the same URL concurrently? The URL uniqueness constraint should deduplicate on merge — the bookmark is kept once.
- What happens when a bookmark is added on Device A and deleted on Device B before sync? Last-write-wins based on timestamp; the most recent operation wins.
- What happens when cloud storage is full or quota is exceeded? The sync command fails with a clear message; local data is preserved.
- What happens when the sync data file in cloud storage is corrupted or invalid? The system refuses to apply the corrupted state, reports the error, and suggests re-initialising sync from the trusted device.
- What happens when the user loses internet connectivity mid-sync? The operation is atomic — either it completes or it does not apply, with no partial state written to the local database. The failed changeset is logged to the pending change log for later reconciliation.
- What happens if the device is offline for an extended period with many pending changes? All changes accumulate in the pending log and are replayed in order when connectivity returns. The log has no expiry — changes are retained until successfully synced.
- What happens on first sync when both devices have independent bookmarks? Both sets are merged, deduplicated by URL, with the cloud record winning on conflict (based on `updated_at` timestamp); the user ends up with the union of all bookmarks.
- What happens when the primary device sets up sync and the cloud has no existing snapshot yet? The local bookmarks are uploaded to create the initial cloud snapshot. Subsequent devices merge against this snapshot.
- What happens when a new device with its own local bookmarks sets up sync? Local and cloud bookmarks are merged; cloud wins on URL conflict. The merged result is immediately available locally and the next auto-push uploads the reconciled state to cloud.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Users MUST be able to configure a cloud storage backend (initially Dropbox) for bookmark sync via `cairn sync setup` or by accepting a first-run prompt.
- **FR-001a**: On first run with no sync configuration present, the CLI MUST prompt the user once to set up sync. Declining MUST suppress future prompts until the user explicitly runs `cairn sync setup`.
- **FR-002**: The system MUST support an extensible backend model so that additional cloud providers (e.g., S3) can be added without altering core sync logic.
- **FR-003**: Users MUST be able to authenticate with Dropbox using the standard OAuth2 device flow from the CLI, without requiring a web browser redirect to localhost.
- **FR-004**: Users MUST be able to push local bookmark changes to the configured cloud backend via a CLI command.
- **FR-005**: Users MUST be able to pull bookmark changes from the configured cloud backend via a CLI command.
- **FR-006**: The system MUST merge bookmarks from the cloud with local bookmarks during setup and on every pull, deduplicating by URL with the cloud record winning on conflict (based on `updated_at` timestamp).
- **FR-006a**: When sync setup runs and the cloud snapshot does not yet exist (first device), the system MUST upload the local bookmarks to initialise the cloud snapshot.
- **FR-007**: The system MUST propagate bookmark deletions across devices, so a bookmark deleted on one device is removed on others after sync.
- **FR-008**: The system MUST track a last-sync timestamp per device to determine which records are new since the last sync.
- **FR-009**: The system MUST store sync configuration (backend type, credentials/tokens, last-sync timestamp) in a dedicated local configuration file separate from the bookmark database.
- **FR-010**: Sync credentials (OAuth tokens, API keys) MUST be stored securely and not in plain text in world-readable files.
- **FR-011**: The system MUST provide a clear status command showing sync configuration, last-sync time, and current backend.
- **FR-012**: All sync operations MUST be atomic — a failed sync MUST NOT leave the local database in a partial or inconsistent state.
- **FR-013**: The system MUST handle conflicts between concurrent edits using last-write-wins based on record timestamps.
- **FR-014**: Users MUST be able to unlink a device from cloud sync without deleting their local bookmarks.
- **FR-015**: The system MUST automatically pull the latest bookmarks from the configured cloud backend each time the CLI starts.
- **FR-016**: The system MUST automatically push changes to the configured cloud backend after every modifying operation (bookmark add, delete, or tag edit).
- **FR-017**: When an automatic sync push fails, the system MUST display a brief warning and record the failed changeset in a local pending sync log, then complete the local operation normally.
- **FR-018**: The system MUST attempt to reconcile and replay any pending (unsynced) changesets the next time a successful connection to the cloud backend is established.
- **FR-019**: The system MUST expose the number of pending (unsynced) changes in the sync status command so users can see when local changes are waiting to be pushed.
- **FR-020**: Automatic sync output MUST be event-driven and minimal — the system MUST print a brief status line only when a meaningful sync event occurs (bookmarks pulled, sync queued due to failure). Successful no-change syncs MUST produce no output.
- **FR-021**: The system MUST use long-lived OAuth2 refresh tokens for Dropbox authentication, silently refreshing access tokens when they expire without requiring user interaction.
- **FR-022**: When a refresh token itself expires or is revoked, the system MUST notify the user with a clear message (e.g., "Sync auth expired — run `cairn sync auth` to reconnect") and queue changes locally until re-authentication is complete.

### Key Entities

- **SyncConfig**: Represents the user's sync configuration — backend type, authentication credentials/tokens, device identifier, and last-sync timestamp. Stored locally per device.
- **SyncBackend**: An abstraction representing a cloud storage provider. Defines the operations: upload snapshot, download snapshot, check connectivity. Dropbox and S3 are concrete implementations.
- **SyncRecord**: A point-in-time snapshot of all bookmarks formatted for cloud exchange. Contains device ID, timestamp, and a list of bookmark entries including deletion markers (tombstones).
- **Bookmark** (existing): Already has URL, title, description, tags, archived status, permanent flag, and timestamps. Sync adds a globally-unique identifier and an updated-at timestamp per record.
- **PendingChange**: An entry in the local change log representing an unsynced operation (add, delete, edit). Contains the operation type, the affected bookmark identifier, a timestamp, and retry count. Stored in the SQLite database (`pending_sync` table) so that bookmark writes and their pending log entries are committed atomically. Entries are removed once successfully pushed to the cloud backend.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can set up sync on a new device and have their full bookmark collection available within 60 seconds under normal network conditions.
- **SC-002**: A sync push or pull operation for a collection of up to 10,000 bookmarks completes within 30 seconds under normal network conditions.
- **SC-003**: After any sync operation, both devices have an identical set of bookmarks with zero duplicates by URL.
- **SC-004**: Sync setup requires no more than 3 commands and one browser-based authentication step from the user.
- **SC-005**: A failed or interrupted sync leaves local bookmarks 100% intact with no data loss or corruption.
- **SC-006**: Adding a new cloud backend requires changes to only one isolated module, with no changes to core bookmark or database logic.
- **SC-007**: The sync status command returns current configuration, last-sync time, and pending change count within 1 second, even with no network connectivity.
- **SC-008**: All changes made while offline are automatically reconciled and pushed the next time the CLI operates with network connectivity, with no user action required.

## Clarifications

### Session 2026-03-11

- Q: Should sync be manual only, or also automatic? → A: Both. Manual sync is available via explicit CLI commands; automatic sync runs as part of the CLI's normal operation (e.g., on startup or on bookmark modification).
- Q: When does automatic sync trigger? → A: Pull from cloud on CLI startup; push to cloud on every modifying operation (add, delete, edit tags).
- Q: If automatic sync fails (e.g., no network), should it block the CLI operation or warn and continue? → A: Warn and continue (Option B). Failed sync operations are recorded in a local change log so they can be reconciled and retried automatically when connectivity is restored.
- Q: How visible should automatic sync output be during normal CLI use? → A: Brief, event-driven output only — print a status line when something meaningful happens (new bookmarks pulled, sync queued due to failure), silent on a successful no-op push.
- Q: How should expired Dropbox OAuth2 tokens be handled during auto-sync? → A: Use long-lived refresh tokens. Access tokens are refreshed silently in the background. Re-authentication is only required if the refresh token itself expires or is revoked.
- Q: Where should the pending change log be stored? → A: In the existing SQLite database as a new `pending_sync` table, enabling atomic writes alongside bookmark operations.

### Session 2026-03-11 (continued)

- Q: When a new device has no sync config, does the CLI prompt the user to set up sync on first run, or is it strictly opt-in? → A: On first run with no sync config, the CLI prompts once: "No sync configured — connect to Dropbox? (y/N)". Declining skips sync and the user can set it up later with `cairn sync setup`.
- Q: When setting up sync on a device that already has local bookmarks, how should they be merged with the cloud snapshot? → A: Merge both sets, deduplicate by URL, cloud record wins on conflict (based on most recent `updated_at` timestamp).

## Assumptions

- The first supported backend is Dropbox; S3 is planned but not implemented in this feature.
- Sync supports two modes: manual (explicit `cairn sync` command) and automatic. Automatic sync pulls on CLI startup and pushes after every modifying operation.
- The sync format is a JSON snapshot stored as a single file in cloud storage (not streaming/incremental). This keeps the implementation simple and the backend requirements minimal.
- Bookmark volume is assumed to be manageable for a single-file snapshot (up to tens of thousands of records). Pagination or chunking is out of scope for this feature.
- The device identifier is generated once at setup time and stored in the sync configuration. It does not need to be user-visible or user-settable.
- Conflict resolution uses last-write-wins based on `updated_at` timestamp. More sophisticated merge strategies are out of scope.
- The Dropbox OAuth2 flow uses the offline access type to obtain a long-lived refresh token; access tokens are refreshed silently. A running local web server for the redirect is not required.
- Credentials are stored in the OS config directory with file permissions restricted to the current user (mode 0600). OS-level keychain integration is out of scope for this feature.

## Dependencies

- Existing bookmark store with URL uniqueness constraint and timestamp fields (already present in V2 schema).
- Schema migration (V3) required to add: globally-unique identifiers and `updated_at` tracking per bookmark record, and a new `pending_sync` table for the change log.
- Dropbox API access (OAuth2 application registration required before implementation).
