package backend

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"golang.org/x/oauth2"
)

// DropboxBackend implements SyncBackend for Dropbox.
type DropboxBackend struct {
	client files.Client
}

// NewDropboxBackend creates a new Dropbox backend with OAuth2 token auto-refresh.
func NewDropboxBackend(token *oauth2.Token, appKey string) *DropboxBackend {
	cfg := &oauth2.Config{
		ClientID: appKey,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/oauth2/authorize",
			TokenURL: "https://api.dropboxapi.com/oauth2/token",
		},
	}

	// Create an HTTP client that auto-refreshes the token.
	tokenSource := cfg.TokenSource(context.Background(), token)
	httpClient := oauth2.NewClient(context.Background(), tokenSource)

	dbxCfg := dropbox.Config{
		Token:  token.AccessToken,
		Client: httpClient,
	}

	return &DropboxBackend{
		client: files.New(dbxCfg),
	}
}

// Upload writes data to the given remote path with overwrite mode.
func (d *DropboxBackend) Upload(data []byte, remotePath string) error {
	arg := files.NewUploadArg(remotePath)
	arg.Mode = &files.WriteMode{Tagged: dropbox.Tagged{Tag: "overwrite"}}

	_, err := d.client.Upload(arg, strings.NewReader(string(data)))
	if err != nil {
		return mapDropboxError(err)
	}
	return nil
}

// Download retrieves data from the given remote path.
func (d *DropboxBackend) Download(remotePath string) ([]byte, error) {
	arg := files.NewDownloadArg(remotePath)
	_, content, err := d.client.Download(arg)
	if err != nil {
		return nil, mapDropboxError(err)
	}
	defer func() { _ = content.Close() }()

	data, err := io.ReadAll(content)
	if err != nil {
		return nil, fmt.Errorf("read dropbox response: %w", err)
	}
	return data, nil
}

// Exists checks whether a file exists at the given remote path.
func (d *DropboxBackend) Exists(remotePath string) (bool, error) {
	arg := files.NewGetMetadataArg(remotePath)
	_, err := d.client.GetMetadata(arg)
	if err != nil {
		mapped := mapDropboxError(err)
		if mapped == ErrNotFound {
			return false, nil
		}
		return false, mapped
	}
	return true, nil
}

// mapDropboxError maps Dropbox API errors to sentinel errors.
func mapDropboxError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()

	if strings.Contains(msg, "not_found") || strings.Contains(msg, "path/not_found") {
		return ErrNotFound
	}
	if strings.Contains(msg, "invalid_access_token") || strings.Contains(msg, "expired_access_token") {
		return ErrAuthExpired
	}
	if strings.Contains(msg, "insufficient_space") {
		return ErrQuotaExceeded
	}

	// Check for HTTP transport errors.
	if isNetworkError(err) {
		return ErrNetworkFailure
	}

	return fmt.Errorf("dropbox: %w", err)
}

// isNetworkError checks if the error is a network-level failure.
func isNetworkError(err error) bool {
	msg := err.Error()
	networkPatterns := []string{
		"connection refused",
		"no such host",
		"timeout",
		"EOF",
		"connection reset",
	}
	for _, p := range networkPatterns {
		if strings.Contains(msg, p) {
			return true
		}
	}

	// Check if it wraps an HTTP error.
	if httpErr, ok := err.(*dropbox.APIError); ok {
		_ = httpErr
		return false
	}

	return false
}

// Ensure DropboxBackend implements SyncBackend at compile time.
var _ SyncBackend = (*DropboxBackend)(nil)
