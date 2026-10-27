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
	ToRekey  []UUIDFix       // Local bookmarks that must adopt the remote copy's UUID
}

// UUIDFix renames a local bookmark's UUID to the one used by the remote copy.
type UUIDFix struct {
	From string // Local bookmark UUID
	To   string // Remote bookmark UUID
}

// CreatedAtFix moves a local bookmark's created_at back to an earlier remote value.
type CreatedAtFix struct {
	UUID      string // Local bookmark UUID, after any rekey
	CreatedAt time.Time
}

// MergeBookmarks compares local bookmarks against the remote snapshot and returns
// what needs to be inserted, updated, or deleted locally.
// Logic:
//   - Remote bookmark not in local (by URL) → insert
//   - Remote bookmark with newer updated_at (matched by URL) → update
//   - Remote tombstone matching a local UUID → delete
//   - Matched bookmark with an earlier remote created_at → heal local created_at
//   - Same URL under different UUIDs → both sides keep the UUID of the earliest copy
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
			localUUID := lb.UUID
			// Diverged copies converge on one UUID so tombstones match everywhere.
			if rb.UUID != "" && rb.UUID != lb.UUID && remoteCopyWins(rb, lb) {
				result.ToRekey = append(result.ToRekey, UUIDFix{From: lb.UUID, To: rb.UUID})
				localUUID = rb.UUID
			}
			// Update if remote is newer.
			if rb.UpdatedAt.After(lb.UpdatedAt) {
				result.ToUpdate = append(result.ToUpdate, rb)
			}
			// The earliest known creation date wins.
			if !rb.CreatedAt.IsZero() && rb.CreatedAt.Before(lb.CreatedAt) {
				result.ToHeal = append(result.ToHeal, CreatedAtFix{UUID: localUUID, CreatedAt: rb.CreatedAt})
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

// remoteCopyWins reports whether the remote copy of a bookmark is the earliest
// one, so its UUID should be kept. Ties go to the smaller UUID, so every device
// picks the same winner. A zero remote date never wins.
func remoteCopyWins(rb BookmarkEntry, lb *store.Bookmark) bool {
	if rb.CreatedAt.IsZero() {
		return false
	}
	if rb.CreatedAt.Equal(lb.CreatedAt) {
		return rb.UUID < lb.UUID
	}
	return rb.CreatedAt.Before(lb.CreatedAt)
}
