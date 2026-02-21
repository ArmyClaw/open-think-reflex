package config

import (
	"os"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	
	if cfg.Database.MaxOpenConns != 25 {
		t.Errorf("Expected 25, got %d", cfg.Database.MaxOpenConns)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected 8080, got %d", cfg.Server.Port)
	}
	if !cfg.Features.EnableCache {
		t.Error("Expected EnableCache to be true")
	}
}

func TestLoad(t *testing.T) {
	// Create temp file
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
