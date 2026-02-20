package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// Storage implements contracts.Storage using SQLite
type Storage struct {
	db *Database
}

// NewStorage creates a new SQLite storage
func NewStorage(db *Database) *Storage {
	return &Storage{db: db}
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

	_, err := s.db.db.ExecContext(ctx, `
		INSERT OR REPLACE INTO patterns (
			id, trigger, response, strength, threshold, decay_rate, decay_enabled,
			connections, created_at, updated_at, reinforcement_count, decay_count,
			last_used_at, tags, project, user_id, deleted_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		p.ID, p.Trigger, p.Response, p.Strength, p.Threshold, p.DecayRate, p.DecayEnabled,
		string(connections), p.CreatedAt.Unix(), p.UpdatedAt.Unix(), p.ReinforceCnt, p.DecayCnt,
		int64TimeToPtr(p.LastUsedAt), string(tags), p.Project, p.UserID, int64TimeToPtr(p.DeletedAt),
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
		json.Unmarshal([]byte(connections.String), &p.Connections)
	}
	if tags.Valid {
		json.Unmarshal([]byte(tags.String), &p.Tags)
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
		last_used_at, tags, project, user_id
		FROM patterns WHERE deleted_at IS NULL`
	args := []interface{}{}

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
			json.Unmarshal([]byte(connections.String), &p.Connections)
		}
		if tags.Valid {
			json.Unmarshal([]byte(tags.String), &p.Tags)
		}
		p.CreatedAt = int64ToTime(createdAt)
		p.UpdatedAt = int64ToTime(updatedAt)
		p.LastUsedAt = int64ToTimePtr(lastUsedAt)

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
		INSERT OR REPLACE INTO spaces (id, name, description, is_default, pattern_limit, pattern_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, space.ID, space.Name, space.Description, boolToInt(space.DefaultSpace),
		space.PatternLimit, space.PatternCount, space.CreatedAt.Unix(), space.UpdatedAt.Unix())

	return err
}

// GetSpace retrieves a space by ID
func (s *Storage) GetSpace(ctx context.Context, id string) (*models.Space, error) {
	var space models.Space

	err := s.db.db.QueryRowContext(ctx, `
		SELECT id, name, description, is_default, pattern_limit, pattern_count, created_at, updated_at
		FROM spaces WHERE id = ?
	`, id).Scan(
		&space.ID, &space.Name, &space.Description, &space.DefaultSpace,
		&space.PatternLimit, &space.PatternCount, &space.CreatedAt, &space.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("space not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	return &space, nil
}

// ListSpaces lists all spaces
func (s *Storage) ListSpaces(ctx context.Context) ([]*models.Space, error) {
	rows, err := s.db.db.QueryContext(ctx, `
		SELECT id, name, description, is_default, pattern_limit, pattern_count, created_at, updated_at
		FROM spaces ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spaces []*models.Space
	for rows.Next() {
		var space models.Space
		err := rows.Scan(
			&space.ID, &space.Name, &space.Description, &space.DefaultSpace,
			&space.PatternLimit, &space.PatternCount, &space.CreatedAt, &space.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		spaces = append(spaces, &space)
	}

	return spaces, rows.Err()
}

// BeginTx starts a new transaction
func (s *Storage) BeginTx(ctx context.Context) (contracts.Transaction, error) {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &transaction{tx: tx}, nil
}

// Close closes the storage
func (s *Storage) Close() error {
	return s.db.Close()
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
