package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// AppConfig represents the final resolved configuration after merging all sources.
type AppConfig struct {
	DBPath        string `mapstructure:"db_path"`
	DropboxAppKey string `mapstructure:"dropbox_app_key"`
}

// Manager handles configuration loading and resolution using viper.
type Manager struct {
	v      *viper.Viper
	config *AppConfig
}

// NewManager creates a new configuration manager.
func NewManager() *Manager {
	return &Manager{
		v: viper.New(),
	}
}

// DefaultConfigPath returns the OS-appropriate path for cairn.json.
func DefaultConfigPath() string {
	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, "Library", "Application Support", "cairn", "cairn.json")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return ""
			}
			appData = home
		}
		return filepath.Join(appData, "cairn", "cairn.json")
	default: // linux and others
		xdg := os.Getenv("XDG_CONFIG_HOME")
		if xdg == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return ""
			}
			xdg = filepath.Join(home, ".config")
		}
		return filepath.Join(xdg, "cairn", "cairn.json")
	}
}

// DefaultDBPath returns the OS-appropriate default database path.
func DefaultDBPath() string {
	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, "Library", "Application Support", "cairn", "bookmarks.db")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return ""
			}
			appData = home
		}
		return filepath.Join(appData, "cairn", "bookmarks.db")
	default: // linux and others
		xdg := os.Getenv("XDG_DATA_HOME")
		if xdg == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return ""
			}
			xdg = filepath.Join(home, ".local", "share")
		}
		return filepath.Join(xdg, "cairn", "bookmarks.db")
	}
}

// LegacyDBPath returns the pre-rename default database path (bookmark-manager → cairn).
// Used only by the one-time migration in the CLI startup path.
func LegacyDBPath() string {
	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, "Library", "Application Support", "bookmark-manager", "bookmarks.db")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return ""
			}
			appData = home
		}
		return filepath.Join(appData, "bookmark-manager", "bookmarks.db")
	default: // linux and others
		xdg := os.Getenv("XDG_DATA_HOME")
		if xdg == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return ""
			}
			xdg = filepath.Join(home, ".local", "share")
		}
		return filepath.Join(xdg, "bookmark-manager", "bookmarks.db")
	}
}

// Load initializes and loads configuration from all sources.
// Precedence (highest to lowest): env vars > CLI flags > config file > defaults
func (m *Manager) Load(configPath string, cliDBFlag string) error {
	// Set config name and type
	m.v.SetConfigName("cairn")
	m.v.SetConfigType("json")

	// Set defaults
	m.v.SetDefault("db_path", DefaultDBPath())
	m.v.SetDefault("dropbox_app_key", "")

	// Add config path
	if configPath != "" {
		// Explicit config file path provided
		m.v.SetConfigFile(configPath)
	} else {
		// Use default config path
		defaultPath := DefaultConfigPath()
		if defaultPath != "" {
			configDir := filepath.Dir(defaultPath)
			m.v.AddConfigPath(configDir)
		}

		// Also check current directory
		m.v.AddConfigPath(".")

		// Check home directory
		if home, err := os.UserHomeDir(); err == nil {
			m.v.AddConfigPath(home)
		}
	}

	// Bind environment variables
	m.v.SetEnvPrefix("CAIRN")
	m.v.AutomaticEnv()

	// Map environment variable names
	_ = m.v.BindEnv("db_path", "CAIRN_DB_PATH")
	_ = m.v.BindEnv("dropbox_app_key", "CAIRN_DROPBOX_APP_KEY")

	// Read config file (ignore error if not found)
	if err := m.v.ReadInConfig(); err != nil {
		// Only return error if it's not a "file not found" error
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && configPath != "" {
			return fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is OK - will use defaults and env vars
	}

	// Override with CLI flag if provided
	if cliDBFlag != "" {
		m.v.Set("db_path", cliDBFlag)
	}

	// Unmarshal into struct
	m.config = &AppConfig{}
	if err := m.v.Unmarshal(m.config); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}

// Get returns the loaded configuration.
func (m *Manager) Get() *AppConfig {
	if m.config == nil {
		// Return defaults if Load was never called
		return &AppConfig{
			DBPath:        DefaultDBPath(),
			DropboxAppKey: "",
		}
	}
	return m.config
}

// GetString returns a string configuration value.
func (m *Manager) GetString(key string) string {
	return m.v.GetString(key)
}

// Set sets a configuration value (useful for testing).
func (m *Manager) Set(key string, value interface{}) {
	m.v.Set(key, value)
}

// ConfigFileUsed returns the config file path that was used.
func (m *Manager) ConfigFileUsed() string {
	return m.v.ConfigFileUsed()
}

// AllSettings returns all configuration settings.
func (m *Manager) AllSettings() map[string]interface{} {
	return m.v.AllSettings()
}

// WriteConfig writes the current configuration to file.
func (m *Manager) WriteConfig() error {
	configPath := m.v.ConfigFileUsed()
	if configPath == "" {
		configPath = DefaultConfigPath()
		// Create directory if it doesn't exist (0700 = owner-only for security)
		dir := filepath.Dir(configPath)
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("error creating config directory: %w", err)
		}
		m.v.SetConfigFile(configPath)
	}
	return m.v.WriteConfig()
}

// SaveConfig writes the current configuration to file, creating it if necessary.
func (m *Manager) SaveConfig() error {
	configPath := m.v.ConfigFileUsed()
	if configPath == "" {
		configPath = DefaultConfigPath()
		// Create directory if it doesn't exist (0700 = owner-only for security)
		dir := filepath.Dir(configPath)
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("error creating config directory: %w", err)
		}
		m.v.SetConfigFile(configPath)
	}
	return m.v.SafeWriteConfig()
}

// Resolve is a backward-compatible function that creates a manager and loads config.
// Deprecated: Use NewManager().Load() instead.
func Resolve(fileConfig interface{}, cliDBFlag string, defaultDBPath string) AppConfig {
	manager := NewManager()

	// Override default if provided
	if defaultDBPath != "" {
		manager.v.SetDefault("db_path", defaultDBPath)
	}

	// Load configuration
	if err := manager.Load("", cliDBFlag); err != nil {
		// Return defaults on error
		return AppConfig{
			DBPath:        defaultDBPath,
			DropboxAppKey: "",
		}
	}

	return *manager.Get()
}
