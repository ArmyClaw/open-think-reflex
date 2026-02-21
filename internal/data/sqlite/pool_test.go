package sqlite

import (
	"context"
	"testing"
	"time"
)

// TestConnectionPoolConfig tests custom connection pool configuration
func TestConnectionPoolConfig(t *testing.T) {
	ctx := context.Background()
	
	// Test with custom config
	config := DatabaseConfig{
		MaxOpenConns:    2,
		MaxIdleConns:    1,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}
	
	db, err := NewDatabaseWithConfig(":memory:", config)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()
	
	// Verify connection pool settings
	stats := db.PoolStats()
	if stats.MaxOpenConnections != 2 {
		t.Errorf("expected max open connections 2, got %d", stats.MaxOpenConnections)
	}
	
	// Run migrations
	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	
	// Test health check
	if !db.IsHealthy(ctx) {
		t.Error("database should be healthy")
	}
}

// TestDefaultDatabaseConfig tests the default configuration
func TestDefaultDatabaseConfig(t *testing.T) {
	config := DefaultDatabaseConfig()
	
	if config.MaxOpenConns != 1 {
		t.Errorf("expected max open conns 1, got %d", config.MaxOpenConns)
	}
	if config.MaxIdleConns != 1 {
		t.Errorf("expected max idle conns 1, got %d", config.MaxIdleConns)
	}
	if config.ConnMaxLifetime != time.Hour {
		t.Errorf("expected conn max lifetime 1h, got %v", config.ConnMaxLifetime)
	}
	if config.ConnMaxIdleTime != 5*time.Minute {
		t.Errorf("expected conn max idle time 5m, got %v", config.ConnMaxIdleTime)
	}
}

// TestDatabasePing tests the Ping method
func TestDatabasePing(t *testing.T) {
	ctx := context.Background()
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()
	
	if err := db.Ping(ctx); err != nil {
		t.Errorf("ping failed: %v", err)
	}
}

// TestPoolStats tests that pool statistics are available
func TestPoolStats(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()
	
	stats := db.PoolStats()
	if stats == nil {
		t.Error("expected non-nil stats")
	}
	
	// Should have at least some idle or open connection info
	// SQLite in-memory databases may show 0, but the method should work
	t.Logf("Pool stats: Open=%d, InUse=%d, Idle=%d, WaitCount=%d", 
		stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount)
}
