package fetcher

import (
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
	userAgent    = "BookmarkManager/1.0"
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

	resp, err := client.Do(req)
	if err != nil {
		return fallback, "", err
	}
	defer resp.Body.Close()

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
	// og:title takes priority.
	if t, exists := doc.Find(`meta[property="og:title"]`).Attr("content"); exists && strings.TrimSpace(t) != "" {
		return strings.TrimSpace(t)
	}
	// Standard <title> tag.
	if t := strings.TrimSpace(doc.Find("title").First().Text()); t != "" {
		return t
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
