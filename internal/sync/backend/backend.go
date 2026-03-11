package backend

import "errors"

// Sentinel errors for sync backend operations.
var (
	ErrNotFound       = errors.New("sync: file not found")
	ErrAuthExpired    = errors.New("sync: authentication expired")
	ErrNetworkFailure = errors.New("sync: network failure")
	ErrQuotaExceeded  = errors.New("sync: storage quota exceeded")
)

// SyncBackend defines the contract for cloud storage providers.
type SyncBackend interface {
	// Upload writes data to the given remote path, overwriting any existing file.
	Upload(data []byte, remotePath string) error

	// Download retrieves data from the given remote path.
	// Returns ErrNotFound if the file does not exist.
	Download(remotePath string) ([]byte, error)

	// Exists checks whether a file exists at the given remote path.
	Exists(remotePath string) (bool, error)
}
