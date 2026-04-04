package main

import (
	"errors"
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
		printCommandHelp("list")
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
		printCommandHelp("search")
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

func runEdit(dbPath, idStr, rawURL string, urlSet bool, title string, titleSet bool, tags []string, tagsSet bool) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fatalf(3, "invalid id: %s", idStr)
	}
	if !urlSet && !titleSet && !tagsSet {
		fatalf(3, "edit: specify at least one of --url, --title, or --tags")
	}
	if urlSet {
		rawURL = strings.TrimSpace(rawURL)
		if rawURL == "" {
			fatalf(3, "edit: --url cannot be empty")
		}
	}
	if titleSet {
		title = strings.TrimSpace(title)
		if title == "" {
			fatalf(3, "edit: --title cannot be empty")
		}
	}

	s := openStore(dbPath)
	defer func() { _ = s.Close() }()

	patch := store.BookmarkPatch{}
	if urlSet {
		patch.URL = &rawURL
		// Re-fetch title from the new URL unless the user explicitly provided --title.
		if !titleSet {
			fetchedTitle, _, _ := fetcher.Fetch(rawURL)
			if fetchedTitle != "" {
				patch.Title = &fetchedTitle
			}
		}
	}
	if titleSet {
		patch.Title = &title
	}
	if tagsSet {
		patch.Tags = &tags
	}

	if err := s.UpdateFields(id, patch); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			_, _ = fmt.Fprintln(os.Stderr, "Not found")
			os.Exit(1)
		}
		if errors.Is(err, store.ErrDuplicateURL) {
			_, _ = fmt.Fprintln(os.Stderr, "Duplicate URL")
			os.Exit(1)
		}
		fatalf(3, "edit: %v", err)
	}

	var updated []string
	if urlSet {
		updated = append(updated, "URL")
	}
	if titleSet {
		updated = append(updated, "title")
	}
	if tagsSet {
		updated = append(updated, "tags")
	}
	fmt.Printf("Updated %s\n", strings.Join(updated, " and "))
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
