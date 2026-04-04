package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ndy40/cairn/internal/store"
)

func openStore(dbPath string) *store.Store {
	s, err := store.Open(dbPath)
	if err != nil {
		fatalf(3, "open database: %v", err)
	}
	return s
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
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(code)
}
