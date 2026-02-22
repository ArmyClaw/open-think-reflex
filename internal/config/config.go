// Package config provides configuration loading and management for Open-Think-Reflex.
// Supports YAML configuration files with environment variable overrides.
package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the top-level application configuration.
// All settings can be overridden via environment variables (e.g., OTR_AI_PROVIDER).
type Config struct {
	// Version for config migration tracking
	Version int `mapstructure:"version"`

	// Application-level settings
	App AppConfig `mapstructure:"app"`

	// Storage backend configuration
	Storage StorageConfig `mapstructure:"storage"`

	// AI provider settings
	AI AIConfig `mapstructure:"ai"`

	// Terminal UI settings
	UI UIConfig `mapstructure:"ui"`

	// Security and audit settings
	Security SecurityConfig `mapstructure:"security"`

	// CurrentSpace (v2.0) - the currently active space
	CurrentSpace string `mapstructure:"current_space"`
}

// AppConfig contains application-level settings.
type AppConfig struct {
	Name     string `mapstructure:"name"`     // Application name
	Version  string `mapstructure:"version"`  // Version string
	DataDir  string `mapstructure:"data_dir"` // Data directory path
	LogLevel string `mapstructure:"log_level"` // Log level (debug, info, warn, error)
	Profile  bool   `mapstructure:"profile"`  // Enable profiling
}

// StorageConfig contains storage backend settings.
type StorageConfig struct {
	Type            string `mapstructure:"type"`             // Storage type (sqlite, etc.)
	Path            string `mapstructure:"path"`             // Database file path
	CacheSize       int    `mapstructure:"cache_size"`       // LRU cache capacity
	MaxOpenConns    int    `mapstructure:"max_open_conns"`   // Max open connections
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`   // Max idle connections
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // Connection max lifetime (seconds)
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"` // Connection max idle time (seconds)
}

// AIConfig contains AI provider configuration.
type AIConfig struct {
	Provider     string         `mapstructure:"provider"`      // Primary provider (anthropic, openai, local)
	Providers    ProvidersConfig `mapstructure:"providers"`   // Per-provider settings
	DefaultModel string         `mapstructure:"default_model"` // Default model name
	Timeout      int            `mapstructure:"timeout"`       // Request timeout (seconds)
	RetryMax     int            `mapstructure:"retry_max"`    // Max retry attempts
}

// ProvidersConfig contains per-provider settings.
type ProvidersConfig struct {
	Anthropic AnthropicConfig `mapstructure:"anthropic"` // Anthropic Claude settings
	OpenAI    OpenAIConfig    `mapstructure:"openai"`    // OpenAI settings
	Local     LocalConfig     `mapstructure:"local"`     // Local model settings
}

// AnthropicConfig contains Anthropic Claude-specific settings.
// API key can be set via config or OTR_ANTHROPIC_API_KEY environment variable.
type AnthropicConfig struct {
	APIKey      string  `mapstructure:"api_key"`       // Anthropic API key
	APIURL      string  `mapstructure:"api_url"`       // API endpoint (for proxy/regional)
	Model       string  `mapstructure:"model"`         // Default model
	MaxTokens   int     `mapstructure:"max_tokens"`    // Max tokens per request
	Temperature float64 `mapstructure:"temperature"`    // Temperature (0.0-1.0)
}

// OpenAIConfig contains OpenAI-specific settings.
type OpenAIConfig struct {
	APIKey      string  `mapstructure:"api_key"`       // OpenAI API key
	APIURL      string  `mapstructure:"api_url"`       // API endpoint
	Model       string  `mapstructure:"model"`         // Default model
	MaxTokens   int     `mapstructure:"max_tokens"`    // Max tokens per request
	Temperature float64 `mapstructure:"temperature"`    // Temperature (0.0-1.0)
}

// LocalConfig contains settings for local AI models (Ollama, LM Studio, etc.).
type LocalConfig struct {
	APIURL string `mapstructure:"api_url"` // Local server URL (e.g., http://localhost:11434)
	Model  string `mapstructure:"model"`   // Model name
}

// UIConfig contains terminal UI configuration.
type UIConfig struct {
	Theme      string      `mapstructure:"theme"`       // Theme name (dark, light)
	Colors     ColorsConfig `mapstructure:"colors"`     // Color scheme
	KeyMap     KeyMapConfig `mapstructure:"keymap"`     // Keyboard shortcuts
	OutputMode string     `mapstructure:"output_mode"` // Output format (terminal, json)
}

