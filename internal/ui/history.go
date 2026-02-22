package ui

import (
	"fmt"
	"time"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/rivo/tview"
)

// HistoryEntry represents a single history entry
type HistoryEntry struct {
	ID        string
	Input     string
	Response  string
	Timestamp time.Time
	PatternID string
}

// HistoryPanel displays interaction history
type HistoryPanel struct {
	view        *tview.TextView
	entries     []HistoryEntry
	theme       *Theme
	visible     bool
	onSelect    func(entry HistoryEntry)
}

// NewHistoryPanel creates a new history panel
func NewHistoryPanel(theme *Theme, onSelect func(entry HistoryEntry)) *HistoryPanel {
	hp := &HistoryPanel{
		theme:    theme,
		entries:  make([]HistoryEntry, 0),
		visible:  false,
		onSelect: onSelect,
	}
	
	hp.view = tview.NewTextView()
	hp.view.SetBorder(true)
	hp.view.SetBorderColor(theme.Border)
	hp.view.SetBackgroundColor(theme.Background)
	hp.view.SetTextColor(theme.Text)
	hp.view.SetTitle("📜 History [h] to close")
	hp.view.SetTitleColor(theme.Primary)
	hp.view.SetScrollable(true)
	
	return hp
}

// View returns the panel's view
func (hp *HistoryPanel) View() *tview.TextView {
	return hp.view
}

// SetTheme updates the panel's theme
func (hp *HistoryPanel) SetTheme(theme *Theme) {
	hp.theme = theme
	hp.view.SetBorderColor(theme.Border)
	hp.view.SetBackgroundColor(theme.Background)
	hp.view.SetTextColor(theme.Text)
	hp.view.SetTitleColor(theme.Primary)
}

// SetVisible sets panel visibility
func (hp *HistoryPanel) SetVisible(visible bool) {
	hp.visible = visible
}

// IsVisible checks if panel is visible
func (hp *HistoryPanel) IsVisible() bool {
	return hp.visible
}

// AddEntry adds a new history entry
func (hp *HistoryPanel) AddEntry(entry HistoryEntry) {
	hp.entries = append([]HistoryEntry{entry}, hp.entries...)
	
	// Keep only last 100 entries
	if len(hp.entries) > 100 {
		hp.entries = hp.entries[:100]
	}
	
	hp.render()
}

// SetEntries sets all history entries
func (hp *HistoryPanel) SetEntries(entries []HistoryEntry) {
	hp.entries = entries
	hp.render()
}

// render renders the history list
func (hp *HistoryPanel) render() {
	hp.view.Clear()
	
	if len(hp.entries) == 0 {
		fmt.Fprint(hp.view, "\nNo history yet.\n\nYour interactions will appear here.")
		return
	}
	
	// Header
	hp.view.SetText(fmt.Sprintf("%-6s %-20s %-40s %s\n%s\n",
		"#",
		"Time",
		"Input",
		"Pattern",
		"─────────────────────────────────────────────────────────────────────────"))
	
	// Entries (newest first)
	for i, entry := range hp.entries {
		input := entry.Input
		if len(input) > 38 {
			input = input[:35] + "..."
		}
		
		pattern := entry.PatternID
		if len(pattern) > 10 {
			pattern = pattern[:7] + "..."
		}
		
		timestamp := entry.Timestamp.Format("15:04:05")
		
		// Highlight newest entry
		if i == 0 {
			hp.view.SetTextColor(hp.theme.Primary)
		}
		
		hp.view.Write([]byte(fmt.Sprintf("%-6d %-20s %-40s %s\n",
			i+1,
			timestamp,
			input,
			pattern)))
	}
	
	// Instructions
	hp.view.Write([]byte(fmt.Sprintf("\n%s\n%s",
		"─────────────────────────────────────────────────────────────────────────",
		"[↑/↓] Navigate | [Enter] View Details | [d] Delete | [c] Clear All")))
}

// GetEntry returns entry at index
func (hp *HistoryPanel) GetEntry(index int) *HistoryEntry {
	if index >= 0 && index < len(hp.entries) {
		return &hp.entries[index]
	}
	return nil
}

// DeleteEntry removes entry at index
func (hp *HistoryPanel) DeleteEntry(index int) {
	if index >= 0 && index < len(hp.entries) {
		hp.entries = append(hp.entries[:index], hp.entries[index+1:]...)
		hp.render()
	}
}

// ClearAll removes all entries
func (hp *HistoryPanel) ClearAll() {
	hp.entries = make([]HistoryEntry, 0)
	hp.render()
}

// SetPatternsHistory converts patterns to history entries (for demo/testing)
func (hp *HistoryPanel) SetPatternsHistory(patterns []*models.Pattern) {
	hp.entries = make([]HistoryEntry, 0)
	
	for i, p := range patterns {
		entry := HistoryEntry{
			ID:        fmt.Sprintf("hist-%d", i),
			Input:     p.Trigger,
			Response:  p.Response,
			Timestamp: p.UpdatedAt,
			PatternID: p.ID,
		}
		hp.entries = append(hp.entries, entry)
	}
	
	// Sort by timestamp descending (newest first)
	for i := 0; i < len(hp.entries)-1; i++ {
		for j := i + 1; j < len(hp.entries); j++ {
			if hp.entries[j].Timestamp.After(hp.entries[i].Timestamp) {
				hp.entries[i], hp.entries[j] = hp.entries[j], hp.entries[i]
			}
		}
	}
	
	hp.render()
}
