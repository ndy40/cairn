package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ndy40/cairn/internal/config"
	"github.com/ndy40/cairn/internal/display"
	"github.com/ndy40/cairn/internal/fetcher"
	"github.com/ndy40/cairn/internal/model"
	"github.com/ndy40/cairn/internal/search"
	"github.com/ndy40/cairn/internal/store"
	csync "github.com/ndy40/cairn/internal/sync"
)

var version = "dev"

const syncLockEnv = "CAIRN_SYNC_LOCK"
const syncLockTTL = 10 * time.Minute

func main() {
	// Intercept root-level -h/--help before flag.Parse() to ensure exit 0.
	for _, a := range os.Args[1:] {
		if a == "-h" || a == "--help" {
			printHelp()
			os.Exit(0)
		}
		if !strings.HasPrefix(a, "-") {
			break // stop at first non-flag (subcommand name)
		}
	}

	// Global flags.
	dbPath := flag.String("db", "", "path to bookmark database")
	flag.Parse()

	args := flag.Args()

	// Initialize configuration manager.
	cfgManager := config.NewManager()

	// Load configuration from all sources.
	if err := cfgManager.Load("", *dbPath); err != nil {
		fatalf(3, "failed to load configuration: %v", err)
	}

	// Get the resolved configuration.
	appCfg := cfgManager.Get()
	resolvedDB := appCfg.DBPath

	if shouldAutoSync(args) {
		// First-run sync prompt.
		checkFirstRunSync()

		// Auto-pull on startup if sync is configured (background, non-blocking).
		backgroundSyncPull()
	}

	if len(args) == 0 {
		// No subcommand — launch TUI.
		runTUI(resolvedDB)
		return
	}

	switch args[0] {
	case "add":
		if len(args) > 1 && (args[1] == "--help" || args[1] == "-h") {
			printAddHelp()
			os.Exit(0)
		}
		if len(args) < 2 {
			fatalf(3, "usage: bm add <url>")
		}
		fs := flag.NewFlagSet("add", flag.ContinueOnError)
		tagsFlag := fs.String("tags", "", "comma-separated tags")
		if err := fs.Parse(args[2:]); err != nil {
			fatalf(3, "bm add: %v", err)
		}
		runAdd(resolvedDB, args[1], store.NormaliseTagsFromString(*tagsFlag))
	case "list":
		runList(resolvedDB, args[1:])
	case "search":
		if len(args) > 1 && (args[1] == "--help" || args[1] == "-h") {
			printSearchHelp()
			os.Exit(0)
		}
		if len(args) < 2 {
			fatalf(3, "usage: bm search <query>")
		}
		runSearch(resolvedDB, args[1], args[2:])
	case "delete":
		if len(args) > 1 && (args[1] == "--help" || args[1] == "-h") {
			printDeleteHelp()
			os.Exit(0)
		}
		if len(args) < 2 {
			fatalf(3, "usage: bm delete <id>")
		}
		runDelete(resolvedDB, args[1])
	case "pin":
		if len(args) < 2 {
			fatalf(3, "usage: cairn pin <id>")
		}
		runPin(resolvedDB, args[1])
	case "sync":
		if len(args) < 2 {
			printSyncHelp()
			os.Exit(0)
		}
		if args[1] == "--help" || args[1] == "-h" {
			printSyncHelp()
			os.Exit(0)
		}
		runSync(resolvedDB, cfgManager, args[1:])
	case "version":
		if len(args) > 1 && (args[1] == "--help" || args[1] == "-h") {
			printVersionHelp()
			os.Exit(0)
		}
		fmt.Printf("bm version %s\n", version)
	case "help":
		printHelp()
	case "config":
		fmt.Printf("CAIRN_DB_PATH=%s\n", appCfg.DBPath)
		if appCfg.DropboxAppKey != "" {
			fmt.Println("CAIRN_DROPBOX_APP_KEY=(set)")
		}
	default:
		_, _ = fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		printHelp()
		os.Exit(3)
	}
}

func runTUI(dbPath string) {
	// US1: prerequisite check before opening the database or TUI.
	result := display.CheckPrerequisites()
	if result.ShouldBlock {
		_, _ = fmt.Fprintln(os.Stderr, result.InstallHint)
		os.Exit(1)
	}
	if result.DisplayType == display.Unknown && !result.ToolFound {
		_, _ = fmt.Fprintln(os.Stderr, result.InstallHint)
	}

	s, err := store.Open(dbPath)
	if err != nil {
		fatalf(3, "open database: %v", err)
	}
	defer func() { _ = s.Close() }()

	// US5: archive stale bookmarks on every startup.
	archiveCount, err := s.ArchiveStale()
	if err != nil {
		// Non-fatal: log and continue.
		_, _ = fmt.Fprintf(os.Stderr, "warning: archive check failed: %v\n", err)
	}

	app := model.New(s, archiveCount)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fatalf(3, "TUI error: %v", err)
	}
}

