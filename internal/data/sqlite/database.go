package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

// Database represents a SQLite database connection
type Database struct {
	db     *sql.DB
	path   string
}

// DatabaseConfig holds connection pool configuration
type DatabaseConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultDatabaseConfig returns default connection pool settings
func DefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		MaxOpenConns:    1,  // SQLite single-writer model
		MaxIdleConns:    1,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 5 * time.Minute,
	}
}

// NewDatabase creates a new database connection with default settings
func NewDatabase(path string) (*Database, error) {
	return NewDatabaseWithConfig(path, DefaultDatabaseConfig())
}

// NewDatabaseWithConfig creates a new database connection with custom config
func NewDatabaseWithConfig(path string, config DatabaseConfig) (*Database, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	if config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(config.ConnMaxLifetime)
	}
	if config.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{
		db:   db,
		path: path,
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// DB returns the underlying sql.DB
func (d *Database) DB() *sql.DB {
	return d.db
}

// PoolStats returns current connection pool statistics
func (d *Database) PoolStats() *sql.DBStats {
	stats := d.db.Stats()
	return &stats
}

// Ping checks database connectivity
func (d *Database) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// IsHealthy checks if the database connection is healthy
func (d *Database) IsHealthy(ctx context.Context) bool {
	if err := d.db.PingContext(ctx); err != nil {
		return false
	}
	stats := d.db.Stats()
	// Consider unhealthy if too many connections are in use
	return stats.InUse < stats.MaxOpenConnections
}

// Migrate runs database migrations
func (d *Database) Migrate(ctx context.Context) error {
	migrations := []string{
		// Patterns table
		`CREATE TABLE IF NOT EXISTS patterns (
			id TEXT PRIMARY KEY,
			trigger TEXT NOT NULL,
			response TEXT NOT NULL,
			strength REAL NOT NULL DEFAULT 0,
			threshold REAL NOT NULL DEFAULT 50,
			decay_rate REAL NOT NULL DEFAULT 0.01,
			decay_enabled INTEGER NOT NULL DEFAULT 1,
			connections TEXT,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			reinforcement_count INTEGER NOT NULL DEFAULT 0,
			decay_count INTEGER NOT NULL DEFAULT 0,
			last_used_at INTEGER,
			tags TEXT,
			project TEXT,
			user_id TEXT,
			space_id TEXT DEFAULT 'global',
			deleted_at INTEGER
		)`,

		// Spaces table
		`CREATE TABLE IF NOT EXISTS spaces (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			owner TEXT,
			is_default INTEGER NOT NULL DEFAULT 0,
			pattern_limit INTEGER NOT NULL DEFAULT 0,
			pattern_count INTEGER NOT NULL DEFAULT 0,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		)`,

		// Events table (for audit/logging)
		`CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			timestamp INTEGER NOT NULL,
			payload TEXT,
			source TEXT,
			trace_id TEXT,
			pattern_id TEXT,
			user_id TEXT
		)`,

		// Notes table (Phase 10: 思绪整理)
		`CREATE TABLE IF NOT EXISTS notes (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			space_id TEXT DEFAULT 'global',
			tags TEXT,
			is_pinned INTEGER NOT NULL DEFAULT 0,
			category TEXT DEFAULT 'note',
			word_count INTEGER NOT NULL DEFAULT 0,
			char_count INTEGER NOT NULL DEFAULT 0,
			last_viewed_at INTEGER,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		)`,

		// Indices - Basic
		`CREATE INDEX IF NOT EXISTS idx_patterns_trigger ON patterns(trigger)`,
		`CREATE INDEX IF NOT EXISTS idx_patterns_strength ON patterns(strength)`,
		`CREATE INDEX IF NOT EXISTS idx_patterns_project ON patterns(project)`,
		`CREATE INDEX IF NOT EXISTS idx_patterns_tags ON patterns(tags)`,
		`CREATE INDEX IF NOT EXISTS idx_patterns_deleted ON patterns(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_events_type ON events(type)`,
		`CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp)`,

		// Indices - Performance optimization (Iter 43)
		`CREATE INDEX IF NOT EXISTS idx_patterns_last_used_at ON patterns(last_used_at)`,
		`CREATE INDEX IF NOT EXISTS idx_patterns_updated_at ON patterns(updated_at)`,
		`CREATE INDEX IF NOT EXISTS idx_patterns_project_deleted ON patterns(project, deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_patterns_strength_threshold ON patterns(strength, threshold)`,
		`CREATE INDEX IF NOT EXISTS idx_patterns_decay_enabled ON patterns(decay_enabled)`,
		`CREATE INDEX IF NOT EXISTS idx_patterns_space_id ON patterns(space_id)`,

		// Notes indices (Phase 10)
		`CREATE INDEX IF NOT EXISTS idx_notes_space_id ON notes(space_id)`,
		`CREATE INDEX IF NOT EXISTS idx_notes_category ON notes(category)`,
		`CREATE INDEX IF NOT EXISTS idx_notes_is_pinned ON notes(is_pinned)`,
		`CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_notes_updated_at ON notes(updated_at)`,
	}

	for _, migration := range migrations {
		if _, err := d.db.ExecContext(ctx, migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

// InitDefaultSpaces initializes default spaces if they don't exist
func (d *Database) InitDefaultSpaces(ctx context.Context) error {
	spaces := models.DefaultSpaces()
	for _, space := range spaces {
		_, err := d.db.ExecContext(ctx, `
			INSERT OR IGNORE INTO spaces (id, name, description, owner, is_default, pattern_limit, pattern_count, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, space.ID, space.Name, space.Description, space.Owner, boolToInt(space.DefaultSpace),
			space.PatternLimit, space.PatternCount, space.CreatedAt.Unix(), space.UpdatedAt.Unix())
		if err != nil {
			return fmt.Errorf("failed to init space %s: %w", space.ID, err)
		}
	}
	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Path returns the database file path
func (d *Database) Path() string {
	return d.path
}
