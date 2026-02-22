package ui

import (
	"github.com/rivo/tview"
)

// HelpPanel displays keyboard shortcuts and help information
type HelpPanel struct {
	view       *tview.TextView
	theme      *Theme
	visible    bool
}

// NewHelpPanel creates a new help panel
func NewHelpPanel(theme *Theme) *HelpPanel {
	h := &HelpPanel{
		theme:   theme,
		visible: false,
	}
	h.view = h.createView()
	return h
}

func (h *HelpPanel) createView() *tview.TextView {
	view := tview.NewTextView()
	view.SetBorder(true)
	view.SetBorderColor(h.theme.Primary)
	view.SetBackgroundColor(h.theme.Background)
	view.SetTextColor(h.theme.Text)
	view.SetTitle(" Help (press ? to close) ")
	
	// Set help content
	view.SetText(h.getHelpContent())
	
	return view
}

func (h *HelpPanel) getHelpContent() string {
	return `
┌─────────────────────────────────────────────────────────────────┐
│                      KEYBOARD SHORTCUTS                         │
├─────────────────────────────────────────────────────────────────┤
│  Global                                                          │
│  ─────────────────────────────────────────────────────────────  │
│  [Tab]       Switch between Input/Navigation mode               │
│  [↑/↓]       Navigate thought branches (Navigation mode)         │
│  [←]         Collapse branch / Go back                          │
│  [→]         Expand branch / Select                             │
│  [Enter]     Use selected response                              │
│  [Space]     Toggle selection                                   │
│  [t]         Toggle Light/Dark theme                            │
│  [/]         Toggle Filter panel                                │
│  [s]         Toggle Statistics panel                            │
│  [S]         Toggle Space Statistics panel                      │
│  [,]         Toggle Settings panel                             │
│  [y]         Toggle History panel                              │
│  [c]         Create new pattern (Input mode)                    │
│  [e]         Edit selected pattern (Navigation mode)            │
│  [d]         Delete selected pattern (Navigation mode)         │
│  [h/?]       Show/Hide this help                                │
│  [q/Esc]     Quit application                                   │
├─────────────────────────────────────────────────────────────────┤
│  Vim-style (Navigation mode)                                    │
│  ─────────────────────────────────────────────────────────────  │
│  [k]         Move up                                            │
│  [j]         Move down                                           │
│  [h]         Collapse / Go back                                  │
│  [l]         Expand / Select                                     │
│  [Esc]       Return to input mode                               │
├─────────────────────────────────────────────────────────────────┤
│  Modes                                                           │
│  ─────────────────────────────────────────────────────────────  │
│  INPUT:       Type queries in the input area                     │
│  NAVIGATE:    Use arrow keys to browse results                  │
└─────────────────────────────────────────────────────────────────┘

Tip: Press ? or h to close this panel
`
}

// View returns the underlying tview primitive
func (h *HelpPanel) View() *tview.TextView {
	return h.view
}

// Show displays the help panel
func (h *HelpPanel) Show() {
	h.visible = true
}

// Hide hides the help panel
func (h *HelpPanel) Hide() {
	h.visible = false
}

// Toggle toggles visibility
func (h *HelpPanel) Toggle() {
	h.visible = !h.visible
}

// IsVisible returns whether the panel is visible
func (h *HelpPanel) IsVisible() bool {
	return h.visible
}

// SetTheme updates the theme
func (h *HelpPanel) SetTheme(theme *Theme) {
	h.theme = theme
	h.view.SetBorderColor(theme.Primary)
	h.view.SetBackgroundColor(theme.Background)
	h.view.SetTextColor(theme.Text)
	h.view.SetTitleColor(theme.Primary)
}

// ShortcutBar displays contextual shortcuts at the bottom of the screen
type ShortcutBar struct {
	view  *tview.TextView
	theme *Theme
}

// NewShortcutBar creates a new shortcut bar
func NewShortcutBar(theme *Theme) *ShortcutBar {
	s := &ShortcutBar{
		theme: theme,
	}
	s.view = s.createView()
	return s
}

func (s *ShortcutBar) createView() *tview.TextView {
	view := tview.NewTextView()
	view.SetBackgroundColor(s.theme.Background)
	view.SetTextColor(s.theme.Secondary)
	view.SetTextAlign(tview.AlignCenter)
	
	// Set initial shortcuts
	view.SetText(s.getShortcuts(ModeInput))
	
	return view
}

func (s *ShortcutBar) getShortcuts(mode AppMode) string {
	if mode == ModeNavigation {
		return " [↑/↓] Navigate | [Enter] Select | [←/→] Expand/Collapse | [e] Edit | [d] Delete | [/] Filter | [s] Stats | [S] Spaces | [y] History | [,] Settings | [?] Help | [q] Quit "
	}
	return " [Tab] Switch | [c] Create | [t] Theme | [/] Filter | [s] Stats | [S] Spaces | [y] History | [,] Settings | [?] Help | [q] Quit "
}

// View returns the underlying tview primitive
func (s *ShortcutBar) View() *tview.TextView {
	return s.view
}

// SetMode updates shortcuts based on current mode
func (s *ShortcutBar) SetMode(mode AppMode) {
	s.view.SetText(s.getShortcuts(mode))
}

// SetTheme updates the theme
func (s *ShortcutBar) SetTheme(theme *Theme) {
	s.theme = theme
	s.view.SetBackgroundColor(theme.Background)
	s.view.SetTextColor(theme.Secondary)
}
