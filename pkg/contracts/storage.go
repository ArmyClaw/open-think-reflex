// Package contracts defines the core interfaces for Open-Think-Reflex.
// These interfaces establish the contracts between different layers of the application.
package contracts

import (
	"context"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// Storage defines the interface for pattern and space persistence.
// Implementations must be thread-safe and handle concurrent access.
type Storage interface {
	// ==================== Pattern Operations ====================

	// SavePattern creates or updates a pattern in storage.
	// If the pattern has an empty ID, a new ID will be generated.
	SavePattern(ctx context.Context, p *models.Pattern) error

	// GetPattern retrieves a pattern by its ID.
	// Returns ErrNotFound if no pattern exists with the given ID.
	GetPattern(ctx context.Context, id string) (*models.Pattern, error)

	// ListPatterns retrieves patterns matching the given filter options.
	// Results are ordered by creation time (newest first).
	ListPatterns(ctx context.Context, opts ListOptions) ([]*models.Pattern, error)

	// DeletePattern removes a pattern by its ID.
	// Returns ErrNotFound if no pattern exists with the given ID.
	DeletePattern(ctx context.Context, id string) error

	// MovePatternToSpace moves a pattern to a different space.
	MovePatternToSpace(ctx context.Context, patternID, newSpaceID string) error

	// UpdatePattern updates an existing pattern.
	// Returns ErrNotFound if no pattern exists with the given ID.
	UpdatePattern(ctx context.Context, p *models.Pattern) error

	// ==================== Query Optimization Methods (Iter 46) ====================

	// GetPatternByTrigger retrieves a pattern by its trigger (exact match).
	// Uses cached statement for better performance.
	GetPatternByTrigger(ctx context.Context, trigger string) (*models.Pattern, error)

	// CountPatterns returns the total count of patterns matching filters.
	// More efficient than len(ListPatterns(...)).
	CountPatterns(ctx context.Context, opts ListOptions) (int, error)

	// GetRecentlyUsedPatterns retrieves patterns ordered by last_used_at.
	// Useful for "frequently used" features.
	GetRecentlyUsedPatterns(ctx context.Context, limit int) ([]*models.Pattern, error)

	// SearchPatterns performs a full-text search on trigger and response.
	SearchPatterns(ctx context.Context, query string, opts ListOptions) ([]*models.Pattern, error)

	// GetTopPatterns retrieves the strongest patterns (for matching priority).
	GetTopPatterns(ctx context.Context, limit int) ([]*models.Pattern, error)

	// ==================== Batch Operations (Iter 44) ====================

	// SavePatternsBatch saves multiple patterns in a single transaction.
	// More efficient than calling SavePattern multiple times.
	// Returns error if any pattern fails validation.
	SavePatternsBatch(ctx context.Context, patterns []*models.Pattern) error

	// DeletePatternsBatch deletes multiple patterns by their IDs.
	// Returns error if any deletion fails.
	DeletePatternsBatch(ctx context.Context, ids []string) error

	// UpdatePatternsBatch updates multiple patterns in a single transaction.
	// More efficient than calling UpdatePattern multiple times.
	UpdatePatternsBatch(ctx context.Context, patterns []*models.Pattern) error

	// ==================== Space Operations ====================

	// CreateSpace creates a new space for organizing patterns.
	CreateSpace(ctx context.Context, s *models.Space) error

	// GetSpace retrieves a space by its ID.
	// Returns ErrNotFound if no space exists with the given ID.
	GetSpace(ctx context.Context, id string) (*models.Space, error)

	// ListSpaces retrieves all spaces in the system.
	ListSpaces(ctx context.Context) ([]*models.Space, error)

	// UpdateSpace updates an existing space.
	UpdateSpace(ctx context.Context, s *models.Space) error

	// DeleteSpace deletes a space by ID.
	DeleteSpace(ctx context.Context, id string) error

	// SetDefaultSpace sets a space as the default space.
	SetDefaultSpace(ctx context.Context, id string) error

	// GetDefaultSpace returns the default space.
	GetDefaultSpace(ctx context.Context) (*models.Space, error)

	// ==================== Transaction Support ====================

	// BeginTx starts a new database transaction.
	// The caller must explicitly commit or rollback the transaction.
	BeginTx(ctx context.Context) (Transaction, error)

	// Close releases all resources held by the storage.
	// After Close returns, the storage should not be used.
	Close() error
}

// ListOptions contains filtering options for pattern listing.
// All fields are optional - zero values are ignored.
type ListOptions struct {
	// Tags filters patterns containing any of these tags.
	// Default: nil (no filter)
	Tags []string

	// Project filters patterns by project name.
	// Default: "" (no filter)
	Project string

	// SpaceID filters patterns by space (v2.0).
	// Default: "" (no filter)
	SpaceID string

	// MinStrength filters patterns with strength >= this value.
	// Default: 0.0 (no filter)
	MinStrength float64

	// Limit restricts the maximum number of results.
	// Default: 0 (no limit)
	Limit int

	// Offset specifies the number of results to skip.
	// Use with Limit for pagination.
	// Default: 0
	Offset int

	// IncludeDeleted specifies whether to include deleted patterns.
	// Default: false
	IncludeDeleted bool
}

// Transaction defines the interface for database transactions.
// Transactions provide atomicity - either all operations succeed
// or none are applied.
type Transaction interface {
	// Commit applies all changes made within the transaction.
	// After Commit returns successfully, the transaction is closed.
	Commit() error

	// Rollback undoes all changes made within the transaction.
	// After Rollback returns, the transaction is closed.
	Rollback() error
}
