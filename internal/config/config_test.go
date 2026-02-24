package config

import (
	"os"
	"testing"
)

func TestConfig_GetTimeout(t *testing.T) {
	cfg := &Config{
		AI: AIConfig{
			Timeout: 30,
		},
	}
	
	result := cfg.GetTimeout()
	if result.Seconds() != 30 {
		t.Errorf("expected 30s, got %v", result)
	}
}

func TestConfig_GetCacheSize(t *testing.T) {
	cfg := &Config{
		Storage: StorageConfig{
			CacheSize: 1000,
		},
	}
	
	result := cfg.GetCacheSize()
	if result != 1000 {
		t.Errorf("expected 1000, got %d", result)
	}
}

func TestConfig_GetAIProvider(t *testing.T) {
	cfg := &Config{
		AI: AIConfig{
			Provider: "anthropic",
		},
	}
	
	result := cfg.GetAIProvider()
	if result != "anthropic" {
		t.Errorf("expected anthropic, got %s", result)
	}
}

func TestConfig_GetDefaultModel(t *testing.T) {
	cfg := &Config{
		AI: AIConfig{
			DefaultModel: "claude-sonnet-4-20250514",
		},
	}
	
	result := cfg.GetDefaultModel()
	if result != "claude-sonnet-4-20250514" {
		t.Errorf("expected claude-sonnet-4-20250514, got %s", result)
	}
}

func TestNewLoader(t *testing.T) {
	loader := NewLoader("/tmp", "config")
	if loader.configPath != "/tmp" {
		t.Errorf("expected /tmp, got %s", loader.configPath)
	}
	if loader.configName != "config" {
		t.Errorf("expected config, got %s", loader.configName)
	}
	if loader.configType != "yaml" {
		t.Errorf("expected yaml, got %s", loader.configType)
	}
}

func TestLoader_LoadDefaults(t *testing.T) {
	// Create a temp directory for config
	tmpDir := t.TempDir()
	
	loader := NewLoader(tmpDir, "nonexistent-config")
	cfg, err := loader.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	
	// Check defaults
	if cfg.App.Name != "open-think-reflex" {
		t.Errorf("expected open-think-reflex, got %s", cfg.App.Name)
	}
	if cfg.Storage.Type != "sqlite" {
		t.Errorf("expected sqlite, got %s", cfg.Storage.Type)
	}
	if cfg.AI.Provider != "anthropic" {
		t.Errorf("expected anthropic, got %s", cfg.AI.Provider)
	}
	if cfg.UI.Theme != "dark" {
		t.Errorf("expected dark, got %s", cfg.UI.Theme)
	}
}

func TestLoader_LoadWithEnvOverride(t *testing.T) {
	// Set environment variable
	os.Setenv("OTR_ANTHROPIC_API_KEY", "test-key-123")
	defer os.Unsetenv("OTR_ANTHROPIC_API_KEY")
	
	tmpDir := t.TempDir()
	loader := NewLoader(tmpDir, "nonexistent-config")
	cfg, err := loader.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	
	// Check env override
	if cfg.AI.Providers.Anthropic.APIKey != "test-key-123" {
		t.Errorf("expected test-key-123, got %s", cfg.AI.Providers.Anthropic.APIKey)
	}
}

func TestLoader_getConfigFilePath(t *testing.T) {
	loader := NewLoader("/home/user/.config/otr", "config")
	path := loader.getConfigFilePath()
	expected := "/home/user/.config/otr/config.yaml"
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}

func TestLoader_resolvePaths(t *testing.T) {
	home, _ := os.UserHomeDir()
	
	loader := NewLoader("/tmp", "config")
	cfg := &Config{
		App: AppConfig{
			DataDir: "$HOME/.otr",
		},
		Storage: StorageConfig{
			Path: "$HOME/.otr/data.db",
		},
		Security: SecurityConfig{
			AuditLog: AuditConfig{
				Path: "$HOME/.otr/audit.log",
			},
		},
	}
	
	loader.resolvePaths(cfg)
	
	// Check paths are resolved
	if cfg.App.DataDir == "$HOME/.otr" {
		t.Error("DataDir should be resolved")
	}
	if cfg.Storage.Path == "$HOME/.otr/data.db" {
		t.Error("Storage.Path should be resolved")
	}
	
	// Verify they contain home directory
	expectedPrefix := home
	if cfg.App.DataDir[:len(home)] != expectedPrefix {
		t.Errorf("expected DataDir to start with %s, got %s", expectedPrefix, cfg.App.DataDir[:len(home)])
	}
}
