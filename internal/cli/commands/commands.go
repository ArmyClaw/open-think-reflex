// Package commands provides CLI command implementations for the otr tool.
// Uses urfave/cli v2 for command-line interface construction.
package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/ArmyClaw/open-think-reflex/internal/config"
	"github.com/ArmyClaw/open-think-reflex/internal/data/sqlite"
	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/export"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/urfave/cli/v2"
)

// BuildCommands constructs all CLI commands and returns them as a slice.
// Each command is registered with appropriate flags and action handlers.
//
// Commands include:
//   - pattern: Manage reflex patterns (list, create, show, delete)
//   - version: Display version information
func BuildCommands(storage *sqlite.Storage, cfg *config.Config, loader *config.Loader) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "pattern",
			Usage: "Manage reflex patterns",
			Subcommands: []*cli.Command{
				{
					Name:  "list",
					Usage: "List all patterns",
					Action: func(c *cli.Context) error {
						return listPatterns(storage)
					},
				},
				{
					Name:  "create",
					Usage: "Create a new pattern",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "trigger",
							Required: true,
							Usage:    "Pattern trigger (keyword that activates the reflex)",
						},
						&cli.StringFlag{
							Name:     "response",
							Required: true,
							Usage:    "Pattern response (content returned when triggered)",
						},
						&cli.StringFlag{
							Name:  "project",
							Usage: "Project name for organization",
						},
					},
					Action: func(c *cli.Context) error {
						return createPattern(storage, c.String("trigger"), c.String("response"), c.String("project"))
					},
				},
				{
					Name:      "show",
					Usage:     "Show pattern details",
					ArgsUsage: "<pattern_id>",
					Action: func(c *cli.Context) error {
						return showPattern(storage, c.Args().First())
					},
				},
				{
					Name:      "delete",
					Usage:     "Delete a pattern",
					ArgsUsage: "<pattern_id>",
					Action: func(c *cli.Context) error {
						return deletePattern(storage, c.Args().First())
					},
				},
			},
		},
		{
			Name:  "space",
			Usage: "Manage pattern spaces",
			Subcommands: []*cli.Command{
				{
					Name:  "list",
					Usage: "List all spaces",
					Action: func(c *cli.Context) error {
						return ListSpaces(storage)
					},
				},
				{
					Name:  "create",
					Usage: "Create a new space",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "name",
							Required: true,
							Usage:    "Space name",
						},
						&cli.StringFlag{
							Name:  "description",
							Usage: "Space description",
						},
					},
					Action: func(c *cli.Context) error {
						return CreateSpace(storage, c.String("name"), c.String("description"))
					},
				},
				{
					Name:      "show",
					Usage:     "Show space details",
					ArgsUsage: "<space_id>",
					Action: func(c *cli.Context) error {
						return ShowSpace(storage, c.Args().First())
					},
				},
				{
					Name:      "delete",
					Usage:     "Delete a space",
					ArgsUsage: "<space_id>",
					Action: func(c *cli.Context) error {
						return DeleteSpace(storage, c.Args().First())
					},
				},
				{
					Name:      "use",
					Usage:     "Switch to a space",
					ArgsUsage: "<space_id>",
					Action: func(c *cli.Context) error {
						return UseSpace(storage, cfg, loader, c.Args().First())
					},
				},
			},
		},
		{
			Name:  "version",
			Usage: "Show version information",
			Action: func(c *cli.Context) error {
				fmt.Println("Open-Think-Reflex v1.0")
				return nil
			},
		},
		{
			Name:  "export",
			Usage: "Export patterns to a file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "output",
					Required: true,
					Usage:    "Output file path (JSON format)",
				},
				&cli.StringFlag{
					Name:  "project",
					Usage: "Filter by project name",
				},
			},
			Action: func(c *cli.Context) error {
				return exportPatterns(storage, c.String("output"), c.String("project"))
			},
		},
		{
			Name:  "import",
			Usage: "Import patterns from a file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "input",
					Required: true,
					Usage:    "Input file path (JSON format)",
				},
				&cli.BoolFlag{
					Name:  "force",
					Usage: "Overwrite existing patterns with same ID",
				},
			},
			Action: func(c *cli.Context) error {
				return importPatterns(storage, c.String("input"), c.Bool("force"))
			},
		},
	}
}

// listPatterns retrieves and displays all patterns.
func listPatterns(storage *sqlite.Storage) error {
	ctx := context.Background()
	patterns, err := storage.ListPatterns(ctx, contracts.ListOptions{Limit: 100})
	if err != nil {
		return fmt.Errorf("failed to list patterns: %w", err)
	}

	if len(patterns) == 0 {
		fmt.Println("No patterns found")
		return nil
	}

	fmt.Printf("Found %d patterns:\n\n", len(patterns))
	for _, p := range patterns {
		// Format: ID(8chars) Trigger Strength%
		fmt.Printf("  %s  %s (strength: %.1f)\n",
			p.ID[:min(8, len(p.ID))],
			p.Trigger,
			p.Strength)
	}

	return nil
}

