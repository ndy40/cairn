package sync

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIsConfigured_NilConfig(t *testing.T) {
	if IsConfigured(nil) {
		t.Error("expected false for nil config")
	}
}

func TestIsConfigured_SyncDeclined(t *testing.T) {
	cfg := &SyncConfig{SyncDeclined: true, Backend: "dropbox", Dropbox: &DropboxConfig{
		AccessToken: "tok",
		AppKey:      "key",
	}}
	if IsConfigured(cfg) {
		t.Error("expected false when SyncDeclined=true")
	}
}

func TestIsConfigured_EmptyBackend(t *testing.T) {
	cfg := &SyncConfig{Backend: ""}
	if IsConfigured(cfg) {
		t.Error("expected false when Backend is empty")
	}
}

func TestIsConfigured_DropboxMissingAccessToken(t *testing.T) {
	cfg := &SyncConfig{
		Backend: "dropbox",
		Dropbox: &DropboxConfig{AppKey: "key", AccessToken: ""},
	}
	if IsConfigured(cfg) {
		t.Error("expected false when Dropbox AccessToken is empty")
	}
}

func TestIsConfigured_DropboxMissingAppKey(t *testing.T) {
	cfg := &SyncConfig{
		Backend: "dropbox",
		Dropbox: &DropboxConfig{AccessToken: "tok", AppKey: ""},
	}
	if IsConfigured(cfg) {
		t.Error("expected false when Dropbox AppKey is empty")
	}
}

func TestIsConfigured_DropboxNilDropboxField(t *testing.T) {
	cfg := &SyncConfig{Backend: "dropbox", Dropbox: nil}
	if IsConfigured(cfg) {
		t.Error("expected false when Dropbox field is nil")
	}
}

func TestIsConfigured_DropboxFullyConfigured(t *testing.T) {
	cfg := &SyncConfig{
		Backend: "dropbox",
		Dropbox: &DropboxConfig{AccessToken: "tok", AppKey: "key"},
	}
	if !IsConfigured(cfg) {
		t.Error("expected true for fully configured Dropbox")
	}
}

func TestIsConfigured_UnknownBackend(t *testing.T) {
	cfg := &SyncConfig{Backend: "s3"}
	if IsConfigured(cfg) {
		t.Error("expected false for unknown backend")
	}
}

func TestLoadConfig_FileNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sync.json")
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if cfg != nil {
		t.Error("expected nil config for missing file")
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sync.json")

	expiry := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	want := &SyncConfig{
		Backend:  "dropbox",
		DeviceID: "device-123",
		Dropbox: &DropboxConfig{
			AccessToken:  "acc",
			RefreshToken: "ref",
			TokenExpiry:  expiry,
			AppKey:       "appkey",
		},
	}
	data, _ := json.Marshal(want)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.Backend != want.Backend {
		t.Errorf("Backend: got %q, want %q", cfg.Backend, want.Backend)
	}
	if cfg.DeviceID != want.DeviceID {
		t.Errorf("DeviceID: got %q, want %q", cfg.DeviceID, want.DeviceID)
	}
	if cfg.Dropbox == nil {
		t.Fatal("expected non-nil Dropbox config")
	}
	if cfg.Dropbox.AccessToken != want.Dropbox.AccessToken {
		t.Errorf("AccessToken: got %q, want %q", cfg.Dropbox.AccessToken, want.Dropbox.AccessToken)
	}
}

func TestLoadConfig_DeclinedConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sync.json")

	declined := &SyncConfig{SyncDeclined: true}
	data, _ := json.Marshal(declined)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if !cfg.SyncDeclined {
		t.Error("expected SyncDeclined=true")
	}
	if IsConfigured(cfg) {
		t.Error("IsConfigured should return false for declined config")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sync.json")
	if err := os.WriteFile(path, []byte("not-json"), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
	if cfg != nil {
		t.Error("expected nil config on parse error")
	}
}
