package sync

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ndy40/cairn/internal/store"
	"github.com/ndy40/cairn/internal/sync/backend"
)

const remotePath = "/cairn/sync.json"

// Engine orchestrates sync operations.
type Engine struct {
	Store      *store.Store
	Backend    backend.SyncBackend
	Config     *SyncConfig
	ConfigPath string
}

// NewEngine creates a new sync engine.
func NewEngine(s *store.Store, b backend.SyncBackend, cfg *SyncConfig, cfgPath string) *Engine {
	return &Engine{
		Store:      s,
		Backend:    b,
		Config:     cfg,
		ConfigPath: cfgPath,
	}
}

// Setup runs the initial sync setup: OAuth2 flow, device ID generation,
// and initial sync (merge or upload).
func (e *Engine) Setup(appKey string) (int, error) {
	// Run OAuth2 flow.
	token, err := RunOAuth2Flow(appKey)
	if err != nil {
		return 0, fmt.Errorf("oauth2 flow: %w", err)
	}

	// Generate device ID.
	deviceID := uuid.New().String()

	// Save config.
	e.Config = &SyncConfig{
		Backend:  "dropbox",
		DeviceID: deviceID,
		Dropbox: &DropboxConfig{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenExpiry:  token.Expiry,
			AppKey:       appKey,
		},
	}

	if err := SaveConfig(e.ConfigPath, e.Config); err != nil {
		return 0, fmt.Errorf("save config: %w", err)
	}

	// Create the Dropbox backend with the new token.
	e.Backend = backend.NewDropboxBackend(token, appKey)

	// Check if cloud snapshot exists.
	exists, err := e.Backend.Exists(remotePath)
	if err != nil {
		return 0, fmt.Errorf("check cloud snapshot: %w", err)
	}

	var count int
	if exists {
		count, err = e.pullAndMerge()
	} else {
		count, err = e.pushInitial()
	}
	if err != nil {
		return 0, err
	}

	// Update last_sync_at.
	now := time.Now().UTC()
	e.Config.LastSyncAt = &now
	if err := SaveConfig(e.ConfigPath, e.Config); err != nil {
		return 0, fmt.Errorf("update last_sync_at: %w", err)
	}

	return count, nil
}

// pullAndMerge downloads the cloud snapshot and merges it with local data.
func (e *Engine) pullAndMerge() (int, error) {
	data, err := e.Backend.Download(remotePath)
	if err != nil {
		return 0, fmt.Errorf("download snapshot: %w", err)
	}

	remote, err := UnmarshalSyncRecord(data)
	if err != nil {
		// Handle corrupted snapshot: upload local data as fresh snapshot.
		fmt.Fprintf(os.Stderr, "warning: corrupted cloud snapshot, re-uploading local data\n")
		return e.pushInitial()
	}

	local, err := e.Store.ExportAll()
	if err != nil {
		return 0, fmt.Errorf("export local bookmarks: %w", err)
	}

	result := MergeBookmarks(local, remote)

	// Apply inserts.
	for _, entry := range result.ToInsert {
		tags := entry.Tags
		if tags == nil {
			tags = []string{}
		}
		_, err := e.Store.Insert(entry.URL, entry.Title, entry.Description, tags)
		if err != nil {
			// Skip duplicates silently.
			continue
		}
	}

	// Apply updates (update tags and metadata for matched bookmarks by URL).
	for _, entry := range result.ToUpdate {
		existing, err := e.Store.GetByUUID(entry.UUID)
		if err != nil {
			// Try by URL if UUID doesn't match.
			continue
		}
		if err := e.Store.UpdateTags(existing.ID, entry.Tags); err != nil {
			continue
		}
	}

	// Apply deletes.
	for _, bookmarkUUID := range result.ToDelete {
		bm, err := e.Store.GetByUUID(bookmarkUUID)
		if err != nil {
			continue
		}
		if _, err := e.Store.DeleteByID(bm.ID); err != nil {
			continue
		}
	}

	// Also upload merged snapshot back to cloud.
	if err := e.uploadSnapshot(); err != nil {
		return 0, fmt.Errorf("upload merged snapshot: %w", err)
	}

	all, err := e.Store.ExportAll()
	if err != nil {
		return 0, err
	}
	return len(all), nil
}

