package commands

import (
	"context"
	"fmt"

	"github.com/ArmyClaw/open-think-reflex/internal/config"
	"github.com/ArmyClaw/open-think-reflex/internal/data/sqlite"
	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/urfave/cli/v2"
)

// BuildCommands builds all CLI commands
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
							Usage:    "Pattern trigger",
						},
						&cli.StringFlag{
							Name:     "response",
							Required: true,
							Usage:    "Pattern response",
						},
						&cli.StringFlag{
							Name:  "project",
							Usage: "Project name",
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
			Name:  "version",
			Usage: "Show version information",
			Action: func(c *cli.Context) error {
				fmt.Println("Open-Think-Reflex v1.0-dev")
				return nil
			},
		},
	}
}

// listPatterns lists all patterns
func listPatterns(storage *sqlite.Storage) error {
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
		fmt.Printf("  %s  %s (strength: %.1f)\n", p.ID[:8], p.Trigger, p.Strength)
	}

	return nil
}

// createPattern creates a new pattern
func createPattern(storage *sqlite.Storage, trigger, response, project string) error {
	ctx := context.Background()
	pattern := models.NewPattern(trigger, response)
	pattern.Project = project

	if err := storage.SavePattern(ctx, pattern); err != nil {
		return err
	}

	fmt.Printf("Pattern created: %s\n", pattern.ID)
	return nil
}

// showPattern shows pattern details
func showPattern(storage *sqlite.Storage, id string) error {
	if id == "" {
		return fmt.Errorf("pattern ID required")
	}

	ctx := context.Background()
	pattern, err := storage.GetPattern(ctx, id)
	if err != nil {
		return err
	}

	fmt.Printf("Pattern: %s\n", pattern.ID)
	fmt.Printf("  Trigger: %s\n", pattern.Trigger)
	fmt.Printf("  Response: %s\n", pattern.Response)
	fmt.Printf("  Strength: %.1f / %.1f\n", pattern.Strength, pattern.Threshold)
	fmt.Printf("  Project: %s\n", pattern.Project)
	fmt.Printf("  Created: %s\n", pattern.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Updated: %s\n", pattern.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

// deletePattern deletes a pattern
func deletePattern(storage *sqlite.Storage, id string) error {
	if id == "" {
		return fmt.Errorf("pattern ID required")
	}

	ctx := context.Background()
	if err := storage.DeletePattern(ctx, id); err != nil {
		return err
	}

	fmt.Printf("Pattern deleted: %s\n", id)
	return nil
}
