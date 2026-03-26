E# Feature Specification: Interactive Setup Configuration Prompts

**Feature Branch**: `012-interactive-setup`
**Created**: 2026-03-26
**Status**: Draft
**Input**: User description: "During setup, we should start prompting the user to key in their DROPBOX_APP_KEY if the environment variable CAIRN_DROPBOX_APP_KEY is not set. After the user keys this in, we can then create the config file in the right location. Also optionally ask for database path. If not set, we maintain the default location."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Guided Dropbox Setup Without Pre-Set Environment Variable (Priority: P1)

A new user runs `cairn sync setup` without having set `CAIRN_DROPBOX_APP_KEY` in their environment. Instead of receiving a hard error and being told to set an environment variable manually, the CLI prompts them to enter their Dropbox App Key interactively. They type it in and setup proceeds. The app key is persisted to `cairn.json` so they do not need to enter it again on subsequent runs.

**Why this priority**: This is the primary user-facing improvement. Without it, users must know to set an environment variable before running setup — a hidden prerequisite that causes unhelpful error messages. Removing that friction makes the tool accessible to all users.

**Independent Test**: Run `cairn sync setup` in an environment where `CAIRN_DROPBOX_APP_KEY` is unset and no `cairn.json` exists. Verify that the CLI prompts for the App Key, accepts typed input, and writes it to `cairn.json` in the OS-appropriate config directory before proceeding with OAuth.

**Acceptance Scenarios**:

1. **Given** `CAIRN_DROPBOX_APP_KEY` is not set and no `cairn.json` contains `dropbox_app_key`, **When** the user runs `cairn sync setup`, **Then** the CLI displays a prompt asking the user to enter their Dropbox App Key.
2. **Given** the CLI is prompting for the Dropbox App Key, **When** the user types a non-empty value and presses Enter, **Then** the value is saved to `cairn.json` and setup continues to the OAuth authentication step.
3. **Given** the CLI is prompting for the Dropbox App Key, **When** the user submits an empty value, **Then** the CLI shows an error message ("App Key cannot be empty") and re-prompts.
4. **Given** `CAIRN_DROPBOX_APP_KEY` is already set in the environment, **When** the user runs `cairn sync setup`, **Then** the CLI does not prompt for the App Key and proceeds directly.
5. **Given** `cairn.json` already contains a non-empty `dropbox_app_key`, **When** the user runs `cairn sync setup`, **Then** the CLI does not prompt for the App Key and proceeds directly.

---

### User Story 2 - Optional Custom Database Path During Setup (Priority: P2)

During `cairn sync setup`, after the App Key is resolved, the CLI optionally prompts the user to specify a custom database path. If the user presses Enter without entering a path, the default OS-appropriate location is used. If the user enters a custom path, it is persisted to `cairn.json`.

**Why this priority**: Most users will use the default database path, making this optional. However, power users who manage multiple cairn instances or use non-standard paths benefit from being guided through this during initial setup rather than discovering they need to set it manually.

**Independent Test**: Run `cairn sync setup` with no config, enter a Dropbox App Key when prompted, then press Enter at the database path prompt without typing anything. Verify setup completes and the default database path is used.

**Acceptance Scenarios**:

1. **Given** the setup is prompting for the database path, **When** the user presses Enter without entering a value, **Then** the default OS-appropriate database path is used and no `db_path` entry is written to `cairn.json`.
2. **Given** the setup is prompting for the database path, **When** the user types a custom path and presses Enter, **Then** the entered path is saved to `cairn.json` under `db_path` and subsequent operations use that path.
3. **Given** `CAIRN_DB_PATH` is already set in the environment, **When** setup prompts are shown, **Then** the database path prompt is skipped (the env var takes precedence and the user is informed of the current value).
4. **Given** `cairn.json` already contains a `db_path`, **When** setup prompts are shown, **Then** the database path prompt is skipped (existing config is respected).

---

### User Story 3 - Config File Written to OS-Appropriate Location (Priority: P3)

After the user provides their App Key (and optionally a custom database path), the CLI creates or updates `cairn.json` in the OS-appropriate config directory. The user is shown the path where the config was written so they know where to find or edit it later.

**Why this priority**: Users need confidence that the config was persisted and know where to find it. This completes the setup feedback loop.

**Independent Test**: After completing a setup flow where a new App Key and custom db path were entered, verify that `cairn.json` in the OS config directory contains the expected keys and values, and that the CLI printed the config file path.

**Acceptance Scenarios**:

