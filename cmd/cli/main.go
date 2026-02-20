package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ArmyClaw/open-think-reflex/internal/config"
	"github.com/ArmyClaw/open-think-reflex/internal/core/matcher"
	"github.com/ArmyClaw/open-think-reflex/internal/data/sqlite"
	"github.com/ArmyClaw/open-think-reflex/internal/ui"
	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/urfave/cli/v2"
)

// Version is set by build flags
var Version = "1.0-dev"

func main() {
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Run runs the CLI application
func Run() error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	storage, err := initStorage(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer storage.Close()

	app := &cli.App{
		Name:    "otr",
		Version: Version,
		Usage:   "Open-Think-Reflex: AI Input Accelerator",
		Commands: buildCommands(storage),
		Action: func(c *cli.Context) error {
			fmt.Println("Open-Think-Reflex v" + Version)
			fmt.Println("\nUse 'otr --help' to see available commands")
			fmt.Println("Use 'otr pattern --help' to manage patterns")
			return nil
		},
	}

	return app.Run(os.Args)
}

func loadConfig() (*config.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := fmt.Sprintf("%s/.openclaw/reflex", homeDir)
	loader := config.NewLoader(configPath, "config")
	return loader.Load()
}

func initStorage(cfg *config.Config) (*sqlite.Storage, error) {
	db, err := sqlite.NewDatabase(cfg.Storage.Path)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if err := db.Migrate(ctx); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := db.InitDefaultSpaces(ctx); err != nil {
		return nil, fmt.Errorf("failed to init default spaces: %w", err)
	}

	return sqlite.NewStorage(db), nil
}

func buildCommands(storage *sqlite.Storage) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "interactive",
			Usage: "Launch interactive TUI mode",
			Aliases: []string{"tui", "ui"},
			Action: func(c *cli.Context) error {
				return runInteractive(storage)
			},
		},
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
				{
					Name:  "update",
					Usage: "Update a pattern",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "id",
							Required: true,
							Usage:    "Pattern ID",
						},
						&cli.StringFlag{
							Name:  "trigger",
							Usage: "New trigger",
						},
						&cli.StringFlag{
							Name:  "response",
							Usage: "New response",
						},
						&cli.StringFlag{
							Name:  "project",
							Usage: "New project",
						},
						&cli.Float64Flag{
							Name:  "strength",
							Usage: "Set strength (0-100)",
						},
						&cli.Float64Flag{
							Name:  "threshold",
							Usage: "Set threshold (0-100)",
						},
					},
					Action: func(c *cli.Context) error {
						return updatePattern(storage, c.String("id"), c.String("trigger"), c.String("response"), c.String("project"), c.Float64("strength"), c.Float64("threshold"))
					},
				},
				{
					Name:  "reinforce",
					Usage: "Reinforce a pattern (increase strength)",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "id",
							Required: true,
							Usage:    "Pattern ID",
						},
						&cli.Float64Flag{
							Name:  "amount",
							Usage: "Reinforce amount (default 5)",
						},
					},
					Action: func(c *cli.Context) error {
						return reinforcePattern(storage, c.String("id"), c.Float64("amount"))
					},
				},
				{
					Name:  "decay",
					Usage: "Apply decay to all patterns",
					Action: func(c *cli.Context) error {
						return decayPatterns(storage)
					},
				},
			},
		},
		{
			Name:  "version",
			Usage: "Show version information",
			Action: func(c *cli.Context) error {
				fmt.Printf("Open-Think-Reflex v%s\n", Version)
				return nil
			},
		},
		{
			Name:  "run",
			Usage: "Run a query against patterns",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "query",
					Required: true,
					Usage:    "Query string",
				},
				&cli.Float64Flag{
					Name:  "threshold",
					Usage: "Minimum confidence threshold (default 30)",
				},
			},
			Action: func(c *cli.Context) error {
				return runQuery(storage, c.String("query"), c.Float64("threshold"))
			},
		},
	}
}

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

func updatePattern(storage *sqlite.Storage, id, trigger, response, project string, strength, threshold float64) error {
	if id == "" {
		return fmt.Errorf("pattern ID required")
	}

	ctx := context.Background()
	pattern, err := storage.GetPattern(ctx, id)
	if err != nil {
		return err
	}

	// Update fields if provided
	if trigger != "" {
		pattern.Trigger = trigger
	}
	if response != "" {
		pattern.Response = response
	}
	if project != "" {
		pattern.Project = project
	}
	if strength > 0 {
		pattern.Strength = strength
	}
	if threshold > 0 {
		pattern.Threshold = threshold
	}

	if err := storage.UpdatePattern(ctx, pattern); err != nil {
		return err
	}

	fmt.Printf("Pattern updated: %s\n", id)
	return nil
}