// pushInitial exports all local bookmarks and uploads them as the initial snapshot.
func (e *Engine) pushInitial() (int, error) {
	if err := e.uploadSnapshot(); err != nil {
		return 0, fmt.Errorf("upload initial snapshot: %w", err)
	}
	all, err := e.Store.ExportAll()
	if err != nil {
		return 0, err
	}
	return len(all), nil
}

// uploadSnapshot builds a SyncRecord from local data and uploads it.
func (e *Engine) uploadSnapshot() error {
	bookmarks, err := e.Store.ExportAll()
	if err != nil {
		return fmt.Errorf("export bookmarks: %w", err)
	}

	record := NewSyncRecord(e.Config.DeviceID)
	for _, b := range bookmarks {
		entry := BookmarkEntry{
			UUID:        b.UUID,
			URL:         b.URL,
			Domain:      b.Domain,
			Title:       b.Title,
			Description: b.Description,
			Tags:        b.Tags,
			IsPermanent: b.IsPermanent,
			IsArchived:  b.IsArchived,
			ArchivedAt:  b.ArchivedAt,
			CreatedAt:   b.CreatedAt,
			UpdatedAt:   b.UpdatedAt,
			Deleted:     false,
		}
		record.Bookmarks = append(record.Bookmarks, entry)
	}

	data, err := record.Marshal()
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}

	return e.Backend.Upload(data, remotePath)
}

// Push uploads the local bookmarks snapshot to cloud and clears pending changes.
func (e *Engine) Push() error {
	if err := e.uploadSnapshot(); err != nil {
		return fmt.Errorf("push: %w", err)
	}

	// Clear pending changes after successful push.
	if err := e.Store.ClearPendingChanges(); err != nil {
		return fmt.Errorf("clear pending changes: %w", err)
	}

	now := time.Now().UTC()
	e.Config.LastSyncAt = &now
	return SaveConfig(e.ConfigPath, e.Config)
}

// Pull downloads the cloud snapshot and merges with local data.
func (e *Engine) Pull() (int, error) {
	count, err := e.pullAndMerge()
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	e.Config.LastSyncAt = &now
	if err := SaveConfig(e.ConfigPath, e.Config); err != nil {
		return 0, fmt.Errorf("update last_sync_at: %w", err)
	}

	return count, nil
}

// AutoPush is a non-fatal version of Push for use in auto-sync hooks.
// Returns an error message string if sync failed, empty string on success.
func (e *Engine) AutoPush() string {
	if err := e.Push(); err != nil {
		if isAuthError(err) {
			return "sync: auth expired, run 'cairn sync auth' to re-authenticate"
		}
		return fmt.Sprintf("sync push warning: %v", err)
	}
	return ""
}

// AutoPull is a non-fatal version of Pull for use in auto-sync hooks.
// Returns bookmark count, and an error message string if sync failed.
func (e *Engine) AutoPull() (int, string) {
	count, err := e.Pull()
	if err != nil {
		if isAuthError(err) {
			return 0, "sync: auth expired, run 'cairn sync auth' to re-authenticate"
		}
		return 0, fmt.Sprintf("sync pull warning: %v", err)
	}
	return count, ""
}

// Status returns sync status information.
type SyncStatus struct {
	Configured   bool
	Backend      string
	DeviceID     string
	LastSyncAt   *time.Time
	PendingCount int
}

// Status returns the current sync status.
func (e *Engine) Status() (*SyncStatus, error) {
	status := &SyncStatus{
		Configured: IsConfigured(e.Config),
	}

	if e.Config != nil {
		status.Backend = e.Config.Backend
		status.DeviceID = e.Config.DeviceID
		status.LastSyncAt = e.Config.LastSyncAt
	}

	count, err := e.Store.PendingChangeCount()
	if err != nil {
		return nil, err
	}
	status.PendingCount = count

	return status, nil
}

// Unlink removes the sync configuration and stops syncing.
func (e *Engine) Unlink() error {
	configPath := e.ConfigPath
	e.Config = nil
	e.Backend = nil
	return os.Remove(configPath)
}

// isAuthError checks if an error chain contains an auth-expired error.
func isAuthError(err error) bool {
	return err != nil && (errors.Is(err, backend.ErrAuthExpired) ||
		strings.Contains(err.Error(), "auth expired") ||
		strings.Contains(err.Error(), "invalid_access_token"))
}

// BookmarkToJSON serializes a bookmark to JSON for pending change payload.
func BookmarkToJSON(b *store.Bookmark) string {
	data, err := json.Marshal(b)
	if err != nil {
		return "{}"
	}
	return string(data)
}
