package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/google/uuid"
)

// Storage implements contracts.Storage using SQLite
type Storage struct {
	db        *Database
	stmtCache map[string]*sql.Stmt
	mu        sync.RWMutex
	queryCache *QueryCache  // Iter 49

	// Concurrency statistics (Iter 47)
	stats StorageStats
}

// StorageStats holds concurrency statistics for monitoring
type StorageStats struct {
	ReadOps        int64 // Total read operations
	WriteOps       int64 // Total write operations
	ActiveReaders  int64 // Current active read operations
	ActiveWriters  int64 // Current active write operations
	ReadWaitTime   int64 // Nanoseconds spent waiting for read locks
	WriteWaitTime  int64 // Nanoseconds spent waiting for write locks
}

// NewStorage creates a new SQLite storage
func NewStorage(db *Database) *Storage {
	return &Storage{
		db:        db,
		stmtCache: make(map[string]*sql.Stmt),
	}
}

// SavePattern saves a pattern to the database
func (s *Storage) SavePattern(ctx context.Context, p *models.Pattern) error {
	// Validate pattern
	if err := p.Validate(); err != nil {
		return err
	}

	connections, _ := json.Marshal(p.Connections)
	tags, _ := json.Marshal(p.Tags)

	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now
	
	// Default space_id to "global" if not set
	if p.SpaceID == "" {
		p.SpaceID = "global"
	}

	_, err := s.db.db.ExecContext(ctx, `
		INSERT OR REPLACE INTO patterns (
			id, trigger, response, strength, threshold, decay_rate, decay_enabled,
			connections, created_at, updated_at, reinforcement_count, decay_count,
			last_used_at, tags, project, user_id, space_id, deleted_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		p.ID, p.Trigger, p.Response, p.Strength, p.Threshold, p.DecayRate, p.DecayEnabled,
		string(connections), p.CreatedAt.Unix(), p.UpdatedAt.Unix(), p.ReinforceCnt, p.DecayCnt,
		int64TimeToPtr(p.LastUsedAt), string(tags), p.Project, p.UserID, p.SpaceID, int64TimeToPtr(p.DeletedAt),
	)

	return err
}

// GetPattern retrieves a pattern by ID
func (s *Storage) GetPattern(ctx context.Context, id string) (*models.Pattern, error) {
	var p models.Pattern
	var connections, tags sql.NullString
	var lastUsedAt, deletedAt, createdAt, updatedAt sql.NullInt64

	err := s.db.db.QueryRowContext(ctx, `
		SELECT id, trigger, response, strength, threshold, decay_rate, decay_enabled,
			connections, created_at, updated_at, reinforcement_count, decay_count,
			last_used_at, tags, project, user_id, deleted_at
		FROM patterns WHERE id = ? AND deleted_at IS NULL
	`, id).Scan(
		&p.ID, &p.Trigger, &p.Response, &p.Strength, &p.Threshold, &p.DecayRate, &p.DecayEnabled,
		&connections, &createdAt, &updatedAt, &p.ReinforceCnt, &p.DecayCnt,
		&lastUsedAt, &tags, &p.Project, &p.UserID, &deletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("pattern not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	if connections.Valid {
		if err := json.Unmarshal([]byte(connections.String), &p.Connections); err != nil {
			return nil, fmt.Errorf("failed to parse connections for pattern %s: %w", id, err)
		}
	}
	if tags.Valid {
		if err := json.Unmarshal([]byte(tags.String), &p.Tags); err != nil {
			return nil, fmt.Errorf("failed to parse tags for pattern %s: %w", id, err)
		}
	}
	p.CreatedAt = int64ToTime(createdAt)
	p.UpdatedAt = int64ToTime(updatedAt)
	p.LastUsedAt = int64ToTimePtr(lastUsedAt)
	p.DeletedAt = int64ToTimePtr(deletedAt)

	return &p, nil
}

// ListPatterns lists patterns with filtering options
func (s *Storage) ListPatterns(ctx context.Context, opts contracts.ListOptions) ([]*models.Pattern, error) {
	query := `SELECT id, trigger, response, strength, threshold, decay_rate, decay_enabled,
		connections, created_at, updated_at, reinforcement_count, decay_count,
		last_used_at, tags, project, user_id, space_id
		FROM patterns WHERE deleted_at IS NULL`
	args := []interface{}{}

	// Filter by space (v2.0)
	if opts.SpaceID != "" {
		query += " AND space_id = ?"
		args = append(args, opts.SpaceID)
	}

	// Apply filters
	if opts.Project != "" {
		query += " AND project = ?"
		args = append(args, opts.Project)
	}

	if opts.MinStrength > 0 {
		query += " AND strength >= ?"
		args = append(args, opts.MinStrength)
	}

	// Add ORDER BY for consistency
	query += " ORDER BY updated_at DESC"

	// Apply limit/offset
	if opts.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, opts.Limit)
	}
	if opts.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, opts.Offset)
	}

	rows, err := s.db.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patterns []*models.Pattern
	for rows.Next() {
		var p models.Pattern
		var connections, tags, spaceID sql.NullString
		var lastUsedAt, createdAt, updatedAt sql.NullInt64

		err := rows.Scan(
			&p.ID, &p.Trigger, &p.Response, &p.Strength, &p.Threshold, &p.DecayRate, &p.DecayEnabled,
			&connections, &createdAt, &updatedAt, &p.ReinforceCnt, &p.DecayCnt,
			&lastUsedAt, &tags, &p.Project, &p.UserID, &spaceID,
		)
		if err != nil {
			return nil, err
		}

		if connections.Valid {
			if err := json.Unmarshal([]byte(connections.String), &p.Connections); err != nil {
				return nil, fmt.Errorf("failed to parse connections: %w", err)
			}
		}
		if tags.Valid {
			if err := json.Unmarshal([]byte(tags.String), &p.Tags); err != nil {
				return nil, fmt.Errorf("failed to parse tags: %w", err)
			}
		}
		p.CreatedAt = int64ToTime(createdAt)
		p.UpdatedAt = int64ToTime(updatedAt)
		p.LastUsedAt = int64ToTimePtr(lastUsedAt)
		p.SpaceID = spaceID.String // Default to empty string if NULL, will use global

		patterns = append(patterns, &p)
	}

	return patterns, rows.Err()
}

// DeletePattern soft deletes a pattern
func (s *Storage) DeletePattern(ctx context.Context, id string) error {
	now := time.Now()
	_, err := s.db.db.ExecContext(ctx, `
		UPDATE patterns SET deleted_at = ? WHERE id = ?
	`, now.Unix(), id)
	return err
}

// MovePatternToSpace moves a pattern to a different space
func (s *Storage) MovePatternToSpace(ctx context.Context, patternID, newSpaceID string) error {
	now := time.Now()
	_, err := s.db.db.ExecContext(ctx, `
		UPDATE patterns SET space_id = ?, updated_at = ? WHERE id = ?
	`, newSpaceID, now.Unix(), patternID)
	return err
}

// UpdatePattern updates an existing pattern
func (s *Storage) UpdatePattern(ctx context.Context, p *models.Pattern) error {
	p.UpdatedAt = time.Now()

	connections, _ := json.Marshal(p.Connections)
	tags, _ := json.Marshal(p.Tags)

	_, err := s.db.db.ExecContext(ctx, `
		UPDATE patterns SET
			trigger = ?, response = ?, strength = ?, threshold = ?,
			decay_rate = ?, decay_enabled = ?, connections = ?,
			updated_at = ?, reinforcement_count = ?, decay_count = ?,
			last_used_at = ?, tags = ?, project = ?
		WHERE id = ?
	`,
		p.Trigger, p.Response, p.Strength, p.Threshold,
		p.DecayRate, p.DecayEnabled, string(connections),
		p.UpdatedAt.Unix(), p.ReinforceCnt, p.DecayCnt,
		int64TimeToPtr(p.LastUsedAt), string(tags), p.Project,
		p.ID,
	)

	return err
}

// CreateSpace creates a new space
func (s *Storage) CreateSpace(ctx context.Context, space *models.Space) error {
	now := time.Now()
	if space.CreatedAt.IsZero() {
		space.CreatedAt = now
	}
	space.UpdatedAt = now

	_, err := s.db.db.ExecContext(ctx, `
		INSERT OR REPLACE INTO spaces (id, name, description, owner, is_default, pattern_limit, pattern_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, space.ID, space.Name, space.Description, space.Owner, boolToInt(space.DefaultSpace),
		space.PatternLimit, space.PatternCount, space.CreatedAt.Unix(), space.UpdatedAt.Unix())

	return err
}

// GetSpace retrieves a space by ID
func (s *Storage) GetSpace(ctx context.Context, id string) (*models.Space, error) {
	var space models.Space
	var ownerNull sql.NullString
	var createdAt, updatedAt sql.NullInt64

	err := s.db.db.QueryRowContext(ctx, `
		SELECT id, name, description, owner, is_default, pattern_limit, pattern_count, created_at, updated_at
		FROM spaces WHERE id = ?
	`, id).Scan(
		&space.ID, &space.Name, &space.Description, &ownerNull, &space.DefaultSpace,
		&space.PatternLimit, &space.PatternCount, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("space not found: %s", id)
	}
	if err != nil {
		return nil, err
	}
	space.Owner = ownerNull.String
	space.CreatedAt = int64ToTime(createdAt)
	space.UpdatedAt = int64ToTime(updatedAt)

	return &space, nil
}

// ListSpaces lists all spaces
func (s *Storage) ListSpaces(ctx context.Context) ([]*models.Space, error) {
	rows, err := s.db.db.QueryContext(ctx, `
		SELECT id, name, description, owner, is_default, pattern_limit, pattern_count, created_at, updated_at
		FROM spaces ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spaces []*models.Space
	for rows.Next() {
		var space models.Space
		var ownerNull sql.NullString
		var createdAt, updatedAt sql.NullInt64
		err := rows.Scan(
			&space.ID, &space.Name, &space.Description, &ownerNull, &space.DefaultSpace,
			&space.PatternLimit, &space.PatternCount, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}
		space.Owner = ownerNull.String
		space.CreatedAt = int64ToTime(createdAt)
		space.UpdatedAt = int64ToTime(updatedAt)
		spaces = append(spaces, &space)
	}

	return spaces, rows.Err()
}

// UpdateSpace updates an existing space.
func (s *Storage) UpdateSpace(ctx context.Context, space *models.Space) error {
	space.UpdatedAt = time.Now()

	_, err := s.db.db.ExecContext(ctx, `
		UPDATE spaces 
		SET name = ?, description = ?, owner = ?, is_default = ?, 
		    pattern_limit = ?, pattern_count = ?, updated_at = ?
		WHERE id = ?
	`, space.Name, space.Description, space.Owner, boolToInt(space.DefaultSpace),
		space.PatternLimit, space.PatternCount, space.UpdatedAt.Unix(), space.ID)

	if err != nil {
		return fmt.Errorf("failed to update space: %w", err)
	}

	return nil
}

// DeleteSpace deletes a space by ID.
func (s *Storage) DeleteSpace(ctx context.Context, id string) error {
	result, err := s.db.db.ExecContext(ctx, `DELETE FROM spaces WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete space: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("space not found: %s", id)
	}

	return nil
}

// SetDefaultSpace sets a space as the default space.
func (s *Storage) SetDefaultSpace(ctx context.Context, id string) error {
	// First, unset all other defaults
	_, err := s.db.db.ExecContext(ctx, `UPDATE spaces SET is_default = 0`)
	if err != nil {
		return fmt.Errorf("failed to unset default spaces: %w", err)
	}

	// Then set the specified space as default
	_, err = s.db.db.ExecContext(ctx, `UPDATE spaces SET is_default = 1 WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to set default space: %w", err)
	}

	return nil
}

// GetDefaultSpace returns the default space.
func (s *Storage) GetDefaultSpace(ctx context.Context) (*models.Space, error) {
	var space models.Space
	var ownerNull sql.NullString
	var createdAt, updatedAt sql.NullInt64

	err := s.db.db.QueryRowContext(ctx, `
		SELECT id, name, description, owner, is_default, pattern_limit, pattern_count, created_at, updated_at
		FROM spaces WHERE is_default = 1
	`).Scan(
		&space.ID, &space.Name, &space.Description, &ownerNull, &space.DefaultSpace,
		&space.PatternLimit, &space.PatternCount, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		// Fallback to "global" if no default is set
		return s.GetSpace(ctx, "global")
	}
	if err != nil {
		return nil, err
	}
	space.Owner = ownerNull.String
	space.CreatedAt = int64ToTime(createdAt)
	space.UpdatedAt = int64ToTime(updatedAt)

	return &space, nil
}

// BeginTx starts a new transaction
func (s *Storage) BeginTx(ctx context.Context) (contracts.Transaction, error) {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &transaction{tx: tx}, nil
}

// Close closes the storage and releases cached statements
func (s *Storage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Close all cached statements
	for _, stmt := range s.stmtCache {
		stmt.Close()
	}
	s.stmtCache = make(map[string]*sql.Stmt)
	return s.db.Close()
}

// getStmt returns a cached prepared statement or creates a new one
func (s *Storage) getStmt(ctx context.Context, query string) (*sql.Stmt, error) {
	s.mu.RLock()
	stmt, exists := s.stmtCache[query]
	s.mu.RUnlock()

	if exists {
		return stmt, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check after acquiring write lock
	if stmt, exists := s.stmtCache[query]; exists {
		return stmt, nil
	}

	prepared, err := s.db.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	s.stmtCache[query] = prepared
	return prepared, nil
}

// SavePatternsBatch saves multiple patterns in a single transaction (Iter 44)
// This is more efficient than calling SavePattern multiple times.
func (s *Storage) SavePatternsBatch(ctx context.Context, patterns []*models.Pattern) error {
	if len(patterns) == 0 {
		return nil
	}

	// Validate all patterns first
	for _, p := range patterns {
		if err := p.Validate(); err != nil {
			return fmt.Errorf("validation failed for pattern %s: %w", p.ID, err)
		}
	}

	// Use transaction for atomicity
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT OR REPLACE INTO patterns (
			id, trigger, response, strength, threshold, decay_rate, decay_enabled,
			connections, created_at, updated_at, reinforcement_count, decay_count,
			last_used_at, tags, project, user_id, deleted_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	for _, p := range patterns {
		connections, _ := json.Marshal(p.Connections)
		tags, _ := json.Marshal(p.Tags)

		if p.CreatedAt.IsZero() {
			p.CreatedAt = now
		}
		p.UpdatedAt = now

		_, err := stmt.ExecContext(ctx,
			p.ID, p.Trigger, p.Response, p.Strength, p.Threshold, p.DecayRate, p.DecayEnabled,
			string(connections), p.CreatedAt.Unix(), p.UpdatedAt.Unix(), p.ReinforceCnt, p.DecayCnt,
			int64TimeToPtr(p.LastUsedAt), string(tags), p.Project, p.UserID, int64TimeToPtr(p.DeletedAt),
		)
		if err != nil {
			return fmt.Errorf("failed to insert pattern %s: %w", p.ID, err)
		}
	}

	return tx.Commit()
}

// DeletePatternsBatch deletes multiple patterns by their IDs (Iter 44)
func (s *Storage) DeletePatternsBatch(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `UPDATE patterns SET deleted_at = ? WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now().Unix()
	for _, id := range ids {
		_, err := stmt.ExecContext(ctx, now, id)
		if err != nil {
			return fmt.Errorf("failed to delete pattern %s: %w", id, err)
		}
	}

	return tx.Commit()
}

// UpdatePatternsBatch updates multiple patterns in a single transaction (Iter 44)
func (s *Storage) UpdatePatternsBatch(ctx context.Context, patterns []*models.Pattern) error {
	if len(patterns) == 0 {
		return nil
	}

	for _, p := range patterns {
		if err := p.Validate(); err != nil {
			return fmt.Errorf("validation failed for pattern %s: %w", p.ID, err)
		}
	}

	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE patterns SET trigger = ?, response = ?, strength = ?, threshold = ?,
			decay_rate = ?, decay_enabled = ?, connections = ?, updated_at = ?,
			reinforcement_count = ?, decay_count = ?, last_used_at = ?, tags = ?, project = ?
		WHERE id = ?
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	for _, p := range patterns {
		connections, _ := json.Marshal(p.Connections)
		tags, _ := json.Marshal(p.Tags)
		p.UpdatedAt = now

		_, err := stmt.ExecContext(ctx,
			p.Trigger, p.Response, p.Strength, p.Threshold,
			p.DecayRate, p.DecayEnabled, string(connections), p.UpdatedAt.Unix(),
			p.ReinforceCnt, p.DecayCnt, int64TimeToPtr(p.LastUsedAt), string(tags), p.Project,
			p.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update pattern %s: %w", p.ID, err)
		}
	}

	return tx.Commit()
}

// ==================== Query Optimization Methods (Iter 46) ====================

// GetPatternByTrigger retrieves a pattern by its trigger (exact match)
// Uses cached statement for better performance
func (s *Storage) GetPatternByTrigger(ctx context.Context, trigger string) (*models.Pattern, error) {
	stmt, err := s.getStmt(ctx, `
		SELECT id, trigger, response, strength, threshold, decay_rate, decay_enabled,
			connections, created_at, updated_at, reinforcement_count, decay_count,
			last_used_at, tags, project, user_id, deleted_at
		FROM patterns WHERE trigger = ? AND deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}

	var p models.Pattern
	var connections, tags sql.NullString
	var lastUsedAt, deletedAt, createdAt, updatedAt sql.NullInt64

	err = stmt.QueryRowContext(ctx, trigger).Scan(
		&p.ID, &p.Trigger, &p.Response, &p.Strength, &p.Threshold, &p.DecayRate, &p.DecayEnabled,
		&connections, &createdAt, &updatedAt, &p.ReinforceCnt, &p.DecayCnt,
		&lastUsedAt, &tags, &p.Project, &p.UserID, &deletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("pattern not found with trigger: %s", trigger)
	}
	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	if connections.Valid {
			if err := json.Unmarshal([]byte(connections.String), &p.Connections); err != nil {
				return nil, fmt.Errorf("failed to parse connections: %w", err)
			}
	}
	if tags.Valid {
			if err := json.Unmarshal([]byte(tags.String), &p.Tags); err != nil {
				return nil, fmt.Errorf("failed to parse tags: %w", err)
			}
	}
	p.CreatedAt = int64ToTime(createdAt)
	p.UpdatedAt = int64ToTime(updatedAt)
	p.LastUsedAt = int64ToTimePtr(lastUsedAt)
	p.DeletedAt = int64ToTimePtr(deletedAt)

	return &p, nil
}

// CountPatterns returns the total count of patterns matching the given filters
// More efficient than len(ListPatterns(...))
func (s *Storage) CountPatterns(ctx context.Context, opts contracts.ListOptions) (int, error) {
	query := "SELECT COUNT(*) FROM patterns WHERE deleted_at IS NULL"
	args := []interface{}{}

	if opts.Project != "" {
		query += " AND project = ?"
		args = append(args, opts.Project)
	}

	if opts.MinStrength > 0 {
		query += " AND strength >= ?"
		args = append(args, opts.MinStrength)
	}

	var count int
	err := s.db.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count patterns: %w", err)
	}

	return count, nil
}

// GetRecentlyUsedPatterns retrieves patterns ordered by last_used_at
// Useful for "frequently used" features
func (s *Storage) GetRecentlyUsedPatterns(ctx context.Context, limit int) ([]*models.Pattern, error) {
	stmt, err := s.getStmt(ctx, `
		SELECT id, trigger, response, strength, threshold, decay_rate, decay_enabled,
			connections, created_at, updated_at, reinforcement_count, decay_count,
			last_used_at, tags, project, user_id
		FROM patterns 
		WHERE deleted_at IS NULL AND last_used_at IS NOT NULL
		ORDER BY last_used_at DESC
		LIMIT ?
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPatternsRows(rows)
}

// SearchPatterns performs a full-text search on trigger and response
// Uses LIKE for simplicity - can be upgraded to FTS5 for production
func (s *Storage) SearchPatterns(ctx context.Context, query string, opts contracts.ListOptions) ([]*models.Pattern, error) {
	searchPattern := "%" + query + "%"

	baseQuery := `
		SELECT id, trigger, response, strength, threshold, decay_rate, decay_enabled,
			connections, created_at, updated_at, reinforcement_count, decay_count,
			last_used_at, tags, project, user_id
		FROM patterns 
		WHERE deleted_at IS NULL 
		AND (trigger LIKE ? OR response LIKE ?)
	`

	args := []interface{}{searchPattern, searchPattern}

	if opts.Project != "" {
		baseQuery += " AND project = ?"
		args = append(args, opts.Project)
	}

	baseQuery += " ORDER BY strength DESC"

	if opts.Limit > 0 {
		baseQuery += " LIMIT ?"
		args = append(args, opts.Limit)
	} else {
		baseQuery += " LIMIT 100" // Default limit for search
	}

	rows, err := s.db.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPatternsRows(rows)
}

// GetTopPatterns retrieves the strongest patterns (for matching priority)
func (s *Storage) GetTopPatterns(ctx context.Context, limit int) ([]*models.Pattern, error) {
	stmt, err := s.getStmt(ctx, `
		SELECT id, trigger, response, strength, threshold, decay_rate, decay_enabled,
			connections, created_at, updated_at, reinforcement_count, decay_count,
			last_used_at, tags, project, user_id
		FROM patterns 
		WHERE deleted_at IS NULL
		ORDER BY strength DESC
		LIMIT ?
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPatternsRows(rows)
}

// scanPatternsRows is a helper to scan pattern rows
func scanPatternsRows(rows *sql.Rows) ([]*models.Pattern, error) {
	var patterns []*models.Pattern
	for rows.Next() {
		var p models.Pattern
		var connections, tags sql.NullString
		var lastUsedAt, createdAt, updatedAt sql.NullInt64

		err := rows.Scan(
			&p.ID, &p.Trigger, &p.Response, &p.Strength, &p.Threshold, &p.DecayRate, &p.DecayEnabled,
			&connections, &createdAt, &updatedAt, &p.ReinforceCnt, &p.DecayCnt,
			&lastUsedAt, &tags, &p.Project, &p.UserID,
		)
		if err != nil {
			return nil, err
		}

		if connections.Valid {
			if err := json.Unmarshal([]byte(connections.String), &p.Connections); err != nil {
				return nil, fmt.Errorf("failed to parse connections: %w", err)
			}
		}
		if tags.Valid {
			if err := json.Unmarshal([]byte(tags.String), &p.Tags); err != nil {
				return nil, fmt.Errorf("failed to parse tags: %w", err)
			}
		}
		p.CreatedAt = int64ToTime(createdAt)
		p.UpdatedAt = int64ToTime(updatedAt)
		p.LastUsedAt = int64ToTimePtr(lastUsedAt)

		patterns = append(patterns, &p)
	}
	return patterns, rows.Err()
}

// Helper functions
func int64TimeToPtr(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	i := t.Unix()
	return &i
}

func int64ToTimePtr(i sql.NullInt64) *time.Time {
	if !i.Valid {
		return nil
	}
	t := time.Unix(i.Int64, 0)
	return &t
}

// ==================== Concurrency Statistics (Iter 47) ====================

// Stats returns the current storage concurrency statistics
func (s *Storage) Stats() StorageStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stats
}

// ResetStats resets the concurrency statistics counters
func (s *Storage) ResetStats() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats = StorageStats{}
}

// BatchGetPatterns retrieves multiple patterns by IDs efficiently (Iter 52)
func (s *Storage) BatchGetPatterns(ctx context.Context, ids []string) ([]*models.Pattern, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// Build query with placeholders
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := "SELECT id, trigger, response, strength, threshold, decay_rate, decay_enabled, connections, created_at, updated_at, reinforcement_count, decay_count, last_used_at, tags, project, user_id, deleted_at FROM patterns WHERE id IN (" + strings.Join(placeholders, ",") + ") AND deleted_at IS NULL"

	rows, err := s.db.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patterns []*models.Pattern
	for rows.Next() {
		var p models.Pattern
		var connections, tags sql.NullString
		var lastUsedAt, deletedAt, createdAt, updatedAt sql.NullInt64
		err := rows.Scan(
			&p.ID, &p.Trigger, &p.Response, &p.Strength, &p.Threshold, &p.DecayRate, &p.DecayEnabled,
			&connections, &createdAt, &updatedAt, &p.ReinforceCnt, &p.DecayCnt,
			&lastUsedAt, &tags, &p.Project, &p.UserID, &deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if connections.Valid {
			json.Unmarshal([]byte(connections.String), &p.Connections)
		}
		if tags.Valid {
			json.Unmarshal([]byte(tags.String), &p.Tags)
		}
		p.CreatedAt = int64ToTime(createdAt)
		p.UpdatedAt = int64ToTime(updatedAt)
		p.LastUsedAt = int64ToTimePtr(lastUsedAt)
		p.DeletedAt = int64ToTimePtr(deletedAt)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, &p)
		
		// Update cache if available
		if s.queryCache != nil {
			s.queryCache.Set(p.ID, &p)
		}
	}

	return patterns, rows.Err()
}

// ==================== Note Operations (Phase 10) ====================

// SaveNote creates or updates a note in storage.
func (s *Storage) SaveNote(ctx context.Context, note *models.Note) error {
	if note.ID == "" {
		note.ID = uuid.New().String()
	}
	note.CalculateStats()
	now := time.Now()
	note.CreatedAt = now
	note.UpdatedAt = now

	tagsJSON, _ := json.Marshal(note.Tags)

	_, err := s.db.db.ExecContext(ctx, `
		INSERT INTO notes (id, title, content, space_id, tags, is_pinned, category, word_count, char_count, last_viewed_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			content = excluded.content,
			space_id = excluded.space_id,
			tags = excluded.tags,
			is_pinned = excluded.is_pinned,
			category = excluded.category,
			word_count = excluded.word_count,
			char_count = excluded.char_count,
			last_viewed_at = excluded.last_viewed_at,
			updated_at = excluded.updated_at
	`,
		note.ID, note.Title, note.Content, note.SpaceID, tagsJSON, note.IsPinned, note.Category,
		note.WordCount, note.CharCount, note.LastViewed, note.CreatedAt.Unix(), note.UpdatedAt.Unix())

	return err
}

// GetNote retrieves a note by its ID.
func (s *Storage) GetNote(ctx context.Context, id string) (*models.Note, error) {
	var note models.Note
	var tagsJSON []byte
	var lastViewed sql.NullInt64

	err := s.db.db.QueryRowContext(ctx, `
		SELECT id, title, content, space_id, tags, is_pinned, category, word_count, char_count, last_viewed_at, created_at, updated_at
		FROM notes WHERE id = ?
	`, id).Scan(
		&note.ID, &note.Title, &note.Content, &note.SpaceID, &tagsJSON,
		&note.IsPinned, &note.Category, &note.WordCount, &note.CharCount, &lastViewed,
		&note.CreatedAt, &note.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("note not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(tagsJSON, &note.Tags)
	if lastViewed.Valid {
		t := time.Unix(lastViewed.Int64, 0)
		note.LastViewed = &t
	}

	return &note, nil
}

// ListNotes retrieves notes matching the given filter options.
func (s *Storage) ListNotes(ctx context.Context, opts contracts.ListOptions) ([]*models.Note, error) {
	query := "SELECT id, title, content, space_id, tags, is_pinned, category, word_count, char_count, last_viewed_at, created_at, updated_at FROM notes WHERE 1=1"
	args := []interface{}{}

	if opts.SpaceID != "" {
		query += " AND space_id = ?"
		args = append(args, opts.SpaceID)
	}

	if opts.Project != "" { // Use Project as category filter
		query += " AND category = ?"
		args = append(args, opts.Project)
	}

	query += " ORDER BY is_pinned DESC, updated_at DESC"

	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}
	if opts.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", opts.Offset)
	}

	rows, err := s.db.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		var note models.Note
		var tagsJSON []byte
		var lastViewed sql.NullInt64
		var createdAt, updatedAt sql.NullInt64

		err := rows.Scan(
			&note.ID, &note.Title, &note.Content, &note.SpaceID, &tagsJSON,
			&note.IsPinned, &note.Category, &note.WordCount, &note.CharCount, &lastViewed,
			&createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(tagsJSON, &note.Tags)
		if lastViewed.Valid {
			t := time.Unix(lastViewed.Int64, 0)
			note.LastViewed = &t
		}
		note.CreatedAt = time.Unix(createdAt.Int64, 0)
		note.UpdatedAt = time.Unix(updatedAt.Int64, 0)

		notes = append(notes, &note)
	}

	return notes, rows.Err()
}

// DeleteNote removes a note by its ID.
func (s *Storage) DeleteNote(ctx context.Context, id string) error {
	result, err := s.db.db.ExecContext(ctx, "DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}

// UpdateNote updates an existing note.
func (s *Storage) UpdateNote(ctx context.Context, note *models.Note) error {
	note.CalculateStats()
	note.UpdatedAt = time.Now()

	tagsJSON, _ := json.Marshal(note.Tags)

	result, err := s.db.db.ExecContext(ctx, `
		UPDATE notes SET
			title = ?,
			content = ?,
			space_id = ?,
			tags = ?,
			is_pinned = ?,
			category = ?,
			word_count = ?,
			char_count = ?,
			last_viewed_at = ?,
			updated_at = ?
		WHERE id = ?
	`,
		note.Title, note.Content, note.SpaceID, tagsJSON, note.IsPinned, note.Category,
		note.WordCount, note.CharCount, note.LastViewed, note.UpdatedAt.Unix(), note.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}

// SearchNotes performs a full-text search on title and content.
func (s *Storage) SearchNotes(ctx context.Context, query string, opts contracts.ListOptions) ([]*models.Note, error) {
	searchQuery := "%" + query + "%"
	sqlQuery := "SELECT id, title, content, space_id, tags, is_pinned, category, word_count, char_count, last_viewed_at, created_at, updated_at FROM notes WHERE (title LIKE ? OR content LIKE ?)"
	args := []interface{}{searchQuery, searchQuery}

	if opts.SpaceID != "" {
		sqlQuery += " AND space_id = ?"
		args = append(args, opts.SpaceID)
	}

	sqlQuery += " ORDER BY is_pinned DESC, updated_at DESC"

	if opts.Limit > 0 {
		sqlQuery += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}

	rows, err := s.db.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		var note models.Note
		var tagsJSON []byte
		var lastViewed sql.NullInt64

		err := rows.Scan(
			&note.ID, &note.Title, &note.Content, &note.SpaceID, &tagsJSON,
			&note.IsPinned, &note.Category, &note.WordCount, &note.CharCount, &lastViewed,
			&note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(tagsJSON, &note.Tags)
		if lastViewed.Valid {
			t := time.Unix(lastViewed.Int64, 0)
			note.LastViewed = &t
		}

		notes = append(notes, &note)
	}

	return notes, rows.Err()
}


// transaction implements contracts.Transaction
type transaction struct {
	tx *sql.Tx
}

func (t *transaction) Commit() error {
	return t.tx.Commit()
}

func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}

func int64ToTime(i sql.NullInt64) time.Time {
	if !i.Valid {
		return time.Time{}
	}
	return time.Unix(i.Int64, 0)
}
