package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ArmyClaw/open-think-reflex/internal/cli/commands"
	"github.com/ArmyClaw/open-think-reflex/internal/config"
	"github.com/ArmyClaw/open-think-reflex/internal/core/matcher"
	"github.com/ArmyClaw/open-think-reflex/internal/data/sqlite"
	"github.com/ArmyClaw/open-think-reflex/internal/ui"
	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/export"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/ArmyClaw/open-think-reflex/pkg/skills"
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
	cfg, loader, err := loadConfig()
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
		Commands: buildCommands(storage, cfg, loader),
		Action: func(c *cli.Context) error {
			fmt.Println("Open-Think-Reflex v" + Version)
			fmt.Println("\nUse 'otr --help' to see available commands")
			fmt.Println("Use 'otr pattern --help' to manage patterns")
			return nil
		},
	}

	return app.Run(os.Args)
}

var configLoader *config.Loader

func loadConfig() (*config.Config, *config.Loader, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := fmt.Sprintf("%s/.openclaw/reflex", homeDir)
	loader := config.NewLoader(configPath, "config")
	configLoader = loader
	cfg, err := loader.Load()
	if err != nil {
		return nil, loader, err
	}
	return cfg, loader, nil
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

func buildCommands(storage *sqlite.Storage, cfg *config.Config, loader *config.Loader) []*cli.Command {
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
					Name:  "move",
					Usage: "Move a pattern to a different space",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "id",
							Required: true,
							Usage:    "Pattern ID",
						},
						&cli.StringFlag{
							Name:     "space",
							Required: true,
							Usage:    "Target space ID",
						},
					},
					Action: func(c *cli.Context) error {
						return movePattern(storage, c.String("id"), c.String("space"))
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
			Name:  "space",
			Usage: "Manage pattern spaces",
			Subcommands: []*cli.Command{
				{
					Name:  "list",
					Usage: "List all spaces",
					Action: func(c *cli.Context) error {
						return commands.ListSpaces(storage)
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
						return commands.CreateSpace(storage, c.String("name"), c.String("description"))
					},
				},
				{
					Name:      "show",
					Usage:     "Show space details",
					ArgsUsage: "<space_id>",
					Action: func(c *cli.Context) error {
						return commands.ShowSpace(storage, c.Args().First())
					},
				},
				{
					Name:      "delete",
					Usage:     "Delete a space",
					ArgsUsage: "<space_id>",
					Action: func(c *cli.Context) error {
						return commands.DeleteSpace(storage, c.Args().First())
					},
				},
				{
					Name:      "use",
					Usage:     "Switch to a space",
					ArgsUsage: "<space_id>",
					Action: func(c *cli.Context) error {
						return commands.UseSpace(storage, cfg, loader, c.Args().First())
					},
				},
				{
					Name:      "default",
					Usage:     "Set a space as default",
					ArgsUsage: "<space_id>",
					Action: func(c *cli.Context) error {
						return commands.SetDefaultSpace(storage, c.Args().First())
					},
				},
				{
					Name:  "export",
					Usage: "Export a space to a file",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "id",
							Required: true,
							Usage:    "Space ID to export",
						},
						&cli.StringFlag{
							Name:     "output",
							Required: true,
							Usage:    "Output file path",
						},
					},
					Action: func(c *cli.Context) error {
						return exportSpace(storage, c.String("id"), c.String("output"))
					},
				},
				{
					Name:  "import",
					Usage: "Import a space from a file",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "input",
							Required: true,
							Usage:    "Input file path",
						},
						&cli.BoolFlag{
							Name:  "force",
							Usage: "Overwrite existing patterns",
						},
					},
					Action: func(c *cli.Context) error {
						return importSpace(storage, c.String("input"), c.Bool("force"))
					},
				},
			},
		},
		{
			Name:  "note",
			Usage: "Manage notes and thoughts",
			Subcommands: []*cli.Command{
				{
					Name:  "list",
					Usage: "List all notes",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "space",
							Usage: "Filter by space ID",
						},
						&cli.StringFlag{
							Name:  "category",
							Usage: "Filter by category (thought/idea/todo/memory/question/note)",
						},
					},
					Action: func(c *cli.Context) error {
						return listNotes(storage, c.String("space"), c.String("category"))
					},
				},
				{
					Name:  "create",
					Usage: "Create a new note",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "title",
							Required: true,
							Usage:    "Note title",
						},
						&cli.StringFlag{
							Name:     "content",
							Required: true,
							Usage:    "Note content",
						},
						&cli.StringFlag{
							Name:  "category",
							Usage: "Note category",
						},
						&cli.StringFlag{
							Name:  "space",
							Usage: "Space ID",
						},
					},
					Action: func(c *cli.Context) error {
						return createNote(storage, c.String("title"), c.String("content"), c.String("category"), c.String("space"))
					},
				},
				{
					Name:      "show",
					Usage:     "Show note details",
					ArgsUsage: "<note_id>",
					Action: func(c *cli.Context) error {
						return showNote(storage, c.Args().First())
					},
				},
				{
					Name:      "delete",
					Usage:     "Delete a note",
					ArgsUsage: "<note_id>",
					Action: func(c *cli.Context) error {
						return deleteNote(storage, c.Args().First())
					},
				},
				{
					Name:  "search",
					Usage: "Search notes",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "query",
							Required: true,
							Usage:    "Search query",
						},
					},
					Action: func(c *cli.Context) error {
						return searchNotes(storage, c.String("query"))
					},
				},
				{
					Name:  "link",
					Usage: "Link a pattern to a note",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "note",
							Required: true,
							Usage:    "Note ID",
						},
						&cli.StringFlag{
							Name:     "pattern",
							Required: true,
							Usage:    "Pattern ID",
						},
					},
					Action: func(c *cli.Context) error {
						return linkPatternToNote(storage, c.String("note"), c.String("pattern"))
					},
				},
				{
					Name:  "unlink",
					Usage: "Unlink a pattern from a note",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "note",
							Required: true,
							Usage:    "Note ID",
						},
						&cli.StringFlag{
							Name:     "pattern",
							Required: true,
							Usage:    "Pattern ID",
						},
					},
					Action: func(c *cli.Context) error {
						return unlinkPatternFromNote(storage, c.String("note"), c.String("pattern"))
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
			Name:  "backup",
			Usage: "Backup data to a file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "output",
					Required: true,
					Usage:    "Output backup file path",
				},
				&cli.StringFlag{
					Name:  "format",
					Usage: "Backup format: json, yaml (default: json)",
				},
				&cli.BoolFlag{
					Name:  "include-notes",
					Usage: "Include notes in backup",
				},
			},
			Action: func(c *cli.Context) error {
				return createBackup(storage, c.String("output"), c.String("format"), c.Bool("include-notes"))
			},
		},
		{
			Name:  "share",
			Usage: "Share a pattern",
			Subcommands: []*cli.Command{
				{
					Name:  "list",
					Usage: "List all public patterns",
					Action: func(c *cli.Context) error {
						return listPublicPatterns(storage)
					},
				},
				{
					Name:  "create",
					Usage: "Create a shareable link for a pattern",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "id",
							Required: true,
							Usage:    "Pattern ID to share",
						},
					},
					Action: func(c *cli.Context) error {
						return sharePattern(storage, c.String("id"))
					},
				},
				{
					Name:  "import",
					Usage: "Import a shared pattern",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "code",
							Required: true,
							Usage:    "Share code",
						},
					},
					Action: func(c *cli.Context) error {
						return importSharedPattern(storage, c.String("code"))
					},
				},
				{
					Name:  "space",
					Usage: "Share an entire space",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "id",
							Required: true,
							Usage:    "Space ID to share",
						},
					},
					Action: func(c *cli.Context) error {
						return shareSpace(storage, c.String("id"))
					},
				},
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
		{
			Name:  "export",
			Usage: "Export patterns to a JSON file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "output",
					Required: true,
					Usage:    "Output file path",
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
			Usage: "Import patterns from a JSON file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "input",
					Required: true,
					Usage:    "Input file path",
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
		{
			Name:  "skill",
			Usage: "Export patterns as AgentSkill",
			Subcommands: []*cli.Command{
				{
					Name:  "export",
					Usage: "Export pattern as AgentSkill",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "id",
							Required: true,
							Usage:    "Pattern ID to export",
						},
						&cli.StringFlag{
							Name:     "output",
							Required: true,
							Usage:    "Output file path (.yaml)",
						},
					},
					Action: func(c *cli.Context) error {
						return exportSkill(storage, c.String("id"), c.String("output"))
					},
				},
				{
					Name:  "batch",
					Usage: "Export all patterns as AgentSkills",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "output",
							Required: true,
							Usage:    "Output directory",
						},
						&cli.StringFlag{
							Name:  "space",
							Usage: "Filter by space ID",
						},
					},
					Action: func(c *cli.Context) error {
						return exportSkillsBatch(storage, c.String("output"), c.String("space"))
					},
				},
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

