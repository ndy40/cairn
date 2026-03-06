# Research: TUI Bookmark Manager

**Feature**: 001-tui-bookmark-manager
**Date**: 2026-03-06

---

## Decision 1: TUI Framework

**Decision**: `charmbracelet/bubbletea` + `charmbracelet/bubbles` + `charmbracelet/lipgloss`

**Rationale**:
Bubbletea uses an Elm-inspired Model-Update-View (MVU) architecture. It is the most actively maintained Go TUI framework (~28k GitHub stars) with a cohesive ecosystem of pre-built components via `bubbles`. Key components available out-of-the-box:
- `bubbles/list` — scrollable, filterable list with keyboard navigation. Purpose-built for a bookmark list.
- `bubbles/textinput` — styled text input with cursor support, ideal for the search bar.
- `bubbles/key` — keyboard map abstraction for clean shortcut documentation.
- `lipgloss` — declarative layout and styling for status bar, footer, and modal overlay.

The MVU architecture maps naturally to the application's discrete states: browsing, searching, adding. Each state transition is explicit and testable via pure functions.

**Alternatives Considered**:
- `rivo/tview`: Mature widget-tree library with built-in modal support and cleaner `SetInputCapture()` for global shortcuts. Rejected because its styling capabilities are weaker (no lipgloss), and mutable widget state makes unit testing harder.
- `gdamore/tcell` (raw): Requires building all UI components from scratch. Cost is disproportionate for an application-level project.

**Gaps to Address in Implementation**:
- No built-in modal component in bubbletea. Use a `showModal bool` state flag and conditionally render a lipgloss-bordered overlay in `View()`.
- Global keyboard shortcuts (Ctrl+P) must be handled at the top-level `Update()` before delegating to sub-models.

---

## Decision 2: Fuzzy Search

**Decision**: `sahilm/fuzzy` with a thin multi-field wrapper

**Rationale**:
`sahilm/fuzzy` uses a Sublime Text-style scoring algorithm that rewards consecutive character matches and word-boundary hits. This produces intuitive ranking (typing "goog" ranks "Google" above "a blog about googling") with minimal allocation. Performance at 1,000 items: well under 1ms per search invocation.

Multi-field search (title, domain, description) is achieved by calling `FindFrom` three times — once per field — with field-weight multipliers (title: 3×, domain: 2×, description: 1×), then merging and deduplicating results by bookmark ID. This wrapper adds ~30 lines of code.

**Alternatives Considered**:
- `lithammer/fuzzysearch`: Simpler API, but ranks by full-string Levenshtein distance. Produces counterintuitive rankings when field lengths vary (e.g., a short title vs. a long description). Rejected for ranking quality.
- `junegunn/fzf` (library): Best-in-class scoring, but internal packages have no stable public API. Breakage on fzf upgrades is a known maintenance hazard. Rejected.
- Custom scorer (Levenshtein + substring): Only advantage is native field weighting, which the sahilm wrapper achieves without custom scoring logic. Rejected as premature.

---

## Decision 3: Local Storage

**Decision**: `modernc.org/sqlite` (pure Go SQLite, no CGO)

**Rationale**:
SQLite via `modernc.org/sqlite` provides full SQL query capability including FTS5 full-text search extension. It is pure Go (mechanical translation of SQLite C source), requires no CGO, and produces a single binary with no system library dependencies. At 1,000+ records, read performance is sub-millisecond.

Key advantages for this application:
- FTS5 enables indexed full-text search across title/description columns — complementing the fuzzy UI search with precise SQL filtering as a pre-filter step.
- Schema migrations are handled by a versioned `schema_version` table checked at startup.
- The `.db` file is portable and inspectable with standard SQLite tooling.
- WAL mode (`PRAGMA journal_mode=WAL`) enables concurrent reads without blocking writes.

**Alternatives Considered**:
- JSON file: Legitimate runner-up. Zero binary size overhead, human-readable. Rejected because all-or-nothing file writes become awkward for concurrent access and FTS5 cannot be replicated without significant custom code.
- `go.etcd.io/bbolt`: Pure Go key-value store. No query language — search requires full scan like JSON. Adding secondary indexes (domain lookup) requires manual bucket management. More complex than JSON with no corresponding benefit. Rejected.
- CSV file: Fragile for URLs and meta descriptions (which frequently contain commas, quotes, newlines). Rejected.

