// Package commands provides CLI command implementations for pattern management.
package commands

import (
	"context"
	"fmt"

	"github.com/ArmyClaw/open-think-reflex/internal/data/sqlite"
	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
)

// ListPatterns lists all patterns
func ListPatterns(storage *sqlite.Storage) error {
	ctx := context.Background()
	patterns, err := storage.ListPatterns(ctx, contracts.ListOptions{Limit: 100})
	if err != nil {
		return err
	}

	if len(patterns) == 0 {
		fmt.Println("No patterns found")
		return nil
	}

	fmt.Printf("Found %d patterns:\n\n", len(patterns))
	for _, p := range patterns {
		fmt.Printf("  %s  %s (strength: %.1f / %.1f)\n", 
			p.ID[:min(8, len(p.ID))], p.Trigger, p.Strength, p.Threshold)
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
