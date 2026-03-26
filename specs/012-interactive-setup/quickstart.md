# Quickstart: Interactive Setup Configuration Prompts

## What's Changing

`cairn sync setup` now guides users through configuration interactively. Instead of failing with an error when `CAIRN_DROPBOX_APP_KEY` is not set, the command prompts for the missing value and writes it to `cairn.json` automatically.

---

## Affected Files

| File | Change Type | Description |
|------|------------|-------------|
| `cmd/cairn/main.go` | Modify | Add `promptForSetupConfig(cfgManager)` call inside `runSyncSetup`; add `promptForSetupConfig` helper function |
| `internal/config/config.go` | No change | Existing `Manager.Set`, `WriteConfig`, `SaveConfig`, `DefaultConfigPath`, `DefaultDBPath` are reused as-is |

---

## Implementation Guide

### 1. Add helper function in `cmd/cairn/main.go`

Add a new function `promptForSetupConfig` that:

1. Checks if `cfgManager.Get().DropboxAppKey` is empty.
2. If empty, loops prompting stdin until a non-empty, trimmed value is entered.
3. Calls `cfgManager.Set("dropbox_app_key", value)`.
4. Checks if `CAIRN_DB_PATH` env var is unset AND `cfgManager.Get().DBPath == config.DefaultDBPath()`.
5. If both true, prompts once for a database path (empty = accept default, no write needed).
6. If non-empty path entered, calls `cfgManager.Set("db_path", path)`.
7. If any `Set` was called, calls `cfgManager.WriteConfig()` (creates dir + file) and prints the path.

```go
// Pseudocode only — implementation fills in real details
func promptForSetupConfig(cfgManager *config.Manager) {
    appKey := cfgManager.Get().DropboxAppKey
    if appKey == "" {
        for {
            fmt.Print("Enter your Dropbox App Key: ")
            key, _ := bufio.NewReader(os.Stdin).ReadString('\n')
            key = strings.TrimSpace(key)
            if key != "" {
                cfgManager.Set("dropbox_app_key", key)
                break
            }
            fmt.Fprintln(os.Stderr, "Error: App Key cannot be empty. Please try again.")
        }
    }

    if os.Getenv("CAIRN_DB_PATH") == "" && cfgManager.Get().DBPath == config.DefaultDBPath() {
        defaultPath := config.DefaultDBPath()
        fmt.Printf("Enter database path (press Enter for default: %s): ", defaultPath)
        path, _ := bufio.NewReader(os.Stdin).ReadString('\n')
        path = strings.TrimSpace(path)
        if path != "" {
            cfgManager.Set("db_path", path)
        }
    }

    // Write to file if we changed anything
    if cfgManager.Get().DropboxAppKey != "" {
        if err := cfgManager.WriteConfig(); err != nil {
            fatalf(3, "cannot write config file: %v", err)
        }
        fmt.Printf("Config written to %s\n", config.DefaultConfigPath())
    }
}
```

### 2. Call helper at top of `runSyncSetup`

```go
func runSyncSetup(dbPath string, appCfg *config.AppConfig) {
    // NEW: prompt for missing config values
    promptForSetupConfig(cfgManager) // cfgManager must be passed in or accessible

    appKey := appCfg.DropboxAppKey
    // ... rest unchanged
}
```

> **Note**: `cfgManager` will need to be accessible in `runSyncSetup`. Either pass it as a parameter (preferred) or use the package-level approach already used in main.

---

## Testing

### Manual test (no env var set, no cairn.json)

```bash
unset CAIRN_DROPBOX_APP_KEY
rm -f ~/.config/cairn/cairn.json
cairn sync setup
# Expect: prompted for App Key, then db path
# After: cairn.json contains dropbox_app_key
```

### Manual test (env var already set)

```bash
export CAIRN_DROPBOX_APP_KEY=mykey
cairn sync setup
# Expect: no App Key prompt; proceeds to OAuth directly
```

### Manual test (empty App Key input)

```bash
unset CAIRN_DROPBOX_APP_KEY
cairn sync setup
# At prompt: press Enter
# Expect: "Error: App Key cannot be empty. Please try again." then re-prompt
```

---

## Key Constraints

- `bufio.NewReader(os.Stdin)` must be created fresh each time or shared carefully to avoid consuming characters meant for OAuth flows.
- The `WriteConfig` call in `config.Manager` uses `viper.WriteConfig` which serialises all currently-set keys — existing keys are automatically preserved.
- No new dependencies. No new packages.
