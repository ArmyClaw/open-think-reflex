package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

// SettingsPanel displays application settings
type SettingsPanel struct {
	view       *tview.Flex
	theme      *Theme
	visible    bool
	app        *App
}

// NewSettingsPanel creates a new settings panel
func NewSettingsPanel(theme *Theme, app *App) *SettingsPanel {
	s := &SettingsPanel{
		theme:   theme,
		visible: false,
		app:     app,
	}
	s.view = s.createView()
	return s
}

func (s *SettingsPanel) createView() *tview.Flex {
	flex := tview.NewFlex()
	flex.SetBorder(true)
	flex.SetBorderColor(s.theme.Primary)
	flex.SetBackgroundColor(s.theme.Background)
	flex.SetTitle(" Settings (press , to close) ")

	// Create content
	content := s.createContent()
	flex.AddItem(content, 0, 1, true)

	return flex
}

func (s *SettingsPanel) createContent() *tview.Flex {
	content := tview.NewFlex()
	content.SetDirection(tview.FlexRow)

	// Title
	title := tview.NewTextView()
	title.SetText("⚙️ Application Settings")
	title.SetTextColor(s.theme.Primary)
	title.SetTextAlign(tview.AlignCenter)
	content.AddItem(title, 1, 0, false)

	// Separator
	sep := tview.NewTextView()
	sep.SetText("\n────────────────────────────────────────\n")
	sep.SetTextColor(s.theme.Border)
	content.AddItem(sep, 1, 0, false)

	// Settings sections
	general := s.createSection("General",
		[]string{
			"Database Path: ~/.otr/otr.db",
			"Log Level: info",
			"Theme: " + s.app.theme.Name,
		})
	content.AddItem(general, 0, 1, false)

	ai := s.createSection("AI Configuration",
		[]string{
			"Provider: Anthropic (Claude)",
			"Model: claude-3-5-sonnet-20241022",
			"Max Tokens: 4096",
			"Temperature: 0.7",
		})
	content.AddItem(ai, 0, 1, false)

	matching := s.createSection("Matching",
		[]string{
			"Default Threshold: 30",
			"Max Results: 10",
			"Cache TTL: 1 minute",
			"Exact Match First: enabled",
		})
	content.AddItem(matching, 0, 1, false)

	pattern := s.createSection("Pattern Management",
		[]string{
			"Strength Reinforce: +5 per use",
			"Decay Rate: -1 per day",
			"Threshold: 80 for activation",
			"Auto-cleanup: disabled",
		})
	content.AddItem(pattern, 0, 1, false)

	// Keyboard shortcuts info
	shortcuts := tview.NewTextView()
	shortcuts.SetText("\n────────────────────────────────────────\n\nKeyboard Shortcuts:\n  [,]     Close settings\n  [t]     Toggle theme\n  [Tab]   Switch mode\n  [q]     Quit")
	shortcuts.SetTextColor(s.theme.Secondary)
	shortcuts.SetTextAlign(tview.AlignLeft)
	content.AddItem(shortcuts, 0, 1, false)

	return content
}

func (s *SettingsPanel) createSection(title string, items []string) *tview.TextView {
	view := tview.NewTextView()
	view.SetBackgroundColor(s.theme.Background)
	view.SetTextColor(s.theme.Text)

	text := fmt.Sprintf("📁 %s\n", title)
	for _, item := range items {
		text += fmt.Sprintf("   • %s\n", item)
	}
	view.SetText(text)

	return view
}

// View returns the underlying tview primitive
func (s *SettingsPanel) View() *tview.Flex {
	return s.view
}

// Show displays the settings panel
func (s *SettingsPanel) Show() {
	s.visible = true
}

// Hide hides the settings panel
func (s *SettingsPanel) Hide() {
	s.visible = false
}

// Toggle toggles visibility
func (s *SettingsPanel) Toggle() {
	s.visible = !s.visible
}

// IsVisible returns whether the panel is visible
func (s *SettingsPanel) IsVisible() bool {
	return s.visible
}

// SetVisible sets visibility
func (s *SettingsPanel) SetVisible(visible bool) {
	s.visible = visible
}

// SetTheme updates the theme
func (s *SettingsPanel) SetTheme(theme *Theme) {
	s.theme = theme
	s.view.SetBorderColor(theme.Primary)
	s.view.SetBackgroundColor(theme.Background)
	s.view.SetTitleColor(theme.Primary)
}
