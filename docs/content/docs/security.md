---
title: "Security"
weight: 70
---

# Security

## Credential Storage

### Dropbox App Key

- **Store in `cairn.json`** (config directory), not in environment variables
- Environment variables are visible in process listings (`ps aux`), logs, and CI/CD output
- The app key is a shared secret for the entire Cairn application; keep it secure

```json
{
  "dropbox_app_key": "your-app-key"
}
```

Ensure restrictive permissions:

```sh
chmod 600 ~/.config/cairn/cairn.json
```

### OAuth2 Tokens (`sync.json`)

- Stored in plaintext in the OS config directory
- Created with `0600` permissions (owner-only) automatically
- **If `sync.json` is compromised**, an attacker can read and write your Dropbox bookmarks

| OS | Path |
|----|------|
| Linux | `~/.config/cairn/sync.json` |
| macOS | `~/Library/Application Support/cairn/sync.json` |
| Windows | `%APPDATA%\cairn\sync.json` |

Run `cairn sync unlink` to revoke access before selling or donating your machine.

### Bookmark Database

- Bookmarks are stored in plaintext SQLite (`bookmarks.db`)
- No encryption at rest
- **Mitigation**: use full-disk encryption (LUKS, FileVault, BitLocker)

## Best Practices

1. **Enable disk encryption** — LUKS (Linux), BitLocker (Windows), or FileVault (macOS)
2. **Restrict config file permissions** — `chmod 600 ~/.config/cairn/cairn.json`
3. **Never commit config files** — add `cairn.json` and `sync.json` to `.gitignore`
4. **Prefer config file over env vars** for the Dropbox app key in long-running sessions
5. **Unlink before losing access** — `cairn sync unlink` before selling or wiping your machine
6. **Monitor Dropbox activity** — check your Dropbox account for unexpected access
7. **Keep Cairn updated** — run the installer again to get the latest release

## Reporting Vulnerabilities

If you discover a security vulnerability, please report it responsibly by opening a [GitHub Security Advisory](https://github.com/ndy40/cairn/security/advisories/new) instead of a public issue. Include:

1. **Description** — what is the vulnerability?
2. **Impact** — what could an attacker do?
3. **Steps to reproduce** — how can it be triggered?
4. **Affected versions** — which versions are vulnerable?

We will acknowledge reports within 48 hours and work to release a patch promptly.

## Known Limitations

| Limitation | Mitigation |
|------------|------------|
| OAuth2 tokens stored in plaintext | Use disk encryption; restrict file permissions |
| No keychain/credential manager integration | Planned for a future release |
| Bookmarks stored in plaintext SQLite | Use full-disk encryption |
| Sync snapshots stored in plaintext Dropbox JSON | Use Dropbox with client-side encryption |

## Contributor Security Checklist

When submitting a pull request, please verify:

- No hardcoded credentials or secrets
- Proper file permissions for sensitive files (`0600` for tokens, `0700` for config dirs)
- Input validation at system boundaries (user input, URLs)
- Use of `crypto/*` stdlib for cryptographic operations
- No deprecated or insecure libraries
- Tests for security-sensitive code paths
