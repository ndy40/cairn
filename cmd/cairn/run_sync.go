package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ndy40/cairn/internal/config"
	csync "github.com/ndy40/cairn/internal/sync"
)

const syncLockEnv = "CAIRN_SYNC_LOCK"
const syncLockTTL = 10 * time.Minute

// checkFirstRunSync prompts the user to set up sync on first run.
func checkFirstRunSync() {
	cfgPath := csync.DefaultConfigPath()
	cfg, err := csync.LoadConfig(cfgPath)
	if err != nil {
		return
	}

	// If config exists (either configured or declined), skip prompt.
	if cfg != nil {
		return
	}

	fmt.Print("No sync configured — connect to Dropbox? (y/N) ")
	var answer string
	_, _ = fmt.Scanln(&answer)
	answer = strings.ToLower(strings.TrimSpace(answer))

	if answer == "y" || answer == "yes" {
		fmt.Println("Run 'cairn sync setup' to configure sync.")
	} else {
		// Record that the user declined.
		declined := &csync.SyncConfig{SyncDeclined: true}
		_ = csync.SaveConfig(cfgPath, declined)
	}
}

// backgroundSyncPull spawns a detached background process to run "cairn sync pull".
// The subprocess inherits no stdout/stderr, so it cannot interfere with the user's
// terminal. If sync is not configured or the binary path cannot be resolved, it
// returns silently. A lockfile prevents concurrent syncs.
func backgroundSyncPull() {
	cfgPath := csync.DefaultConfigPath()
	cfg, err := csync.LoadConfig(cfgPath)
	if err != nil || !csync.IsConfigured(cfg) {
		return
	}

	lockPath, lockFile, ok := acquireSyncLock("pull")
	defer func() {
		if lockFile != nil {
			_ = lockFile.Close()
		}
	}()
	if !ok {
		return
	}

	self, err := os.Executable()
	if err != nil {
		_ = os.Remove(lockPath)
		return
	}

	cmd := exec.Command(self, "sync", "pull")
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		_ = os.Remove(lockPath)
		return
	}
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	cmd.Stdin = nil
	cmd.Env = append(os.Environ(), syncLockEnv+"="+lockPath)

	if err := cmd.Start(); err != nil {
		_ = os.Remove(lockPath)
	}
	_ = devNull.Close()
}

// backgroundSyncPush spawns a detached background process to run "cairn sync push".
// The subprocess inherits no stdout/stderr, so it cannot interfere with the user's
// terminal. If sync is not configured or the binary path cannot be resolved, it
// returns silently. A lockfile prevents concurrent syncs.
func backgroundSyncPush() {
	cfgPath := csync.DefaultConfigPath()
	cfg, err := csync.LoadConfig(cfgPath)
	if err != nil || !csync.IsConfigured(cfg) {
		return
	}

	lockPath, lockFile, ok := acquireSyncLock("push")
	defer func() {
		if lockFile != nil {
			_ = lockFile.Close()
		}
	}()
	if !ok {
		return
	}

	self, err := os.Executable()
	if err != nil {
		_ = os.Remove(lockPath)
		return
	}

	cmd := exec.Command(self, "sync", "push")
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		_ = os.Remove(lockPath)
		return
	}
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	cmd.Stdin = nil
	cmd.Env = append(os.Environ(), syncLockEnv+"="+lockPath)

	if err := cmd.Start(); err != nil {
		_ = os.Remove(lockPath)
	}
	_ = devNull.Close()
}

func acquireSyncLock(kind string) (string, *os.File, bool) {
	cfgPath := csync.DefaultConfigPath()
	lockPath := filepath.Join(filepath.Dir(cfgPath), fmt.Sprintf("cairn-sync-%s.lock", kind))
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if os.IsExist(err) {
			if cleanupStaleLock(lockPath) {
				lockFile, err = os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
				if err == nil {
					_, _ = fmt.Fprintf(lockFile, "%d\n", os.Getpid())
					return lockPath, lockFile, true
				}
			}
			return lockPath, nil, false
		}
		return lockPath, nil, false
	}
	_, _ = fmt.Fprintf(lockFile, "%d\n", os.Getpid())
	return lockPath, lockFile, true
}

func cleanupStaleLock(lockPath string) bool {
	info, err := os.Stat(lockPath)
	if err != nil {
		return false
	}
	if time.Since(info.ModTime()) > syncLockTTL {
		_ = os.Remove(lockPath)
		return true
	}
	pid, err := readLockPID(lockPath)
	if err != nil || pid <= 0 {
		_ = os.Remove(lockPath)
		return true
	}
	if !processAlive(pid) {
		_ = os.Remove(lockPath)
		return true
	}
	return false
}

