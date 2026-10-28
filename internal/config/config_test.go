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

func TestDesktopHotkeyDefault(t *testing.T) {
	path := writeConfigFile(t, `{}`)
	m := NewManager()
	if err := m.Load(path, ""); err != nil {
		t.Fatalf("load: %v", err)
	}
	if got := m.Get().DesktopHotkey; got != "" {
		t.Errorf("DesktopHotkey = %q, want empty default", got)
	}
}

func TestDesktopHotkeyFromFile(t *testing.T) {
	path := writeConfigFile(t, `{"desktop_hotkey": "Ctrl+Alt+B"}`)
	m := NewManager()
	if err := m.Load(path, ""); err != nil {
		t.Fatalf("load: %v", err)
	}
	if got := m.Get().DesktopHotkey; got != "Ctrl+Alt+B" {
		t.Errorf("DesktopHotkey = %q, want Ctrl+Alt+B", got)
	}
}

func TestDesktopHotkeyFromEnv(t *testing.T) {
	// Env vars take precedence over the config file value.
	path := writeConfigFile(t, `{"desktop_hotkey": "Ctrl+Alt+B"}`)
	t.Setenv("CAIRN_DESKTOP_HOTKEY", "off")
	m := NewManager()
	if err := m.Load(path, ""); err != nil {
		t.Fatalf("load: %v", err)
	}
	if got := m.Get().DesktopHotkey; got != "off" {
		t.Errorf("DesktopHotkey = %q, want CAIRN_DESKTOP_HOTKEY=off to override config file", got)
	}
}

func TestSetUpdatesResolvedConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("HOME", tmp)
	t.Setenv("CAIRN_DROPBOX_APP_KEY", "")
	t.Setenv("CAIRN_DB_PATH", "")

	m := NewManager()
	if err := m.Load("", ""); err != nil {
		t.Fatalf("Load: %v", err)
	}
	cfg := m.Get()
	if cfg.DropboxAppKey != "" {
		t.Fatalf("expected empty app key, got %q", cfg.DropboxAppKey)
	}

	m.Set("dropbox_app_key", "new-key")
	m.Set("db_path", "/tmp/custom.db")

	if got := m.Get().DropboxAppKey; got != "new-key" {
		t.Errorf("DropboxAppKey = %q, want %q", got, "new-key")
	}
	if got := m.Get().DBPath; got != "/tmp/custom.db" {
		t.Errorf("DBPath = %q, want %q", got, "/tmp/custom.db")
	}
	// Previously returned pointers observe the update too.
	if cfg.DropboxAppKey != "new-key" {
		t.Errorf("earlier Get() pointer DropboxAppKey = %q, want %q", cfg.DropboxAppKey, "new-key")
	}
}
