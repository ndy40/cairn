package sync

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/ndy40/cairn/internal/store"
	"github.com/ndy40/cairn/internal/sync/backend"
)

// memBackend is an in-memory SyncBackend shared between test devices.
type memBackend struct {
	files map[string][]byte
}

func newMemBackend() *memBackend {
	return &memBackend{files: make(map[string][]byte)}
}

func (m *memBackend) Upload(data []byte, remotePath string) error {
	m.files[remotePath] = append([]byte(nil), data...)
	return nil
}

func (m *memBackend) Download(remotePath string) ([]byte, error) {
	data, ok := m.files[remotePath]
	if !ok {
		return nil, backend.ErrNotFound
	}
	return data, nil
}

func (m *memBackend) Exists(remotePath string) (bool, error) {
	_, ok := m.files[remotePath]
	return ok, nil
}

// newTestDevice returns an engine with its own store and config, sharing b.
func newTestDevice(t *testing.T, deviceID string, b backend.SyncBackend) *Engine {
	t.Helper()
	dir := t.TempDir()
	s, err := store.Open(filepath.Join(dir, "bookmarks.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	cfg := &SyncConfig{Backend: "dropbox", DeviceID: deviceID}
	return NewEngine(s, b, cfg, filepath.Join(dir, "sync.json"))
}

func cloudEntry(t *testing.T, b *memBackend, url string) *BookmarkEntry {
	t.Helper()
	record, err := UnmarshalSyncRecord(b.files[remotePath])
	if err != nil {
		t.Fatalf("unmarshal cloud snapshot: %v", err)
	}
	for i := range record.Bookmarks {
		if record.Bookmarks[i].URL == url {
			return &record.Bookmarks[i]
		}
	}
	return nil
}

// TestPullPreservesMetadata covers a two-device round trip: the pulled copy
// keeps its UUID, creation date, and pin, and a later push from the second
// device does not rewrite the creation date.
func TestPullPreservesMetadata(t *testing.T) {
	cloud := newMemBackend()
	a := newTestDevice(t, "device-a", cloud)
	b := newTestDevice(t, "device-b", cloud)

	const url = "https://example.com/article"
	original := time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)

	bm, err := a.Store.Insert(url, "Article", "desc", []string{"read"})
	if err != nil {
		t.Fatalf("insert on A: %v", err)
	}
	if err := a.Store.SetCreatedAt(bm.UUID, original); err != nil {
		t.Fatalf("backdate on A: %v", err)
	}
	if err := a.Store.SetPermanent(bm.ID, true); err != nil {
		t.Fatalf("pin on A: %v", err)
	}
	if err := a.Push(); err != nil {
		t.Fatalf("push A: %v", err)
	}

	if _, err := b.Pull(); err != nil {
		t.Fatalf("pull B: %v", err)
	}
	got, err := b.Store.GetByUUID(bm.UUID)
	if err != nil {
		t.Fatalf("bookmark missing on B by original UUID: %v", err)
	}
	if !got.CreatedAt.Equal(original) {
		t.Errorf("B CreatedAt = %v, want %v", got.CreatedAt, original)
	}
	if !got.IsPermanent {
		t.Error("B lost the pin")
	}
	if n, _ := b.Store.PendingChangeCount(); n != 0 {
		t.Errorf("B pending changes after pull = %d, want 0", n)
	}

	// B pushing must not overwrite the cloud's creation date with its pull time.
	if err := b.Push(); err != nil {
		t.Fatalf("push B: %v", err)
	}
	entry := cloudEntry(t, cloud, url)
	if entry == nil || !entry.CreatedAt.Equal(original) || entry.UUID != bm.UUID {
		t.Fatalf("cloud entry after B push = %+v, want uuid %s created %v", entry, bm.UUID, original)
	}

}

// TestPullHealsLaterCreatedAt checks that a device holding a later creation
// date for a bookmark adopts the cloud's earlier one and uploads it back.
func TestPullHealsLaterCreatedAt(t *testing.T) {
	cloud := newMemBackend()
	a := newTestDevice(t, "device-a", cloud)
	b := newTestDevice(t, "device-b", cloud)

	const url = "https://example.com/healed"
	original := time.Date(2022, 3, 4, 5, 6, 7, 0, time.UTC)

	bmA, err := a.Store.Insert(url, "Healed", "", nil)
	if err != nil {
		t.Fatalf("insert on A: %v", err)
	}
	if err := a.Store.SetCreatedAt(bmA.UUID, original); err != nil {
		t.Fatalf("backdate on A: %v", err)
	}
	if err := a.Push(); err != nil {
		t.Fatalf("push A: %v", err)
	}

	// B already has the same URL under its own UUID with a later date, as a
	// device that pulled before the fix would.
	bmB, err := b.Store.Insert(url, "Healed", "", nil)
	if err != nil {
		t.Fatalf("insert on B: %v", err)
	}
	if err := b.Store.ClearPendingChanges(); err != nil {
		t.Fatalf("clear pending on B: %v", err)
	}
	if _, err := b.Pull(); err != nil {
		t.Fatalf("pull B: %v", err)
	}

	got, err := b.Store.GetByUUID(bmB.UUID)
	if err != nil {
		t.Fatalf("GetByUUID on B: %v", err)
	}
	if !got.CreatedAt.Equal(original) {
		t.Errorf("B CreatedAt = %v, want %v", got.CreatedAt, original)
	}
	if !got.UpdatedAt.Equal(bmB.UpdatedAt.Truncate(time.Second)) {
		t.Errorf("B UpdatedAt changed: got %v, want %v", got.UpdatedAt, bmB.UpdatedAt)
	}
	if n, _ := b.Store.PendingChangeCount(); n != 0 {
		t.Errorf("B pending changes after heal = %d, want 0", n)
	}
	if entry := cloudEntry(t, cloud, url); entry == nil || !entry.CreatedAt.Equal(original) {
		t.Errorf("cloud entry after B pull = %+v, want created %v", entry, original)
	}
}
