# Contract: SyncBackend Interface

**Feature**: 001-bookmark-sync
**Date**: 2026-03-11

## Interface Definition

The `SyncBackend` interface defines the contract that all cloud storage providers must implement. This is the extensibility point for adding S3, GCS, or other backends in the future.

### Operations

| Method | Input | Output | Description |
|--------|-------|--------|-------------|
| `Upload` | snapshot (byte slice), remote path (string) | error | Upload a sync snapshot to cloud storage, overwriting any existing file at the path |
| `Download` | remote path (string) | byte slice, error | Download the sync snapshot from cloud storage; returns specific error if file does not exist |
| `Exists` | remote path (string) | bool, error | Check whether the snapshot file exists in cloud storage |

### Error Conditions

| Error | When | Caller Action |
|-------|------|---------------|
| `ErrNotFound` | Download/Exists called and file does not exist | Treat as first-time setup (no cloud data yet) |
| `ErrAuthExpired` | Access token expired and refresh failed | Prompt user to re-authenticate (`cairn sync auth`) |
| `ErrNetworkFailure` | Cloud storage unreachable | Queue changes in pending_sync; warn user |
| `ErrQuotaExceeded` | Cloud storage quota full | Warn user; queue changes |

### Implementation Requirements

1. Each backend receives its credentials/config at construction time (no global state).
2. Upload must be atomic from the caller's perspective — either the full file is written or the previous version is preserved.
3. Download returns the full file contents; no streaming required for the expected file sizes (<10MB).
4. Backends must not cache responses — each call must hit the remote storage.

### Dropbox Implementation Notes

- `Upload` → calls Dropbox `files/upload` endpoint with `WriteMode.Overwrite`
- `Download` → calls Dropbox `files/download` endpoint
- `Exists` → calls Dropbox `files/get_metadata` endpoint; `ErrNotFound` on 409/path_not_found
- Remote path: `/cairn/sync.json`
- Auth: Uses `golang.org/x/oauth2.TokenSource` wrapping the Dropbox OAuth2 config for automatic token refresh

### Future S3 Implementation Notes (not in scope)

- `Upload` → `s3:PutObject`
- `Download` → `s3:GetObject`
- `Exists` → `s3:HeadObject`
- Config: bucket name, region, access key ID, secret access key
