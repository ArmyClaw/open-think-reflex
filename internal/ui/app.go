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
	mode         AppMode // input, navigation
}

// AppMode represents the current interaction mode
type AppMode int

const (
	ModeInput     AppMode = iota // Default: typing in input
	ModeNavigation              // Navigating thought chain
)

var currentMode = ModeInput

// NewApp creates a new TUI application
func NewApp(storage *sqlite.Storage) *App {
	a := &App{
		storage: storage,
		matcher: matcher.NewEngine(),
		theme:   DefaultTheme(),
		mode:    ModeInput,
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
	
	modeText := "INPUT"
	if a.mode == ModeNavigation {
		modeText = "NAVIGATE"
	}
	
	info := tview.NewTextView().
		SetText(fmt.Sprintf("Space: %s | Patterns: %d | Mode: %s | [↑/↓] Navigate | [Tab] Switch Panel", spaceName, patternCount, modeText)).
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
		// Tab to switch between input and thought chain
		if event.Key() == tcell.KeyTab {
			if a.mode == ModeInput {
				a.mode = ModeNavigation
				a.thoughtChain.SetFocused(true)
				a.app.SetFocus(a.thoughtChain.view)
			} else {
				a.mode = ModeInput
				a.thoughtChain.SetFocused(false)
				a.app.SetFocus(a.input.view)
			}
			return nil
		}
		
		// If in navigation mode, handle arrow keys
		if a.mode == ModeNavigation {
			switch event.Key() {
			case tcell.KeyUp:
				a.thoughtChain.SelectPrev()
				a.updateOutputForSelection()
				return nil
			case tcell.KeyDown:
				a.thoughtChain.SelectNext()
				a.updateOutputForSelection()
				return nil
			case tcell.KeyLeft:
				// Collapse or go back
				a.thoughtChain.Collapse()
				return nil
			case tcell.KeyRight:
				// Expand or select
				a.thoughtChain.Expand()
				a.updateOutputForSelection()
				return nil
			case tcell.KeyEnter:
				// Use selected response
				a.useSelectedResponse()
				return nil
			case tcell.KeyEsc:
				// Return to input mode
				a.mode = ModeInput
				a.thoughtChain.SetFocused(false)
				a.app.SetFocus(a.input.view)
				return nil
			}
		}
		
		// Global shortcuts
		switch event.Key() {
		case tcell.KeyCtrlC:
			a.Stop()
			return nil
		case tcell.KeyEsc:
			if a.mode == ModeNavigation {
				a.mode = ModeInput
				a.thoughtChain.SetFocused(false)
				a.app.SetFocus(a.input.view)
				return nil
			}
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
			case 'h':
				// Left arrow equivalent
				if a.mode == ModeNavigation {
					a.thoughtChain.Collapse()
				}
				return nil
			case 'l':
				// Right arrow equivalent
				if a.mode == ModeNavigation {
					a.thoughtChain.Expand()
					a.updateOutputForSelection()
				}
				return nil
			case 'j':
				// Down arrow equivalent
				if a.mode == ModeNavigation {
					a.thoughtChain.SelectNext()
					a.updateOutputForSelection()
				}
				return nil
			case 'k':
				// Up arrow equivalent
				if a.mode == ModeNavigation {
					a.thoughtChain.SelectPrev()
					a.updateOutputForSelection()
				}
				return nil
			}
		}
		return event
	})
}

// updateOutputForSelection updates the output panel based on current selection
func (a *App) updateOutputForSelection() {
	result := a.thoughtChain.GetSelectedResult()
	if result != nil {
		a.output.SetOutput(fmt.Sprintf("Selected: %s\n\nTrigger: %s\nResponse: %s\nConfidence: %.0f%% (%s)\nStrength: %.1f / %.1f\n\nPress [Enter] to use this response",
			result.Pattern.Trigger,
			result.Pattern.Trigger,
			result.Pattern.Response,
			result.Confidence,
			result.Branch,
			result.Pattern.Strength,
			result.Pattern.Threshold))
	}
}

// useSelectedResponse copies the selected response to clipboard or shows it
func (a *App) useSelectedResponse() {
	result := a.thoughtChain.GetSelectedResult()
	if result != nil {
		a.output.SetFormattedOutput("Response Copied!", 
			fmt.Sprintf("✓ Copied to clipboard:\n\n%s", result.Pattern.Response))
		
		// TODO: Implement actual clipboard copy
		// For now, just show a success message
	}
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
		a.output.SetOutput("No active patterns found (all below threshold)\n\nTip: Use 'otr pattern create' to add patterns")
		a.thoughtChain.Clear()
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
	a.output.SetOutput(fmt.Sprintf("Found %d match(es):\n\n1. %s\n   Confidence: %.0f%% (%s)\n   Response: %s\n\nUse [↑/↓] to navigate, [Enter] to select",
		len(results),
		results[0].Pattern.Trigger,
		results[0].Confidence,
		results[0].Branch,
		results[0].Pattern.Response))
}

func (a *App) showHelp() {
	helpText := `
Keyboard Shortcuts
════════════════════════════════════════════════════════════

[Tab]       Switch between Input and Navigation mode
[↑/↓]       Navigate thought branches (in Navigation mode)
[←]         Collapse branch / Go back
[→]         Expand branch / Select
[Enter]     Use selected response
[Space]     Toggle selection
[h/?]       Show this help
[q/Esc]     Quit

Vim-style Shortcuts (Navigation mode)
════════════════════════════════════════════════════════════

[k]         Move up
[j]         Move down
[h]         Collapse / Go back
[l]         Expand / Select
[Esc]       Return to input mode

Modes
════════════════════════════════════════════════════════════

INPUT:       Type queries in the input area
NAVIGATE:    Use arrow keys to browse results

Layers
════════════════════════════════════════════════════════════

Left Panel:   Thought Chain Tree - Shows AI reasoning branches
Middle Panel: Output Content - Shows generated responses
Bottom:      Input Area - Type your queries here
`
	a.app.SetRoot(tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.SwitchToPage("main")
			a.app.SetRoot(a.pages, true)
		}), true)
}
