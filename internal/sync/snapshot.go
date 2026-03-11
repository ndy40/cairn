package sync

import (
	"encoding/json"
	"time"
)

// SyncRecord is the cloud snapshot format stored in Dropbox.
type SyncRecord struct {
	Version       int              `json:"version"`
	LastUpdatedBy string           `json:"last_updated_by"`
	LastUpdatedAt time.Time        `json:"last_updated_at"`
	Bookmarks     []BookmarkEntry  `json:"bookmarks"`
	Tombstones    []TombstoneEntry `json:"tombstones"`
}

// BookmarkEntry represents a bookmark in the cloud snapshot.
type BookmarkEntry struct {
	UUID        string     `json:"uuid"`
	URL         string     `json:"url"`
	Domain      string     `json:"domain"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Tags        []string   `json:"tags"`
	IsPermanent bool       `json:"is_permanent"`
	IsArchived  bool       `json:"is_archived"`
	ArchivedAt  *time.Time `json:"archived_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Deleted     bool       `json:"deleted"`
}

// TombstoneEntry represents a deleted bookmark marker in the cloud snapshot.
type TombstoneEntry struct {
	UUID      string    `json:"uuid"`
	URL       string    `json:"url"`
	DeletedAt time.Time `json:"deleted_at"`
	DeletedBy string    `json:"deleted_by"`
}

// NewSyncRecord creates an empty SyncRecord with version 1.
func NewSyncRecord(deviceID string) *SyncRecord {
	return &SyncRecord{
		Version:       1,
		LastUpdatedBy: deviceID,
		LastUpdatedAt: time.Now().UTC(),
		Bookmarks:     []BookmarkEntry{},
		Tombstones:    []TombstoneEntry{},
	}
}

// Marshal serializes the SyncRecord to JSON.
func (r *SyncRecord) Marshal() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// UnmarshalSyncRecord deserializes a SyncRecord from JSON.
func UnmarshalSyncRecord(data []byte) (*SyncRecord, error) {
	var r SyncRecord
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
