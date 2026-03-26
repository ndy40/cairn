# Research: Interactive Setup Configuration Prompts

## Interactive Prompting in Go CLI

**Decision**: Use `bufio.NewReader(os.Stdin)` with `ReadString('\n')` for multi-line-safe input reading.

**Rationale**: The codebase already uses `fmt.Scanln` for single-word answers (e.g., `checkFirstRunSync` in `cmd/cairn/main.go:577`). For App Key input (which can contain arbitrary characters) and database paths (which can contain spaces), `bufio.Reader.ReadString('\n')` is more robust — it reads until newline regardless of whitespace within the value.

**Alternatives considered**:
- `fmt.Scanln` — stops at whitespace; fine for y/N but not for paths or keys with special characters
- Third-party prompt libraries (e.g., `charmbracelet/huh`) — not currently in go.mod; adding a dependency purely for two prompts violates the "no unnecessary abstractions" principle
- `os.Stdin.Read` directly — too low-level

---

## Config File Write Path

**Decision**: Reuse the existing `config.Manager` methods `WriteConfig()` / `SaveConfig()` and `DefaultConfigPath()`.

**Rationale**: `internal/config/config.go` already provides:
- `DefaultConfigPath()` — returns the OS-appropriate `cairn.json` path
- `SaveConfig()` — calls `viper.SafeWriteConfig()`, which creates the file if it doesn't exist
- `WriteConfig()` — overwrites existing file
- `MkdirAll(dir, 0700)` — already used inside both methods to create directories

The setup command must preserve existing keys (e.g., `db_path` must not be lost when writing `dropbox_app_key`). `viper` holds all loaded settings in memory; calling `v.Set("dropbox_app_key", value)` then `WriteConfig()` serialises all currently-loaded keys, naturally preserving them.

**Approach**: After the user enters the App Key prompt, call `cfgManager.Set("dropbox_app_key", appKey)`, then `cfgManager.WriteConfig()` (or `SaveConfig()` when no file exists). This writes all settings — including previously loaded db_path if any — to disk.

---

## Database Path Prompt Conditions

**Decision**: Show the db path prompt only when BOTH of these are true:
1. `CAIRN_DB_PATH` environment variable is not set
2. The config manager's resolved `db_path` equals `config.DefaultDBPath()` (i.e., no explicit override exists in cairn.json)

**Rationale**: The spec says the prompt is skipped when a higher-precedence source already provides the value. Since `config.Manager` resolves precedence, comparing the resolved value to the default cleanly captures "no explicit source set it".

---

## Prompt Placement: `runSyncSetup` in `cmd/cairn/main.go`

**Decision**: Add a `promptForSetupConfig(cfgManager *config.Manager)` helper called at the top of `runSyncSetup`, before the `appKey == ""` guard.

**Rationale**: This keeps all setup-specific prompt logic in a single function and makes `runSyncSetup` trivially testable — callers can pre-populate cfgManager with test values to bypass prompts. No changes to other sync subcommands are needed.

---

## No New Dependencies Required

**Decision**: All changes use `bufio`, `fmt`, `os`, and `strings` from the Go standard library only.

**Rationale**: The feature requires two sequential stdin prompts and one config file write — all achievable with stdlib. Adding a TUI prompt library would violate the constitution's "no unnecessary abstractions" principle.

---

## Constitution Check Results

| Gate | Status | Notes |
|------|--------|-------|
| No CGO | PASS | Only stdlib + existing pure-Go dependencies |
| Single binary | PASS | No new runtime or process |
| Task management | PASS | Tasks will be created via Backlog CLI after `/speckit.tasks` |
| Backward-compatible migrations | PASS | No schema changes |
