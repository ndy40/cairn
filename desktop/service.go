package main

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/ndy40/cairn/internal/fetcher"
	"github.com/ndy40/cairn/internal/search"
	"github.com/ndy40/cairn/internal/store"
)

// ErrInvalidURL is returned by Add when the input is not an http(s) URL.
var ErrInvalidURL = errors.New("enter a valid http(s) URL")

// ErrAlreadyBookmarked is returned by Add when the URL is already saved.
var ErrAlreadyBookmarked = errors.New("already bookmarked")

// maxResults caps how many rows the popup renders for a single query.
const maxResults = 50

// Bookmark is the frontend-facing view of a store.Bookmark.
type Bookmark struct {
	ID     int64    `json:"id"`
	Title  string   `json:"title"`
	URL    string   `json:"url"`
	Domain string   `json:"domain"`
	Tags   []string `json:"tags"`
	Pinned bool     `json:"pinned"`
}

// AddResult reports the outcome of Add.
type AddResult struct {
	Bookmark Bookmark `json:"bookmark"`
	// TitleFetched is false when the page could not be fetched and the
	// hostname was used as the title, mirroring `cairn add` exit code 2.
	TitleFetched bool `json:"titleFetched"`
}

// BookmarkService exposes bookmark operations to the frontend.
// Every call reads from the store so changes made by the CLI/TUI show up immediately.
type BookmarkService struct {
	store         *store.Store
	openURL       func(url string) error
	readClipboard func() string
	fetch         func(rawURL string) (title, description string, err error)
	// afterWrite, when set, is called after a bookmark is saved (used to push sync).
	afterWrite func()
	// writeClipboard puts text on the system clipboard, reporting success.
	writeClipboard func(text string) bool
	// changed, when set, tells the frontend bookmarks changed (e.g. "pin").
	changed func(reason string)
	// syncNow, when set, runs a full sync for SyncNow.
	syncNow func() error
}

// NewBookmarkService returns a service backed by s that opens links with
// openURL and reads the system clipboard with readClipboard.
func NewBookmarkService(s *store.Store, openURL func(url string) error, readClipboard func() string) *BookmarkService {
	return &BookmarkService{
		store:         s,
		openURL:       openURL,
		readClipboard: readClipboard,
		fetch:         fetcher.Fetch,
	}
}

// Search returns bookmarks matching query using the same ranking as the TUI.
// An empty query returns pinned bookmarks first, then the most recently added.
func (b *BookmarkService) Search(query string) ([]Bookmark, error) {
	all, err := b.store.List()
	if err != nil {
		return nil, fmt.Errorf("list bookmarks: %w", err)
	}

	var results []*store.Bookmark
	if query == "" {
		results = all
		sort.SliceStable(results, func(i, j int) bool {
			return results[i].IsPermanent && !results[j].IsPermanent
		})
	} else {
		results = search.Search(query, all)
	}

	if len(results) > maxResults {
		results = results[:maxResults]
	}
	return toDTOs(results), nil
}

// Open launches the bookmark's URL in the default browser.
func (b *BookmarkService) Open(id int64) error {
	bm, err := b.store.GetByID(id)
	if err != nil {
		return fmt.Errorf("get bookmark %d: %w", id, err)
	}
	if err := b.openURL(bm.URL); err != nil {
		return fmt.Errorf("open %s: %w", bm.URL, err)
	}
	return nil
}

// Add saves rawURL with comma-separated tags, fetching the page title and
// description the same way `cairn add` does. A fetch failure is not an error:
// the bookmark is saved with the hostname as its title.
func (b *BookmarkService) Add(rawURL, tags string) (*AddResult, error) {
	normalised, err := normaliseURL(rawURL)
	if err != nil {
		return nil, err
	}

	exists, err := b.store.ExistsByURL(normalised)
	if err != nil {
		return nil, fmt.Errorf("check existing bookmark: %w", err)
	}
	if exists {
		return nil, ErrAlreadyBookmarked
	}

	title, description, fetchErr := b.fetch(normalised)
	bm, err := b.store.Insert(normalised, title, description, store.NormaliseTagsFromString(tags))
	if errors.Is(err, store.ErrDuplicate) {
		return nil, ErrAlreadyBookmarked
	}
	if err != nil {
		return nil, fmt.Errorf("save bookmark: %w", err)
	}

	if b.afterWrite != nil {
		b.afterWrite()
	}
	return &AddResult{Bookmark: toDTO(bm), TitleFetched: fetchErr == nil}, nil
}

// FindByURL normalises rawURL like Add and returns the saved bookmark with
// that URL, or nil when there is none, so a duplicate add can select it.
// Archived bookmarks are not listed, so they are reported as not found.
func (b *BookmarkService) FindByURL(rawURL string) (*Bookmark, error) {
	normalised, err := normaliseURL(rawURL)
	if err != nil {
		return nil, err
	}
	all, err := b.store.List()
	if err != nil {
		return nil, fmt.Errorf("list bookmarks: %w", err)
	}
	for _, bm := range all {
		if bm.URL == normalised {
			dto := toDTO(bm)
			return &dto, nil
		}
	}
	return nil, nil
}

// SetPinned pins or unpins a bookmark, like `cairn pin`. Like the CLI it
// does not queue a sync change.
func (b *BookmarkService) SetPinned(id int64, pinned bool) error {
	if err := b.store.SetPermanent(id, pinned); err != nil {
		return fmt.Errorf("pin bookmark %d: %w", id, err)
	}
	if b.changed != nil {
		b.changed("pin")
	}
	return nil
}

// CopyLink puts the bookmark's URL on the clipboard.
func (b *BookmarkService) CopyLink(id int64) error {
	bm, err := b.store.GetByID(id)
	if err != nil {
		return fmt.Errorf("get bookmark %d: %w", id, err)
	}
	if b.writeClipboard == nil || !b.writeClipboard(bm.URL) {
		return errors.New("could not copy the link")
	}
	return nil
}

// SyncNow pulls then pushes, ignoring the pull throttle, and returns once
// done. The error says why nothing synced or what went wrong.
func (b *BookmarkService) SyncNow() error {
	if b.syncNow == nil {
		return errSyncNotSetUp
	}
	return b.syncNow()
}

// ClipboardURL returns the clipboard contents if they are an http(s) URL, or "".
func (b *BookmarkService) ClipboardURL() string {
	text := strings.TrimSpace(b.readClipboard())
	if !strings.HasPrefix(text, "http://") && !strings.HasPrefix(text, "https://") {
		return ""
	}
	normalised, err := normaliseURL(text)
	if err != nil {
		return ""
	}
	return normalised
}

// normaliseURL trims input, defaults a missing scheme to https and rejects
// anything that is not an http(s) URL with a host. The URL is otherwise kept
// as typed (not re-encoded) so duplicate detection matches what `cairn add` stores.
func normaliseURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ErrInvalidURL
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || strings.ContainsAny(u.Host, " \t") {
		return "", ErrInvalidURL
	}
	return raw, nil
}

func toDTOs(bookmarks []*store.Bookmark) []Bookmark {
	out := make([]Bookmark, 0, len(bookmarks))
	for _, bm := range bookmarks {
		out = append(out, toDTO(bm))
	}
	return out
}

func toDTO(bm *store.Bookmark) Bookmark {
	tags := bm.Tags
	if tags == nil {
		tags = []string{}
	}
	return Bookmark{
		ID:     bm.ID,
		Title:  bm.Title,
		URL:    bm.URL,
		Domain: bm.Domain,
		Tags:   tags,
		Pinned: bm.IsPermanent,
	}
}