// ColorsConfig contains terminal color settings (hex color codes).
type ColorsConfig struct {
	Root        string `mapstructure:"root"`         // Root node color
	Branch      string `mapstructure:"branch"`       // Branch color
	Selected    string `mapstructure:"selected"`     // Selected item color
	Unmatched   string `mapstructure:"unmatched"`     // Unmatched item color
	Background  string `mapstructure:"background"`    // Background color
}

// KeyMapConfig contains keyboard shortcut mappings.
// Default is vim-style navigation.
type KeyMapConfig struct {
	Up      string `mapstructure:"up"`       // Move up
	Down    string `mapstructure:"down"`     // Move down
	Left    string `mapstructure:"left"`     // Move left
	Right   string `mapstructure:"right"`    // Move right
	Select  string `mapstructure:"select"`   // Select item
	Confirm string `mapstructure:"confirm"`  // Confirm action
	Cancel  string `mapstructure:"cancel"`   // Cancel action
	Quit    string `mapstructure:"quit"`     // Quit application
	Help    string `mapstructure:"help"`     // Show help
}

// SecurityConfig contains security-related settings.
type SecurityConfig struct {
	APIKeysEnvPrefix string     `mapstructure:"api_keys_env_prefix"` // Env var prefix for API keys
	ConfigFileMode   string     `mapstructure:"config_file_mode"`   // Config file permissions
	AuditLog         AuditConfig `mapstructure:"audit_log"`         // Audit logging settings
}

// AuditConfig contains audit logging configuration.
type AuditConfig struct {
	Enabled bool   `mapstructure:"enabled"` // Enable audit logging
	Path    string `mapstructure:"path"`   // Audit log file path
}

// Loader handles configuration loading from files with environment variable overrides.
// Uses Viper for flexible configuration management.
type Loader struct {
	configPath string     // Directory containing config file
	configName string     // Config file name (without extension)
	configType string     // Config file type (yaml, json, toml)
	v          *viper.Viper // Viper instance for config management
}

// NewLoader creates a new configuration loader.
//
// Example:
//
//	loader := NewLoader("/home/user/.config/reflex", "config")
//	cfg, err := loader.Load()
func NewLoader(configPath, configName string) *Loader {
	return &Loader{
		configPath: configPath,
		configName: configName,
		configType: "yaml",
	}
}

// Load loads the configuration from file
func (l *Loader) Load() (*Config, error) {
	l.v = viper.New()

	// Set defaults
	l.setDefaults()

	// Set config path and name
	l.v.AddConfigPath(l.configPath)
	l.v.SetConfigName(l.configName)
	l.v.SetConfigType(l.configType)

	// Enable environment variable override
	l.v.AutomaticEnv()
	l.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := l.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, use defaults
			fmt.Printf("Config file not found, using defaults: %s\n", err)
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Apply environment variable overrides for sensitive data
	l.applyEnvOverrides()

	// Unmarshal config
	var cfg Config
	if err := l.v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Resolve paths
	l.resolvePaths(&cfg)

	return &cfg, nil
}

