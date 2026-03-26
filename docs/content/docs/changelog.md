---
title: "Changelog"
weight: 90
---

# Changelog

All notable changes to Cairn are automatically generated from [conventional commits](https://www.conventionalcommits.org/) and published as [GitHub Releases](https://github.com/ndy40/cairn/releases).

## Release Process

When a version tag is pushed (`git push origin v0.2.0`), GitHub Actions automatically:

1. Runs tests and linting
2. Builds binaries for all platforms (Linux, macOS, Windows — amd64 and arm64)
3. Generates release notes from commits since the last tag
4. Creates a GitHub Release with binaries and SHA256 checksums

## Commit Message Format

```
<type>(<scope>): <description>

<body>

<footer>
```

### Types

| Type | Appears in Release Notes | Description |
|------|--------------------------|-------------|
| `feat` | Yes | New feature |
| `fix` | Yes | Bug fix |
| `perf` | Yes | Performance improvement |
| `docs` | No | Documentation changes |
| `refactor` | No | Code refactoring |
| `test` | No | Tests |
| `ci` | No | CI/CD changes |
| `chore` | No | Build, dependencies |

### Examples

```sh
# New feature
git commit -m "feat(sync): Add auto-sync on startup"

# Bug fix with scope
git commit -m "fix(tui): Resolve bookmark list rendering issue"

# Breaking change (! suffix)
git commit -m "feat!: Change bookmark ID format"

# With body
git commit -m "feat(extension): Add Vicinae bookmark integration

- Search bookmarks from launcher
- Save current page as bookmark
- Browse all bookmarks in Vicinae"
```

## Full Changelog

See the [GitHub Releases page](https://github.com/ndy40/cairn/releases) for the complete, auto-generated changelog.