// createPattern creates a new pattern with the given trigger and response.
func createPattern(storage *sqlite.Storage, trigger, response, project string) error {
	ctx := context.Background()
	pattern := models.NewPattern(trigger, response)
	pattern.Project = project

	if err := storage.SavePattern(ctx, pattern); err != nil {
		return fmt.Errorf("failed to create pattern: %w", err)
	}

	fmt.Printf("Pattern created: %s\n", pattern.ID)
	return nil
}

// showPattern displays detailed information about a specific pattern.
func showPattern(storage *sqlite.Storage, id string) error {
	if id == "" {
		return fmt.Errorf("pattern ID required")
	}

	ctx := context.Background()
	pattern, err := storage.GetPattern(ctx, id)
	if err != nil {
		return fmt.Errorf("pattern not found: %w", err)
	}

	timeFormat := "2006-01-02 15:04:05"

	fmt.Printf("Pattern: %s\n", pattern.ID)
	fmt.Printf("  Trigger:   %s\n", pattern.Trigger)
	fmt.Printf("  Response:  %s\n", pattern.Response)
	fmt.Printf("  Strength:  %.1f / %.1f\n", pattern.Strength, pattern.Threshold)
	fmt.Printf("  Project:   %s\n", pattern.Project)
	fmt.Printf("  Created:   %s\n", pattern.CreatedAt.Format(timeFormat))
	fmt.Printf("  Updated:   %s\n", pattern.UpdatedAt.Format(timeFormat))
	fmt.Printf("  Reinforced: %d times\n", pattern.ReinforceCnt)
	fmt.Printf("  Decayed:    %d times\n", pattern.DecayCnt)

	if pattern.LastUsedAt != nil {
		fmt.Printf("  Last Used: %s\n", pattern.LastUsedAt.Format(timeFormat))
	}

	return nil
}

// deletePattern removes a pattern from storage.
func deletePattern(storage *sqlite.Storage, id string) error {
	if id == "" {
		return fmt.Errorf("pattern ID required")
	}

	ctx := context.Background()
	if err := storage.DeletePattern(ctx, id); err != nil {
		return fmt.Errorf("failed to delete pattern: %w", err)
	}

	fmt.Printf("Pattern deleted: %s\n", id)
	return nil
}

// ListSpaces retrieves and displays all spaces.
func ListSpaces(storage *sqlite.Storage) error {
	ctx := context.Background()
	spaces, err := storage.ListSpaces(ctx)
	if err != nil {
		return fmt.Errorf("failed to list spaces: %w", err)
	}

	if len(spaces) == 0 {
		fmt.Println("No spaces found")
		return nil
	}

	fmt.Printf("Found %d spaces:\n\n", len(spaces))
	for _, s := range spaces {
		fmt.Printf("  %s  %s\n", s.ID[:min(8, len(s.ID))], s.Name)
	}

	return nil
}

