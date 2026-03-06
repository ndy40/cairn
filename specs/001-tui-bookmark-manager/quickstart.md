# Quickstart: TUI Bookmark Manager

**Feature**: 001-tui-bookmark-manager
**Date**: 2026-03-06

This guide covers the prerequisites, project setup, and development workflow for the bookmark manager.

---

## Prerequisites

| Tool | Minimum Version | Install |
|------|----------------|---------|
| Go | 1.22+ | https://go.dev/dl/ |
| Git | Any recent | System package manager |
| `xclip` or `xsel` | Any | **Linux only**: `sudo apt install xclip` |

> **Note**: All Go dependencies are pure Go (no CGO). No C compiler or C development headers are required.

---

## Project Initialisation

```bash
# From the repository root
go mod init github.com/<your-username>/bookmark-manager
go mod tidy
```

---

## Dependency Acquisition

After running `go mod tidy` the following dependencies will be resolved:

```
charmbracelet/bubbletea     — TUI runtime
charmbracelet/bubbles       — List, textinput components
charmbracelet/lipgloss      — Styling and layout
modernc.org/sqlite          — Embedded SQLite (pure Go)
PuerkitoBio/goquery         — HTML meta tag parsing
golang.org/x/net            — HTML charset handling (transitive)
sahilm/fuzzy                — Fuzzy search scoring
atotto/clipboard            — Clipboard read
```

---

## Recommended Source Layout

```
bookmark-manager/
├── cmd/
│   └── bm/
│       └── main.go              # Entry point: parses CLI args, launches TUI or runs subcommand
├── internal/
│   ├── model/
│   │   ├── app.go               # Root bubbletea Model (state, Update, View)
│   │   ├── browse.go            # Browse mode sub-model
│   │   ├── search.go            # Search mode sub-model
│   │   └── add.go               # Add modal sub-model
│   ├── store/
│   │   ├── store.go             # SQLite open, migrate, close
│   │   ├── bookmark.go          # CRUD operations
│   │   └── search.go            # FTS5 pre-filter query
│   ├── fetcher/
│   │   └── fetcher.go           # HTTP fetch + goquery meta tag extraction
│   ├── search/
│   │   └── fuzzy.go             # Multi-field fuzzy wrapper around sahilm/fuzzy
│   └── clipboard/
│       └── clipboard.go         # Thin wrapper around atotto/clipboard
├── specs/                       # Feature specifications (this directory)
├── go.mod
└── go.sum
```

---

## Build

```bash
# Development build
go build ./cmd/bm

# Run directly
go run ./cmd/bm

# Release build (single static binary)
CGO_ENABLED=0 go build -ldflags="-s -w" -o bm ./cmd/bm
```

---

## Run

```bash
# Launch interactive TUI
./bm

# Add a bookmark non-interactively
./bm add https://example.com

# Search from the command line
./bm search "golang tui"

# List all bookmarks as JSON
./bm list --json
```

---

## Test

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Database Location

The SQLite database is created automatically on first run:

| Platform | Default Path |
|----------|-------------|
| Linux | `~/.local/share/bookmark-manager/bookmarks.db` |
| macOS | `~/Library/Application Support/bookmark-manager/bookmarks.db` |
| Windows | `%APPDATA%\bookmark-manager\bookmarks.db` |

Override with `--db <path>` flag or `BM_DB_PATH` environment variable.

---

## Linux Clipboard Setup

The clipboard integration on Linux requires `xclip` or `xsel`:

```bash
# Debian/Ubuntu
sudo apt install xclip

# Fedora/RHEL
sudo dnf install xclip

# Arch
sudo pacman -S xclip
```

> **Wayland users**: The application uses the X11 clipboard via XWayland, which is active by default on most desktop environments (GNOME, KDE, etc.). If running a pure Wayland session without XWayland, clipboard paste (Ctrl+P) will not function. A Wayland-native path is planned for a future version.