func runAdd(dbPath, rawURL string, tags []string) {
	s := openStore(dbPath)
	defer func() { _ = s.Close() }()

	title, description, fetchErr := fetcher.Fetch(rawURL)
	_, insertErr := s.Insert(rawURL, title, description, tags)
	if insertErr != nil {
		if insertErr == store.ErrDuplicate {
			_, _ = fmt.Fprintln(os.Stderr, "Already bookmarked")
			os.Exit(1)
		}
		fatalf(3, "save bookmark: %v", insertErr)
	}

	domain := domainFromURL(rawURL)
	if fetchErr != nil {
		fmt.Printf("Saved (title unavailable): (%s)\n", domain)
		backgroundSyncPush()
		os.Exit(2)
	}
	fmt.Printf("Saved: %q (%s)\n", title, domain)
	backgroundSyncPush()
}

func runList(dbPath string, args []string) {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.Usage = func() {
		printListHelp()
		os.Exit(0)
	}
	jsonOut := fs.Bool("json", false, "output as JSON")
	order := fs.String("order", "desc", "sort order: asc or desc (default desc)")
	if err := fs.Parse(args); err != nil {
		os.Exit(3)
	}

	var asc bool
	switch strings.ToLower(*order) {
	case "asc":
		asc = true
	case "desc":
		asc = false
	default:
		fatalf(3, "list: --order must be asc or desc, got %q", *order)
	}

	s := openStore(dbPath)
	defer func() { _ = s.Close() }()

	bookmarks, err := s.ListOrdered(asc)
	if err != nil {
		fatalf(3, "list: %v", err)
	}
	printBookmarks(bookmarks, *jsonOut)
}

func runSearch(dbPath, query string, args []string) {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.Usage = func() {
		printSearchHelp()
		os.Exit(0)
	}
	jsonOut := fs.Bool("json", false, "output as JSON")
	limit := fs.Int("limit", 10, "maximum results")
	if err := fs.Parse(args); err != nil {
		os.Exit(3)
	}

	s := openStore(dbPath)
	defer func() { _ = s.Close() }()

	all, err := s.List()
	if err != nil {
		fatalf(3, "list: %v", err)
	}

	ids, _ := s.FTSSearch(query)
	var candidates []*store.Bookmark
	if len(ids) > 0 {
		idSet := make(map[int64]bool, len(ids))
		for _, id := range ids {
			idSet[id] = true
		}
		for _, b := range all {
			if idSet[b.ID] {
				candidates = append(candidates, b)
			}
		}
	} else {
		candidates = all
	}

	results := search.Search(query, candidates)
	if *limit > 0 && len(results) > *limit {
		results = results[:*limit]
	}
	printBookmarks(results, *jsonOut)
}

func runDelete(dbPath, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fatalf(3, "invalid id: %s", idStr)
	}

	s := openStore(dbPath)
	defer func() { _ = s.Close() }()

	if _, err := s.DeleteByID(id); err != nil {
		if err == store.ErrNotFound {
			_, _ = fmt.Fprintln(os.Stderr, "Not found")
			os.Exit(1)
		}
		fatalf(3, "delete: %v", err)
	}
	fmt.Println("Deleted")
	backgroundSyncPush()
}

func runPin(dbPath, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fatalf(3, "invalid id: %s", idStr)
	}

	s := openStore(dbPath)
	defer func() { _ = s.Close() }()

	b, err := s.GetByID(id)
	if err != nil {
		if err == store.ErrNotFound {
			_, _ = fmt.Fprintln(os.Stderr, "Bookmark not found")
			os.Exit(1)
		}
		fatalf(3, "pin: %v", err)
	}

	newState := !b.IsPermanent
	if err := s.SetPermanent(id, newState); err != nil {
		fatalf(3, "pin: %v", err)
	}

	verb := "Pinned"
	if !newState {
		verb = "Unpinned"
	}
	fmt.Printf("%s: %q (%s)\n", verb, b.Title, domainFromURL(b.URL))
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
		printSyncHelp()
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