---

## Decision 4: HTTP Fetching + HTML Meta Tag Parsing

**Decision**: `net/http` (stdlib) + `github.com/PuerkitoBio/goquery` + `golang.org/x/net/html/charset`

**Rationale**:
goquery provides a jQuery-style CSS selector API over `golang.org/x/net/html`. Extracting all required meta tags reduces to concise `Find` + `AttrOr` calls, compared to writing a custom recursive node-walker with `golang.org/x/net/html` alone. Both goquery and its `cascadia` CSS selector dependency are pure Go.

`golang.org/x/net/html/charset` is applied unconditionally as a transform around the HTTP response body before parsing. This correctly handles non-UTF-8 pages (ISO-8859-1, Windows-1252, Shift-JIS, GBK, etc.) via HTML5 charset sniffing.

**HTTP Client Configuration**:
- Timeout: 8 seconds (total round-trip)
- Body limit: 512 KB via `io.LimitReader` (meta tags appear in `<head>`, well within first 512 KB of any page)
- User-Agent: `BookmarkManager/1.0` (avoids Go default UA blocks)
- Redirects: default `http.Client` behavior (follows up to 10 redirects)

**Meta Tag Priority**:
- Title: `og:title` > `<title>` > URL hostname
- Description: `og:description` > `meta[name=description]` > empty string

Note: Open Graph tags use `property` attribute; standard meta tags use `name`. Both use `content` for value. Extraction logic must check both attribute names.

**Alternatives Considered**:
- `net/http` + `golang.org/x/net/html` only: Requires custom recursive tree walker. More verbose and error-prone for attribute matching. Rejected in favour of goquery's selector API.
- `mvdan.cc/xurls`: URL extraction from text, not HTML parsing. Not applicable.

---

## Decision 5: Clipboard

**Decision**: `github.com/atotto/clipboard`

**Rationale**:
This application only needs to read plain text (a URL) from the clipboard on Ctrl+P. `atotto/clipboard` handles this with zero CGO on all platforms:
- macOS: `pbcopy`/`pbpaste` via `exec.Command`
- Windows: `syscall` Win32 API
- Linux X11/XWayland: `xclip` or `xsel` external binary

Zero CGO means clean cross-compilation and no C header dependencies at build time — important for a Go CLI tool distributed as a single binary.

**Linux Caveat**: On Linux, `xclip` or `xsel` must be installed on the host system. This is a one-time setup step documented in the quickstart.

**Alternatives Considered**:
- `golang.design/x/clipboard`: Explicitly supports Wayland and image clipboard. Rejected because it requires CGO on all platforms, complicating cross-compilation and requiring C development headers. The additional capabilities (image clipboard, Wayland native) are not needed for this application.

**Pure Wayland Limitation**: On systems running pure Wayland without XWayland (uncommon for interactive terminal users), `atotto/clipboard` will not function. This is an acceptable limitation for an initial version. If user demand warrants it, a Wayland-native fallback using `wl-paste` via `exec.Command` can be added with a `WAYLAND_DISPLAY` environment variable check.

---

## Dependency Summary

| Package | Purpose | CGO | Source |
|---------|---------|-----|--------|
| `charmbracelet/bubbletea` | TUI runtime (MVU) | No | github.com/charmbracelet/bubbletea |
| `charmbracelet/bubbles` | Pre-built TUI components (list, textinput) | No | github.com/charmbracelet/bubbles |
| `charmbracelet/lipgloss` | TUI layout and styling | No | github.com/charmbracelet/lipgloss |
| `modernc.org/sqlite` | Embedded SQLite storage | No | modernc.org/sqlite |
| `PuerkitoBio/goquery` | HTML meta tag extraction (CSS selectors) | No | github.com/PuerkitoBio/goquery |
| `golang.org/x/net/html/charset` | Non-UTF-8 HTML charset detection | No | golang.org/x/net |
| `sahilm/fuzzy` | Fuzzy search scoring | No | github.com/sahilm/fuzzy |
| `atotto/clipboard` | Clipboard read (Ctrl+P paste) | No | github.com/atotto/clipboard |

All dependencies are pure Go (zero CGO). The final binary will be a single statically linked executable.