func readLockPID(lockPath string) (int, error) {
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

func runSync(dbPath string, cfgManager *config.Manager, args []string) {
	subcmd := args[0]
	switch subcmd {
	case "setup":
		runSyncSetup(dbPath, cfgManager)
	case "push":
		runSyncPush(dbPath)
	case "pull":
		runSyncPull(dbPath)
	case "status":
		runSyncStatus(dbPath)
	case "auth":
		runSyncAuth(dbPath, cfgManager.Get())
	case "unlink":
		runSyncUnlink(dbPath)
	default:
		_, _ = fmt.Fprintf(os.Stderr, "unknown sync command: %s\n", subcmd)
		printCommandHelp("sync")
		os.Exit(3)
	}
}

// promptForSetupConfig interactively prompts the user for any missing configuration
// values needed by sync setup, then persists them to cairn.json.
func promptForSetupConfig(cfgManager *config.Manager) {
	reader := bufio.NewReader(os.Stdin)
	var promptedAppKey, promptedDBPath string

	// Prompt for Dropbox App Key only when not already resolved.
	if cfgManager.Get().DropboxAppKey == "" {
		for {
			fmt.Print("Enter your Dropbox App Key: ")
			key, err := reader.ReadString('\n')
			if err != nil {
				// EOF or read error — exit cleanly without writing.
				_, _ = fmt.Fprintln(os.Stderr, "\nSetup cancelled.")
				os.Exit(0)
			}
			key = strings.TrimSpace(key)
			if key != "" {
				promptedAppKey = key
				cfgManager.Set("dropbox_app_key", key)
				break
			}
			_, _ = fmt.Fprintln(os.Stderr, "Error: App Key cannot be empty. Please try again.")
		}
	}

	// Prompt for database path only when no higher-precedence source has set it.
	if os.Getenv("CAIRN_DB_PATH") == "" && cfgManager.Get().DBPath == config.DefaultDBPath() {
		defaultPath := config.DefaultDBPath()
		fmt.Printf("Enter database path (press Enter for default: %s): ", defaultPath)
		path, err := reader.ReadString('\n')
		if err == nil {
			path = strings.TrimSpace(path)
			if path != "" {
				promptedDBPath = path
				cfgManager.Set("db_path", path)
			}
		}
	}

	// Write only the values that were explicitly provided during this prompt session.
	if promptedAppKey != "" || promptedDBPath != "" {
		if err := writePromptedConfig(promptedAppKey, promptedDBPath); err != nil {
			fatalf(3, "%v", err)
		}
		fmt.Printf("Config written to %s\n", config.DefaultConfigPath())
	}
}

// writePromptedConfig writes explicitly-provided config values to cairn.json,
// preserving any existing keys. Only non-empty values are written.
func writePromptedConfig(appKey, dbPath string) error {
	configPath := config.DefaultConfigPath()
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("cannot create config directory: %w", err)
	}

	// Load existing config or start with an empty map.
	data := map[string]interface{}{}
	if raw, err := os.ReadFile(configPath); err == nil {
		_ = json.Unmarshal(raw, &data)
	}

	if appKey != "" {
		data["dropbox_app_key"] = appKey
	}
	if dbPath != "" {
		data["db_path"] = dbPath
	}

	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal config: %w", err)
	}
	if err := os.WriteFile(configPath, append(raw, '\n'), 0600); err != nil {
		return fmt.Errorf("cannot write config file %s: %w", configPath, err)
	}
	return nil
}

func runSyncSetup(dbPath string, cfgManager *config.Manager) {
	promptForSetupConfig(cfgManager)

	appKey := cfgManager.Get().DropboxAppKey

	s := openStore(dbPath)
	defer func() { _ = s.Close() }()

	cfgPath := csync.DefaultConfigPath()
	cfg, err := csync.LoadConfig(cfgPath)
	if err != nil {
		fatalf(3, "load sync config: %v", err)
	}
	if csync.IsConfigured(cfg) {
		fmt.Println("Sync is already configured. Use 'cairn sync unlink' first to reconfigure.")
		return
	}

	engine := csync.NewEngine(s, nil, cfg, cfgPath)
	count, err := engine.Setup(appKey)
	if err != nil {
		fatalf(3, "sync setup: %v", err)
	}
	fmt.Printf("Sync configured. %d bookmarks synced.\n", count)
}