func reinforcePattern(storage *sqlite.Storage, id string, amount float64) error {
	if id == "" {
		return fmt.Errorf("pattern ID required")
	}
	if amount <= 0 {
		amount = 5.0 // Default reinforce amount
	}
	if amount > 20 {
		amount = 20 // Cap at 20
	}

	ctx := context.Background()
	pattern, err := storage.GetPattern(ctx, id)
	if err != nil {
		return err
	}

	// Increase strength
	oldStrength := pattern.Strength
	pattern.Strength = pattern.Strength + amount
	if pattern.Strength > 100 {
		pattern.Strength = 100
	}
	pattern.ReinforceCnt++

	if err := storage.UpdatePattern(ctx, pattern); err != nil {
		return err
	}

	fmt.Printf("Pattern reinforced: %s\n", id)
	fmt.Printf("  Strength: %.1f → %.1f\n", oldStrength, pattern.Strength)
	return nil
}

func decayPatterns(storage *sqlite.Storage) error {
	ctx := context.Background()
	patterns, err := storage.ListPatterns(ctx, contracts.ListOptions{Limit: 1000})
	if err != nil {
		return err
	}

	now := time.Now()
	decayedCount := 0

	for _, p := range patterns {
		if !p.DecayEnabled {
			continue
		}

		// Calculate time since last update
		hoursSinceUpdate := now.Sub(p.UpdatedAt).Hours()
		if hoursSinceUpdate < 24 {
			continue // Only decay patterns older than 24 hours
		}

		// Calculate decay
		days := hoursSinceUpdate / 24
		decayAmount := p.DecayRate * days * 100 // Convert to strength points
		oldStrength := p.Strength
		p.Strength = p.Strength - decayAmount
		p.DecayCnt++

		if p.Strength < 0 {
			p.Strength = 0
		}

		if err := storage.UpdatePattern(ctx, p); err != nil {
			continue
		}

		decayedCount++
		fmt.Printf("Pattern decayed: %s (%.1f → %.1f)\n", p.ID[:8], oldStrength, p.Strength)
	}

	fmt.Printf("Total patterns decayed: %d\n", decayedCount)
	return nil
}

func runQuery(storage *sqlite.Storage, query string, threshold float64) error {
	if threshold <= 0 {
		threshold = 30.0
	}

	ctx := context.Background()
	
	// Get all patterns
	patterns, err := storage.ListPatterns(ctx, contracts.ListOptions{Limit: 1000})
	if err != nil {
		return err
	}

	// Filter by threshold (only active patterns)
	var activePatterns []*models.Pattern
	for _, p := range patterns {
		if p.Strength >= p.Threshold {
			activePatterns = append(activePatterns, p)
		}
	}

	if len(activePatterns) == 0 {
		fmt.Println("No active patterns found (all below threshold)")
		return nil
	}

	// Create matcher and run
	engine := matcher.NewEngine()
	opts := contracts.MatchOptions{
		Threshold:  threshold,
		Limit:      10,
		ExactFirst: true,
	}
	
	results := engine.Match(ctx, query, activePatterns, opts)
	
	if len(results) == 0 {
		fmt.Printf("No matches found for: %s\n", query)
		fmt.Println("\nTip: Use 'otr pattern create' to add patterns")
		return nil
	}

	fmt.Printf("Found %d match(es):\n\n", len(results))
	for i, r := range results {
		fmt.Printf("%d. %s\n", i+1, r.Pattern.Trigger)
		fmt.Printf("   Confidence: %.0f%% (%s)\n", r.Confidence, r.Branch)
		fmt.Printf("   Response: %s\n", truncate(r.Pattern.Response, 60))
		fmt.Printf("   Strength: %.1f / %.1f\n", r.Pattern.Strength, r.Pattern.Threshold)
		fmt.Println()
	}

	return nil
}

func runInteractive(storage *sqlite.Storage) error {
	fmt.Println("Starting interactive mode...")
	
	app := ui.NewApp(storage)
	ctx := context.Background()
	
	return app.Run(ctx)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
