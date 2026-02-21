package config

import (
	"os"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	
	if cfg.Database.MaxOpenConns != 25 {
		t.Errorf("Expected 25, got %d", cfg.Database.MaxOpenConns)
	}
	if cfg.Database.MaxIdleConns != 10 {
		t.Errorf("Expected 10, got %d", cfg.Database.MaxIdleConns)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected 8080, got %d", cfg.Server.Port)
	}
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected 0.0.0.0, got %s", cfg.Server.Host)
	}
	if !cfg.Features.EnableCache {
		t.Error("Expected EnableCache to be true")
	}
	if !cfg.Features.EnableCompression {
		t.Error("Expected EnableCompression to be true")
	}
	if !cfg.Features.EnableMetrics {
		t.Error("Expected EnableMetrics to be true")
	}
}

func TestLoad(t *testing.T) {
	f, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	
	f.WriteString(`{"server":{"port":9090}}`)
	f.Close()
	
	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	
	if cfg.Server.Port != 9090 {
		t.Errorf("Expected 9090, got %d", cfg.Server.Port)
	}
}

func TestLoadNotExist(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.json")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", cfg.Server.Port)
	}
}

func TestSave(t *testing.T) {
	cfg := Default()
	cfg.Server.Port = 3000
	
	f, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	
	cfg.filepath = f.Name()
	if err := cfg.Save(); err != nil {
		t.Fatal(err)
	}
	
	loaded, err := Load(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	
	if loaded.Server.Port != 3000 {
		t.Errorf("Expected 3000, got %d", loaded.Server.Port)
	}
}

func TestGetDatabase(t *testing.T) {
	cfg := Default()
	cfg.Database.MaxOpenConns = 50
	
	db := cfg.GetDatabase()
	if db.MaxOpenConns != 50 {
		t.Errorf("Expected 50, got %d", db.MaxOpenConns)
	}
}

func TestGetCache(t *testing.T) {
	cfg := Default()
	cfg.Cache.TTL = time.Hour
	
	cache := cfg.GetCache()
	if cache.TTL != time.Hour {
		t.Errorf("Expected 1h, got %v", cache.TTL)
	}
}

func TestGetServer(t *testing.T) {
	cfg := Default()
	cfg.Server.Host = "localhost"
	
	server := cfg.GetServer()
	if server.Host != "localhost" {
		t.Errorf("Expected localhost, got %s", server.Host)
	}
}

func TestSetDatabase(t *testing.T) {
	cfg := Default()
	cfg.SetDatabase(DatabaseConfig{MaxOpenConns: 100})
	
	if cfg.Database.MaxOpenConns != 100 {
		t.Errorf("Expected 100, got %d", cfg.Database.MaxOpenConns)
	}
}

func TestSetCache(t *testing.T) {
	cfg := Default()
	cfg.SetCache(CacheConfig{Enabled: false})
	
	if cfg.Cache.Enabled {
		t.Error("Expected cache to be disabled")
	}
}

func TestSetServer(t *testing.T) {
	cfg := Default()
	cfg.SetServer(ServerConfig{Port: 9000})
	
	if cfg.Server.Port != 9000 {
		t.Errorf("Expected 9000, got %d", cfg.Server.Port)
	}
}

func TestSaveEmptyFilepath(t *testing.T) {
	cfg := Default()
	err := cfg.Save()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
