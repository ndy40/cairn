# Research: Bookmark Cloud Sync

**Feature**: 001-bookmark-sync
**Date**: 2026-03-11

## Decision 1: Dropbox Go Library

**Decision**: Use `github.com/dropbox/dropbox-sdk-go-unofficial/v6`

**Rationale**: The only actively maintained Go library for Dropbox. Hosted under the Dropbox GitHub org. Auto-generated from Dropbox's public API spec. Pure Go (no CGO) — satisfies the Zero CGO constitution gate. Provides typed wrappers for `files.Upload()` and `files.Download()`.

**Alternatives considered**:
- `github.com/tj/go-dropbox` — less maintained, limited feature set
- Raw HTTP calls to Dropbox API — unnecessary effort when SDK abstracts headers and error handling
- `github.com/ncw/dropbox` — legacy v1 API, effectively abandoned

## Decision 2: OAuth2 Flow for CLI

**Decision**: Use PKCE authorization code flow with no-redirect (manual code entry)

**Rationale**: Dropbox does **not** support RFC 8628 (device authorization grant). Two options exist: (a) PKCE with localhost redirect, or (b) PKCE with no redirect (user copies code from browser). Option (b) is chosen because it avoids binding a local port (which can fail due to port conflicts or restricted environments), and aligns with the spec requirement (FR-003) of "no localhost redirect." The user visits a URL, approves, copies a code, and pastes it into the CLI.

**Alternatives considered**:
- PKCE with localhost redirect — better UX (auto-capture), but requires a local HTTP server and can fail on headless machines or port conflicts. Also contradicts FR-003.
- Long-lived access tokens — Dropbox no longer issues these; all tokens are now short-lived with mandatory refresh tokens.

## Decision 3: OAuth2 Library

**Decision**: Use `golang.org/x/oauth2`

**Rationale**: Official Go project. Pure Go, no CGO. Production-grade. Provides `TokenSource` abstraction that automatically refreshes expired access tokens — directly satisfies FR-021 (silent token refresh). Well-documented. Already used by the Dropbox unofficial SDK internally.

**Alternatives considered**:
- Manual HTTP token exchange — unnecessary complexity when `x/oauth2` handles PKCE, token exchange, and refresh natively.

## Decision 4: Sync Data Format

**Decision**: Single JSON file stored in Dropbox at a fixed path (e.g., `/cairn/sync.json`)

**Rationale**: The spec assumes a JSON snapshot approach (Assumptions section). A single file keeps the Dropbox interactions simple: one upload, one download. The file contains all bookmarks plus tombstones for deletions. At 10,000 bookmarks with average ~500 bytes each, the file is ~5MB — well within Dropbox's upload limit (150MB for simple upload endpoint).

**Alternatives considered**:
- One file per bookmark — complex to manage, many API calls, harder to make atomic
- Delta/changelog file — more complex to implement, requires merging multiple files, harder to reason about correctness
- Binary format (protobuf, msgpack) — adds a dependency and reduces debuggability for minimal size benefit

## Decision 5: Sync Configuration Storage

**Decision**: JSON config file at `$XDG_CONFIG_HOME/cairn/sync.json` (Linux) or `~/Library/Application Support/cairn/sync.json` (macOS), with file permissions mode 0600.

**Rationale**: FR-009 requires sync config to be separate from the bookmark database. FR-010 requires credentials not to be in world-readable files. JSON is human-readable for debugging. Mode 0600 restricts access to the current user. Using the OS config directory follows platform conventions and keeps config separate from data.

**Alternatives considered**:
- Store in SQLite alongside bookmarks — mixes concerns; config changes don't need transactional guarantees with bookmark writes; harder to inspect
- Environment variables for tokens — not persistent across sessions
- OS keychain (macOS Keychain, libsecret) — adds CGO dependencies or complex IPC; explicitly out of scope per spec assumptions

## Decision 6: Declined Sync Prompt Persistence

**Decision**: Store a `sync_declined` boolean flag in the sync config file. When present and true, suppress the first-run sync prompt.

**Rationale**: FR-001a requires that declining the first-run prompt suppresses future prompts. The simplest durable signal is a flag in the sync config file. If the file doesn't exist at all, the prompt is shown. If it exists with `sync_declined: true`, the prompt is suppressed. If it exists with backend config, sync is already set up.

**Alternatives considered**:
- SQLite flag — adds schema coupling for a one-time preference
- Separate dot file — yet another file to manage
