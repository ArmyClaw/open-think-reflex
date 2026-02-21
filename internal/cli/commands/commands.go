// Package commands provides CLI command implementations for the otr tool.
// Uses urfave/cli v2 for command-line interface construction.
package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/ArmyClaw/open-think-reflex/internal/data/sqlite"
	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/urfave/cli/v2"
)

// BuildCommands constructs all CLI commands and returns them as a slice.
// Each command is registered with appropriate flags and action handlers.
//
// Commands include:
//   - pattern: Manage reflex patterns (list, create, show, delete)
//   - version: Display version information
func BuildCommands(storage *sqlite.Storage) []*cli.Command {
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
						return listSpaces(storage)
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
						return createSpace(storage, c.String("name"), c.String("description"))
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

// listSpaces retrieves and displays all spaces.
func listSpaces(storage *sqlite.Storage) error {
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

// createSpace creates a new space for organizing patterns.
func createSpace(storage *sqlite.Storage, name, description string) error {
	ctx := context.Background()
	space := &models.Space{
		ID:          models.NewSpaceID(),
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
