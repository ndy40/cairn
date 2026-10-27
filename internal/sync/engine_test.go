package sync

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/ndy40/cairn/internal/store"
	"github.com/ndy40/cairn/internal/sync/backend"
)

// memBackend is an in-memory SyncBackend shared between test devices.
type memBackend struct {
	files      map[string][]byte
	failUpload bool
}

func newMemBackend() *memBackend {
	return &memBackend{files: make(map[string][]byte)}
}

func (m *memBackend) Upload(data []byte, remotePath string) error {
	if m.failUpload {
		return errors.New("upload failed")
	}
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

// TestPullPreservesMetadataAndPropagatesDeletes covers a two-device round trip:
// the pulled copy keeps its UUID, creation date, and pin, a later push from the
// second device does not rewrite the creation date, and a delete on the first
// device reaches the second.
func TestPullPreservesMetadataAndPropagatesDeletes(t *testing.T) {
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

	// Delete on A propagates to B.
	if _, err := a.Pull(); err != nil {
		t.Fatalf("pull A: %v", err)
	}
	if _, err := a.Store.DeleteByID(bm.ID); err != nil {
		t.Fatalf("delete on A: %v", err)
	}
	if err := a.Push(); err != nil {
		t.Fatalf("push A after delete: %v", err)
	}
	if _, err := b.Pull(); err != nil {
		t.Fatalf("pull B after delete: %v", err)
	}
	if _, err := b.Store.GetByUUID(bm.UUID); err != store.ErrNotFound {
		t.Errorf("bookmark still on B after delete: err = %v", err)
	}
	if cloudEntry(t, cloud, url) != nil {
		t.Error("bookmark still in cloud snapshot after B merged the delete")
	}
}

// TestPullHealsLaterCreatedAt checks that a device holding a later copy of a
// bookmark adopts the cloud's earlier creation date and UUID and uploads them
// back, after which a delete on the other device reaches it.
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

	got, err := b.Store.GetByUUID(bmA.UUID)
	if err != nil {
		t.Fatalf("B did not adopt A's UUID: %v", err)
	}
	if _, err := b.Store.GetByUUID(bmB.UUID); err != store.ErrNotFound {
		t.Errorf("B's old UUID still present: err = %v", err)
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
	if entry := cloudEntry(t, cloud, url); entry == nil || !entry.CreatedAt.Equal(original) || entry.UUID != bmA.UUID {
		t.Errorf("cloud entry after B pull = %+v, want uuid %s created %v", entry, bmA.UUID, original)
	}

	// With the UUIDs converged, a delete on A now reaches B.
	if _, err := a.Store.DeleteByID(bmA.ID); err != nil {
		t.Fatalf("delete on A: %v", err)
	}
	if err := a.Push(); err != nil {
		t.Fatalf("push A after delete: %v", err)
	}
	if _, err := b.Pull(); err != nil {
		t.Fatalf("pull B after delete: %v", err)
	}
	if n, _ := b.Store.Count(); n != 0 {
		t.Errorf("B has %d bookmarks after delete, want 0", n)
	}
}

// TestPullKeepsEarlierLocalUUID checks the other direction: when the local copy
// is the earliest, it keeps its UUID, uploads it, and the other device adopts it.
func TestPullKeepsEarlierLocalUUID(t *testing.T) {
	cloud := newMemBackend()
	a := newTestDevice(t, "device-a", cloud)
	b := newTestDevice(t, "device-b", cloud)

	const url = "https://example.com/older-on-b"
	original := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	bmA, err := a.Store.Insert(url, "Older on B", "", nil)
	if err != nil {
		t.Fatalf("insert on A: %v", err)
	}
	if err := a.Push(); err != nil {
		t.Fatalf("push A: %v", err)
	}

	bmB, err := b.Store.Insert(url, "Older on B", "", nil)
	if err != nil {
		t.Fatalf("insert on B: %v", err)
	}
	if err := b.Store.SetCreatedAt(bmB.UUID, original); err != nil {
		t.Fatalf("backdate on B: %v", err)
	}
	if _, err := b.Pull(); err != nil {
		t.Fatalf("pull B: %v", err)
	}
	if _, err := b.Store.GetByUUID(bmB.UUID); err != nil {
		t.Fatalf("B lost its earlier UUID: %v", err)
	}

	if _, err := a.Pull(); err != nil {
		t.Fatalf("pull A: %v", err)
	}
	got, err := a.Store.GetByUUID(bmB.UUID)
	if err != nil {
		t.Fatalf("A did not adopt B's UUID: %v", err)
	}
	if got.ID != bmA.ID || !got.CreatedAt.Equal(original) {
		t.Errorf("A row = id %d created %v, want id %d created %v", got.ID, got.CreatedAt, bmA.ID, original)
	}
}

func cloudRecord(t *testing.T, b *memBackend) *SyncRecord {
	t.Helper()
	record, err := UnmarshalSyncRecord(b.files[remotePath])
	if err != nil {
		t.Fatalf("unmarshal cloud snapshot: %v", err)
	}
	return record
}

// TestPushKeepsOtherDevicesBookmarks checks that a push from a device that has
// not pulled merges the cloud first instead of overwriting it.
func TestPushKeepsOtherDevicesBookmarks(t *testing.T) {
	cloud := newMemBackend()
	a := newTestDevice(t, "device-a", cloud)
	b := newTestDevice(t, "device-b", cloud)

	if _, err := a.Store.Insert("https://example.com/from-a", "A", "", nil); err != nil {
		t.Fatalf("insert on A: %v", err)
	}
	if err := a.Push(); err != nil {
		t.Fatalf("push A: %v", err)
	}
	if _, err := b.Store.Insert("https://example.com/from-b", "B", "", nil); err != nil {
		t.Fatalf("insert on B: %v", err)
	}
	if err := b.Push(); err != nil {
		t.Fatalf("push B: %v", err)
	}

	if cloudEntry(t, cloud, "https://example.com/from-a") == nil || cloudEntry(t, cloud, "https://example.com/from-b") == nil {
		t.Errorf("cloud lost a bookmark: %+v", cloudRecord(t, cloud).Bookmarks)
	}
	if ok, _ := b.Store.ExistsByURL("https://example.com/from-a"); !ok {
		t.Error("B did not merge A's bookmark during push")
	}
	if n, _ := b.Store.PendingChangeCount(); n != 0 {
		t.Errorf("B pending changes after push = %d, want 0", n)
	}
}

// TestFailedPushDoesNotResurrectDelete checks that a delete whose push failed is
// not undone by the next pull, and still reaches the other device afterwards.
func TestFailedPushDoesNotResurrectDelete(t *testing.T) {
	cloud := newMemBackend()
	a := newTestDevice(t, "device-a", cloud)
	b := newTestDevice(t, "device-b", cloud)

	const url = "https://example.com/doomed"
	bm, err := a.Store.Insert(url, "Doomed", "", nil)
	if err != nil {
		t.Fatalf("insert on A: %v", err)
	}
	if err := a.Push(); err != nil {
		t.Fatalf("push A: %v", err)
	}
	if _, err := b.Pull(); err != nil {
		t.Fatalf("pull B: %v", err)
	}

	if _, err := a.Store.DeleteByID(bm.ID); err != nil {
		t.Fatalf("delete on A: %v", err)
	}
	cloud.failUpload = true
	if err := a.Push(); err == nil {
		t.Fatal("expected push to fail")
	}
	cloud.failUpload = false

	if _, err := a.Pull(); err != nil {
		t.Fatalf("pull A: %v", err)
	}
	if ok, _ := a.Store.ExistsByURL(url); ok {
		t.Fatal("pull brought the deleted bookmark back on A")
	}

	if _, err := b.Pull(); err != nil {
		t.Fatalf("pull B: %v", err)
	}
	if ok, _ := b.Store.ExistsByURL(url); ok {
		t.Error("delete did not reach B")
	}
}

// TestPullAppliesRemoteEdits checks that a newer title, description, tags, and
// URL from another device are applied without queuing local changes, while the
// pin stays per-device.
func TestPullAppliesRemoteEdits(t *testing.T) {
	cloud := newMemBackend()
	a := newTestDevice(t, "device-a", cloud)
	b := newTestDevice(t, "device-b", cloud)

	old := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	seed := &store.Bookmark{
		UUID: "edit-me", URL: "https://example.com/before", Title: "Before",
		CreatedAt: old, UpdatedAt: old,
	}
	if err := a.Store.InsertSynced(seed); err != nil {
		t.Fatalf("seed A: %v", err)
	}
	if err := a.Push(); err != nil {
		t.Fatalf("push A: %v", err)
	}
	if _, err := b.Pull(); err != nil {
		t.Fatalf("pull B: %v", err)
	}
	onB, err := b.Store.GetByUUID("edit-me")
	if err != nil {
		t.Fatalf("GetByUUID on B: %v", err)
	}
	if err := b.Store.SetPermanent(onB.ID, true); err != nil {
		t.Fatalf("pin on B: %v", err)
	}

	onA, err := a.Store.GetByUUID("edit-me")
	if err != nil {
		t.Fatalf("GetByUUID on A: %v", err)
	}
	title, newURL, tags := "After", "https://example.com/after", []string{"edited"}
	if err := a.Store.UpdateFields(onA.ID, store.BookmarkPatch{Title: &title, URL: &newURL, Tags: &tags}); err != nil {
		t.Fatalf("edit on A: %v", err)
	}
	if err := a.Push(); err != nil {
		t.Fatalf("push A after edit: %v", err)
	}
	if _, err := b.Pull(); err != nil {
		t.Fatalf("pull B after edit: %v", err)
	}

	got, err := b.Store.GetByUUID("edit-me")
	if err != nil {
		t.Fatalf("GetByUUID on B after edit: %v", err)
	}
	if got.Title != title || got.URL != newURL || len(got.Tags) != 1 || got.Tags[0] != "edited" {
		t.Errorf("B row = %+v, want A's edit", got)
	}
	if !got.IsPermanent {
		t.Error("B's pin was overwritten by the remote edit")
	}
	if !got.CreatedAt.Equal(old) {
		t.Errorf("B CreatedAt = %v, want %v", got.CreatedAt, old)
	}
	if n, _ := b.Store.PendingChangeCount(); n != 0 {
		t.Errorf("B pending changes = %d, want 0", n)
	}
}

// TestUploadPrunesOldTombstones checks that tombstones past the TTL are dropped
// and undated ones are stamped instead of dropped.
func TestUploadPrunesOldTombstones(t *testing.T) {
	cloud := newMemBackend()
	a := newTestDevice(t, "device-a", cloud)

	now := time.Now().UTC()
	seed := NewSyncRecord("device-x")
	seed.Tombstones = []TombstoneEntry{
		{UUID: "expired", DeletedAt: now.Add(-tombstoneTTL - time.Hour)},
		{UUID: "recent", DeletedAt: now.Add(-time.Hour)},
		{UUID: "undated"},
	}
	data, err := seed.Marshal()
	if err != nil {
		t.Fatalf("marshal seed: %v", err)
	}
	cloud.files[remotePath] = data

	if _, err := a.Pull(); err != nil {
		t.Fatalf("pull A: %v", err)
	}

	got := make(map[string]TombstoneEntry)
	for _, ts := range cloudRecord(t, cloud).Tombstones {
		got[ts.UUID] = ts
	}
	if _, ok := got["expired"]; ok {
		t.Error("expired tombstone was kept")
	}
	if _, ok := got["recent"]; !ok {
		t.Error("recent tombstone was dropped")
	}
	if ts, ok := got["undated"]; !ok || ts.DeletedAt.IsZero() {
		t.Errorf("undated tombstone = %+v, want kept with a date", ts)
	}
}