1. **Given** the user has provided a Dropbox App Key during setup, **When** setup writes the config, **Then** `cairn.json` is created (or updated) in the OS-appropriate config directory containing `dropbox_app_key`.
2. **Given** setup successfully writes or updates `cairn.json`, **When** the write completes, **Then** the CLI prints a confirmation message including the full path to `cairn.json`.
3. **Given** the OS config directory does not exist, **When** the CLI attempts to write `cairn.json`, **Then** the directory is created automatically and the file is written successfully.
4. **Given** an existing `cairn.json` already has other settings, **When** setup writes the new app key, **Then** existing settings in `cairn.json` are preserved (only the new/updated keys are changed).

---

### Edge Cases

- What happens when the user enters a Dropbox App Key containing only whitespace? The CLI trims whitespace and re-prompts if the trimmed value is empty.
- What happens when `cairn.json` exists but is read-only? The CLI reports a clear error explaining it cannot write to the config file and shows the path.
- What happens when the config directory path cannot be determined (e.g., missing home directory)? The CLI reports an error and instructs the user to set `CAIRN_DROPBOX_APP_KEY` manually.
- What happens when the user runs `cairn sync setup` and sync is already configured? The existing early-return behaviour is preserved; no prompts are shown.
- What happens if the user interrupts the prompt (Ctrl+C)? The CLI exits cleanly without writing a partial config.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: During `cairn sync setup`, the system MUST detect whether a Dropbox App Key is available (via environment variable `CAIRN_DROPBOX_APP_KEY` or existing `cairn.json` `dropbox_app_key`).
- **FR-002**: When no Dropbox App Key is available, the system MUST interactively prompt the user to enter it before proceeding with OAuth authentication.
- **FR-003**: The system MUST reject empty or whitespace-only App Key input and re-prompt the user with a clear error message.
- **FR-004**: After the user provides a valid App Key, the system MUST persist it to `cairn.json` in the OS-appropriate config directory before proceeding.
- **FR-005**: During `cairn sync setup`, the system MUST offer an optional prompt for a custom database path, clearly indicating that pressing Enter keeps the default location.
- **FR-006**: If the user enters a custom database path, the system MUST persist it to `cairn.json` under `db_path`.
- **FR-007**: If the user accepts the default database path, the system MUST NOT write a `db_path` entry to `cairn.json` (relying on built-in defaults).
- **FR-008**: Database path and App Key prompts MUST be skipped when their values are already resolved from higher-precedence sources (environment variables or existing config file).
- **FR-009**: After writing to `cairn.json`, the system MUST display the full path of the config file to the user.
- **FR-010**: When writing to `cairn.json`, the system MUST preserve any existing keys not modified during setup.
- **FR-011**: The system MUST create the OS config directory if it does not already exist when writing `cairn.json`.

### Key Entities

- **App Key**: The Dropbox application identifier required to authenticate OAuth flows. Sourced from environment, config file, or interactive prompt during setup.
- **Config File (`cairn.json`)**: A JSON file in the OS-appropriate config directory holding persistent application settings (`dropbox_app_key`, `db_path`). Created or updated by the setup command when interactive prompts are completed.
- **Database Path**: The file path to the SQLite database. Optional — defaults to the OS-appropriate data directory path if not specified.

## Assumptions

- The interactive prompts apply only to the `cairn sync setup` command. Other sync subcommands (`push`, `pull`, `status`, `auth`, `unlink`) continue to error if the App Key is missing.
- Input is read from stdin; the setup command is not intended to run non-interactively (scripts should set environment variables or pre-populate `cairn.json`).
- The config precedence order established in feature 010 is unchanged: CLI flag > environment variable > `cairn.json` > default. Interactive input during setup writes to `cairn.json` (not to the environment).
- The database path prompt is shown only when `CAIRN_DB_PATH` is unset and no `db_path` exists in an existing `cairn.json`.
- Input masking (hidden characters) for the App Key is not required as Dropbox App Keys are not secret credentials in the same sense as passwords — they are safe to display.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user with no prior cairn configuration can complete `cairn sync setup` successfully by responding only to interactive prompts, without setting any environment variables beforehand.
- **SC-002**: After completing the interactive setup flow, `cairn.json` exists in the correct OS config directory and contains the entered App Key value.
- **SC-003**: Users who already have `CAIRN_DROPBOX_APP_KEY` set in their environment experience no change in behaviour — setup proceeds without any new prompts.
- **SC-004**: The setup command displays the path of the written config file, so users can locate it without guessing.
- **SC-005**: An invalid (empty) App Key entry is rejected immediately with a clear message, and the user is re-prompted without restarting the command.
