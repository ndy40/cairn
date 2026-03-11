package sync

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// DropboxConfig holds Dropbox-specific OAuth2 credentials.
type DropboxConfig struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenExpiry  time.Time `json:"token_expiry"`
	AppKey       string    `json:"app_key"`
}

// SyncConfig holds the sync configuration persisted as a JSON file.
type SyncConfig struct {
	Backend      string         `json:"backend"`
	DeviceID     string         `json:"device_id"`
	LastSyncAt   *time.Time     `json:"last_sync_at,omitempty"`
	SyncDeclined bool           `json:"sync_declined"`
	Dropbox      *DropboxConfig `json:"dropbox,omitempty"`
}

// DefaultConfigPath returns the OS-appropriate path for the sync config file.
func DefaultConfigPath() string {
	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, "Library", "Application Support", "cairn", "sync.json")
	default: // linux and others
		xdg := os.Getenv("XDG_CONFIG_HOME")
		if xdg == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return ""
			}
			xdg = filepath.Join(home, ".config")
		}
		return filepath.Join(xdg, "cairn", "sync.json")
	}
}

// LoadConfig reads the sync config from disk.
// Returns nil, nil if the file does not exist.
func LoadConfig(path string) (*SyncConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var cfg SyncConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig writes the sync config to disk with mode 0600.
func SaveConfig(path string, cfg *SyncConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// IsConfigured returns true if the config has an active backend (not just declined).
func IsConfigured(cfg *SyncConfig) bool {
	return cfg != nil && cfg.Backend != "" && !cfg.SyncDeclined
}
