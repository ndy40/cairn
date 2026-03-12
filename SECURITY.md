# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in Cairn, please report it responsibly by emailing security@example.com (replace with your actual contact) instead of using the issue tracker. Please include:

1. **Description** — What is the vulnerability?
2. **Impact** — What could an attacker do?
3. **Steps to reproduce** — How can it be triggered?
4. **Affected versions** — Which versions are vulnerable?

We will acknowledge your report within 48 hours and work to release a patch promptly.

## Known Security Considerations

### Plaintext Token Storage

OAuth2 tokens are stored in plaintext JSON files on disk. This is a known limitation:

- **Where**: `~/.config/cairn/sync.json` (Linux), `~/Library/Application Support/cairn/sync.json` (macOS), `%APPDATA%/cairn/sync.json` (Windows)
- **Risk**: If your system is compromised, an attacker can read the tokens
- **Mitigation**: Use disk encryption (e.g., LUKS, BitLocker, FileVault) and keep your system secure
- **Future**: Consider integrating with OS credential stores (Keychain, Secret Service, Credential Manager) in future versions

### Environment Variables

Do not store the Dropbox app key in environment variables in production:

- **Risk**: `ps aux`, process monitors, CI/CD logs, and `.env` files can expose the key
- **Safe approach**: Store in `cairn.json` with restrictive file permissions
- **Recommended**: Use your system's credential management system

### File Permissions

Config directories are created with `0700` (owner-only) permissions by default:

- **Linux/macOS**: `~/.config/cairn/` and sync config are owner-only
- **Windows**: Inherits ACLs from parent directory
- **Best practice**: Verify permissions with `ls -ld ~/.config/cairn/` — should show `drwx------`

### No Encryption of Stored Data

- Bookmarks are stored in plaintext SQLite (`bookmarks.db`)
- Sync snapshots are stored in plaintext JSON in Dropbox
- **Mitigation**: Use encrypted storage or cloud storage with client-side encryption

## Security Best Practices for Users

1. **Enable disk encryption** — Use LUKS (Linux), BitLocker (Windows), or FileVault (macOS)
2. **Restrict file permissions** — Ensure `~/.config/cairn/` is only readable by your user
3. **Don't share config files** — Never commit or share `cairn.json` or `sync.json`
4. **Unlink before losing access** — Run `cairn sync unlink` before selling/donating your machine
5. **Keep Cairn updated** — Regularly update to the latest version for security patches
6. **Monitor Dropbox activity** — Check your Dropbox account for unauthorized access

## Supported Versions

Security patches will be provided for:

- The latest stable release
- The previous major version (for 3 months after release)

## Disclosure Timeline

We follow responsible disclosure:

1. **Report received** — Acknowledged within 48 hours
2. **Assessment** — We determine severity and impact
3. **Patch development** — We work on a fix
4. **Release** — Patch released; vulnerability disclosed publicly (CVE if applicable)
5. **Follow-up** — Users are encouraged to upgrade

## Security Checklist for Contributors

When submitting code, please consider:

- [ ] No hardcoded credentials or secrets
- [ ] Proper file permissions for sensitive files
- [ ] Input validation and sanitization
- [ ] Use of secure cryptographic functions (from `crypto/*` stdlib)
- [ ] No use of deprecated or insecure libraries
- [ ] Tests for security-sensitive code paths

## Contact

For security-related questions (non-urgent), open a GitHub Discussion.
For vulnerability reports, email security@example.com (replace with your contact).
