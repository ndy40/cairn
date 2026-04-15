# Implementation Plan: Fix Fetcher Title Extraction

## Feature: 015-fix-fetcher-title

---

## Tech Stack

- **Language**: Go 1.25.0
- **Build**: `CGO_ENABLED=0` (pure-Go, static binary)
- **Libraries**:
  - `github.com/PuerkitoBio/goquery` — HTML parsing (existing)
  - `golang.org/x/net/html/charset` — charset detection (existing)
  - `net/http` stdlib — HTTP client (existing)
- **No new dependencies required**

---

## Project Structure (affected files)

```
internal/
  fetcher/
    fetcher.go        ← extractTitle priority swap + header improvements
    fetcher_test.go   ← unit tests for extractTitle and Fetch (existing or new)
```

---

## Implementation Approach

### US1 — Priority swap (`<title>` before `og:title`)

Change `extractTitle` to check the standard `<title>` tag first, then fall back to `og:title`, then the hostname fallback. This is a two-line reorder inside `extractTitle`.

```go
func extractTitle(doc *goquery.Document, fallback string) string {
    // Standard <title> tag takes priority.
    if t := strings.TrimSpace(doc.Find("title").First().Text()); t != "" {
        return t
    }
    // og:title as fallback.
    if t, exists := doc.Find(`meta[property="og:title"]`).Attr("content"); exists && strings.TrimSpace(t) != "" {
        return strings.TrimSpace(t)
    }
    return fallback
}
```

### US2 — Improved request headers for bot-gated pages

Update `Fetch` to send a more browser-like `User-Agent` string and add `Accept` / `Accept-Language` headers. YouTube and similar sites serve richer static HTML (including the `<title>` tag with the video title) when the request looks like a real browser.

```go
const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
req.Header.Set("Accept-Language", "en-US,en;q=0.5")
```

No new dependencies or schema changes are required.

---

## Testing Strategy

- Unit tests for `extractTitle` covering all three priority branches.
- Manual smoke test: `bm add https://www.youtube.com/watch?v=<id>` and verify the stored title matches the video title.
- `go build ./... && go vet ./...` must pass.