// CreateSpace creates a new space for organizing patterns.
func CreateSpace(storage *sqlite.Storage, name, description string) error {
	ctx := context.Background()
	space := &models.Space{
		ID:          generateSpaceID(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := storage.CreateSpace(ctx, space); err != nil {
		return fmt.Errorf("failed to create space: %w", err)
	}

	fmt.Printf("Space created: %s\n", space.ID)
	return nil
}

// ShowSpace displays details of a specific space.
func ShowSpace(storage *sqlite.Storage, spaceID string) error {
	if spaceID == "" {
		return fmt.Errorf("space ID is required")
	}
	
	ctx := context.Background()
	space, err := storage.GetSpace(ctx, spaceID)
	if err != nil {
		return fmt.Errorf("failed to get space: %w", err)
	}

	fmt.Printf("Space Details:\n")
	fmt.Printf("  ID:          %s\n", space.ID)
	fmt.Printf("  Name:        %s\n", space.Name)
	fmt.Printf("  Description: %s\n", space.Description)
	fmt.Printf("  Owner:       %s\n", space.Owner)
	fmt.Printf("  Default:     %v\n", space.DefaultSpace)
	fmt.Printf("  Pattern Limit: %d\n", space.PatternLimit)
	fmt.Printf("  Pattern Count: %d\n", space.PatternCount)
	fmt.Printf("  Created:     %s\n", space.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Updated:     %s\n", space.UpdatedAt.Format("2006-01-02 15:04:05"))
	
	return nil
}

// DeleteSpace deletes a space by ID.
func DeleteSpace(storage *sqlite.Storage, spaceID string) error {
	if spaceID == "" {
		return fmt.Errorf("space ID is required")
	}
	
	ctx := context.Background()
	
	// Check if space exists
	_, err := storage.GetSpace(ctx, spaceID)
	if err != nil {
		return fmt.Errorf("space not found: %s", spaceID)
	}

	// Delete the space
	if err := storage.DeleteSpace(ctx, spaceID); err != nil {
		return fmt.Errorf("failed to delete space: %w", err)
	}
	
	fmt.Printf("Space '%s' deleted\n", spaceID)
	return nil
}

// UseSpace switches to a different space and saves to config if loader is provided.
func UseSpace(storage *sqlite.Storage, cfg *config.Config, loader *config.Loader, spaceID string) error {
	if spaceID == "" {
		return fmt.Errorf("space ID is required")
	}
	
	ctx := context.Background()
	
	// Verify space exists
	space, err := storage.GetSpace(ctx, spaceID)
	if err != nil {
		return fmt.Errorf("failed to get space: %w", err)
	}

	// Save current space to config if loader is provided
	if cfg != nil && loader != nil {
		cfg.CurrentSpace = spaceID
		if err := loader.Save(cfg); err != nil {
			fmt.Printf("Warning: failed to save current space: %v\n", err)
			fmt.Printf("Switched to space: %s (%s)\n", space.Name, space.ID)
		} else {
			fmt.Printf("Switched to space: %s (%s)\n", space.Name, space.ID)
			fmt.Println("Current space saved to config")
		}
		return nil
	}

	fmt.Printf("Switched to space: %s (%s)\n", space.Name, space.ID)
	return nil
}

// SetDefaultSpace sets a space as the default space.
func SetDefaultSpace(storage *sqlite.Storage, spaceID string) error {
	if spaceID == "" {
		return fmt.Errorf("space ID is required")
	}
	
	ctx := context.Background()
	
	// Verify space exists
	_, err := storage.GetSpace(ctx, spaceID)
	if err != nil {
		return fmt.Errorf("failed to get space: %w", err)
	}

	// Set as default
	if err := storage.SetDefaultSpace(ctx, spaceID); err != nil {
		return fmt.Errorf("failed to set default space: %w", err)
	}

	fmt.Printf("Space '%s' is now the default\n", spaceID)
	return nil
}

// generateSpaceID generates a new space ID.
func generateSpaceID() string {
	return fmt.Sprintf("space-%d", time.Now().UnixNano())
}

// exportPatterns exports patterns to a JSON file.
func exportPatterns(storage *sqlite.Storage, outputPath, projectFilter string) error {
	ctx := context.Background()

	var patterns []*models.Pattern
	var err error

	if projectFilter != "" {
		// Filter by project
		allPatterns, err := storage.ListPatterns(ctx, contracts.ListOptions{Limit: 10000})
		if err != nil {
			return fmt.Errorf("failed to list patterns: %w", err)
		}
		for _, p := range allPatterns {
			if p.Project == projectFilter {
				patterns = append(patterns, p)
			}
		}
	} else {
		// Export all
		patterns, err = storage.ListPatterns(ctx, contracts.ListOptions{Limit: 10000})
		if err != nil {
			return fmt.Errorf("failed to list patterns: %w", err)
		}
	}

	exporter := export.NewExporter()
	if err := exporter.ExportToJSON(ctx, patterns, outputPath); err != nil {
		return fmt.Errorf("failed to export: %w", err)
	}

	fmt.Printf("Exported %d patterns to %s\n", len(patterns), outputPath)
	return nil
}

// importPatterns imports patterns from a JSON file.
func importPatterns(storage *sqlite.Storage, inputPath string, force bool) error {
	ctx := context.Background()

	importer := export.NewImporter()
	importData, err := importer.ImportFromJSON(ctx, inputPath)
	if err != nil {
		return fmt.Errorf("failed to import: %w", err)
	}

	imported := 0
	skipped := 0

	for _, p := range importData.Patterns {
		// Check if pattern already exists
		existing, err := storage.GetPattern(ctx, p.ID)
		if err == nil && existing != nil {
			if force {
				// Update existing pattern
				if err := storage.SavePattern(ctx, &p); err != nil {
					return fmt.Errorf("failed to update pattern %s: %w", p.ID, err)
				}
				imported++
			} else {
				fmt.Printf("Skipping existing pattern: %s (use --force to overwrite)\n", p.ID)
				skipped++
			}
			continue
		}

		// Create new pattern
		if err := storage.SavePattern(ctx, &p); err != nil {
			return fmt.Errorf("failed to save pattern %s: %w", p.ID, err)
		}
		imported++
	}

	fmt.Printf("Imported %d patterns from %s\n", imported, inputPath)
	if skipped > 0 {
		fmt.Printf("Skipped %d existing patterns\n", skipped)
	}
	return nil
}
