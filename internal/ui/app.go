package ui

import (
	"context"
	"fmt"

	"github.com/ArmyClaw/open-think-reflex/internal/core/matcher"
	"github.com/ArmyClaw/open-think-reflex/internal/data/sqlite"
	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App represents the main TUI application
type App struct {
	app       *tview.Application
	storage   *sqlite.Storage
	matcher   *matcher.Engine
	pages     *tview.Pages
	theme     *Theme
	
	// Layers
	thoughtChain *ThoughtChainView
	output       *OutputView
	input        *InputView
	
	// State
	currentSpace *models.Space
	patterns     []*models.Pattern
	results      []contracts.MatchResult
}

// NewApp creates a new TUI application
func NewApp(storage *sqlite.Storage) *App {
	a := &App{
		storage: storage,
		matcher: matcher.NewEngine(),
		theme:   DefaultTheme(),
	}
	
	a.app = tview.NewApplication()
	a.setupPages()
	
	return a
}

// Run starts the TUI application
func (a *App) Run(ctx context.Context) error {
	// Load initial data
	if err := a.loadData(ctx); err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}
	
	// Set up the application
	a.app.SetRoot(a.pages, true)
	a.app.SetFocus(a.input.view)
	
	// Enable mouse support
	a.app.EnableMouse(true)
	
	// Set up keyboard shortcuts
	a.setupKeyBindings()
	
	return a.app.Run()
}

// Stop stops the TUI application
func (a *App) Stop() {
	a.app.Stop()
}

func (a *App) loadData(ctx context.Context) error {
	// Load patterns
	patterns, err := a.storage.ListPatterns(ctx, contracts.ListOptions{Limit: 1000})
	if err != nil {
		return err
	}
	a.patterns = patterns
	
	// Load default space
	spaces, err := a.storage.ListSpaces(ctx)
	if err != nil {
		return err
	}
	if len(spaces) > 0 {
		a.currentSpace = spaces[0]
	}
	
	return nil
}

func (a *App) setupPages() {
	// Create the main layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.createHeader(), 3, 0, false).
		AddItem(a.createMainContent(), 0, 1, true).
		AddItem(a.createInputArea(), 5, 0, false)
	
	// Create pages
	a.pages = tview.NewPages()
	a.pages.AddPage("main", flex, true, true)
}

func (a *App) createHeader() tview.Primitive {
	title := tview.NewTextView().
		SetText("🤖 Open-Think-Reflex - AI Input Accelerator").
		SetTextColor(a.theme.Primary).
		SetTextAlign(tview.AlignCenter)
	
	patternCount := len(a.patterns)
	spaceName := "default"
	if a.currentSpace != nil {
		spaceName = a.currentSpace.Name
	}
	info := tview.NewTextView().
		SetText(fmt.Sprintf("Space: %s | Patterns: %d | Press '?' for help | 'q' to quit", spaceName, patternCount)).
		SetTextColor(a.theme.Secondary).
		SetTextAlign(tview.AlignCenter)
	
	header := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 1, 0, false).
		AddItem(info, 1, 0, false)
	
	header.SetBorder(false)
	return header
}

func (a *App) createMainContent() tview.Primitive {
	// Create three-column layout
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	
	// Layer 1: Thought Chain Tree (left panel)
	a.thoughtChain = NewThoughtChainView(a.theme)
	flex.AddItem(a.thoughtChain.view, 0, 1, false)
	
	// Layer 2: Output Content (middle panel)
	a.output = NewOutputView(a.theme)
	flex.AddItem(a.output.view, 0, 2, true)
	
	// Separator
	sep := tview.NewTextView().
		SetText("│").
		SetTextColor(a.theme.Border)
	flex.AddItem(sep, 1, 0, false)
	
	return flex
}

func (a *App) createInputArea() tview.Primitive {
	a.input = NewInputView(a.theme, func(text string) {
		a.handleInput(text)
	})
	
	inputBox := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().
			SetText("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━").
			SetTextColor(a.theme.Border), 1, 0, false).
		AddItem(a.input.view, 0, 1, true)
	
	return inputBox
}

func (a *App) setupKeyBindings() {
	// Global key bindings
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			a.Stop()
			return nil
		case tcell.KeyEsc:
			a.Stop()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				a.Stop()
				return nil
			case '?':
				a.showHelp()
				return nil
			}
		}
		return event
	})
}

func (a *App) handleInput(text string) {
	ctx := context.Background()
	
	// Filter active patterns
	var activePatterns []*models.Pattern
	for _, p := range a.patterns {
		if p.Strength >= p.Threshold {
			activePatterns = append(activePatterns, p)
		}
	}
	
	if len(activePatterns) == 0 {
		a.output.SetOutput("No active patterns found (all below threshold)")
		return
	}
	
	// Run matching
	opts := contracts.MatchOptions{
		Threshold:  30,
		Limit:       10,
		ExactFirst: true,
	}
	
	results := a.matcher.Match(ctx, text, activePatterns, opts)
	a.results = results
	
	if len(results) == 0 {
		a.output.SetOutput(fmt.Sprintf("No matches found for: %s\n\nTip: Use 'otr pattern create' to add patterns", text))
		a.thoughtChain.Clear()
		return
	}
	
	// Update thought chain with results
	a.thoughtChain.SetResults(results)
	
	// Update output with first match
	a.output.SetOutput(fmt.Sprintf("Match found!\n\nTrigger: %s\nResponse: %s\nConfidence: %.0f%%",
		results[0].Pattern.Trigger,
		results[0].Pattern.Response,
		results[0].Confidence))
}

func (a *App) showHelp() {
	helpText := `
Keyboard Shortcuts
═══════════════════════════════════════════════════════════

[↑/↓]     Navigate thought branches
[→]        Expand/select branch
[←]        Go back
[Enter]    Use selected response
[Space]    Toggle selection
[h/?]      Show this help
[q/Esc]    Quit

Layers
═══════════════════════════════════════════════════════════

Left Panel:   Thought Chain Tree - Shows AI reasoning branches
Middle Panel: Output Content - Shows generated responses
Bottom:       Input Area - Type your queries here
`
	a.app.SetRoot(tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.SwitchToPage("main")
			a.app.SetRoot(a.pages, true)
		}), true)
}
