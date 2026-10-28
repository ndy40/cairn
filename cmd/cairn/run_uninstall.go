package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/ndy40/cairn/internal/config"
	csync "github.com/ndy40/cairn/internal/sync"
	"github.com/ndy40/cairn/internal/updater"
)

// uninstallPlan lists everything cairn uninstall can remove.
type uninstallPlan struct {
	binary       string   // resolved path of the running cairn binary
	extension    string   // Vicinae extension directory; empty when not installed
	desktop      string   // desktop app (binary on Linux, .app on macOS); empty when not installed
	desktopFiles []string // desktop app's icon, version file, launcher and login items
	otherDesktop string   // desktop app found elsewhere on PATH, which is left in place
	data         []string // database, config and sync files that exist on disk
	syncFile     string   // sync.json path when Dropbox sync is set up (holds the token)
	ownDirs      []string // cairn's own directories, removed with --purge once empty
}

func runUninstall(dbPath string, cfgManager *config.Manager, purge, yes bool) {
	plan, err := buildUninstallPlan(dbPath, cfgManager.ConfigFileUsed())
	if err != nil {
		fatalf(1, "cairn uninstall: %v", err)
	}

	if runtime.GOOS == "windows" {
		printManualUninstall(plan)
		return
	}

	if blocked := blockedPaths(plan); len(blocked) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "cairn uninstall: permission denied: cannot remove %s\n", strings.Join(blocked, ", "))
		printManualSteps(plan, blocked)
		printKeptData(plan)
		os.Exit(4)
	}

	if user := os.Getenv("SUDO_USER"); user != "" {
		fmt.Printf("Note: running under sudo, so the data paths below may belong to root rather than %s.\n\n", user)
	}

	printUninstallPlan(plan, purge)
	if !yes && !confirm(os.Stdin, "Continue? [y/N] ") {
		fmt.Println("Uninstall cancelled.")
		return
	}

	if plan.desktop != "" {
		updater.StopDesktop(plan.desktop)
	}
	if err := executeUninstall(plan, purge); err != nil {
		fatalf(1, "cairn uninstall: %v", err)
	}

	fmt.Println("cairn has been uninstalled.")
	if !purge {
		printKeptData(plan)
	}
	if plan.syncFile != "" {
		fmt.Println("\nThe copy of your bookmarks in Dropbox has not been deleted.")
	}
}

// buildUninstallPlan resolves the running binary, the Vicinae extension and
// every data file cairn may have written. cfgFile is the config file loaded at
// startup; when empty the default cairn.json location is checked instead.
func buildUninstallPlan(dbPath, cfgFile string) (*uninstallPlan, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("resolve executable path: %w", err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return nil, fmt.Errorf("resolve symlinks: %w", err)
	}

	plan := &uninstallPlan{binary: exe}
	if dir, installed := updater.DetectExtension(); installed {
		plan.extension = dir
	}
	if path, installed := updater.DetectDesktop(); installed {
		// The installer puts the Linux desktop app beside cairn. One found
		// elsewhere on PATH belongs to another install, so leave it alone.
		if runtime.GOOS == "darwin" || filepath.Dir(path) == filepath.Dir(exe) {
			plan.desktop = path
			plan.desktopFiles = updater.DesktopFiles()
		} else {
			plan.otherDesktop = path
		}
	}

	if cfgFile == "" {
		cfgFile = config.DefaultConfigPath()
	}
	syncFile := csync.DefaultConfigPath()
	legacyDB := config.LegacyDBPath()

	plan.data = collectDataPaths(dbPath, legacyDB, cfgFile, syncFile)
	// sync.json also exists when sync was declined, so check it holds a token.
	if syncCfg, err := csync.LoadConfig(syncFile); err == nil && csync.IsConfigured(syncCfg) {
		plan.syncFile = syncFile
	}
	// Only directories cairn creates itself are candidates for removal, so a
	// custom --db location such as ~/Documents is never deleted when empty.
	for _, p := range []string{config.DefaultDBPath(), legacyDB, config.DefaultConfigPath(), syncFile} {
		if p != "" {
			plan.ownDirs = append(plan.ownDirs, filepath.Dir(p))
		}
	}
	return plan, nil
}

// collectDataPaths returns the database, config and sync files that exist on
// disk, in a stable order and without duplicates. The SQLite WAL and shared
// memory files are included alongside each database, as are any sync lock
// files left behind in the sync config directory.
func collectDataPaths(dbPath, legacyDBPath, cfgFile, syncFile string) []string {
	var candidates []string
	for _, db := range []string{dbPath, legacyDBPath} {
		if db != "" {
			candidates = append(candidates, db, db+"-wal", db+"-shm")
		}
	}
	candidates = append(candidates, cfgFile, syncFile)
	if syncFile != "" {
		locks, _ := filepath.Glob(filepath.Join(filepath.Dir(syncFile), "cairn-sync-*.lock"))
		candidates = append(candidates, locks...)
	}

	seen := make(map[string]bool)
	var found []string
	for _, p := range candidates {
		if p == "" || seen[p] || !fileExists(p) {
			continue
		}
		seen[p] = true
		found = append(found, p)
	}
	return found
}

