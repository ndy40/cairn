# CLI Contract: Interactive Setup Configuration Prompts

## Command: `cairn sync setup`

This document describes the updated interactive behaviour of `cairn sync setup`.

---

## Invocation

```
cairn sync setup
```

No new flags are introduced. Existing flags (`--db`) continue to work as before.

---

## Prompt Flow (new behaviour)

### Step 1 — Dropbox App Key

Shown **only when** `CAIRN_DROPBOX_APP_KEY` is unset AND `cairn.json` does not contain a non-empty `dropbox_app_key`.

```
Enter your Dropbox App Key: _
```

- Input is read from stdin (characters visible, no masking required).
- Whitespace is trimmed from both ends.
- Empty or whitespace-only input prints an error and re-prompts:
  ```
  Error: App Key cannot be empty. Please try again.
  Enter your Dropbox App Key: _
  ```
- On valid input, the value is persisted to `cairn.json` before OAuth begins.

### Step 2 — Database Path (optional)

Shown **only when** `CAIRN_DB_PATH` is unset AND no explicit `db_path` is present in `cairn.json` (i.e., resolved path equals the built-in default).

```
Enter database path (press Enter for default: /home/alice/.local/share/bookmark-manager/bookmarks.db): _
```

- The default path shown is the OS-appropriate value for the current user.
- Empty input → use default; do NOT write `db_path` to `cairn.json`.
- Non-empty input → write entered value to `cairn.json` as `db_path`.

### Step 3 — Confirmation

After any write to `cairn.json`:

```
Config written to /home/alice/.config/cairn/cairn.json
```

---

## Post-Prompt Behaviour

After prompts are resolved (or skipped), execution continues to the existing OAuth authentication flow unchanged.

---

## Error Conditions

| Condition | Output | Exit Code |
|-----------|--------|-----------|
| Cannot determine config directory (no home dir) | `Error: cannot determine config file path` | 3 |
| Cannot write `cairn.json` (e.g., read-only) | `Error: cannot write config file <path>: <reason>` | 3 |
| User sends EOF/Ctrl-D at prompt | Exit cleanly, no partial config written | 0 |

---

## Unchanged Behaviour

- When `CAIRN_DROPBOX_APP_KEY` is set in the environment, no prompts appear.
- When sync is already configured (`cairn.json` has a valid access token), the existing early-return message is shown and no prompts appear.
- All other `cairn sync` subcommands (`push`, `pull`, `status`, `auth`, `unlink`) are unaffected.