func openStore(dbPath string) *store.Store {
	s, err := store.Open(dbPath)
	if err != nil {
		fatalf(3, "open database: %v", err)
	}
	return s
}

func printBookmarks(bookmarks []*store.Bookmark, asJSON bool) {
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(bookmarks)
		return
	}
	for _, b := range bookmarks {
		fmt.Printf("%d\t%s\t%s\t%s\t%s\n",
			b.ID, b.Title, b.URL, b.Domain, b.CreatedAt.Format("2006-01-02T15:04:05Z"))
	}
}

func domainFromURL(rawURL string) string {
	// Reuse the same logic as store.extractDomain via URL parse.
	parts := strings.SplitN(rawURL, "//", 2)
	if len(parts) < 2 {
		return rawURL
	}
	host := strings.SplitN(parts[1], "/", 2)[0]
	host = strings.TrimPrefix(host, "www.")
	return host
}

func shouldAutoSync(args []string) bool {
	if len(args) == 0 {
		return true
	}
	switch args[0] {
	case "add", "list", "search", "delete", "pin", "sync":
		return true
	default:
		return false
	}
}

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
	if !ok {
		return
	}
	defer func() {
		if lockFile != nil {
			_ = lockFile.Close()
		}
	}()

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
	if !ok {
		return
	}
	defer func() {
		if lockFile != nil {
			_ = lockFile.Close()
		}
	}()

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

func fatalf(code int, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(code)
}

func printHelp() {
	fmt.Println(`cairn - terminal bookmark manager

Usage:
  cairn                    Launch interactive TUI
  cairn add <url> [--tags <tags>]  Save a bookmark non-interactively
  cairn list [--json] [--order asc|desc]  List all bookmarks
  cairn search <query> [--json] [--limit N]  Search bookmarks
  cairn delete <id>        Delete a bookmark by ID
  cairn pin <id>           Toggle pin (permanent) flag on a bookmark
  cairn sync <command>     Manage bookmark sync
  cairn version            Print version
  cairn help               Show this help

Sync Commands:
  cairn sync setup         Connect to Dropbox and set up sync
  cairn sync push          Push local changes to cloud
  cairn sync pull          Pull remote changes from cloud
  cairn sync status        Show sync status
  cairn sync auth          Re-authenticate with Dropbox
  cairn sync unlink        Disconnect sync (keeps local data)

Flags:
  --db <path>           Override default database path

Environment:
  CAIRN_DB_PATH            Override default database path
  CAIRN_DROPBOX_APP_KEY    Dropbox app key for sync setup`)
}

func printAddHelp() {
	fmt.Println(`Usage: cairn add <url> [--tags <comma-separated>]

Save a bookmark by URL. The page title and description are fetched automatically.

Arguments:
  <url>    The URL to bookmark (required)

Flags:
  --tags   Comma-separated tags (e.g. "work, go, tools") — max 3 tags

Exit codes:
  0  Saved successfully
  1  Already bookmarked (duplicate URL)
  2  Saved but title could not be fetched
  3  Error (invalid arguments, database error)`)
}

func printListHelp() {
	fmt.Println(`Usage: cairn list [--json] [--order asc|desc]

List all bookmarks ordered by date added.

Flags:
  --json         Output as JSON array instead of tab-separated text
  --order asc    Oldest first
  --order desc   Newest first (default)`)
}

func printSearchHelp() {
	fmt.Println(`Usage: cairn search <query> [--json] [--limit N]

Search bookmarks by title, domain, and description.

Arguments:
  <query>  Search query (required)

Flags:
  --json       Output as JSON array
  --limit N    Maximum number of results to return (default: 10)`)
}

func printDeleteHelp() {
	fmt.Println(`Usage: cairn delete <id>

Delete a bookmark by its numeric ID.

Arguments:
  <id>    Bookmark ID (required, use cairn list to find IDs)

Exit codes:
  0  Deleted successfully
  1  Bookmark not found
  3  Error`)
}

func printSyncHelp() {
	fmt.Println(`Usage: cairn sync <command>

Manage bookmark synchronization across devices.

Commands:
  setup     Connect to Dropbox and set up sync
  push      Push local changes to cloud
  pull      Pull remote changes from cloud
  status    Show sync configuration and pending changes
  auth      Re-authenticate with Dropbox
  unlink    Disconnect sync (local bookmarks are preserved)

Environment:
  CAIRN_DROPBOX_APP_KEY    Required for setup and auth commands`)
}

func printVersionHelp() {
	fmt.Println(`Usage: cairn version

Print the application version and exit.`)
}
