# Changelog

All notable changes to Cairn are automatically generated from git commits and published in [GitHub Releases](https://github.com/ndy40/cairn/releases).

This changelog is maintained by [Cocogitto](https://cocogitto.io) based on [Conventional Commits](https://www.conventionalcommits.org/).

## Commit Message Format

To have your changes appear in release notes, use **conventional commits**:

```
<type>(<scope>): <description>

<body>

<footer>
```

### Types

- **feat**: A new feature (bumps minor version)
- **fix**: A bug fix (bumps patch version)
- **perf**: Performance improvement (bumps patch version)
- **refactor**: Code refactoring (no version bump)
- **docs**: Documentation changes (no version bump)
- **doc**: Documentation changes (no version bump)
- **test**: Adding or updating tests (no version bump)
- **tests**: Adding or updating tests (no version bump)
- **ci**: CI/CD changes (no version bump)
- **chore**: Build process, dependencies, etc. (no version bump)
- **lint**: Linting and code style (no version bump)

### Examples

```bash
# New feature
git commit -m "feat(sync): Add auto-sync on startup"

# Bug fix with scope
git commit -m "fix(tui): Resolve bookmark list rendering issue"

# Performance improvement
git commit -m "perf(search): Optimize FTS5 queries for large collections"

# Breaking change (bumps major version)
git commit -m "feat!: Change bookmark ID from int to UUID"

# With detailed description
git commit -m "feat(extension): Add Vicinae bookmark integration

- Search bookmarks from launcher
- Save current page as bookmark
- Browse all bookmarks in Vicinae"
```

## Release Process

1. Merge code to `main` with conventional commits
2. GitHub Actions workflow (`bump.yml`) runs automatically:
   - Analyzes commits since the last release
   - Determines version bump (major, minor, or patch)
   - Creates and pushes a new `v*` tag
3. The `release.yml` workflow triggers on the new tag:
   - Runs tests and linting
   - Builds binaries for all platforms
   - Generates release notes using Cocogitto
   - Creates a GitHub Release with binaries + checksums
4. This changelog is updated automatically via Cocogitto

## View All Releases

See all releases and auto-generated notes: https://github.com/ndy40/cairn/releases
