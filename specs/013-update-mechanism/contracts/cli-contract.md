# CLI Contract: Update Command

**Feature**: 013-update-mechanism
**Date**: 2026-03-26

---

## Commands

### `cairn update`

Check for a newer version of the cairn binary and apply the update if one is available.

**Usage**:
```
cairn update [flags]
```

**Flags**:

| Flag          | Description                                                  |
|---------------|--------------------------------------------------------------|
| `--check`     | Report version status only; do not download or modify files  |
| `--extension` | Target the Vicinae extension instead of (or in addition to) the CLI binary |

**Flag combinations**:

| Invocation                        | Behaviour                                                        |
|-----------------------------------|------------------------------------------------------------------|
| `cairn update`                    | Check and apply CLI binary update                               |
| `cairn update --check`            | Report CLI version status; no changes                           |
| `cairn update --extension`        | Check and apply extension update                                |
| `cairn update --extension --check`| Report extension version status; no changes                     |

---

## Exit Codes

| Code | Meaning                                                       |
|------|---------------------------------------------------------------|
| `0`  | Success (update applied, or already up to date, or --check completed) |
| `1`  | General error (network failure, API error, unknown platform)  |
| `3`  | Checksum verification failed                                  |
| `4`  | Permission denied (cannot write to install directory)         |

---

## Standard Output (human-readable)

All output is written to stdout. Messages are terse, one line per event.

### `cairn update` (update available)
```
cairn: current version v0.1.0, latest v0.2.0
cairn: downloading cairn-linux-amd64...
cairn: verifying checksum...
cairn: updated to v0.2.0
```

### `cairn update` (already up to date)
```
cairn: already up to date (v0.2.0)
```

### `cairn update --check` (update available)
```
cairn: current version v0.1.0, latest v0.2.0 (update available)
```

### `cairn update --check` (up to date)
```
cairn: already up to date (v0.2.0)
```

### `cairn update --extension` (extension not installed)
```
cairn: extension not installed; run the install script with --with-extension to install it
```

### `cairn update --extension` (update available)
```
cairn: extension current v0.1.0, latest v0.2.0
cairn: downloading extension v0.2.0...
cairn: verifying checksum...
cairn: extension updated to v0.2.0
```

---

## Standard Error

Errors are written to stderr with a consistent prefix.

```
cairn update: failed to reach release server: <detail>
cairn update: checksum mismatch for cairn-linux-amd64
cairn update: permission denied: cannot write to /usr/local/bin/cairn
```

---

## Help Text

```
Usage: cairn update [--check] [--extension]

Update cairn to the latest available release.

Flags:
  --check       Check for updates without applying them
  --extension   Update the Vicinae extension instead of the CLI binary

Exit codes:
  0  Success
  1  Error (network, API, or unknown platform)
  3  Checksum verification failed
  4  Permission denied
```

---

## Constraints

- The command MUST NOT make any network requests during normal cairn subcommands (`add`, `list`, `search`, etc.). Network activity is strictly limited to explicit `cairn update` invocations.
- On Windows, the command MUST print a message directing the user to re-run the install script, then exit 0 without attempting a binary replacement.
- The existing binary MUST be backed up before replacement and MUST be restored if replacement fails.
- The command MUST run without requiring elevated privileges; if the install directory is not writable by the current user, exit code 4 is returned with a descriptive error.
