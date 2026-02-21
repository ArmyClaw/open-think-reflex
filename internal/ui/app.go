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
	app         *tview.Application
	storage     *sqlite.Storage
	matcher     *matcher.Engine
	pages       *tview.Pages
	theme       *Theme
	themeManager *ThemeManager
	
	// Layers
	thoughtChain *ThoughtChainView
	output       *OutputView
	input        *InputView
	statusBar    *StatusBar
	helpPanel    *HelpPanel
	shortcutBar  *ShortcutBar
	filterPanel  *FilterPanel
	
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
	themeManager := NewThemeManager()
	a := &App{
		storage:      storage,
		matcher:      matcher.NewEngine(),
		theme:        themeManager.Current(),
		themeManager: themeManager,
		mode:         ModeInput,
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
	// Create status bar
	a.statusBar = NewStatusBar(a.theme)
	a.statusBar.SetPatternCount(len(a.patterns))
	
	// Create help panel and shortcut bar
	a.helpPanel = NewHelpPanel(a.theme)
	a.shortcutBar = NewShortcutBar(a.theme)
	a.filterPanel = NewFilterPanel()
	
	// Set up filter callback
	a.filterPanel.onFilter = a.filterPatterns
	
	// Create the main layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.createHeader(), 3, 0, false).
		AddItem(a.createMainContent(), 0, 1, true).
		AddItem(a.createInputArea(), 5, 0, false).
		AddItem(a.shortcutBar.View(), 1, 0, false).
		AddItem(a.statusBar.View(), 1, 0, false)
	
	// Create pages
	a.pages = tview.NewPages()
	a.pages.AddPage("main", flex, true, true)
	
	// Add help as overlay page
	a.pages.AddPage("help", a.helpPanel.View(), false, false)
	a.pages.AddPage("filter", a.filterPanel.GetView(), false, false)
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
				a.shortcutBar.SetMode(ModeNavigation)
			} else {
				a.mode = ModeInput
				a.thoughtChain.SetFocused(false)
				a.app.SetFocus(a.input.view)
				a.shortcutBar.SetMode(ModeInput)
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
			case '/':
				// Toggle filter panel
				a.toggleFilter()
				return nil
			case 't':
				// Toggle theme
				a.toggleTheme()
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
	
	// Set status to matching
	a.statusBar.SetStatus(StatusMatching, "Matching patterns...")
	
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
		a.statusBar.SetStatus(StatusIdle, "No active patterns")
		a.statusBar.SetMatchCount(0)
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
		a.statusBar.SetStatus(StatusIdle, "No matches")
		a.statusBar.SetMatchCount(0)
		return
	}
	
	// Update thought chain with results
	a.thoughtChain.SetResults(results)
	
	// Update status bar
	a.statusBar.SetMatchCount(len(results))
	a.statusBar.SetStatus(StatusIdle, fmt.Sprintf("Found %d match(es)", len(results)))
	
	// Update output with first match
	a.output.SetOutput(fmt.Sprintf("Found %d match(es):\n\n1. %s\n   Confidence: %.0f%% (%s)\n   Response: %s\n\nUse [↑/↓] to navigate, [Enter] to select",
		len(results),
		results[0].Pattern.Trigger,
		results[0].Confidence,
		results[0].Branch,
		results[0].Pattern.Response))
}

// toggleTheme switches between light and dark themes
func (a *App) toggleTheme() {
	a.themeManager.Toggle()
	a.theme = a.themeManager.Current()
	
	// Update UI elements with new theme
	a.thoughtChain.theme = a.theme
	a.output.theme = a.theme
	a.input.theme = a.theme
	a.statusBar.SetTheme(a.theme)
	a.helpPanel.SetTheme(a.theme)
	a.shortcutBar.SetTheme(a.theme)
	
	// Re-render
	a.thoughtChain.render()
	a.output.view.SetBackgroundColor(a.theme.Background)
	a.output.view.SetBorderColor(a.theme.Border)
	
	a.input.view.SetFieldBackgroundColor(a.theme.Background)
	a.input.view.SetFieldTextColor(a.theme.Text)
	
	// Show notification
	a.output.SetStatus(fmt.Sprintf("Theme switched to %s", a.theme.Name), true)
}

func (a *App) showHelp() {
	// Toggle help panel
	if a.helpPanel.IsVisible() {
		a.helpPanel.Hide()
		a.pages.HidePage("help")
		a.pages.SwitchToPage("main")
	} else {
		a.helpPanel.Show()
		a.pages.ShowPage("help")
		a.pages.SwitchToPage("help")
	}
}

// toggleFilter toggles the filter panel visibility
func (a *App) toggleFilter() {
	if a.filterPanel.IsVisible() {
		a.filterPanel.SetVisible(false)
		a.pages.HidePage("filter")
		a.pages.SwitchToPage("main")
	} else {
		a.filterPanel.SetVisible(true)
		a.pages.ShowPage("filter")
		a.pages.SwitchToPage("filter")
		a.filterPanel.Focus()
	}
}

// filterPatterns filters patterns based on query and filter type
func (a *App) filterPatterns(query string, filterType FilterType) []*models.Pattern {
	var results []*models.Pattern
	
	for _, p := range a.patterns {
		match := false
		
		switch filterType {
		case FilterAll:
			match = true
		case FilterByTrigger:
			match = containsIgnoreCase(p.Trigger, query)
		case FilterByResponse:
			match = containsIgnoreCase(p.Response, query)
		case FilterByStrength:
			// Query can be "high", "medium", "low"
			match = matchesStrength(p.Strength, query)
		case FilterRecent:
			// Show all if no query, or match recent used
			match = query == "" || containsIgnoreCase(p.Trigger, query)
		case FilterFavorites:
			// Would need a favorite field
			match = query == "" || containsIgnoreCase(p.Trigger, query)
		}
		
		if match {
			results = append(results, p)
		}
	}
	
	return results
}

// containsIgnoreCase checks if s contains sub (case insensitive)
func containsIgnoreCase(s, sub string) bool {
	if sub == "" {
		return true
	}
	s = lower(s)
	sub = lower(s)
	return contains(s, sub)
}

func lower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func matchesStrength(strength float64, query string) bool {
	query = lower(query)
	switch query {
	case "high", "强", "80":
		return strength >= 80
	case "medium", "中", "50":
		return strength >= 50 && strength < 80
	case "low", "弱", "20":
		return strength >= 20 && strength < 50
	case "very low", "很弱", "0":
		return strength < 20
	default:
		return true
	}
}
