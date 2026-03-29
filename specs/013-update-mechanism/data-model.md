# Data Model: Cairn Self-Update Mechanism

**Feature**: 013-update-mechanism
**Date**: 2026-03-26

> No persistent storage changes. All entities are transient (in-memory during the update command's lifetime).

---

## Entities

### ReleaseInfo

Represents the latest published release fetched from the release source.

| Field          | Type   | Description                                      |
|----------------|--------|--------------------------------------------------|
| TagName        | string | Version string, e.g. `v0.3.0`                   |
| BinaryAssetURL | string | Download URL for the platform-specific binary    |
| ChecksumURL    | string | Download URL for the checksum manifest           |
| PublishedAt    | string | ISO 8601 timestamp of the release publication   |

**Derived at runtime** from the releases API. Not persisted.

---

### UpdateStatus

Represents the result of a version check for a component (CLI binary or extension).

| Field           | Type   | Description                                                      |
|-----------------|--------|------------------------------------------------------------------|
| Component       | string | `"cli"` or `"extension"`                                         |
| CurrentVersion  | string | Version of the installed component (`"dev"` for dev builds)      |
| LatestVersion   | string | Version reported by the release source                           |
| UpdateAvailable | bool   | True when latest differs from current and current is not already latest |
| Error           | error  | Non-nil if the check could not be completed                      |

---

### DownloadedArtifact

Transient representation of a downloaded binary or archive before installation.

| Field        | Type   | Description                                             |
|--------------|--------|---------------------------------------------------------|
| TempPath     | string | Absolute path to the temporary file on disk             |
| ExpectedHash | string | SHA256 hex string from the checksum manifest            |
| TargetPath   | string | Absolute path where the artifact will be installed      |
| BackupPath   | string | Absolute path of the `.bak` copy of the existing file   |

**Lifecycle**: created during download, discarded (temp file removed) after installation or on error.

---

## State Transitions

### CLI Binary Update Flow

```
[Start]
  │
  ▼
CheckLatestVersion
  │
  ├── error ──────────────────────────────► [Report error, exit 1, no changes]
  │
  ├── no update available ────────────────► [Report "already up to date", exit 0]
  │
  └── update available
        │
        ▼
      DownloadBinary → TempFile created
        │
        ├── error ──────────────────────► [Remove temp file, report error, exit 1]
        │
        ▼
      VerifyChecksum
        │
        ├── mismatch ───────────────────► [Remove temp file, report error, exit 1]
        │
        ▼
      BackupCurrentBinary → .bak created
        │
        ▼
      ReplaceBinary (os.Rename or copy+delete)
        │
        ├── error ──────────────────────► [Restore from .bak, report error, exit 1]
        │
        ▼
      [Report success, remove .bak, exit 0]
```

### Extension Update Flow

```
[Start]
  │
  ▼
DetectExtensionInstall
  │
  ├── not installed ──────────────────────► [Report "not installed", exit 0]
  │
  ▼
CheckExtensionVersion
  │
  ├── error ──────────────────────────────► [Report error, exit 1]
  ├── no update available ────────────────► [Report "already up to date", exit 0]
  │
  └── update available
        │
        ▼
      DownloadExtensionArchive → TempFile
        │
        ├── error ──────────────────────► [Remove temp, report error, exit 1]
        │
        ▼
      VerifyChecksum
        │
        ├── mismatch ───────────────────► [Remove temp, report error, exit 1]
        │
        ▼
      ExtractToExtensionDir
        │
        ├── error ──────────────────────► [Report error, exit 1]
        │
        ▼
      WriteVersionFile
        │
        ▼
      [Report success, exit 0]
```