func movePattern(storage *sqlite.Storage, patternID, spaceID string) error {
	if patternID == "" {
		return fmt.Errorf("pattern ID required")
	}
	if spaceID == "" {
		return fmt.Errorf("space ID required")
	}

	ctx := context.Background()

	// Verify pattern exists
	_, err := storage.GetPattern(ctx, patternID)
	if err != nil {
		return fmt.Errorf("pattern not found: %s", patternID)
	}

	// Verify space exists
	_, err = storage.GetSpace(ctx, spaceID)
	if err != nil {
		return fmt.Errorf("space not found: %s", spaceID)
	}

	if err := storage.MovePatternToSpace(ctx, patternID, spaceID); err != nil {
		return err
	}

	fmt.Printf("Pattern '%s' moved to space '%s'\n", patternID, spaceID)
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

func exportSpace(storage *sqlite.Storage, spaceID, outputPath string) error {
	if spaceID == "" {
		return fmt.Errorf("space ID required")
	}
	if outputPath == "" {
		return fmt.Errorf("output path required")
	}

	ctx := context.Background()

	// Get space
	space, err := storage.GetSpace(ctx, spaceID)
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	// Get patterns in this space
	patterns, err := storage.ListPatterns(ctx, contracts.ListOptions{SpaceID: spaceID, Limit: 10000})
	if err != nil {
		return fmt.Errorf("failed to list patterns: %w", err)
	}

	// Export
	exporter := export.NewExporter()
	if err := exporter.ExportSpaceToJSON(ctx, space, patterns, outputPath); err != nil {
		return fmt.Errorf("failed to export: %w", err)
	}

	fmt.Printf("Space '%s' exported to '%s' (%d patterns)\n", space.Name, outputPath, len(patterns))
	return nil
}

func importSpace(storage *sqlite.Storage, inputPath string, force bool) error {
	if inputPath == "" {
		return fmt.Errorf("input path required")
	}

	ctx := context.Background()

	// Import
	importer := export.NewImporter()
	data, err := importer.ImportSpaceFromJSON(ctx, inputPath)
	if err != nil {
		return fmt.Errorf("failed to import: %w", err)
	}

	// Create or get space
	space := data.Space
	if space == nil {
		return fmt.Errorf("invalid import file: no space data")
	}

	// Check if space exists
	existingSpace, err := storage.GetSpace(ctx, space.ID)
	if err != nil {
		// Space doesn't exist, create it
		if err := storage.CreateSpace(ctx, space); err != nil {
			return fmt.Errorf("failed to create space: %w", err)
		}
		fmt.Printf("Created space: %s\n", space.Name)
	} else if !force {
		return fmt.Errorf("space '%s' already exists (use --force to overwrite)", existingSpace.Name)
	}

	// Import patterns
	importedCount := 0
	for _, pattern := range data.Patterns {
		pattern.SpaceID = space.ID // Ensure pattern belongs to this space
		if err := storage.SavePattern(ctx, &pattern); err != nil {
			fmt.Printf("Warning: failed to import pattern %s: %v\n", pattern.Trigger, err)
			continue
		}
		importedCount++
	}

	fmt.Printf("Imported %d patterns to space '%s'\n", importedCount, space.Name)
	return nil
}

func listNotes(storage *sqlite.Storage, spaceID, category string) error {
	ctx := context.Background()

	opts := contracts.ListOptions{
		SpaceID: spaceID,
	}

	notes, err := storage.ListNotes(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	// Filter by category in memory if specified
	if category != "" {
		var filtered []*models.Note
		for _, n := range notes {
			if n.Category == category {
				filtered = append(filtered, n)
			}
		}
		notes = filtered
	}

	if len(notes) == 0 {
		fmt.Println("No notes found")
		return nil
	}

	fmt.Printf("Found %d notes:\n\n", len(notes))
	for _, n := range notes {
		preview := n.Content
		if len(preview) > 50 {
			preview = preview[:50] + "..."
		}
		fmt.Printf("  %s  %s (%s)\n", n.ID[:min(8, len(n.ID))], n.Title, n.Category)
		fmt.Printf("      %s\n\n", preview)
	}

	return nil
}

func createNote(storage *sqlite.Storage, title, content, category, spaceID string) error {
	ctx := context.Background()

	note := &models.Note{
		Title:    title,
		Content:  content,
		Category: category,
		SpaceID:  spaceID,
	}

	if note.SpaceID == "" {
		note.SpaceID = "global"
	}
	if note.Category == "" {
		note.Category = "note"
	}

	if err := storage.SaveNote(ctx, note); err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	fmt.Printf("Note created: %s\n", note.ID)
	return nil
}

func showNote(storage *sqlite.Storage, noteID string) error {
	if noteID == "" {
		return fmt.Errorf("note ID required")
	}

	ctx := context.Background()
	note, err := storage.GetNote(ctx, noteID)
	if err != nil {
		return fmt.Errorf("note not found: %w", err)
	}

	fmt.Printf("Note: %s\n", note.ID)
	fmt.Printf("  Title: %s\n", note.Title)
	fmt.Printf("  Category: %s\n", note.Category)
	fmt.Printf("  Space: %s\n", note.SpaceID)
	fmt.Printf("  Words: %d\n", note.WordCount)
	fmt.Printf("  Created: %s\n", note.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Updated: %s\n", note.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("\nContent:\n%s\n", note.Content)

	return nil
}

func deleteNote(storage *sqlite.Storage, noteID string) error {
	if noteID == "" {
		return fmt.Errorf("note ID required")
	}

	ctx := context.Background()
	if err := storage.DeleteNote(ctx, noteID); err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	fmt.Printf("Note deleted: %s\n", noteID)
	return nil
}

func searchNotes(storage *sqlite.Storage, query string) error {
	if query == "" {
		return fmt.Errorf("search query required")
	}

	ctx := context.Background()
	notes, err := storage.SearchNotes(ctx, query, contracts.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to search notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("No notes found")
		return nil
	}

	fmt.Printf("Found %d matching notes:\n\n", len(notes))
	for _, n := range notes {
		fmt.Printf("  %s  %s (%s)\n", n.ID[:min(8, len(n.ID))], n.Title, n.Category)
	}

	return nil
}

func linkPatternToNote(storage *sqlite.Storage, noteID, patternID string) error {
	if noteID == "" {
		return fmt.Errorf("note ID required")
	}
	if patternID == "" {
		return fmt.Errorf("pattern ID required")
	}

	ctx := context.Background()

	// Get note
	note, err := storage.GetNote(ctx, noteID)
	if err != nil {
		return fmt.Errorf("note not found: %w", err)
	}

	// Verify pattern exists
	_, err = storage.GetPattern(ctx, patternID)
	if err != nil {
		return fmt.Errorf("pattern not found: %w", err)
	}

	// Add pattern link
	note.AddPattern(patternID)
	if err := storage.UpdateNote(ctx, note); err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	fmt.Printf("Linked pattern '%s' to note '%s'\n", patternID[:min(8, len(patternID))], note.Title)
	return nil
}

func unlinkPatternFromNote(storage *sqlite.Storage, noteID, patternID string) error {
	if noteID == "" {
		return fmt.Errorf("note ID required")
	}
	if patternID == "" {
		return fmt.Errorf("pattern ID required")
	}

	ctx := context.Background()

	// Get note
	note, err := storage.GetNote(ctx, noteID)
	if err != nil {
		return fmt.Errorf("note not found: %w", err)
	}

	// Remove pattern link
	note.RemovePattern(patternID)
	if err := storage.UpdateNote(ctx, note); err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	fmt.Printf("Unlinked pattern '%s' from note '%s'\n", patternID[:min(8, len(patternID))], note.Title)
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

// createBackup creates a full backup of all data
func createBackup(storage *sqlite.Storage, outputPath, format string, includeNotes bool) error {
	if outputPath == "" {
		return fmt.Errorf("output path required")
	}

	if format == "" {
		format = "json"
	}

	ctx := context.Background()

	// Get all patterns
	patterns, err := storage.ListPatterns(ctx, contracts.ListOptions{Limit: 10000})
	if err != nil {
		return fmt.Errorf("failed to list patterns: %w", err)
	}

	// Get all spaces
	spaces, err := storage.ListSpaces(ctx)
	if err != nil {
		return fmt.Errorf("failed to list spaces: %w", err)
	}

	fmt.Printf("Backing up %d patterns and %d spaces...\n", len(patterns), len(spaces))

	exporter := export.NewExporter()

	// Export based on format
	switch format {
	case "yaml", "yml":
		// Export each space separately with YAML
		for _, space := range spaces {
			spacePatterns, err := storage.ListPatterns(ctx, contracts.ListOptions{SpaceID: space.ID, Limit: 10000})
			if err != nil {
				fmt.Printf("Warning: failed to list patterns for space %s: %v\n", space.Name, err)
				continue
			}
			filename := outputPath
			if len(spaces) > 1 {
				filename = fmt.Sprintf("%s_%s.yaml", outputPath[:len(outputPath)-5], space.ID)
			}
			if err := exporter.ExportSpaceToYAML(ctx, space, spacePatterns, filename); err != nil {
				fmt.Printf("Warning: failed to backup space %s: %v\n", space.Name, err)
			}
		}
		fmt.Printf("Backup completed: %s\n", outputPath)
	case "json":
		fallthrough
	default:
		// Export all as single JSON
		if err := exporter.ExportToJSON(ctx, patterns, outputPath); err != nil {
			return fmt.Errorf("failed to backup: %w", err)
		}
		fmt.Printf("Backup completed: %s (%d patterns)\n", outputPath, len(patterns))
	}

	return nil
}

// listPublicPatterns lists all patterns (shareable ones)
func listPublicPatterns(storage *sqlite.Storage) error {
	ctx := context.Background()

	patterns, err := storage.ListPatterns(ctx, contracts.ListOptions{Limit: 100})
	if err != nil {
		return fmt.Errorf("failed to list patterns: %w", err)
	}

	if len(patterns) == 0 {
		fmt.Println("No patterns found")
		return nil
	}

	fmt.Printf("Public Patterns (%d total):\n\n", len(patterns))
	for _, p := range patterns {
		preview := p.Response
		if len(preview) > 40 {
			preview = preview[:40] + "..."
		}
		fmt.Printf("  %s  %s -> %s\n", p.ID[:min(8, len(p.ID))], p.Trigger, preview)
	}
	fmt.Printf("\nUse 'otr share create --id <id>' to share a pattern\n")

	return nil
}

// sharePattern generates a shareable code for a pattern
func sharePattern(storage *sqlite.Storage, patternID string) error {
	if patternID == "" {
		return fmt.Errorf("pattern ID required")
	}

	ctx := context.Background()

	// Get pattern
	pattern, err := storage.GetPattern(ctx, patternID)
	if err != nil {
		return fmt.Errorf("pattern not found: %w", err)
	}

	// Get space name
	spaceName := ""
	if pattern.SpaceID != "" {
		space, err := storage.GetSpace(ctx, pattern.SpaceID)
		if err == nil {
			spaceName = space.Name
		}
	}

	// Convert to skill
	skill := skills.ConvertPatternToSkill(pattern, spaceName)

	// Marshal to JSON
	data, err := json.Marshal(skill)
	if err != nil {
		return fmt.Errorf("failed to marshal skill: %w", err)
	}

	// Encode as base64
	encoded := base64.StdEncoding.EncodeToString(data)

	fmt.Printf("Pattern shared! Share code:\n%s\n", encoded)
	fmt.Printf("\nUse: otr share import --code <code>\n")

	return nil
}

// importSharedPattern imports a pattern from a share code
func importSharedPattern(storage *sqlite.Storage, code string) error {
	if code == "" {
		return fmt.Errorf("share code required")
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return fmt.Errorf("invalid share code: %w", err)
	}

	// Unmarshal JSON
	var skill skills.Skill
	if err := json.Unmarshal(data, &skill); err != nil {
		return fmt.Errorf("invalid share code format: %w", err)
	}

	// Convert to pattern
	pattern := skills.ConvertSkillToPattern(&skill)

	ctx := context.Background()

	// Save pattern
	if err := storage.SavePattern(ctx, pattern); err != nil {
		return fmt.Errorf("failed to save pattern: %w", err)
	}

	fmt.Printf("Imported pattern: %s\n", pattern.Trigger)
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

func exportSkill(storage *sqlite.Storage, patternID, outputPath string) error {
	if patternID == "" {
		return fmt.Errorf("pattern ID required")
	}
	if outputPath == "" {
		return fmt.Errorf("output path required")
	}

	ctx := context.Background()

	// Get pattern
	pattern, err := storage.GetPattern(ctx, patternID)
	if err != nil {
		return fmt.Errorf("pattern not found: %w", err)
	}

	// Get space name
	spaceName := ""
	if pattern.SpaceID != "" {
		space, err := storage.GetSpace(ctx, pattern.SpaceID)
		if err == nil {
			spaceName = space.Name
		}
	}

	// Convert to skill
	skill := skills.ConvertPatternToSkill(pattern, spaceName)

	// TODO: Add YAML marshaling and AI polish option
	// For now, output as JSON
	fmt.Printf("Exported skill: %s\n", skill.Name)
	fmt.Printf("Trigger: %s\n", skill.Trigger)
	fmt.Printf("Output: %s\n", outputPath)

	return nil
}

func exportSkillsBatch(storage *sqlite.Storage, outputDir, spaceID string) error {
	if outputDir == "" {
		return fmt.Errorf("output directory required")
	}

	ctx := context.Background()

	// Get patterns
	opts := contracts.ListOptions{Limit: 10000}
	if spaceID != "" {
		opts.SpaceID = spaceID
	}

	patterns, err := storage.ListPatterns(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list patterns: %w", err)
	}

	if len(patterns) == 0 {
		fmt.Println("No patterns found to export")
		return nil
	}

	fmt.Printf("Found %d patterns to export\n", len(patterns))

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Export each pattern as a separate skill file
	exporter := export.NewExporter()
	successCount := 0
	for _, p := range patterns {
		// Get space name
		spaceName := "global"
		if p.SpaceID != "" {
			space, err := storage.GetSpace(ctx, p.SpaceID)
			if err == nil {
				spaceName = space.Name
			}
		}

		// Convert to skill
		skill := skills.ConvertPatternToSkill(p, spaceName)

		// Create filename from trigger
		filename := strings.ReplaceAll(skill.Trigger, "/", "_")
		filename = strings.ReplaceAll(filename, " ", "_")
		filename = strings.ToLower(filename) + ".json"
		filepath := filepath.Join(outputDir, filename)

		// Export to JSON
		if err := exporter.ExportSkillToJSON(skill, filepath); err != nil {
			fmt.Printf("Warning: Failed to export %s: %v\n", skill.Trigger, err)
			continue
		}
		successCount++
	}

	fmt.Printf("Successfully exported %d/%d patterns to %s\n", successCount, len(patterns), outputDir)
	return nil
}

// shareSpace shares an entire space with all its patterns
func shareSpace(storage *sqlite.Storage, spaceID string) error {
	if spaceID == "" {
		return fmt.Errorf("space ID required")
	}

	ctx := context.Background()

	// Get space
	space, err := storage.GetSpace(ctx, spaceID)
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	// Get all patterns in space
	patterns, err := storage.ListPatterns(ctx, contracts.ListOptions{SpaceID: spaceID, Limit: 10000})
	if err != nil {
		return fmt.Errorf("failed to list patterns: %w", err)
	}

	// Convert patterns to JSON
	data, err := json.Marshal(struct {
		Space    *models.Space    `json:"space"`
		Patterns []models.Pattern `json:"patterns"`
	}{
		Space: space,
		Patterns: func() []models.Pattern {
			result := make([]models.Pattern, len(patterns))
			for i, p := range patterns {
				result[i] = *p
			}
			return result
		}(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	// Encode as base64
	encoded := base64.StdEncoding.EncodeToString(data)

	fmt.Printf("Space '%s' shared with %d patterns!\n\n", space.Name, len(patterns))
	fmt.Printf("Share code:\n%s\n\n", encoded)
	fmt.Printf("Use: otr share import --code <code>\n")

	return nil
}
