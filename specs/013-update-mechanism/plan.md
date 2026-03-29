# Implementation Plan: Cairn Self-Update Mechanism

**Branch**: `013-update-mechanism` | **Date**: 2026-03-26 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/013-update-mechanism/spec.md`

---

## Summary

Add a `cairn update` subcommand that checks GitHub Releases for a newer version of the cairn binary, downloads it, verifies its SHA256 checksum, and atomically replaces the running binary. A `--check` flag reports version status without applying changes. A `--extension` flag targets the Vicinae extension instead. No new external dependencies; uses stdlib `net/http` and `crypto/sha256` only.

---

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**: stdlib only (`net/http`, `crypto/sha256`, `encoding/json`, `os`, `runtime`) — no new external packages
**Storage**: N/A — no schema changes; extension version tracked via a plain `version.txt` file in the extension directory
**Testing**: `go test ./...`
**Target Platform**: Linux (amd64, arm64), macOS (amd64, arm64); Windows excluded from in-process update (locked executables)
**Project Type**: CLI tool
**Performance Goals**: Update completes in under 60 seconds on a standard broadband connection
**Constraints**: CGO_ENABLED=0 (all pure Go); single binary; no background network calls during normal subcommands
**Scale/Scope**: Single binary replacement + optional extension file replacement

---

## Constitution Check

| Gate | Status | Notes |
|------|--------|-------|
| No CGO | PASS | Stdlib `net/http` only; no new external deps |
| Single binary | PASS | No external runtime or service added |
| Task management | PASS | Tasks will be created via Backlog CLI after `/speckit.tasks` |
| Backward-compatible migrations | PASS | No SQLite schema changes |

---

## Project Structure

### Documentation (this feature)

```text
specs/013-update-mechanism/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   └── cli-contract.md  # CLI command schema
└── tasks.md             # Phase 2 output (/speckit.tasks — not yet created)
```

### Source Code (repository root)

```text
cmd/cairn/
└── main.go              # Add "update" subcommand handler (~60 lines)

internal/updater/
├── updater.go           # Core update logic (new package)
└── updater_test.go      # Unit tests with httptest server
```

**Structure Decision**: Single project layout (Option 1). The `internal/updater` package has one clear responsibility — version checking and binary/extension replacement — satisfying the constitution's no-unnecessary-abstractions rule. The subcommand dispatch stays in `main.go` consistent with all existing subcommands.

---

## Implementation Phases

### Phase A — Core CLI Update (FR-001 through FR-006, FR-009, FR-010)

Delivers: `cairn update` and `cairn update --check` for the binary.

**Tasks**:

1. **Create `internal/updater` package**
   - `CheckLatestVersion(currentVersion string) (latest string, available bool, err error)`: query GitHub Releases API, parse `tag_name`, compare with current.
   - `UpdateBinary(currentVersion, latestVersion string) error`: download binary, verify SHA256, backup existing, atomic rename.
   - HTTP client with 8s timeout and `cairn/<version>` User-Agent, mirroring `internal/fetcher`.

2. **Add `update` subcommand to `cmd/cairn/main.go`**
   - Parse `--check` and `--extension` flags.
   - Route to updater functions based on flags.
   - Output messages per CLI contract (`contracts/cli-contract.md`).
   - Map errors to exit codes (0/1/3/4).

3. **Unit tests for `internal/updater`**
   - Use `net/http/httptest` to serve mock GitHub API responses and mock binary downloads.
   - Cover: update available, already up to date, network error, checksum mismatch, permission error.

### Phase B — Extension Update (FR-007, FR-008)

Delivers: `cairn update --extension`.

**Tasks**:

4. **Extend `internal/updater` with extension functions**
   - `DetectExtension() (dir string, installed bool)`: resolve platform-specific extension directory.
   - `CheckExtensionVersion(dir string) (current, latest string, available bool, err error)`: read `version.txt`, query latest release.
   - `UpdateExtension(dir, latestVersion string) error`: download archive, verify checksum, extract to extension dir, write `version.txt`.

5. **Wire extension flags in `main.go`**
   - `cairn update --extension`: check + apply extension update.
   - `cairn update --extension --check`: report extension version status only.

6. **Unit tests for extension functions**
   - Mock filesystem state (extension installed / not installed).
   - Cover: not installed, already up to date, update available, checksum mismatch.

### Phase C — Help and Windows Guard

7. **Help text for `cairn update`**
   - Add help text per `contracts/cli-contract.md`.
   - Integrate with existing `cairn help` dispatch.

8. **Windows guard**
   - On `runtime.GOOS == "windows"`, print guidance to re-run the install script and exit 0 without touching any files.

---

## Complexity Tracking

No constitution violations. No complexity justification required.