func runSyncPush(dbPath string) {
	lockPath := os.Getenv(syncLockEnv)
	if lockPath != "" {
		defer func() { _ = os.Remove(lockPath) }()
	}
	engine := openSyncEngine(dbPath)
	defer func() { _ = engine.Store.Close() }()

	if err := engine.Push(); err != nil {
		if lockPath != "" {
			_, _ = fmt.Fprintf(os.Stderr, "sync push: %v\n", err)
			return
		}
		fatalf(3, "sync push: %v", err)
	}
	fmt.Println("Push complete.")
}

func runSyncPull(dbPath string) {
	lockPath := os.Getenv(syncLockEnv)
	if lockPath != "" {
		defer func() { _ = os.Remove(lockPath) }()
	}
	engine := openSyncEngine(dbPath)
	defer func() { _ = engine.Store.Close() }()

	count, err := engine.Pull()
	if err != nil {
		if lockPath != "" {
			_, _ = fmt.Fprintf(os.Stderr, "sync pull: %v\n", err)
			return
		}
		fatalf(3, "sync pull: %v", err)
	}
	fmt.Printf("Pull complete. %d bookmarks synced.\n", count)
}

func runSyncStatus(dbPath string) {
	s := openStore(dbPath)
	defer func() { _ = s.Close() }()

	cfgPath := csync.DefaultConfigPath()
	cfg, err := csync.LoadConfig(cfgPath)
	if err != nil {
		fatalf(3, "load sync config: %v", err)
	}

	engine := csync.NewEngine(s, nil, cfg, cfgPath)
	status, err := engine.Status()
	if err != nil {
		fatalf(3, "sync status: %v", err)
	}

	if !status.Configured {
		fmt.Println("Sync is not configured. Run 'cairn sync setup' to get started.")
		return
	}

	fmt.Printf("Backend:         %s\n", status.Backend)
	fmt.Printf("Device ID:       %s\n", status.DeviceID)
	if status.LastSyncAt != nil {
		fmt.Printf("Last sync:       %s\n", status.LastSyncAt.Format("2006-01-02 15:04:05 UTC"))
	} else {
		fmt.Println("Last sync:       never")
	}
	fmt.Printf("Pending changes: %d\n", status.PendingCount)
}

func runSyncAuth(_ string, appCfg *config.AppConfig) {
	appKey := appCfg.DropboxAppKey
	if appKey == "" {
		fatalf(3, "CAIRN_DROPBOX_APP_KEY is required (set via environment variable or cairn.json)")
	}

	cfgPath := csync.DefaultConfigPath()
	cfg, err := csync.LoadConfig(cfgPath)
	if err != nil {
		fatalf(3, "load sync config: %v", err)
	}
	if cfg == nil {
		fatalf(3, "sync not configured. Run 'cairn sync setup' first.")
		return
	}

	token, err := csync.RunOAuth2Flow(appKey)
	if err != nil {
		fatalf(3, "oauth2 flow: %v", err)
	}

	cfg.Dropbox = &csync.DropboxConfig{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenExpiry:  token.Expiry,
		AppKey:       appKey,
	}

	if err := csync.SaveConfig(cfgPath, cfg); err != nil {
		fatalf(3, "save config: %v", err)
	}
	fmt.Println("Authentication updated.")
}

func runSyncUnlink(dbPath string) {
	cfgPath := csync.DefaultConfigPath()
	cfg, err := csync.LoadConfig(cfgPath)
	if err != nil {
		fatalf(3, "load sync config: %v", err)
	}
	if cfg == nil {
		fmt.Println("Sync is not configured.")
		return
	}

	s := openStore(dbPath)
	defer func() { _ = s.Close() }()

	engine := csync.NewEngine(s, nil, cfg, cfgPath)
	if err := engine.Unlink(); err != nil {
		fatalf(3, "unlink: %v", err)
	}
	fmt.Println("Sync unlinked. Local bookmarks are preserved.")
}

// openSyncEngine creates a sync engine with an active backend from config.
func openSyncEngine(dbPath string) *csync.Engine {
	s := openStore(dbPath)

	cfgPath := csync.DefaultConfigPath()
	cfg, err := csync.LoadConfig(cfgPath)
	if err != nil {
		_ = s.Close()
		fatalf(3, "load sync config: %v", err)
	}
	if !csync.IsConfigured(cfg) {
		_ = s.Close()
		fatalf(3, "sync not configured. Run 'cairn sync setup' first.")
	}

	b, err := csync.NewBackend(cfg)
	if err != nil {
		_ = s.Close()
		fatalf(3, "create sync backend: %v", err)
	}

	return csync.NewEngine(s, b, cfg, cfgPath)
}
