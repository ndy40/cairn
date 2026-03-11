# Contract: Install Script CLI Interface

**Feature**: 009-install-script
**Date**: 2026-03-11

## Invocation

```sh
# Standard interactive install (latest version)
curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh

# With options (download first, then run)
curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh -o install.sh
sh install.sh [OPTIONS]
```

## Flags

| Flag                  | Short | Default         | Description                                              |
|-----------------------|-------|-----------------|----------------------------------------------------------|
| `--install-dir DIR`   | `-d`  | `~/.local/bin`  | Custom installation directory for the cairn binary       |
| `--version VERSION`   | `-v`  | latest          | Install a specific version (e.g., `v0.0.1`)             |
| `--non-interactive`   | `-y`  | (interactive)   | Skip all prompts, install CLI only                       |
| `--with-extension`    |       | (off)           | Also install Vicinae extension (non-interactive mode)    |
| `--help`              | `-h`  |                 | Show usage information                                   |

## Environment Variables

| Variable              | Default | Description                                              |
|-----------------------|---------|----------------------------------------------------------|
| `CAIRN_INSTALL_DIR`   | (none)  | Override install directory (flag takes precedence)        |

## Exit Codes

| Code | Meaning                                              |
|------|------------------------------------------------------|
| 0    | Installation completed successfully                  |
| 1    | General error (network failure, unexpected error)    |
| 2    | Unsupported platform (OS or architecture)            |
| 3    | Checksum verification failed                         |
| 4    | Permission denied (cannot write to install directory)|

## Output Behavior

- **Interactive mode**: Progress messages, prompts for Vicinae extension, success/failure summary
- **Non-interactive mode**: Minimal output, errors to stderr, no prompts
- **Piped mode** (detected via `! [ -t 0 ]`): Behaves like non-interactive when stdin is not a terminal, unless user explicitly passes flags
