package ui

import (
	"context"
	"fmt"
	"time"

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
	helpPanel     *HelpPanel
	shortcutBar   *ShortcutBar
	filterPanel   *FilterPanel
	statsPanel    *StatsPanel
	spaceStatsPanel *SpaceStatsPanel
	patternForm   *PatternFormPanel
	settingsPanel *SettingsPanel
	historyPanel  *HistoryPanel
	
	// State
	currentSpace *models.Space
	patterns     []*models.Pattern
	results      []contracts.MatchResult
	mode         AppMode // input, navigation
	showStats    bool
	showSpaceStats bool
	showForm     bool
	showSettings bool
	showHistory  bool
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
	a.statsPanel = NewStatsPanel()
	a.spaceStatsPanel = NewSpaceStatsPanel()
	a.showStats = false
	a.showSpaceStats = false
	a.showForm = false
	
	// Create pattern form panel
	a.patternForm = NewPatternFormPanel(a.theme, a.handlePatternSave, a.handlePatternCancel)
	
	// Create settings panel
	a.settingsPanel = NewSettingsPanel(a.theme, a)
	a.showSettings = false
	
	// Create history panel
	a.historyPanel = NewHistoryPanel(a.theme, a.handleHistorySelect)
	a.showHistory = false
	
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
	a.pages.AddPage("stats", a.statsPanel.GetView(), false, false)
	a.pages.AddPage("spaceStats", a.spaceStatsPanel.GetView(), false, false)
	a.pages.AddPage("form", a.patternForm.GetView(), false, false)
	a.pages.AddPage("settings", a.settingsPanel.View(), false, false)
	a.pages.AddPage("history", a.historyPanel.View(), false, false)
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
			case 's':
				// Toggle stats panel
				a.toggleStats()
				return nil
			case 'S':
				// Toggle space stats panel
				a.toggleSpaceStats()
				return nil
			case ',':
				// Toggle settings panel
				a.toggleSettings()
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
			case 'c':
				// Create new pattern (only in input mode)
				if a.mode == ModeInput {
					a.showPatternForm(false, nil)
				}
				return nil
			case 'e':
				// Edit selected pattern (only in navigation mode)
				if a.mode == ModeNavigation {
					result := a.thoughtChain.GetSelectedResult()
					if result != nil {
						a.showPatternForm(true, result.Pattern)
					}
				}
				return nil
			case 'd':
				// Delete selected pattern (only in navigation mode)
				if a.mode == ModeNavigation {
					result := a.thoughtChain.GetSelectedResult()
					if result != nil {
						a.deleteSelectedPattern(result.Pattern)
					}
				}
				return nil
			case 'y':
				// Toggle history panel
				a.toggleHistory()
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
	
	// Record in history (for each matched pattern)
	for _, r := range results {
		entry := HistoryEntry{
			ID:        fmt.Sprintf("hist-%d", len(a.historyPanel.entries)),
			Input:     text,
			Response:  r.Pattern.Response,
			Timestamp: time.Now(),
			PatternID: r.Pattern.ID,
		}
		a.historyPanel.AddEntry(entry)
	}
	
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

// toggleStats toggles the statistics panel visibility
func (a *App) toggleStats() {
	if a.showStats {
		a.showStats = false
		a.statsPanel.SetVisible(false)
		a.pages.HidePage("stats")
		a.pages.SwitchToPage("main")
		a.app.SetFocus(a.input.view)
	} else {
		// Update stats with current patterns
		var patternModels []models.Pattern
		for _, p := range a.patterns {
			patternModels = append(patternModels, *p)
		}
		a.statsPanel.SetPatterns(patternModels)
		
		a.showStats = true
		a.statsPanel.SetVisible(true)
		a.pages.ShowPage("stats")
		a.pages.SwitchToPage("stats")
	}
}

// toggleSpaceStats toggles the space statistics panel visibility
func (a *App) toggleSpaceStats() {
	if a.showSpaceStats {
		a.showSpaceStats = false
		a.spaceStatsPanel.SetVisible(false)
		a.pages.HidePage("spaceStats")
		a.pages.SwitchToPage("main")
		a.app.SetFocus(a.input.view)
	} else {
		// Get spaces from storage
		spaces, err := a.storage.ListSpaces(context.Background())
		if err != nil {
			a.output.SetOutput("Error loading spaces: " + err.Error())
			return
		}
		a.spaceStatsPanel.SetSpaces(spaces)
		
		a.showSpaceStats = true
		a.spaceStatsPanel.SetVisible(true)
		a.pages.ShowPage("spaceStats")
		a.pages.SwitchToPage("spaceStats")
	}
}

// toggleSettings toggles the settings panel visibility
func (a *App) toggleSettings() {
	if a.showSettings {
		a.showSettings = false
		a.settingsPanel.SetVisible(false)
		a.pages.HidePage("settings")
		a.pages.SwitchToPage("main")
		a.app.SetFocus(a.input.view)
	} else {
		a.showSettings = true
		a.settingsPanel.SetVisible(true)
		a.pages.ShowPage("settings")
		a.pages.SwitchToPage("settings")
	}
}

// toggleHistory toggles the history panel visibility
func (a *App) toggleHistory() {
	if a.showHistory {
		a.showHistory = false
		a.historyPanel.SetVisible(false)
		a.pages.HidePage("history")
		a.pages.SwitchToPage("main")
		a.app.SetFocus(a.input.view)
	} else {
		// Populate history from patterns
		a.historyPanel.SetPatternsHistory(a.patterns)
		
		a.showHistory = true
		a.historyPanel.SetVisible(true)
		a.pages.ShowPage("history")
		a.pages.SwitchToPage("history")
	}
}

// handleHistorySelect handles history entry selection
func (a *App) handleHistorySelect(entry HistoryEntry) {
	a.output.SetFormattedOutput("History Detail",
		fmt.Sprintf("Input: %s\n\nResponse: %s\n\nTime: %s",
			entry.Input,
			entry.Response,
			entry.Timestamp.Format("2006-01-02 15:04:05")))
}

// showPatternForm shows the pattern form for creating or editing
func (a *App) showPatternForm(isEdit bool, pattern *models.Pattern) {
	if isEdit {
		a.patternForm.SetEditMode(pattern)
	} else {
		a.patternForm.SetCreateMode()
	}
	
	a.showForm = true
	a.pages.ShowPage("form")
	a.pages.SwitchToPage("form")
}

// handlePatternSave handles pattern save from the form
func (a *App) handlePatternSave(pattern *models.Pattern) {
	ctx := context.Background()
	
	var err error
	if pattern.ID == "" {
		// Create new pattern
		err = a.storage.SavePattern(ctx, pattern)
	} else {
		// Update existing pattern
		err = a.storage.UpdatePattern(ctx, pattern)
	}
	
	if err != nil {
		a.output.SetStatus(fmt.Sprintf("Error saving pattern: %v", err), false)
		return
	}
	
	// Refresh pattern list
	if err := a.loadData(ctx); err != nil {
		a.output.SetStatus(fmt.Sprintf("Error reloading patterns: %v", err), false)
		return
	}
	
	// Update status bar
	a.statusBar.SetPatternCount(len(a.patterns))
	
	// Return to main view
	a.showForm = false
	a.pages.HidePage("form")
	a.pages.SwitchToPage("main")
	a.app.SetFocus(a.input.view)
	
	a.output.SetStatus(fmt.Sprintf("Pattern '%s' saved successfully", pattern.Trigger), true)
}

// handlePatternCancel handles pattern form cancellation
func (a *App) handlePatternCancel() {
	a.showForm = false
	a.pages.HidePage("form")
	a.pages.SwitchToPage("main")
	a.app.SetFocus(a.input.view)
}

// deleteSelectedPattern deletes the selected pattern after confirmation
func (a *App) deleteSelectedPattern(pattern *models.Pattern) {
	ctx := context.Background()
	
	// Show confirmation
	a.output.SetOutput(fmt.Sprintf("Delete pattern '%s'?\n\nPress 'y' to confirm or any other key to cancel",
		pattern.Trigger))
	
	// Set up a one-time input handler for confirmation
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Restore original handler
		a.setupKeyBindings()
		
		if event.Key() == tcell.KeyRune && (event.Rune() == 'y' || event.Rune() == 'Y') {
			// Delete the pattern
			if err := a.storage.DeletePattern(ctx, pattern.ID); err != nil {
				a.output.SetStatus(fmt.Sprintf("Error deleting pattern: %v", err), false)
				return nil
			}
			
			// Refresh pattern list
			if err := a.loadData(ctx); err != nil {
				a.output.SetStatus(fmt.Sprintf("Error reloading patterns: %v", err), false)
				return nil
			}
			
			// Update status bar
			a.statusBar.SetPatternCount(len(a.patterns))
			a.output.SetStatus(fmt.Sprintf("Pattern '%s' deleted", pattern.Trigger), true)
		}
		
		return event
	})
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
