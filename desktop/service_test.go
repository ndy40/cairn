package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/ndy40/cairn/internal/store"
)

// openTestStore creates a Store backed by a temporary SQLite database.
func openTestStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "bookmarks.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func mustInsert(t *testing.T, s *store.Store, url, title string, tags ...string) *store.Bookmark {
	t.Helper()
	b, err := s.Insert(url, title, "", tags)
	if err != nil {
		t.Fatalf("insert %s: %v", url, err)
	}
	return b
}

func TestSearchEmptyQueryListsPinnedFirst(t *testing.T) {
	s := openTestStore(t)
	mustInsert(t, s, "https://a.example.com", "Alpha")
	pinned := mustInsert(t, s, "https://b.example.com", "Bravo")
	mustInsert(t, s, "https://c.example.com", "Charlie")
	if err := s.SetPermanent(pinned.ID, true); err != nil {
		t.Fatalf("pin: %v", err)
	}

	got, err := NewBookmarkService(s, nil, nil).Search("")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("got %d results, want 3", len(got))
	}
	if got[0].ID != pinned.ID || !got[0].Pinned {
		t.Errorf("first result = %+v, want pinned bookmark %d", got[0], pinned.ID)
	}
}

func TestSearchMatchesTitleAndTags(t *testing.T) {
	s := openTestStore(t)
	mustInsert(t, s, "https://go.dev", "The Go Programming Language", "golang")
	mustInsert(t, s, "https://rust-lang.org", "Rust", "systems")

	svc := NewBookmarkService(s, nil, nil)

	got, err := svc.Search("programming")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(got) == 0 || got[0].URL != "https://go.dev" {
		t.Errorf("search by title: got %+v, want go.dev first", got)
	}

	got, err = svc.Search("systems")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(got) == 0 || got[0].URL != "https://rust-lang.org" {
		t.Errorf("search by tag: got %+v, want rust-lang.org first", got)
	}
}

func TestSearchListsEveryBookmark(t *testing.T) {
	s := openTestStore(t)
	for i := range 101 {
		mustInsert(t, s, fmt.Sprintf("https://example.com/%d", i), fmt.Sprintf("Page %d", i))
	}

	got, err := NewBookmarkService(s, nil, nil).Search("")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(got) != 101 {
		t.Errorf("got %d results, want all 101", len(got))
	}
}

func TestSearchIncludesAddedDate(t *testing.T) {
	s := openTestStore(t)
	b := mustInsert(t, s, "https://example.com", "Example")

	got, err := NewBookmarkService(s, nil, nil).Search("")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	added, err := time.Parse(time.RFC3339, got[0].AddedAt)
	if err != nil {
		t.Fatalf("AddedAt %q is not RFC 3339: %v", got[0].AddedAt, err)
	}
	if !added.Equal(b.CreatedAt.Truncate(time.Second)) {
		t.Errorf("AddedAt = %v, want %v", added, b.CreatedAt)
	}
}

func TestSearchReturnsEmptyTagSlice(t *testing.T) {
	s := openTestStore(t)
	mustInsert(t, s, "https://example.com", "Example")

	got, err := NewBookmarkService(s, nil, nil).Search("")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if got[0].Tags == nil {
		t.Error("Tags is nil; frontend expects an array")
	}
}

func TestOpenUsesBookmarkURL(t *testing.T) {
	s := openTestStore(t)
	b := mustInsert(t, s, "https://example.com/page", "Example")

	var opened string
	svc := NewBookmarkService(s, func(url string) error {
		opened = url
		return nil
	}, nil)

	if err := svc.Open(b.ID); err != nil {
		t.Fatalf("open: %v", err)
	}
	if opened != "https://example.com/page" {
		t.Errorf("opened %q, want https://example.com/page", opened)
	}
}

