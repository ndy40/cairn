# Tasks: Bookmark Cloud Sync

**Input**: Design documents from `/specs/001-bookmark-sync/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Not explicitly requested. Test tasks are omitted. Unit and integration tests should be written as part of implementation where appropriate.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **CLI entry**: `cmd/cairn/main.go`
- **New sync package**: `internal/sync/`
- **New backend subpackage**: `internal/sync/backend/`
- **Existing store package**: `internal/store/`
- **Existing TUI models**: `internal/model/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add new dependencies and create package scaffolding

- [ ] T001 Add new dependencies: `go get golang.org/x/oauth2 github.com/dropbox/dropbox-sdk-go-unofficial/v6 github.com/google/uuid`
- [ ] T002 [P] Create directory structure: `internal/sync/` and `internal/sync/backend/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Schema migration, SyncBackend interface, SyncConfig, snapshot format, and merge algorithm — all shared across user stories

**CRITICAL**: No user story work can begin until this phase is complete

- [ ] T003 Implement schema migration V3 in `internal/store/store.go`: add `uuid` and `updated_at` columns to `bookmarks` table with DEFAULT values, backfill existing rows (uuid from google/uuid, updated_at from created_at), create unique index on uuid, create index on updated_at, create `pending_sync` table with id/bookmark_uuid/operation/payload/created_at/retry_count columns, create index on pending_sync.created_at, record version 3 in schema_version. Follow existing migration pattern in the `migrations` slice.
- [ ] T004 Modify `internal/store/bookmark.go`: update `Insert()` to generate UUID v4 and set `updated_at` to current time on every new bookmark. Update `UpdateTags()` to set `updated_at` to current time. Update `DeleteByID()` to accept and return the bookmark UUID before deletion (needed for tombstone recording). Update `scanBookmark()` and `scanBookmarks()` to read the new uuid and updated_at columns. Add `uuid` and `UpdatedAt` fields to the `Bookmark` struct.
- [ ] T005 Create pending sync store methods in `internal/store/sync.go`: implement `AddPendingChange(bookmarkUUID, operation, payload string) error` (insert into pending_sync within same transaction as bookmark write), `ListPendingChanges() ([]*PendingChange, error)`, `ClearPendingChanges() error`, `PendingChangeCount() (int, error)`, and `ExportAllBookmarks() ([]*Bookmark, error)` (all non-archived bookmarks with uuid and updated_at for snapshot building). Define `PendingChange` struct matching pending_sync columns.
- [ ] T006 [P] Define SyncBackend interface and error types in `internal/sync/backend/backend.go`: interface with `Upload(data []byte, remotePath string) error`, `Download(remotePath string) ([]byte, error)`, `Exists(remotePath string) (bool, error)` methods. Define sentinel errors: `ErrNotFound`, `ErrAuthExpired`, `ErrNetworkFailure`, `ErrQuotaExceeded`.
- [ ] T007 [P] Implement SyncRecord snapshot types and JSON serialization in `internal/sync/snapshot.go`: define `SyncRecord` struct with Version int, LastUpdatedBy string, LastUpdatedAt time.Time, Bookmarks []BookmarkEntry, Tombstones []TombstoneEntry. Define `BookmarkEntry` with all bookmark fields plus Deleted bool. Define `TombstoneEntry` with UUID, URL, DeletedAt, DeletedBy. Implement `MarshalSnapshot(record *SyncRecord) ([]byte, error)` and `UnmarshalSnapshot(data []byte) (*SyncRecord, error)`.
- [ ] T008 [P] Implement SyncConfig load/save in `internal/sync/config.go`: define `SyncConfig` struct matching data-model.md (Backend, DeviceID, LastSyncAt, SyncDeclined, Dropbox credentials). Implement `DefaultConfigPath() string` (XDG_CONFIG_HOME/cairn/sync.json on Linux, ~/Library/Application Support/cairn/sync.json on macOS). Implement `LoadConfig() (*SyncConfig, error)` (returns nil if file doesn't exist), `SaveConfig(cfg *SyncConfig) error` (write JSON with mode 0600), `IsConfigured(cfg *SyncConfig) bool` (true if backend is set and not just declined).
- [ ] T009 Implement merge algorithm in `internal/sync/merge.go`: implement `MergeBookmarks(local []*store.Bookmark, remote *SyncRecord) (toInsert []*BookmarkEntry, toUpdate []*BookmarkEntry, toDelete []string, error)`. Logic: iterate remote bookmarks, match by URL against local set; if URL exists locally and remote updated_at is newer, mark for update; if URL not in local, mark for insert; iterate remote tombstones, if UUID exists locally mark for delete. Return three slices for the caller to apply atomically.

**Checkpoint**: Foundation ready — SyncBackend interface defined, snapshot format ready, merge algorithm implemented, store has uuid/updated_at/pending_sync support

---

## Phase 3: User Story 1 — First-Time Sync Setup (Priority: P1) MVP

**Goal**: Users can run `cairn sync setup` (or accept first-run prompt) to authenticate with Dropbox, download existing bookmarks, and configure sync on their device.

**Independent Test**: Run `cairn sync setup` on a fresh device, authenticate with Dropbox, verify bookmarks from cloud appear locally. Run it with existing local bookmarks and verify merge behavior.

### Implementation for User Story 1

- [ ] T010 [US1] Implement Dropbox backend in `internal/sync/backend/dropbox.go`: create `DropboxBackend` struct implementing SyncBackend interface. Constructor takes oauth2 token and app key. Use `dropbox-sdk-go-unofficial/v6` for `Upload()` (files/upload with WriteMode.Overwrite at `/cairn/sync.json`), `Download()` (files/download), and `Exists()` (files/get_metadata, return ErrNotFound on path_not_found). Wrap HTTP client with `golang.org/x/oauth2.TokenSource` for automatic token refresh. Map Dropbox API errors to sentinel errors (ErrAuthExpired, ErrNetworkFailure, ErrQuotaExceeded).
- [ ] T011 [US1] Implement OAuth2 PKCE flow for Dropbox in `internal/sync/config.go` (or separate `internal/sync/auth.go`): implement `RunOAuth2Flow(appKey string) (*oauth2.Token, error)` using PKCE with no-redirect. Generate code_verifier and code_challenge (S256). Build authorization URL with Dropbox OAuth2 endpoint, print to stdout. Read authorization code from stdin. Exchange code + verifier for token via `golang.org/x/oauth2.Config.Exchange()`. Return token containing access_token, refresh_token, and expiry.
- [ ] T012 [US1] Implement sync engine Setup() in `internal/sync/engine.go`: create `Engine` struct holding `*store.Store`, `SyncBackend`, `*SyncConfig`. Implement `Setup(appKey string) error`: run OAuth2 flow (T011), generate device UUID, save config. Check if cloud snapshot exists via backend.Exists(). If exists: download, unmarshal, call merge (T009), apply inserts/updates/deletes to store atomically. If not exists: export all local bookmarks via store.ExportAllBookmarks(), build SyncRecord, marshal, upload via backend.Upload(). Update last_sync_at. Return bookmark count for display.
- [ ] T013 [US1] Add `cairn sync setup` CLI subcommand in `cmd/cairn/main.go`: add "sync" to the subcommand router. Parse `sync setup`, `sync push`, `sync pull`, `sync status`, `sync auth`, `sync unlink` sub-subcommands. For this task, implement only the `setup` path: open store, call engine.Setup(), print "Sync configured. N bookmarks synced." on success, handle errors with appropriate exit codes (0=success, 1=auth failed, 3=error). Parse `--backend` flag (default "dropbox").
- [ ] T014 [US1] Implement first-run sync prompt in `cmd/cairn/main.go`: before executing any command, check if sync config exists (LoadConfig). If config is nil (no file), prompt "No sync configured — connect to Dropbox? (y/N):". On "y": run sync setup flow inline, then continue with original command. On "N"/Enter: create config file with `sync_declined: true` via SaveConfig, continue normally. Skip prompt entirely if config file exists (whether configured or declined).
- [ ] T015 [US1] Implement `cairn sync status` CLI subcommand in `cmd/cairn/main.go`: load sync config, query pending change count from store. Print formatted text output (Sync: configured/not configured, Backend, Device ID, Last sync, Pending changes). Support `--json` flag for JSON output. Exit code always 0.
- [ ] T016 [US1] Implement `cairn sync unlink` CLI subcommand in `cmd/cairn/main.go`: check sync is configured, prompt "Unlink this device from sync? Local bookmarks will be kept. (y/N)". On confirm: delete sync config file (os.Remove), clear pending_sync table via store.ClearPendingChanges(). Print confirmation. Exit codes: 0=success, 1=not configured/cancelled.

**Checkpoint**: User Story 1 complete — `cairn sync setup`, `cairn sync status`, `cairn sync unlink` work. First-run prompt functional. Dropbox auth flow works. Initial merge/upload on setup works.

---

## Phase 4: User Story 2 — Push Local Bookmarks to Cloud (Priority: P2)

**Goal**: Users can run `cairn sync push` to upload local bookmark changes (including deletions) to cloud storage.

**Independent Test**: Add a bookmark locally, run `cairn sync push`, verify the cloud snapshot contains the new bookmark. Delete a bookmark, push again, verify tombstone appears in cloud snapshot.

### Implementation for User Story 2

- [ ] T017 [US2] Modify `internal/store/bookmark.go` to record pending changes atomically: update `Insert()` to call `AddPendingChange(uuid, "add", jsonPayload)` within the same database transaction. Update `DeleteByID()` to call `AddPendingChange(uuid, "delete", "")` within the same transaction. Update `UpdateTags()` to call `AddPendingChange(uuid, "update", jsonPayload)` within the same transaction. This requires changing these methods to use explicit transactions (sql.Tx) instead of direct db.Exec.
- [ ] T018 [US2] Implement sync engine Push() in `internal/sync/engine.go`: implement `Push() (int, error)`. Check if cloud snapshot exists. If exists: download, unmarshal. Export all local bookmarks. Build new SyncRecord from local bookmarks, preserving existing tombstones from cloud and adding any new tombstones from pending_sync delete operations. Marshal and upload. Clear pending_sync table on success. Update last_sync_at in config. Return count of changes pushed.
- [ ] T019 [US2] Add `cairn sync push` CLI subcommand in `cmd/cairn/main.go`: load config, check IsConfigured, construct Dropbox backend and engine. Call engine.Push(). Print "Pushed N changes. Up to date." or "Already up to date." on success. Print "Sync not configured. Run `cairn sync setup` first." if unconfigured. Exit codes: 0=success, 1=not configured, 3=error.

**Checkpoint**: User Story 2 complete — `cairn sync push` works. Pending changes recorded atomically. Tombstones propagated. Cloud snapshot updated.

---

## Phase 5: User Story 3 — Pull Updates from Cloud (Priority: P3)

**Goal**: Users can run `cairn sync pull` to download bookmark changes from cloud and merge them locally.

**Independent Test**: Upload a snapshot with new bookmarks to Dropbox (from another device or manually), run `cairn sync pull`, verify new bookmarks appear locally. Add a tombstone to the cloud snapshot, pull, verify the bookmark is deleted locally.

### Implementation for User Story 3

- [ ] T020 [US3] Implement store batch import methods in `internal/store/sync.go`: implement `ImportBookmarks(bookmarks []*BookmarkEntry) (inserted int, updated int, error)` — for each entry, check if URL exists locally; if yes and remote updated_at is newer, update all fields; if no, insert with provided uuid and updated_at. Implement `DeleteByUUIDs(uuids []string) (int, error)` — delete bookmarks matching given UUIDs (for tombstone application). Both operations run in a single transaction for atomicity.
- [ ] T021 [US3] Implement sync engine Pull() in `internal/sync/engine.go`: implement `Pull() (inserted int, deleted int, error)`. Download cloud snapshot via backend.Download(). Unmarshal. Call merge algorithm (T009) with local bookmarks and remote snapshot. Apply inserts/updates via store.ImportBookmarks(). Apply deletes via store.DeleteByUUIDs(). Update last_sync_at in config. Return counts.
- [ ] T022 [US3] Add `cairn sync pull` CLI subcommand in `cmd/cairn/main.go`: load config, check IsConfigured, construct backend and engine. Call engine.Pull(). Print "Pulled N new, M deleted. Up to date." or "Already up to date." on success. Print "Sync not configured." if unconfigured. Exit codes: 0=success, 1=not configured, 3=error.

**Checkpoint**: User Story 3 complete — `cairn sync pull` works. Merge algorithm applies correctly. Tombstones delete local bookmarks. Bidirectional sync (push + pull) is fully functional.

---

## Phase 6: User Story 4 — Automatic Sync During Normal Use (Priority: P4)

**Goal**: Auto-pull on CLI startup, auto-push after every modifying operation, with graceful failure and pending queue replay.

**Independent Test**: Add a bookmark via `cairn add <url>` — verify it auto-pushes to cloud without running `cairn sync push`. Launch `cairn` TUI — verify it auto-pulls new bookmarks from cloud. Disconnect network, add a bookmark, reconnect, run any command — verify pending change is reconciled automatically.

### Implementation for User Story 4

- [ ] T023 [US4] Implement engine AutoPull() and AutoPush() in `internal/sync/engine.go`: `AutoPull() (int, error)` wraps Pull() but returns 0 on any error instead of failing (non-blocking). `AutoPush() error` wraps Push() but on failure increments retry_count on pending changes instead of propagating error. Both methods also replay any pending changes from the queue on success. Add `ReplayPending() (int, error)` that lists pending changes, applies them to the snapshot, uploads, and clears on success.
- [ ] T024 [US4] Add auto-pull hook to CLI startup in `cmd/cairn/main.go`: after the first-run prompt check and before executing the main command, if sync is configured, call engine.AutoPull(). On success with changes > 0, print `↓ N new bookmarks synced`. On failure, print `⚠ Sync pull failed: <reason>`. Continue with the original command regardless.
- [ ] T025 [US4] Add auto-push hooks after modifying CLI operations in `cmd/cairn/main.go`: after the `add`, `delete` subcommand handlers complete successfully, if sync is configured, call engine.AutoPush(). On failure, print `⚠ Sync push failed: change queued for later`. On success with replayed pending > 0, print `↑ Synced`.
- [ ] T026 [US4] Add auto-pull to TUI startup in `internal/model/app.go`: in the `Init()` function, add a tea.Cmd that calls engine.AutoPull() asynchronously (similar to how loadBookmarks() works). On completion, if new bookmarks were pulled, trigger a loadBookmarks() refresh. Display brief footer message if sync occurred or failed.
- [ ] T027 [US4] Add auto-push after TUI modifying operations in `internal/model/app.go`: after successful bookmark add (fetchAndSave completes), delete confirmation (yes), or tag edit (save), trigger engine.AutoPush() as an async tea.Cmd. Display brief footer message on sync status.
- [ ] T028 [US4] Implement `cairn sync auth` CLI subcommand in `cmd/cairn/main.go`: check sync is configured, run same OAuth2 PKCE flow as setup (reuse RunOAuth2Flow), update tokens in config file, attempt to replay pending changes via engine.ReplayPending(). Print "Re-authenticated. N pending changes synced." Exit codes: 0=success, 1=not configured/auth failed.
- [ ] T029 [US4] Add auth-expired detection to auto-sync in `internal/sync/engine.go`: in AutoPull() and AutoPush(), if the backend returns ErrAuthExpired, print `⚠ Sync auth expired — run 'cairn sync auth' to reconnect` instead of the generic failure message. Queue changes locally as normal.

**Checkpoint**: User Story 4 complete — auto-sync fires on all CLI and TUI operations. Offline changes queue and replay automatically. Auth expiry detected and reported with clear instructions.

---

## Phase 7: User Story 5 — Pluggable Backend Architecture (Priority: P5)

**Goal**: Validate that the SyncBackend interface is clean and extensible. An unsupported backend type produces a clear error. Adding a new backend requires only one new file.

**Independent Test**: Set backend to "s3" in config — verify clear error message. Verify that the Dropbox backend is constructed via a factory function that takes backend type as input.

### Implementation for User Story 5

- [ ] T030 [US5] Implement backend factory function in `internal/sync/backend/backend.go`: implement `NewBackend(cfg *sync.SyncConfig) (SyncBackend, error)` that switches on cfg.Backend ("dropbox" → construct DropboxBackend, unknown → return descriptive error "unsupported sync backend: %s"). Update engine construction in `cmd/cairn/main.go` to use this factory instead of directly constructing DropboxBackend.
- [ ] T031 [US5] Validate backend extensibility: ensure the `internal/sync/engine.go` Engine struct only references the `SyncBackend` interface, never `DropboxBackend` directly. Verify that no sync logic outside `internal/sync/backend/dropbox.go` imports the Dropbox SDK. If any coupling exists, refactor to use only the interface.

**Checkpoint**: User Story 5 complete — backend is fully abstracted. Adding S3 requires only a new file `internal/sync/backend/s3.go` implementing SyncBackend and a new case in the factory function.

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Final validation, edge case hardening, and cleanup

- [ ] T032 Handle edge case: corrupted cloud snapshot in `internal/sync/engine.go` — if UnmarshalSnapshot fails during Pull() or Setup(), return clear error suggesting re-initialising sync from a trusted device. Do not apply partial state.
- [ ] T033 Handle edge case: concurrent URL conflict in `internal/sync/merge.go` — when two devices add same URL, the merge should keep the record with the most recent `updated_at`, not create duplicates. Verify the merge algorithm handles this and add explicit handling if missing.
- [ ] T034 [P] Add sync help text to `cmd/cairn/main.go` help output: update the help subcommand to include sync subcommands (sync setup, sync push, sync pull, sync status, sync auth, sync unlink) with brief descriptions.
- [ ] T035 [P] Verify atomic operations: ensure all store methods that modify bookmarks and write to pending_sync use explicit database transactions. Verify that a crash between bookmark write and pending_sync write cannot leave inconsistent state.
- [ ] T036 Run quickstart.md validation: build the binary, execute each command from quickstart.md, verify expected outputs and exit codes.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Phase 2 — MVP deliverable
- **User Story 2 (Phase 4)**: Depends on Phase 2 — can start in parallel with US1 but logically follows it
- **User Story 3 (Phase 5)**: Depends on Phase 2 — can start in parallel with US1/US2 but logically follows US2
- **User Story 4 (Phase 6)**: Depends on Phases 3, 4, and 5 (uses Setup, Push, Pull, and pending queue from all three)
- **User Story 5 (Phase 7)**: Depends on Phase 2 — can start in parallel with US1-US3 (just refactoring backend construction)
- **Polish (Phase 8)**: Depends on all user stories being complete

### User Story Dependencies

- **US1 (Setup)**: Needs foundational phase only. Delivers: auth flow, config, initial merge, status, unlink
- **US2 (Push)**: Needs foundational phase. Builds on US1's engine and backend setup but is independently implementable
- **US3 (Pull)**: Needs foundational phase. Builds on merge algorithm. Independent of US2
- **US4 (Auto-sync)**: Needs US1 + US2 + US3 all complete (wraps Push/Pull with auto-fire and error handling)
- **US5 (Pluggable backend)**: Needs foundational phase only. Refactors backend construction

### Within Each User Story

- Store methods before engine logic
- Engine logic before CLI routing
- CLI routing before auto-sync hooks

### Parallel Opportunities

**Phase 2 (Foundational)**:
- T006 (SyncBackend interface), T007 (snapshot types), T008 (SyncConfig) can all run in parallel — different files, no dependencies

**Phase 3 (US1)**:
- T010 (Dropbox backend) and T011 (OAuth2 flow) can run in parallel — different concerns

**Phase 8 (Polish)**:
- T034 (help text) and T035 (atomic verification) can run in parallel — independent concerns

---

## Parallel Example: Foundational Phase

```
# These three tasks touch different files and can run simultaneously:
T006: Define SyncBackend interface in internal/sync/backend/backend.go
T007: Implement SyncRecord snapshot types in internal/sync/snapshot.go
T008: Implement SyncConfig load/save in internal/sync/config.go
```

## Parallel Example: User Story 1

```
# After T009 (merge) completes, these two can run simultaneously:
T010: Implement Dropbox backend in internal/sync/backend/dropbox.go
T011: Implement OAuth2 PKCE flow in internal/sync/auth.go
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001–T002)
2. Complete Phase 2: Foundational (T003–T009)
3. Complete Phase 3: User Story 1 (T010–T016)
4. **STOP and VALIDATE**: Test `cairn sync setup` end-to-end with a real Dropbox account
5. User can now set up sync, see status, and unlink

### Incremental Delivery

1. Setup + Foundational → foundation ready
2. Add User Story 1 → Test setup flow → Deploy (MVP!)
3. Add User Story 2 → Test push → Now devices can send changes
4. Add User Story 3 → Test pull → Bidirectional sync complete
5. Add User Story 4 → Test auto-sync → Seamless daily experience
6. Add User Story 5 → Validate extensibility → Future-proofed
7. Polish → Edge cases, help text, atomic verification

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- The Dropbox app key must be registered before T010/T011 can be tested against real Dropbox
