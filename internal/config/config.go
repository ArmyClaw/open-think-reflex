package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	// Version for config migration
	Version int `mapstructure:"version"`

	// Application settings
	App AppConfig `mapstructure:"app"`

	// Storage settings
	Storage StorageConfig `mapstructure:"storage"`

	// AI provider settings
	AI AIConfig `mapstructure:"ai"`

	// UI settings
	UI UIConfig `mapstructure:"ui"`

	// Security settings
	Security SecurityConfig `mapstructure:"security"`
}

// AppConfig contains application-level settings
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	DataDir     string `mapstructure:"data_dir"`
	LogLevel    string `mapstructure:"log_level"`
	Profile     bool   `mapstructure:"profile"`
}

// StorageConfig contains storage-related settings
type StorageConfig struct {
	Type     string `mapstructure:"type"`
	Path     string `mapstructure:"path"`
	CacheSize int   `mapstructure:"cache_size"`
}

// AIConfig contains AI provider settings
type AIConfig struct {
	Provider     string         `mapstructure:"provider"`
	Providers    ProvidersConfig `mapstructure:"providers"`
	DefaultModel string        `mapstructure:"default_model"`
	Timeout      int           `mapstructure:"timeout"`
	RetryMax     int           `mapstructure:"retry_max"`
}

// ProvidersConfig contains per-provider settings
type ProvidersConfig struct {
	Anthropic AnthropicConfig `mapstructure:"anthropic"`
	OpenAI    OpenAIConfig    `mapstructure:"openai"`
	Local     LocalConfig     `mapstructure:"local"`
}

// AnthropicConfig contains Anthropic-specific settings
type AnthropicConfig struct {
	APIKey      string  `mapstructure:"api_key"`
	APIURL      string  `mapstructure:"api_url"`
	Model       string  `mapstructure:"model"`
	MaxTokens   int     `mapstructure:"max_tokens"`
	Temperature float64 `mapstructure:"temperature"`
}

// OpenAIConfig contains OpenAI-specific settings
type OpenAIConfig struct {
	APIKey      string  `mapstructure:"api_key"`
	APIURL      string  `mapstructure:"api_url"`
	Model       string  `mapstructure:"model"`
	MaxTokens   int     `mapstructure:"max_tokens"`
	Temperature float64 `mapstructure:"temperature"`
}

// LocalConfig contains local AI settings
type LocalConfig struct {
	APIURL string `mapstructure:"api_url"`
	Model  string `mapstructure:"model"`
}

// UIConfig contains UI-related settings
type UIConfig struct {
	Theme      string      `mapstructure:"theme"`
	Colors     ColorsConfig `mapstructure:"colors"`
	KeyMap     KeyMapConfig `mapstructure:"keymap"`
	OutputMode string     `mapstructure:"output_mode"`
}

// ColorsConfig contains color settings
type ColorsConfig struct {
	Root        string `mapstructure:"root"`
	Branch      string `mapstructure:"branch"`
	Selected    string `mapstructure:"selected"`
	Unmatched   string `mapstructure:"unmatched"`
	Background  string `mapstructure:"background"`
}

// KeyMapConfig contains keyboard mapping settings
type KeyMapConfig struct {
	Up      string `mapstructure:"up"`
	Down    string `mapstructure:"down"`
	Left    string `mapstructure:"left"`
	Right   string `mapstructure:"right"`
	Select  string `mapstructure:"select"`
	Confirm string `mapstructure:"confirm"`
	Cancel  string `mapstructure:"cancel"`
	Quit    string `mapstructure:"quit"`
	Help    string `mapstructure:"help"`
}

// SecurityConfig contains security-related settings
type SecurityConfig struct {
	APIKeysEnvPrefix string     `mapstructure:"api_keys_env_prefix"`
	ConfigFileMode   string     `mapstructure:"config_file_mode"`
	AuditLog         AuditConfig `mapstructure:"audit_log"`
}

// AuditConfig contains audit log settings
type AuditConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

// Loader handles configuration loading
type Loader struct {
	configPath string
	configName string
	configType string
	v          *viper.Viper
}

// NewLoader creates a new configuration loader
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
