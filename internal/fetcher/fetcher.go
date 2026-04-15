package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

const (
	maxBodyBytes = 512 * 1024 // 512 KB
	timeout      = 8 * time.Second
	// Realistic browser UA so that sites like YouTube serve full HTML (including <title>)
	// rather than a minimal bot-response.
	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
)

var client = &http.Client{Timeout: timeout}

// Fetch retrieves the web page at rawURL and extracts its title and
// meta description. On any error it returns the URL hostname as a
// fallback title and an empty description alongside the error.
func Fetch(rawURL string) (title, description string, err error) {
	fallback := hostname(rawURL)

	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return fallback, "", err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return fallback, "", err
	}
	defer func() { _ = resp.Body.Close() }()

	// Cloudflare (and similar) bot challenges return a 403 with cf-mitigated: challenge.
	// Detect this before reading the body so we don't record "Just a moment..." as a title.
	if resp.Header.Get("cf-mitigated") == "challenge" {
		return fallback, "", fmt.Errorf("fetch %s: blocked by bot-protection challenge", rawURL)
	}

	// Treat non-2xx responses as errors; parsing an error page body would
	// produce a misleading title (e.g. "403 Forbidden", "Just a moment...").
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fallback, "", fmt.Errorf("fetch %s: unexpected status %s", rawURL, resp.Status)
	}

	// Detect charset and wrap reader so goquery sees UTF-8.
	contentType := resp.Header.Get("Content-Type")
	utf8Reader, err := charset.NewReader(io.LimitReader(resp.Body, maxBodyBytes), contentType)
	if err != nil {
		return fallback, "", err
	}

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		return fallback, "", err
	}

	title = extractTitle(doc, fallback)
	description = extractDescription(doc)
	return title, description, nil
}

func extractTitle(doc *goquery.Document, fallback string) string {
	// Standard <title> tag takes priority — it matches what the user sees in the browser tab.
	if t := strings.TrimSpace(doc.Find("title").First().Text()); t != "" {
		return t
	}
	// og:title as fallback when no <title> is present.
	if t, exists := doc.Find(`meta[property="og:title"]`).Attr("content"); exists && strings.TrimSpace(t) != "" {
		return strings.TrimSpace(t)
	}
	return fallback
}

func extractDescription(doc *goquery.Document) string {
	// og:description takes priority.
	if d, exists := doc.Find(`meta[property="og:description"]`).Attr("content"); exists && strings.TrimSpace(d) != "" {
		return strings.TrimSpace(d)
	}
	// Standard meta description.
	if d, exists := doc.Find(`meta[name="description"]`).Attr("content"); exists && strings.TrimSpace(d) != "" {
		return strings.TrimSpace(d)
	}
	return ""
}

// hostname extracts the hostname from a URL string for use as a fallback title.
func hostname(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	h := u.Hostname()
	if h == "" {
		return rawURL
	}
	// Strip www. prefix.
	h = strings.TrimPrefix(h, "www.")
	return h
}