// executeUninstall removes the extension, the desktop app, the data files when
// purge is set, and the binary last, so a failure part-way leaves cairn
// runnable for a retry.
func executeUninstall(plan *uninstallPlan, purge bool) error {
	if plan.extension != "" {
		if err := os.RemoveAll(plan.extension); err != nil {
			return fmt.Errorf("remove extension: %w", err)
		}
	}

	for _, p := range plan.desktopFiles {
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove %s: %w", p, err)
		}
	}
	if plan.desktop != "" {
		if err := os.RemoveAll(plan.desktop); err != nil {
			return fmt.Errorf("remove desktop app: %w", err)
		}
	}

	if purge {
		for _, p := range plan.data {
			if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("remove %s: %w", p, err)
			}
		}
		for _, dir := range plan.ownDirs {
			removeEmptyDir(dir)
		}
	}

	if err := os.Remove(plan.binary); err != nil {
		return fmt.Errorf("remove binary: %w", err)
	}
	return nil
}

func printUninstallPlan(plan *uninstallPlan, purge bool) {
	fmt.Println("This will remove:")
	fmt.Printf("  %s\n", plan.binary)
	if plan.extension != "" {
		fmt.Printf("  %s (Vicinae extension)\n", plan.extension)
	}
	if plan.desktop != "" {
		fmt.Printf("  %s (desktop app; stopped first if running)\n", plan.desktop)
	}
	for _, p := range plan.desktopFiles {
		fmt.Printf("  %s\n", p)
	}
	if purge {
		for _, p := range plan.data {
			fmt.Printf("  %s\n", p)
		}
		if len(plan.data) > 0 {
			fmt.Println("\nYour bookmarks and settings will be permanently deleted.")
		}
	} else if len(plan.data) > 0 {
		fmt.Println("\nYour bookmarks and settings will be kept. Use --purge to delete them too.")
	}
	if plan.otherDesktop != "" {
		fmt.Printf("\nThe desktop app at %s is left in place: it was not installed beside this cairn.\n", plan.otherDesktop)
	}
	fmt.Println()
}

// blockedPaths returns the binaries cairn lacks permission to remove.
func blockedPaths(plan *uninstallPlan) []string {
	var blocked []string
	for _, p := range []string{plan.binary, plan.desktop} {
		if p != "" && !dirWritable(filepath.Dir(p)) {
			blocked = append(blocked, p)
		}
	}
	return blocked
}

// printManualSteps prints the commands that remove cairn by hand. It suggests
// sudo only for the blocked paths rather than "sudo cairn uninstall": under
// sudo, HOME may point at root's home and the wrong data paths would be used.
func printManualSteps(plan *uninstallPlan, blocked []string) {
	_, _ = fmt.Fprintln(os.Stderr, "Run these commands instead:")
	if plan.desktop != "" {
		_, _ = fmt.Fprintln(os.Stderr, "  pkill -x cairn-desktop")
	}
	_, _ = fmt.Fprintf(os.Stderr, "  sudo %s\n", rmCommand(blocked))

	var rest []string
	for _, p := range []string{plan.binary, plan.desktop, plan.extension} {
		if p != "" && !slices.Contains(blocked, p) {
			rest = append(rest, p)
		}
	}
	rest = append(rest, plan.desktopFiles...)
	if len(rest) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "  %s\n", rmCommand(rest))
	}
}

// printKeptData lists the data files left on disk and a command to delete them.
func printKeptData(plan *uninstallPlan) {
	if len(plan.data) == 0 {
		return
	}
	fmt.Println("\nYour bookmarks and settings were kept:")
	for _, p := range plan.data {
		if p == plan.syncFile {
			fmt.Printf("  %s (contains your Dropbox access token)\n", p)
		} else {
			fmt.Printf("  %s\n", p)
		}
	}
	fmt.Println("\nTo delete them, run:")
	fmt.Printf("  %s\n", rmCommand(plan.data))
}

// rmCommand returns an rm command for paths, recursive when any of them is a
// directory such as the Vicinae extension or the macOS .app bundle.
func rmCommand(paths []string) string {
	cmd := "rm"
	quoted := make([]string, len(paths))
	for i, p := range paths {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			cmd = "rm -r"
		}
		quoted[i] = shellQuote(p)
	}
	return cmd + " " + strings.Join(quoted, " ")
}

// printManualUninstall is used on Windows, where a running program cannot
// delete its own executable.
func printManualUninstall(plan *uninstallPlan) {
	fmt.Println("cairn uninstall: Windows cannot delete a running program. Remove cairn manually:")
	fmt.Printf("  %s\n", plan.binary)
	if len(plan.data) > 0 {
		fmt.Println("\nYour bookmarks and settings (delete these only if you no longer need them):")
		for _, p := range plan.data {
			fmt.Printf("  %s\n", p)
		}
	}
}

// confirm prints prompt and reports whether the reply was "y" or "yes".
// End of input counts as "no".
func confirm(r io.Reader, prompt string) bool {
	fmt.Print(prompt)
	answer, err := bufio.NewReader(r).ReadString('\n')
	if err != nil && answer == "" {
		fmt.Println()
		return false
	}
	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "y", "yes":
		return true
	}
	return false
}

// shellQuote wraps p in single quotes for POSIX shells, so paths such as
// "~/Library/Application Support" can be copied straight into a terminal.
func shellQuote(p string) string {
	return "'" + strings.ReplaceAll(p, "'", `'\''`) + "'"
}

func dirWritable(dir string) bool {
	f, err := os.CreateTemp(dir, ".cairn-perm-check-*")
	if err != nil {
		return false
	}
	_ = f.Close()
	_ = os.Remove(f.Name())
	return true
}
