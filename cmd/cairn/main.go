package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ndy40/cairn/internal/display"
	"github.com/ndy40/cairn/internal/fetcher"
	"github.com/ndy40/cairn/internal/model"
	"github.com/ndy40/cairn/internal/search"
	"github.com/ndy40/cairn/internal/store"
)

var version = "dev"

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
	dbPath := flag.String("db", "", "path to bookmark database (overrides CAIRN_DB_PATH)")
	flag.Parse()

	args := flag.Args()

	// Resolve database path.
	resolvedDB, err := resolveDBPath(*dbPath)
	if err != nil {
		fatalf(3, "resolve DB path: %v", err)
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
	case "version":
		if len(args) > 1 && (args[1] == "--help" || args[1] == "-h") {
			printVersionHelp()
			os.Exit(0)
		}
		fmt.Printf("bm version %s\n", version)
	case "help":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		printHelp()
		os.Exit(3)
	}
}

func resolveDBPath(flag string) (string, error) {
	if flag != "" {
		return flag, nil
	}
	if env := os.Getenv("CAIRN_DB_PATH"); env != "" {
		return env, nil
	}
	return store.DefaultPath()
}

func runTUI(dbPath string) {
	// US1: prerequisite check before opening the database or TUI.
	result := display.CheckPrerequisites()
	if result.ShouldBlock {
		fmt.Fprintln(os.Stderr, result.InstallHint)
		os.Exit(1)
	}
	if result.DisplayType == display.Unknown && !result.ToolFound {
		fmt.Fprintln(os.Stderr, result.InstallHint)
	}

	s, err := store.Open(dbPath)
	if err != nil {
		fatalf(3, "open database: %v", err)
	}
	defer s.Close()

	// US5: archive stale bookmarks on every startup.
	archiveCount, err := s.ArchiveStale()
	if err != nil {
		// Non-fatal: log and continue.
		fmt.Fprintf(os.Stderr, "warning: archive check failed: %v\n", err)
	}

	app := model.New(s, archiveCount)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fatalf(3, "TUI error: %v", err)
	}
}

func runAdd(dbPath, rawURL string, tags []string) {
	s := openStore(dbPath)
	defer s.Close()

	title, description, fetchErr := fetcher.Fetch(rawURL)
	_, insertErr := s.Insert(rawURL, title, description, tags)
	if insertErr != nil {
		if insertErr == store.ErrDuplicate {
			fmt.Fprintln(os.Stderr, "Already bookmarked")
			os.Exit(1)
		}
		fatalf(3, "save bookmark: %v", insertErr)
	}

	domain := domainFromURL(rawURL)
	if fetchErr != nil {
		fmt.Printf("Saved (title unavailable): (%s)\n", domain)
		os.Exit(2)
	}
	fmt.Printf("Saved: %q (%s)\n", title, domain)
}

func runList(dbPath string, args []string) {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.Usage = func() {
		printListHelp()
		os.Exit(0)
	}
	jsonOut := fs.Bool("json", false, "output as JSON")
	if err := fs.Parse(args); err != nil {
		os.Exit(3)
	}

	s := openStore(dbPath)
	defer s.Close()

	bookmarks, err := s.List()
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
	defer s.Close()

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
	defer s.Close()

	if err := s.DeleteByID(id); err != nil {
		if err == store.ErrNotFound {
			fmt.Fprintln(os.Stderr, "Not found")
			os.Exit(1)
		}
		fatalf(3, "delete: %v", err)
	}
	fmt.Println("Deleted")
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
		enc.Encode(bookmarks)
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

func fatalf(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(code)
}

func printHelp() {
	fmt.Println(`cairn - terminal bookmark manager

Usage:
  cairn                    Launch interactive TUI
  cairn add <url> [--tags <tags>]  Save a bookmark non-interactively
  cairn list [--json]      List all bookmarks
  cairn search <query> [--json] [--limit N]  Search bookmarks
  cairn delete <id>        Delete a bookmark by ID
  cairn version            Print version
  cairn help               Show this help

Flags:
  --db <path>           Override default database path

Environment:
  CAIRN_DB_PATH            Override default database path`)
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
	fmt.Println(`Usage: cairn list [--json]

List all bookmarks ordered by date added (newest first).

Flags:
  --json    Output as JSON array instead of tab-separated text`)
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

func printVersionHelp() {
	fmt.Println(`Usage: cairn version

Print the application version and exit.`)
}
