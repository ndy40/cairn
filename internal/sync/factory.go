package sync

import (
	"fmt"

	"github.com/ndy40/cairn/internal/sync/backend"
	"golang.org/x/oauth2"
)

// NewBackend creates a SyncBackend based on the config's backend type.
func NewBackend(cfg *SyncConfig) (backend.SyncBackend, error) {
	switch cfg.Backend {
	case "dropbox":
		if cfg.Dropbox == nil {
			return nil, fmt.Errorf("dropbox config is missing")
		}
		token := &oauth2.Token{
			AccessToken:  cfg.Dropbox.AccessToken,
			RefreshToken: cfg.Dropbox.RefreshToken,
			Expiry:       cfg.Dropbox.TokenExpiry,
		}
		return backend.NewDropboxBackend(token, cfg.Dropbox.AppKey), nil
	default:
		return nil, fmt.Errorf("unsupported backend: %s", cfg.Backend)
	}
}
