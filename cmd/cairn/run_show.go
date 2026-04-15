package main

import (
	"fmt"

	"github.com/ndy40/cairn/internal/config"
	csync "github.com/ndy40/cairn/internal/sync"
)

func runShow(db string, cfgManager *config.Manager) {
	// --- database path & bookmark count ---
	s := openStore(db)
	defer func() { _ = s.Close() }()

	count, err := s.Count()
	if err != nil {
		fatalf(3, "show: count bookmarks: %v", err)
	}

	fmt.Printf("%-12s %s\n", "Database:", db)
	fmt.Printf("%-12s %d\n", "Bookmarks:", count)

	// --- sync status ---
	cfgPath := csync.DefaultConfigPath()
	cfg, err := csync.LoadConfig(cfgPath)
	if err != nil || !csync.IsConfigured(cfg) {
		fmt.Printf("%-12s %s\n", "Sync:", "not configured")
		return
	}

	fmt.Printf("%-12s configured (%s)\n", "Sync:", cfg.Backend)

	if cfg.LastSyncAt != nil {
		fmt.Printf("%-12s %s\n", "Last sync:", cfg.LastSyncAt.UTC().Format("2006-01-02 15:04:05 UTC"))
	} else {
		fmt.Printf("%-12s %s\n", "Last sync:", "never")
	}
}
