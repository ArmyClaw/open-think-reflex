package sqlite

import (
	"context"
	"testing"
)

// TestIndexOptimization verifies that the new indexes are created correctly
// and can be used by the query planner.
func TestIndexOptimization(t *testing.T) {
	ctx := context.Background()
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	// Run migrations (includes index creation)
	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// Verify indexes exist
	indexes := []string{
		"idx_patterns_trigger",
		"idx_patterns_strength",
		"idx_patterns_project",
		"idx_patterns_tags",
		"idx_patterns_deleted",
		// New indexes from Iter 43
		"idx_patterns_last_used_at",
		"idx_patterns_updated_at",
		"idx_patterns_project_deleted",
		"idx_patterns_strength_threshold",
		"idx_patterns_decay_enabled",
	}

	for _, idx := range indexes {
		var count int
		err := db.db.QueryRowContext(ctx, 
			"SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name=?", idx).Scan(&count)
		if err != nil {
			t.Errorf("failed to check index %s: %v", idx, err)
			continue
		}
		if count == 0 {
			t.Errorf("index %s was not created", idx)
		} else {
			t.Logf("✅ Index %s exists", idx)
		}
	}
}

// TestQueryPlanAnalysis shows the query execution plan for common queries
func TestQueryPlanAnalysis(t *testing.T) {
	ctx := context.Background()
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// Test query plans
	queries := []struct {
		name  string
		query string
	}{
		{
			name: "ListPatterns with project filter",
			query: "EXPLAIN QUERY PLAN SELECT * FROM patterns WHERE project = 'test' AND deleted_at IS NULL ORDER BY updated_at DESC LIMIT 10",
		},
		{
			name: "Filter by strength threshold",
			query: "EXPLAIN QUERY PLAN SELECT * FROM patterns WHERE strength >= 50 AND threshold <= 50 AND deleted_at IS NULL",
		},
		{
			name: "Filter by decay_enabled",
			query: "EXPLAIN QUERY PLAN SELECT * FROM patterns WHERE decay_enabled = 1 AND deleted_at IS NULL",
		},
	}

	for _, q := range queries {
		rows, err := db.db.QueryContext(ctx, q.query)
		if err != nil {
			t.Errorf("failed to explain query %s: %v", q.name, err)
			continue
		}

		t.Logf("=== %s ===", q.name)
		for rows.Next() {
			var id, parent, notused, detail string
			rows.Scan(&id, &parent, &notused, &detail)
			t.Logf("  %s | %s | %s", id, parent, detail)
		}
		rows.Close()
	}
}