func TestOpenUnknownID(t *testing.T) {
	s := openTestStore(t)
	svc := NewBookmarkService(s, func(string) error {
		t.Fatal("openURL should not be called")
		return nil
	}, nil)

	err := svc.Open(999)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

// newAddService returns a service whose fetcher returns the given title, or an
// error when fetchErr is non-nil, and records the URLs it was asked to fetch.
func newAddService(s *store.Store, title string, fetchErr error) (*BookmarkService, *[]string) {
	svc := NewBookmarkService(s, nil, nil)
	var fetched []string
	svc.fetch = func(rawURL string) (string, string, error) {
		fetched = append(fetched, rawURL)
		if fetchErr != nil {
			return "example.com", "", fetchErr
		}
		return title, "a description", nil
	}
	return svc, &fetched
}

func TestAddSavesFetchedTitleAndTags(t *testing.T) {
	s := openTestStore(t)
	svc, _ := newAddService(s, "Example Page", nil)

	res, err := svc.Add("  https://example.com/page  ", "Go, reading ,go")
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if !res.TitleFetched || res.Bookmark.Title != "Example Page" {
		t.Errorf("result = %+v, want fetched title", res)
	}

	saved, err := s.GetByID(res.Bookmark.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if saved.URL != "https://example.com/page" {
		t.Errorf("URL = %q, want trimmed URL", saved.URL)
	}
	if saved.Description != "a description" {
		t.Errorf("Description = %q, want fetched description", saved.Description)
	}
	if len(saved.Tags) != 2 {
		t.Errorf("Tags = %v, want 2 normalised tags", saved.Tags)
	}
}

func TestAddDefaultsSchemeToHTTPS(t *testing.T) {
	s := openTestStore(t)
	svc, fetched := newAddService(s, "Go", nil)

	res, err := svc.Add("go.dev/doc", "")
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if res.Bookmark.URL != "https://go.dev/doc" {
		t.Errorf("URL = %q, want https://go.dev/doc", res.Bookmark.URL)
	}
	if len(*fetched) != 1 || (*fetched)[0] != "https://go.dev/doc" {
		t.Errorf("fetched %v, want the normalised URL", *fetched)
	}
}

func TestAddFetchFailureStillSaves(t *testing.T) {
	s := openTestStore(t)
	svc, _ := newAddService(s, "", errors.New("timeout"))

	res, err := svc.Add("https://example.com", "")
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if res.TitleFetched {
		t.Error("TitleFetched = true, want false after fetch failure")
	}
	if res.Bookmark.Title != "example.com" {
		t.Errorf("Title = %q, want hostname fallback", res.Bookmark.Title)
	}
}

func TestAddRejectsDuplicateWithoutFetching(t *testing.T) {
	s := openTestStore(t)
	mustInsert(t, s, "https://example.com", "Example")
	svc, fetched := newAddService(s, "Example", nil)

	_, err := svc.Add("https://example.com", "")
	if !errors.Is(err, ErrAlreadyBookmarked) {
		t.Errorf("err = %v, want ErrAlreadyBookmarked", err)
	}
	if len(*fetched) != 0 {
		t.Errorf("fetched %v for a duplicate, want no network call", *fetched)
	}
}

func TestAddRejectsInvalidURLs(t *testing.T) {
	s := openTestStore(t)
	svc, fetched := newAddService(s, "x", nil)

	for _, in := range []string{"", "   ", "ftp://example.com", "javascript:alert(1)", "https://", "not a url"} {
		if _, err := svc.Add(in, ""); !errors.Is(err, ErrInvalidURL) {
			t.Errorf("Add(%q) err = %v, want ErrInvalidURL", in, err)
		}
	}
	if len(*fetched) != 0 {
		t.Errorf("fetched %v for invalid input, want no network call", *fetched)
	}
}

func TestClipboardURL(t *testing.T) {
	s := openTestStore(t)
	cases := map[string]string{
		"https://example.com/a": "https://example.com/a",
		"  http://example.com ": "http://example.com",
		"example.com":           "",
		"some copied text":      "",
		"":                      "",
	}
	for clip, want := range cases {
		svc := NewBookmarkService(s, nil, func() string { return clip })
		if got := svc.ClipboardURL(); got != want {
			t.Errorf("ClipboardURL() with %q = %q, want %q", clip, got, want)
		}
	}
}

func TestAddCallsAfterWriteOnlyOnSave(t *testing.T) {
	s := openTestStore(t)
	svc, _ := newAddService(s, "Example", nil)
	writes := 0
	svc.afterWrite = func() { writes++ }

	if _, err := svc.Add("https://example.com", ""); err != nil {
		t.Fatalf("add: %v", err)
	}
	_, _ = svc.Add("https://example.com", "") // duplicate
	_, _ = svc.Add("not a url", "")           // invalid

	if writes != 1 {
		t.Errorf("afterWrite called %d times, want 1 (only the successful save)", writes)
	}
}

func TestFindByURL(t *testing.T) {
	s := openTestStore(t)
	saved := mustInsert(t, s, "https://go.dev/doc", "Go docs")
	mustInsert(t, s, "https://example.com", "Example")
	svc := NewBookmarkService(s, nil, nil)

	for _, in := range []string{"https://go.dev/doc", "  go.dev/doc "} {
		got, err := svc.FindByURL(in)
		if err != nil {
			t.Fatalf("FindByURL(%q): %v", in, err)
		}
		if got == nil || got.ID != saved.ID {
			t.Errorf("FindByURL(%q) = %+v, want bookmark %d", in, got, saved.ID)
		}
	}

	got, err := svc.FindByURL("https://go.dev/blog")
	if err != nil || got != nil {
		t.Errorf("FindByURL(unknown) = %+v, %v; want nil, nil", got, err)
	}

	if _, err := svc.FindByURL("ftp://go.dev"); !errors.Is(err, ErrInvalidURL) {
		t.Errorf("FindByURL(ftp) err = %v, want ErrInvalidURL", err)
	}
}

func TestSetPinnedPinsAndNotifies(t *testing.T) {
	s := openTestStore(t)
	b := mustInsert(t, s, "https://go.dev", "Go")

	svc := NewBookmarkService(s, nil, nil)
	var reasons []string
	svc.changed = func(reason string) { reasons = append(reasons, reason) }

	if err := svc.SetPinned(b.ID, true); err != nil {
		t.Fatalf("SetPinned(true): %v", err)
	}
	if got, _ := s.GetByID(b.ID); !got.IsPermanent {
		t.Error("bookmark not pinned")
	}
	if err := svc.SetPinned(b.ID, false); err != nil {
		t.Fatalf("SetPinned(false): %v", err)
	}
	if got, _ := s.GetByID(b.ID); got.IsPermanent {
		t.Error("bookmark still pinned")
	}
	if len(reasons) != 2 || reasons[0] != "pin" {
		t.Errorf("changed reasons = %v, want [pin pin]", reasons)
	}
}

func TestCopyLinkWritesURL(t *testing.T) {
	s := openTestStore(t)
	b := mustInsert(t, s, "https://go.dev/doc", "Go docs")

	svc := NewBookmarkService(s, nil, nil)
	var copied string
	svc.writeClipboard = func(text string) bool {
		copied = text
		return true
	}
	if err := svc.CopyLink(b.ID); err != nil {
		t.Fatalf("CopyLink: %v", err)
	}
	if copied != "https://go.dev/doc" {
		t.Errorf("copied %q, want the bookmark URL", copied)
	}

	svc.writeClipboard = func(string) bool { return false }
	if err := svc.CopyLink(b.ID); err == nil {
		t.Error("CopyLink succeeded although the clipboard refused the text")
	}
	if err := svc.CopyLink(b.ID + 100); err == nil {
		t.Error("CopyLink succeeded for a missing bookmark")
	}
}

func TestSyncNowWithoutSyncer(t *testing.T) {
	svc := NewBookmarkService(openTestStore(t), nil, nil)
	if err := svc.SyncNow(); !errors.Is(err, errSyncNotSetUp) {
		t.Errorf("SyncNow() = %v, want errSyncNotSetUp", err)
	}
}