// setDefaults sets the default configuration values
func (l *Loader) setDefaults() {
	// Version
	l.v.SetDefault("version", 1)

	// App defaults
	l.v.SetDefault("app.name", "open-think-reflex")
	l.v.SetDefault("app.version", "dev")
	l.v.SetDefault("app.data_dir", "$HOME/.openclaw/reflex")
	l.v.SetDefault("app.log_level", "info")
	l.v.SetDefault("app.profile", false)

	// Storage defaults
	l.v.SetDefault("storage.type", "sqlite")
	l.v.SetDefault("storage.path", "$HOME/.openclaw/reflex/data.db")
	l.v.SetDefault("storage.cache_size", 1000)
	l.v.SetDefault("storage.max_open_conns", 1)   // SQLite single-writer model
	l.v.SetDefault("storage.max_idle_conns", 1)
	l.v.SetDefault("storage.conn_max_lifetime", 3600)  // 1 hour
	l.v.SetDefault("storage.conn_max_idle_time", 300)  // 5 minutes

	// AI defaults
	l.v.SetDefault("ai.provider", "anthropic")
	l.v.SetDefault("ai.default_model", "claude-sonnet-4-20250514")
	l.v.SetDefault("ai.timeout", 30)
	l.v.SetDefault("ai.retry_max", 3)

	// AI Providers defaults
	l.v.SetDefault("ai.providers.anthropic.api_url", "https://api.anthropic.com/v1")
	l.v.SetDefault("ai.providers.anthropic.max_tokens", 4096)
	l.v.SetDefault("ai.providers.anthropic.temperature", 0.7)

	l.v.SetDefault("ai.providers.openai.api_url", "https://api.openai.com/v1")
	l.v.SetDefault("ai.providers.openai.model", "gpt-4")
	l.v.SetDefault("ai.providers.openai.max_tokens", 4096)
	l.v.SetDefault("ai.providers.openai.temperature", 0.7)

	l.v.SetDefault("ai.providers.local.api_url", "http://localhost:11434/v1")
	l.v.SetDefault("ai.providers.local.model", "llama2")

	// UI defaults
	l.v.SetDefault("ui.theme", "dark")
	l.v.SetDefault("ui.output_mode", "terminal")

	// KeyMap defaults (vim-style)
	l.v.SetDefault("ui.keymap.up", "k")
	l.v.SetDefault("ui.keymap.down", "j")
	l.v.SetDefault("ui.keymap.left", "h")
	l.v.SetDefault("ui.keymap.right", "l")
	l.v.SetDefault("ui.keymap.select", "space")
	l.v.SetDefault("ui.keymap.confirm", "enter")
	l.v.SetDefault("ui.keymap.cancel", "esc")
	l.v.SetDefault("ui.keymap.quit", "q")
	l.v.SetDefault("ui.keymap.help", "?")

	// Security defaults
	l.v.SetDefault("security.api_keys_env_prefix", "OTR_")
	l.v.SetDefault("security.config_file_mode", "0600")
	l.v.SetDefault("security.audit_log.enabled", true)
	l.v.SetDefault("security.audit_log.path", "$HOME/.openclaw/reflex/audit.log")
}

// applyEnvOverrides applies environment variable overrides for sensitive data
func (l *Loader) applyEnvOverrides() {
	prefix := l.v.GetString("security.api_keys_env_prefix")

	// Anthropic API key
	if key := os.Getenv(prefix + "ANTHROPIC_API_KEY"); key != "" {
		l.v.Set("ai.providers.anthropic.api_key", key)
	}

	// OpenAI API key
	if key := os.Getenv(prefix + "OPENAI_API_KEY"); key != "" {
		l.v.Set("ai.providers.openai.api_key", key)
	}
}

// resolvePaths resolves relative paths to absolute paths
func (l *Loader) resolvePaths(cfg *Config) {
	home, _ := os.UserHomeDir()

	// Resolve data_dir
	if strings.HasPrefix(cfg.App.DataDir, "$HOME") {
		cfg.App.DataDir = strings.Replace(cfg.App.DataDir, "$HOME", home, 1)
	}

	// Resolve storage path
	if strings.HasPrefix(cfg.Storage.Path, "$HOME") {
		cfg.Storage.Path = strings.Replace(cfg.Storage.Path, "$HOME", home, 1)
	}

	// Resolve audit log path
	if strings.HasPrefix(cfg.Security.AuditLog.Path, "$HOME") {
		cfg.Security.AuditLog.Path = strings.Replace(cfg.Security.AuditLog.Path, "$HOME", home, 1)
	}
}

// Save saves the configuration to file
func (l *Loader) Save(cfg *Config) error {
	// Ensure config directory exists
	configDir := l.configPath
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write config file
	l.v.SetConfigFile(l.getConfigFilePath())
	if err := l.v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigFilePath returns the full config file path
func (l *Loader) getConfigFilePath() string {
	return fmt.Sprintf("%s/%s.%s", l.configPath, l.configName, l.configType)
}

// GetTimeout returns the timeout as a duration
func (c *Config) GetTimeout() time.Duration {
	return time.Duration(c.AI.Timeout) * time.Second
}

// GetCacheSize returns the cache size as int
func (c *Config) GetCacheSize() int {
	return c.Storage.CacheSize
}

// GetAIProvider returns the configured AI provider
func (c *Config) GetAIProvider() string {
	return c.AI.Provider
}

// GetDefaultModel returns the default AI model
func (c *Config) GetDefaultModel() string {
	return c.AI.DefaultModel
}

// GetCurrentSpace returns the current space ID
func (c *Config) GetCurrentSpace() string {
	return c.CurrentSpace
}

// SetCurrentSpace sets the current space ID
func (c *Config) SetCurrentSpace(spaceID string) {
	c.CurrentSpace = spaceID
}
