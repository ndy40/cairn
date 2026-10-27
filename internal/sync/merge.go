package sync

import (
	"time"

	"github.com/ndy40/cairn/internal/store"
)

// MergeResult contains the outcome of merging remote data with local data.
type MergeResult struct {
	ToInsert []BookmarkEntry // Remote bookmarks not present locally
	ToUpdate []BookmarkEntry // Remote bookmarks newer than local
	ToDelete []string        // UUIDs of local bookmarks that have remote tombstones
	ToHeal   []CreatedAtFix  // Local bookmarks whose created_at is later than the remote copy's
}

// CreatedAtFix moves a local bookmark's created_at back to an earlier remote value.
type CreatedAtFix struct {
	UUID      string // Local bookmark UUID
	CreatedAt time.Time
}

// MergeBookmarks compares local bookmarks against the remote snapshot and returns
// what needs to be inserted, updated, or deleted locally.
// Logic:
//   - Remote bookmark not in local (by URL) → insert
//   - Remote bookmark with newer updated_at (matched by URL) → update
//   - Remote tombstone matching a local UUID → delete
//   - Matched bookmark with an earlier remote created_at → heal local created_at
func MergeBookmarks(local []*store.Bookmark, remote *SyncRecord) *MergeResult {
	result := &MergeResult{}

	// Index local bookmarks by URL and UUID for fast lookup.
	localByURL := make(map[string]*store.Bookmark, len(local))
	localByUUID := make(map[string]*store.Bookmark, len(local))
	for _, b := range local {
		localByURL[b.URL] = b
		if b.UUID != "" {
			localByUUID[b.UUID] = b
		}
	}

	// Process remote bookmarks.
	for _, rb := range remote.Bookmarks {
		if rb.Deleted {
			continue
		}

		// Try to match by UUID first, then by URL.
		lb, ok := localByUUID[rb.UUID]
		if !ok {
			// Same URL, different UUID — dedup by URL.
			lb, ok = localByURL[rb.URL]
		}
		if ok {
			// Update if remote is newer.
			if rb.UpdatedAt.After(lb.UpdatedAt) {
				result.ToUpdate = append(result.ToUpdate, rb)
			}
			// The earliest known creation date wins.
			if !rb.CreatedAt.IsZero() && rb.CreatedAt.Before(lb.CreatedAt) {
				result.ToHeal = append(result.ToHeal, CreatedAtFix{UUID: lb.UUID, CreatedAt: rb.CreatedAt})
			}
			continue
		}

		// Not present locally — insert.
		result.ToInsert = append(result.ToInsert, rb)
	}

	// Process remote tombstones.
	for _, ts := range remote.Tombstones {
		if _, ok := localByUUID[ts.UUID]; ok {
			result.ToDelete = append(result.ToDelete, ts.UUID)
		}
	}

	return result
}
