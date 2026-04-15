package fetcher

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// --- extractTitle unit tests ---

func makeDoc(t *testing.T, html string) *goquery.Document {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("makeDoc: %v", err)
	}
	return doc
}

func TestExtractTitle_PrefersTitleOverOgTitle(t *testing.T) {
	doc := makeDoc(t, `<html><head>
		<title>Real Page Title</title>
		<meta property="og:title" content="OG Title">
	</head></html>`)
	got := extractTitle(doc, "fallback.com")
	if got != "Real Page Title" {
		t.Errorf("expected %q, got %q", "Real Page Title", got)
	}
}

func TestExtractTitle_OgTitleWhenNoTitleTag(t *testing.T) {
	doc := makeDoc(t, `<html><head>
		<meta property="og:title" content="OG Only Title">
	</head></html>`)
	got := extractTitle(doc, "fallback.com")
	if got != "OG Only Title" {
		t.Errorf("expected %q, got %q", "OG Only Title", got)
	}
}

func TestExtractTitle_FallbackWhenNeither(t *testing.T) {
	doc := makeDoc(t, `<html><head></head><body></body></html>`)
	got := extractTitle(doc, "fallback.com")
	if got != "fallback.com" {
		t.Errorf("expected %q, got %q", "fallback.com", got)
	}
}

func TestExtractTitle_TrimsWhitespace(t *testing.T) {
	doc := makeDoc(t, `<html><head><title>  Spaced Title  </title></head></html>`)
	got := extractTitle(doc, "fallback.com")
	if got != "Spaced Title" {
		t.Errorf("expected %q, got %q", "Spaced Title", got)
	}
}

func TestExtractTitle_EmptyTitleFallsBackToOg(t *testing.T) {
	// A <title> tag that contains only whitespace should be ignored.
	doc := makeDoc(t, `<html><head>
		<title>   </title>
		<meta property="og:title" content="OG Fallback">
	</head></html>`)
	got := extractTitle(doc, "fallback.com")
	if got != "OG Fallback" {
		t.Errorf("expected %q, got %q", "OG Fallback", got)
	}
}

// --- Fetch integration tests against a local test server ---

// youtubeStyleHTML mimics the static HTML that YouTube returns to a browser UA:
// it has a descriptive <title> tag and no og:title in the served HTML.
const youtubeStyleHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Never Gonna Give You Up - Rick Astley - YouTube</title>
</head>
<body></body>
</html>`

func TestFetch_YouTubeStylePage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify browser-like headers are sent.
		if ua := r.Header.Get("User-Agent"); !strings.Contains(ua, "Mozilla") {
			t.Errorf("expected browser-like User-Agent, got %q", ua)
		}
		if r.Header.Get("Accept") == "" {
			t.Error("expected Accept header to be set")
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(youtubeStyleHTML))
	}))
	defer srv.Close()

	title, _, err := Fetch(srv.URL + "/watch?v=dQw4w9WgXcQ")
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}
	if title == "localhost" || title == "" {
		t.Errorf("expected video title, got fallback %q", title)
	}
	if !strings.Contains(title, "Rick Astley") {
		t.Errorf("expected title to contain video info, got %q", title)
	}
}

func TestFetch_Non2xxReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`<html><head><title>403 Forbidden</title></head></html>`))
	}))
	defer srv.Close()

	title, _, err := Fetch(srv.URL)
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
	if title != "127.0.0.1" && !strings.Contains(title, "localhost") && title != "::1" {
		// title should be the fallback hostname, not "403 Forbidden"
		if strings.Contains(title, "403") {
			t.Errorf("title should be fallback hostname, got %q", title)
		}
	}
}

func TestFetch_CloudflareChallengeReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("cf-mitigated", "challenge")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`<html><head><title>Just a moment...</title></head></html>`))
	}))
	defer srv.Close()

	title, _, err := Fetch(srv.URL)
	if err == nil {
		t.Fatal("expected error for Cloudflare challenge, got nil")
	}
	if strings.Contains(title, "Just a moment") {
		t.Errorf("title should be fallback hostname, got challenge page title %q", title)
	}
}

func TestFetch_RegularPage(t *testing.T) {
	const html = `<html><head><title>My Blog Post</title></head><body></body></html>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}))
	defer srv.Close()

	title, _, err := Fetch(srv.URL)
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}
	if title != "My Blog Post" {
		t.Errorf("expected %q, got %q", "My Blog Post", title)
	}
}
