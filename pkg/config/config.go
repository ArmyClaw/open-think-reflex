package config

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Config holds application configuration (Iter 54)
type Config struct {
	mu sync.RWMutex
	
	// Database
	Database DatabaseConfig `json:"database"`
	
	// Cache
	Cache CacheConfig `json:"cache"`
	
	// Server
	Server ServerConfig `json:"server"`
	
	// Logger
	Logger LoggerConfig `json:"logger"`
	
	// Features
	Features FeaturesConfig `json:"features"`
	
	// loaded from file
	filepath string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path            string `json:"path"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"`
	ConnMaxIdleTime int    `json:"conn_max_idle_time"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled    bool          `json:"enabled"`
	TTL        time.Duration `json:"ttl"`
	MaxSize    int           `json:"max_size"`
	RedisAddr  string        `json:"redis_addr"`
	RedisPassword string     `json:"redis_password"`
	RedisDB    int           `json:"redis_db"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Timeout   int    `json:"timeout"`
	EnableTLS bool   `json:"enable_tls"`
	CertFile  string `json:"cert_file"`
	KeyFile   string `json:"key_file"`
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level    string `json:"level"`
	Format   string `json:"format"`
	Output   string `json:"output"`
	FilePath string `json:"file_path"`
	MaxSize  int    `json:"max_size"`
}

// FeaturesConfig holds feature flags
type FeaturesConfig struct {
	EnableCache       bool `json:"enable_cache"`
	EnableCompression bool `json:"enable_compression"`
	EnableMetrics     bool `json:"enable_metrics"`
	EnableTracing     bool `json:"enable_tracing"`
}

// Default returns default configuration
func Default() *Config {
	return &Config{
		Database: DatabaseConfig{
			Path:            "./data/otr.db",
			MaxOpenConns:    25,
			MaxIdleConns:    10,
			ConnMaxLifetime: 300,
			ConnMaxIdleTime: 60,
		},
		Cache: CacheConfig{
			Enabled: true,
			TTL:     time.Minute,
			MaxSize: 1000,
		},
		Server: ServerConfig{
			Host:    "0.0.0.0",
			Port:    8080,
			Timeout: 30,
		},
		Logger: LoggerConfig{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		},
		Features: FeaturesConfig{
			EnableCache:       true,
			EnableCompression: true,
			EnableMetrics:     true,
			EnableTracing:     false,
		},
	}
}

// Load loads configuration from file
func Load(filepath string) (*Config, error) {
	cfg := Default()
	cfg.filepath = filepath
	
	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

// Save saves configuration to file
func (c *Config) Save() error {
	if c.filepath == "" {
		return nil
	}
	
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(c.filepath, data, 0644)
}

// GetDatabase returns database config
func (c *Config) GetDatabase() DatabaseConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Database
}

// GetCache returns cache config
func (c *Config) GetCache() CacheConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Cache
}

// GetServer returns server config
func (c *Config) GetServer() ServerConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Server
}
