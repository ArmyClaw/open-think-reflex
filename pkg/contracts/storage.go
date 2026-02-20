package contracts

import (
	"context"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
)

// Storage defines the interface for pattern persistence
type Storage interface {
	// Pattern operations
	SavePattern(ctx context.Context, p *models.Pattern) error
	GetPattern(ctx context.Context, id string) (*models.Pattern, error)
	ListPatterns(ctx context.Context, opts ListOptions) ([]*models.Pattern, error)
	DeletePattern(ctx context.Context, id string) error
	UpdatePattern(ctx context.Context, p *models.Pattern) error

	// Space operations
	CreateSpace(ctx context.Context, s *models.Space) error
	GetSpace(ctx context.Context, id string) (*models.Space, error)
	ListSpaces(ctx context.Context) ([]*models.Space, error)

	// Transaction support
	BeginTx(ctx context.Context) (Transaction, error)

	// Close
	Close() error
}

// ListOptions contains filtering options for pattern listing
type ListOptions struct {
	Tags         []string
	Project      string
	MinStrength  float64
	Limit        int
	Offset       int
	IncludeDeleted bool
}

// Transaction defines the interface for database transactions
type Transaction interface {
	Commit() error
	Rollback() error
}
