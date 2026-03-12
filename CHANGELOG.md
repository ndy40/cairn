# Changelog

All notable changes to Cairn are automatically generated from git commits and published as [GitHub Releases](https://github.com/ndy40/cairn/releases).

## Commit Message Format

To have your changes appear in release notes, use **conventional commits**:

```
<type>(<scope>): <description>

<body>

<footer>
```

### Types

- **feat**: A new feature (appears in release notes)
- **fix**: A bug fix (appears in release notes)
- **docs**: Documentation changes
- **perf**: Performance improvement
- **refactor**: Code refactoring
- **test**: Adding or updating tests
- **ci**: CI/CD changes
- **chore**: Build process, dependencies, etc.

### Examples

```bash
# New feature
git commit -m "feat(sync): Add auto-sync on startup"

# Bug fix with scope
git commit -m "fix(tui): Resolve bookmark list rendering issue"

# Performance improvement
git commit -m "perf(search): Optimize FTS5 queries for large collections"

# Breaking change (requires ! after type)
git commit -m "feat!: Change bookmark ID from int to UUID"

# With detailed description
git commit -m "feat(extension): Add Vicinae bookmark integration

- Search bookmarks from launcher
- Save current page as bookmark
- Browse all bookmarks in Vicinae"
```

## Release Process

1. When you push a tag (`git push origin v0.2.0`), GitHub Actions automatically:
   - Runs tests and linting
   - Builds binaries for all platforms
   - **Generates release notes** from commits since the last tag
   - Creates a GitHub Release with binaries + checksums

2. Release notes are auto-generated using [conventional-changelog](https://github.com/conventional-changelog/conventional-changelog)

3. No manual CHANGELOG.md editing needed — commits are the source of truth

## View Releases

See all releases and auto-generated notes here: https://github.com/ndy40/cairn/releases
