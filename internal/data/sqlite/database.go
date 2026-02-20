package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

// Database represents a SQLite database connection
type Database struct {
	db     *sql.DB
	path   string
}

// NewDatabase creates a new database connection
func NewDatabase(path string) (*Database, error) {
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

	// Configure connection
	db.SetMaxOpenConns(1) // SQLite single-writer model
	db.SetMaxIdleConns(1)

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
			deleted_at INTEGER
		)`,

		// Spaces table
		`CREATE TABLE IF NOT EXISTS spaces (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
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

		// Indices
		`CREATE INDEX IF NOT EXISTS idx_patterns_trigger ON patterns(trigger)`,
		`CREATE INDEX IF NOT EXISTS idx_patterns_strength ON patterns(strength)`,
		`CREATE INDEX IF NOT EXISTS idx_patterns_project ON patterns(project)`,
		`CREATE INDEX IF NOT EXISTS idx_patterns_tags ON patterns(tags)`,
		`CREATE INDEX IF NOT EXISTS idx_patterns_deleted ON patterns(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_events_type ON events(type)`,
		`CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp)`,
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
			INSERT OR IGNORE INTO spaces (id, name, description, is_default, pattern_limit, pattern_count, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, space.ID, space.Name, space.Description, boolToInt(space.DefaultSpace),
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
