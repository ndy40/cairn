package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/ndy40/cairn/internal/config"
)

//go:embed text
var helpFS embed.FS

var helpTexts = map[string]string{}

func init() {
	entries, _ := fs.ReadDir(helpFS, "text")
	for _, e := range entries {
		data, _ := helpFS.ReadFile("text/" + e.Name())
		key := strings.TrimSuffix(e.Name(), ".txt") // "help-add" or "help"
		key = strings.TrimPrefix(key, "help-")      // "add" or "help"
		if key == "help" {
			key = ""
		}
		helpTexts[key] = string(data)
	}
}

var version = "dev"

// cmdContext carries the resolved runtime values passed to every command handler.
type cmdContext struct {
	args       []string
	db         string
	appCfg     *config.AppConfig
	cfgManager *config.Manager
}

// command describes a registered subcommand.
type command struct {
	run      func(ctx cmdContext)
	autoSync bool // whether to trigger background sync before running
}

// commands is the central registry mapping subcommand names to their handlers.
var commands = map[string]command{
	"add":     {run: cmdAdd, autoSync: true},
	"edit":    {run: cmdEdit, autoSync: true},
	"list":    {run: cmdList, autoSync: true},
	"search":  {run: cmdSearch, autoSync: true},
	"delete":  {run: cmdDelete, autoSync: true},
	"pin":     {run: cmdPin, autoSync: true},
	"sync":    {run: cmdSyncCmd, autoSync: true},
	"update":  {run: cmdUpdate, autoSync: false},
	"version": {run: cmdVersion, autoSync: false},
	"config":  {run: cmdConfig, autoSync: false},
	"help":    {run: cmdHelpCmd, autoSync: false},
}

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

	if len(args) == 0 {
		// No subcommand — launch TUI (auto-sync happens inside runTUI path).
		checkFirstRunSync()
		backgroundSyncPull()
		runTUI(resolvedDB)
		return
	}

	cmd, ok := commands[args[0]]
	if !ok {
		_, _ = fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		printHelp()
		os.Exit(3)
	}

	ctx := cmdContext{
		args:       args[1:],
		db:         resolvedDB,
		appCfg:     appCfg,
		cfgManager: cfgManager,
	}

	if cmd.autoSync {
		checkFirstRunSync()
		backgroundSyncPull()
	}

	cmd.run(ctx)
}
