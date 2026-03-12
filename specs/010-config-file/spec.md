# Feature Specification: Configuration File Support

**Feature Branch**: `010-config-file`
**Created**: 2026-03-12
**Status**: Draft
**Input**: User description: "Add optional cairn.json configuration file with precedence over environment variables and defaults. Update README with config file location and env variable documentation."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Load Configuration from File (Priority: P1)

As a user, I want cairn to read settings from an optional `cairn.json` file so that I can configure the application without setting environment variables every time.

The configuration file is located in the OS-appropriate config directory:
- **Linux**: `$XDG_CONFIG_HOME/cairn/cairn.json` (defaults to `~/.config/cairn/cairn.json`)
- **macOS**: `~/Library/Application Support/cairn/cairn.json`
- **Windows**: `%APPDATA%/cairn/cairn.json`

**Why this priority**: This is the core feature — without config file loading, nothing else matters.

**Independent Test**: Run `cairn config` with a `cairn.json` file present containing a custom `db_path` value (and no `CAIRN_DB_PATH` env var set) and verify the output reflects the file's value instead of the default.

**Acceptance Scenarios**:

1. **Given** a `cairn.json` file exists with `{"db_path": "/tmp/test.db"}` and no `CAIRN_DB_PATH` env var is set, **When** the user runs `cairn config`, **Then** the output shows `CAIRN_DB_PATH=/tmp/test.db`.
2. **Given** no `cairn.json` file exists, **When** the user runs any cairn command, **Then** the application uses defaults as it does today (no error, no warning).
3. **Given** a `cairn.json` file exists with invalid JSON, **When** the user runs any cairn command, **Then** the application prints a clear error message indicating the config file has a syntax error and exits.

---

### User Story 2 - Configuration Precedence (Priority: P2)

As a user, I want the configuration precedence to be clearly defined so that I can predictably override settings at different levels.

The precedence order (highest to lowest):
1. **Environment variables** (`CAIRN_*`) — highest priority
2. **CLI flags** (e.g., `--db`)
3. **Config file** (`cairn.json`)
4. **Defaults** (OS-appropriate paths built into the code)

**Why this priority**: Users need predictable behavior when multiple configuration sources exist. Environment variables take top priority to support deployment/CI overrides without modifying files or commands.

**Independent Test**: Set `CAIRN_DB_PATH` env var to one path, `cairn.json` `db_path` to another, and pass `--db` flag with a third — verify the env var wins.

**Acceptance Scenarios**:

1. **Given** `CAIRN_DB_PATH=/env/path.db` is set AND `cairn.json` contains `{"db_path": "/file/path.db"}` AND the user runs `cairn --db /flag/path.db config`, **When** the config is resolved, **Then** the output shows `CAIRN_DB_PATH=/env/path.db` (env var wins over CLI flag and config file).
2. **Given** no `CAIRN_DB_PATH` env var is set AND `cairn.json` contains `{"db_path": "/file/path.db"}` AND the user runs `cairn --db /flag/path.db config`, **When** the config is resolved, **Then** the output shows `CAIRN_DB_PATH=/flag/path.db` (CLI flag wins over config file).
3. **Given** no `CAIRN_DB_PATH` env var is set AND no `--db` flag AND `cairn.json` contains `{"db_path": "/file/path.db"}`, **When** the user runs `cairn config`, **Then** the output shows `CAIRN_DB_PATH=/file/path.db` (config file wins over default).
4. **Given** no env var, no CLI flag, and no config file, **When** the user runs `cairn config`, **Then** the output shows the OS-appropriate default path.

---

### User Story 3 - README Documentation (Priority: P3)

As a user, I want the README to document the config file location and supported environment variables so that I can set up cairn without reading the source code.

**Why this priority**: Documentation is important but the feature must work first.

**Independent Test**: Read the README and verify it contains a section explaining config file location, supported settings, and environment variables.

**Acceptance Scenarios**:

1. **Given** the README exists, **When** a user reads it, **Then** they find a "Configuration" section that lists the config file path for each OS, the supported JSON keys, the supported environment variables, and the precedence order.

---

### Edge Cases

- What happens when the config file contains unknown keys? The application silently ignores them (forward compatibility).
- What happens when the config file has correct JSON but wrong value types (e.g., `db_path` is a number)? The application prints a clear error and exits.
- What happens when the config file exists but is empty (`{}`)? The application treats it as no overrides and falls back to env vars / defaults.
- What happens when the config file directory exists but the file does not? No error — the file is optional.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST look for `cairn.json` in the OS-appropriate config directory on startup.
- **FR-002**: System MUST NOT require the config file to exist — it is optional.
- **FR-003**: System MUST apply configuration precedence: environment variables > CLI flags > config file > defaults.
- **FR-004**: System MUST support the following config keys in `cairn.json`: `db_path` (string), `dropbox_app_key` (string).
- **FR-005**: System MUST report a clear error and exit if the config file exists but contains invalid JSON or wrong value types.
- **FR-006**: System MUST silently ignore unknown keys in the config file.
- **FR-007**: The `cairn config` command MUST display the resolved configuration values reflecting the correct precedence.
- **FR-008**: The README MUST document config file location (per OS), supported JSON keys, supported environment variables, and precedence order.

### Key Entities

- **AppConfig**: Represents the resolved configuration after merging all sources. Contains `db_path` and `dropbox_app_key`.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can configure cairn via a JSON file without setting environment variables.
- **SC-002**: Configuration precedence behaves predictably — environment variables always win, then CLI flags, then config file, then defaults.
- **SC-003**: The application starts without error when no config file exists.
- **SC-004**: Users can find all configuration options documented in the README within 30 seconds.

## Clarifications

### Session 2026-03-12

- Q: What is the configuration precedence order? → A: defaults < config file < CLI flags < env vars (environment variables highest priority).

## Assumptions

- The config file name is `cairn.json` (not `config.json` as mentioned in the user's initial description — using `cairn.json` for consistency with the project name and to avoid ambiguity).
- The config directory follows the same OS conventions already used by the sync config (`$XDG_CONFIG_HOME/cairn/` on Linux, `~/Library/Application Support/cairn/` on macOS).
- Only two config keys are needed initially (`db_path`, `dropbox_app_key`), matching the current environment variables. More can be added later.
- Precedence clarified by user: defaults < config file < CLI flags < env vars. Environment variables have the highest priority to support deployment/CI override scenarios.
