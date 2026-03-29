# Research: Cairn Self-Update Mechanism

**Feature**: 013-update-mechanism
**Date**: 2026-03-26

---

## 1. Version Discovery

**Decision**: Query the GitHub Releases API (`/repos/ndy40/cairn/releases/latest`) to retrieve the latest tag name.

**Rationale**: The install script (`install.sh`) already uses this exact endpoint. Reusing the same pattern keeps behaviour consistent and avoids a second source of truth.

**Alternatives considered**:
- Polling a static `version.txt` file: requires extra CI step, more moving parts.
- RSS feed from GitHub releases: more complex parsing, no clear advantage.

---

## 2. Current Version Source

**Decision**: Read the `version` package-level variable in `cmd/cairn/main.go` (injected via `-X main.version=$TAG` ldflags at build time). The update command will receive it as a parameter from `main`.

**Rationale**: The variable already exists and is the canonical version. No file on disk needed.

**Alternatives considered**:
- Read a `~/.local/share/cairn/version` file: redundant, can drift from actual binary.
- Parse `os.Executable` headers: complex, fragile across platforms.

---

## 3. Version Comparison

**Decision**: Simple lexicographic string comparison on the full tag string (e.g., `v0.1.0` < `v0.2.0`). Treat `dev` as "always out of date" so development builds always offer an update.

**Rationale**: All tags in the repository follow the `vMAJOR.MINOR.PATCH` format. Lexicographic ordering on this format is equivalent to semantic version ordering for single-digit components. Bringing in a semver library would violate the no-unnecessary-abstractions rule.

**Alternatives considered**:
- `golang.org/x/mod/semver`: correct for all cases but adds a dependency.
- Numeric component parsing: correct and zero-dependency; acceptable fallback if multi-digit version components appear in future.

---

## 4. Binary Download and Replacement

**Decision**: Download `cairn-{GOOS}-{GOARCH}` from the release asset URL, stream to a temp file in the same directory as the installed binary, verify SHA256 against `checksums.txt`, then use `os.Rename` for atomic replacement. If `os.Rename` fails (cross-device), fall back to copy-then-delete.

**Rationale**:
- `os.Rename` is atomic on POSIX when source and destination are on the same filesystem — the most common case (same install dir).
- Temp file in the same directory avoids cross-filesystem issues in the common path.
- Checksum verification mirrors `install.sh` behaviour.

**Alternatives considered**:
- Write directly to binary path: not atomic, leaves broken binary if interrupted.
- Backup to `.bak` before rename: the install script does this for rollback; included as a safety step.

---

## 5. Checksum Verification

**Decision**: Download `checksums.txt` from the same release, parse lines as `<sha256>  <filename>`, compute SHA256 of the downloaded binary using stdlib `crypto/sha256`, compare in constant time.

**Rationale**: Mirrors `install.sh` exactly. Stdlib only — no new dependencies.

**Alternatives considered**:
- Skip checksum for "simple" update: violates the integrity requirement in FR-004.
- GPG signature: no GPG infrastructure in the release pipeline.

---

## 6. Current Binary Location

**Decision**: Use `os.Executable()` (with `filepath.EvalSymlinks` to resolve symlinks) to determine where to write the updated binary.

**Rationale**: Guaranteed to return the real path of the running process. Symlink resolution handles the case where cairn is symlinked from a bin directory.

**Alternatives considered**:
- `os.Args[0]`: not reliable — may be a relative path or shell alias.
- Hard-coded default path: breaks custom install directories.

---

## 7. Extension Location and Update

**Decision**: Mirror the `install.sh` extension directory logic:
- macOS: `~/Library/Application Support/vicinae/extensions/cairn`
- Linux: `$XDG_DATA_HOME/vicinae/extensions/cairn` → `~/.local/share/vicinae/extensions/cairn`

Extension "version" is determined by a `version.txt` file written into the extension directory at install time. The update checks a `vicinae-extension-{version}.tar.gz` (or zip) from the same GitHub release.

**Rationale**: Consistency with the install script's path logic.

**Alternatives considered**:
- Delegate to `vici` toolchain for extension updates: `vici` is not on public npm and may not be available.
- Skip extension versioning: then `--check` cannot report extension version, which reduces usefulness.

---

## 8. New Dependencies

**Decision**: Zero new external dependencies. Use stdlib only:
- `net/http` — HTTP requests (already used by `internal/fetcher`)
- `crypto/sha256`, `encoding/hex` — checksum
- `os`, `io`, `runtime` — file ops, platform detection
- `encoding/json` — parse GitHub API response

**Rationale**: Keeps the zero-CGO, single-binary constraints satisfied. The existing `internal/fetcher` pattern shows stdlib HTTP is sufficient.

---

## 9. Package Structure

**Decision**: New package `internal/updater/` with a single file `updater.go`. Exposes:
- `CheckLatestVersion(currentVersion string) (latestVersion string, updateAvailable bool, err error)`
- `UpdateBinary(latestVersion string) error`
- `CheckExtensionVersion() (installed string, latest string, updateAvailable bool, err error)`
- `UpdateExtension(latestVersion string) error`

The subcommand handler lives in `cmd/cairn/main.go` (consistent with all other subcommands).

**Rationale**: Single clear responsibility (version checking + binary/extension replacement). Keeps subcommand dispatch in main where all others live.

---

## 10. Platform Scope

**Decision**: Support Linux (amd64, arm64) and macOS (amd64, arm64). Windows is out of scope for the in-process self-update because Windows locks executables that are currently running.

**Rationale**: Windows users can re-run the install script for updates. The install script already handles Windows. Adding Windows in-process update requires elevated complexity (rename running exe, service restart) with minimal user base.

**Alternatives considered**:
- Windows update via temp process: significantly more complex, deferred to future feature.
