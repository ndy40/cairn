# Quickstart: Implementing the Update Mechanism

**Feature**: 013-update-mechanism
**Date**: 2026-03-26

---

## What Gets Built

1. **`internal/updater/updater.go`** — Core update logic (version check, download, verify, replace).
2. **`cmd/cairn/main.go`** — `update` subcommand handler (≈60 lines, consistent with existing subcommands).
3. **`internal/updater/updater_test.go`** — Unit tests using an `http.Handler` test server.

No schema changes. No new external dependencies.

---

## Key Integration Points

### Version variable
`cmd/cairn/main.go:25`
```go
var version = "dev"
```
Pass this to `updater.CheckLatestVersion(version)`.

### Subcommand dispatch
`cmd/cairn/main.go` — add `"update"` to the switch block following the same pattern as `"sync"` or `"config"`.

### Binary location
```go
exe, _ := os.Executable()
exe, _ = filepath.EvalSymlinks(exe)
```

### HTTP client pattern
Follow `internal/fetcher/fetcher.go`: 8-second timeout, `cairn/<version>` User-Agent.

### Platform string
```go
// GOOS → "linux" | "darwin"
// GOARCH → "amd64" | "arm64"
binaryName := fmt.Sprintf("cairn-%s-%s", runtime.GOOS, runtime.GOARCH)
```

### GitHub API response (relevant fields only)
```json
{
  "tag_name": "v0.3.0",
  "assets": [
    { "name": "cairn-linux-amd64", "browser_download_url": "https://..." },
    { "name": "checksums.txt",     "browser_download_url": "https://..." }
  ]
}
```

### Checksum file format
```
abc123...  cairn-linux-amd64
def456...  cairn-darwin-arm64
```
Parse: split each line on two-or-more spaces; match by filename.

### Extension directory (Linux)
```go
dataHome := os.Getenv("XDG_DATA_HOME")
if dataHome == "" {
    dataHome = filepath.Join(os.UserHomeDir(), ".local", "share")
}
extDir := filepath.Join(dataHome, "vicinae", "extensions", "cairn")
```

### Extension version file
```
extDir/version.txt  →  contains plain version string, e.g. "v0.2.0"
```

---

## Atomic Binary Replacement

```go
// 1. Download to temp file in same dir as target
tmp, _ := os.CreateTemp(filepath.Dir(targetPath), ".cairn-update-*")
// 2. Stream body to tmp
// 3. Verify SHA256
// 4. Backup existing
os.Rename(targetPath, targetPath+".bak")
// 5. Move temp into place (atomic on same filesystem)
err := os.Rename(tmp.Name(), targetPath)
if err != nil {
    os.Rename(targetPath+".bak", targetPath) // restore
}
// 6. Remove backup on success
os.Remove(targetPath + ".bak")
```

---

## Running Tests

```bash
go test ./internal/updater/...
go build ./... && go vet ./...
```

Test the `--check` flag manually:
```bash
go run ./cmd/cairn update --check
```

---

## Acceptance Checklist (quick reference)

- [ ] `cairn update` downloads and replaces the binary when a newer version exists
- [ ] `cairn update` reports "already up to date" when current == latest
- [ ] `cairn update --check` prints version diff without modifying any files
- [ ] Checksum mismatch exits with code 3, no files modified
- [ ] Write permission error exits with code 4
- [ ] `cairn update --extension` reports "not installed" when extension dir is absent
- [ ] `cairn update --extension` updates extension files when available
- [ ] No network calls occur during `cairn add`, `cairn list`, or any other subcommand
- [ ] `go build ./... && go vet ./...` pass with zero errors
