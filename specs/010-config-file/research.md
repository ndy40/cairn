# Research: Configuration File Support

**Feature**: 010-config-file
**Date**: 2026-03-12

## Config File Location Convention

- **Decision**: Use OS-appropriate config directories, matching the existing sync config pattern in `internal/sync/config.go`.
- **Rationale**: The project already resolves config paths per-OS (XDG on Linux, Library/Application Support on macOS, APPDATA on Windows). Reusing this pattern keeps behavior consistent.
- **Alternatives considered**: Home directory dotfile (`~/.cairn.json`) — rejected because it doesn't follow XDG and the project already uses XDG.

## Config File Format

- **Decision**: JSON (`cairn.json`).
- **Rationale**: The project already uses JSON for sync config. Go stdlib `encoding/json` requires no new dependencies. JSON is familiar to users.
- **Alternatives considered**: TOML (requires external dependency, violates "no unnecessary dependencies"), YAML (same issue), env file (doesn't support structured config).

## Precedence Order

- **Decision**: env vars > CLI flags > config file > defaults.
- **Rationale**: User explicitly specified this order. Environment variables at top priority supports CI/deployment overrides.
- **Alternatives considered**: CLI flags > config file > env vars > defaults (more conventional but user explicitly chose env vars first).

## Config Package Placement

- **Decision**: `internal/config/config.go` — new package.
- **Rationale**: Separates config concerns from store and sync. Single file, single responsibility. Does not bloat existing packages.
- **Alternatives considered**: Adding to `cmd/cairn/main.go` directly — rejected because config resolution logic would be untestable and mixed with CLI wiring.

## Unknown Keys Handling

- **Decision**: Silently ignore unknown keys using standard `json.Unmarshal` behavior (only populates known struct fields).
- **Rationale**: Forward compatibility — users can add keys for future versions without breaking current version.
- **Alternatives considered**: Strict validation with error on unknown keys — rejected per spec requirement FR-006.
