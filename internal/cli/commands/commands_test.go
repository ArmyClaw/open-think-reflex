package commands

import (
	"context"
	"testing"

	"github.com/ArmyClaw/open-think-reflex/internal/data/sqlite"
	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestStorage creates a test in-memory SQLite storage
func setupTestStorage(t *testing.T) *sqlite.Storage {
	db, err := sqlite.NewDatabase(":memory:")
	require.NoError(t, err)
	
	// Run migrations
	err = db.Migrate(context.Background())
	require.NoError(t, err)
	
	storage := sqlite.NewStorage(db)
	return storage
}

func TestListPatterns_Empty(t *testing.T) {
	storage := setupTestStorage(t)
	err := listPatterns(storage)
	assert.NoError(t, err)
}

func TestListPatterns_WithData(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	// Create test patterns
	patterns := []*models.Pattern{
		models.NewPattern("test-trigger-1", "test-response-1"),
		models.NewPattern("test-trigger-2", "test-response-2"),
		models.NewPattern("test-trigger-3", "test-response-3"),
	}

	for _, p := range patterns {
		p.Project = "test-project"
		err := storage.SavePattern(ctx, p)
		require.NoError(t, err)
	}

	// List patterns
	err := listPatterns(storage)
	assert.NoError(t, err)
}

func TestCreatePattern(t *testing.T) {
	storage := setupTestStorage(t)

	err := createPattern(storage, "new-trigger", "new-response", "test-project")
	assert.NoError(t, err)

	// Verify pattern was created
	ctx := context.Background()
	patterns, err := storage.ListPatterns(ctx, contracts.ListOptions{Limit: 10})
	assert.NoError(t, err)
	assert.Len(t, patterns, 1)
	assert.Equal(t, "new-trigger", patterns[0].Trigger)
	assert.Equal(t, "new-response", patterns[0].Response)
	assert.Equal(t, "test-project", patterns[0].Project)
}

func TestCreatePattern_EmptyTrigger(t *testing.T) {
	storage := setupTestStorage(t)

	err := createPattern(storage, "", "response", "")
	// Empty trigger triggers validation error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trigger cannot be empty")
}

func TestShowPattern(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	// Create a test pattern
	pattern := models.NewPattern("show-trigger", "show-response")
	pattern.Project = "show-project"
	err := storage.SavePattern(ctx, pattern)
	require.NoError(t, err)

	// Show pattern
	err = showPattern(storage, pattern.ID)
	assert.NoError(t, err)
}

func TestShowPattern_EmptyID(t *testing.T) {
	storage := setupTestStorage(t)

	err := showPattern(storage, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pattern ID required")
}

func TestShowPattern_NotFound(t *testing.T) {
	storage := setupTestStorage(t)

	err := showPattern(storage, "non-existent-id")
	assert.Error(t, err)
}

func TestDeletePattern(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	// Create a test pattern
	pattern := models.NewPattern("delete-trigger", "delete-response")
	err := storage.SavePattern(ctx, pattern)
	require.NoError(t, err)

	// Delete pattern
	err = deletePattern(storage, pattern.ID)
	assert.NoError(t, err)

	// Verify pattern was deleted
	_, err = storage.GetPattern(ctx, pattern.ID)
	assert.Error(t, err)
}

func TestDeletePattern_EmptyID(t *testing.T) {
	storage := setupTestStorage(t)

	err := deletePattern(storage, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pattern ID required")
}

func TestDeletePattern_NotFound(t *testing.T) {
	storage := setupTestStorage(t)

	// DeletePattern does not check if pattern exists
	err := deletePattern(storage, "non-existent-id")
	assert.NoError(t, err)
}
