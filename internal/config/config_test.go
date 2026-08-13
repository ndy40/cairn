package config

import (
	"os"
	"path/filepath"
	"testing"
)

// writeConfigFile creates a cairn.json with the given JSON body and returns its path.
func writeConfigFile(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "cairn.json")
	if err := os.WriteFile(path, []byte(body), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func TestDisableAutoArchiveDefault(t *testing.T) {
	path := writeConfigFile(t, `{}`)
	m := NewManager()
	if err := m.Load(path, ""); err != nil {
		t.Fatalf("load: %v", err)
	}
	if m.Get().DisableAutoArchive {
		t.Error("expected DisableAutoArchive to default to false")
	}
}

func TestDisableAutoArchiveFromFile(t *testing.T) {
	path := writeConfigFile(t, `{"disable_auto_archive": true}`)
	m := NewManager()
	if err := m.Load(path, ""); err != nil {
		t.Fatalf("load: %v", err)
	}
	if !m.Get().DisableAutoArchive {
		t.Error("expected DisableAutoArchive to be true from config file")
	}
}

func TestDisableAutoArchiveFromEnv(t *testing.T) {
	// Env vars take precedence over the config file value.
	path := writeConfigFile(t, `{"disable_auto_archive": false}`)
	t.Setenv("CAIRN_DISABLE_AUTO_ARCHIVE", "true")
	m := NewManager()
	if err := m.Load(path, ""); err != nil {
		t.Fatalf("load: %v", err)
	}
	if !m.Get().DisableAutoArchive {
		t.Error("expected CAIRN_DISABLE_AUTO_ARCHIVE=true to override config file")
	}
}
