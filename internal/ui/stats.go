package ui

import (
	"fmt"
	"time"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// StatsPanel displays pattern statistics
type StatsPanel struct {
	view      *tview.TextView
	patterns  []models.Pattern
	title     string
	visible   bool
}

// NewStatsPanel creates a new stats panel
func NewStatsPanel() *StatsPanel {
	statsView := tview.NewTextView()
	statsView.SetDynamicColors(true)
	statsView.SetScrollable(true)
	statsView.SetBorder(true)
	statsView.SetTitle(" Pattern Statistics ")

	panel := &StatsPanel{
		view:    statsView,
		title:   "Statistics",
		visible: false,
	}

	return panel
}

// SetPatterns updates the patterns to show stats for
func (p *StatsPanel) SetPatterns(patterns []models.Pattern) {
	p.patterns = patterns
	p.updateStats()
}

// updateStats refreshes the statistics display
func (p *StatsPanel) updateStats() {
	if len(p.patterns) == 0 {
		p.view.SetText("[yellow]No patterns available[white]")
		return
	}

	var totalStrength float64
	activeCount := 0
	decayedCount := 0
	highStrength := 0
	mediumStrength := 0
	lowStrength := 0
	recentlyUsed := 0

	oneDayAgo := time.Now().Add(-24 * time.Hour)

	for _, pattern := range p.patterns {
		totalStrength += pattern.Strength

		// Count by strength level
		if pattern.Strength >= 80 {
			highStrength++
		} else if pattern.Strength >= 50 {
			mediumStrength++
		} else {
			lowStrength++
		}

		// Count by status (assuming threshold is 60)
		if pattern.DecayEnabled && pattern.Strength < 60 {
			decayedCount++
		} else {
			activeCount++
		}

		// Count recently used
		if pattern.LastUsedAt != nil && pattern.LastUsedAt.After(oneDayAgo) {
			recentlyUsed++
		}
	}

	avgStrength := totalStrength / float64(len(p.patterns))

	// Build stats display
	statsText := fmt.Sprintf(`[bold]📊 Pattern Statistics[white]

[cyan]Total Patterns:[white] %d

[cyan]Status:[white]
  ✓ Active: %d
  ○ Decayed: %d

[cyan]Strength Distribution:[white]
  🔴 High (≥80): %d
  🟡 Medium (50-79): %d
  🟢 Low (<50): %d
  [yellow]Avg Strength: %.1f[white]

[cyan]Usage (24h):[white]
  Recently Used: %d patterns

[cyan]Strength Bar:[white]
%s

[dim]Last updated: %s[white]`,
		len(p.patterns),
		activeCount,
		decayedCount,
		highStrength,
		mediumStrength,
		lowStrength,
		avgStrength,
		recentlyUsed,
		p.buildStrengthBar(highStrength, mediumStrength, lowStrength),
		time.Now().Format("15:04:05"),
	)

	p.view.SetText(statsText)
}

// buildStrengthBar creates a visual bar showing strength distribution
func (p *StatsPanel) buildStrengthBar(high, medium, low int) string {
	total := high + medium + low
	if total == 0 {
		return "[red]░░░░░░░░░░[white] (no data)"
	}

	const barLength = 12
	highLen := (high * barLength) / total
	medLen := (medium * barLength) / total
	lowLen := barLength - highLen - medLen

	if lowLen < 0 {
		lowLen = 0
	}

	bar := ""
	for i := 0; i < highLen; i++ {
		bar += "[red]█[white]"
	}
	for i := 0; i < medLen; i++ {
		bar += "[yellow]█[white]"
	}
	for i := 0; i < lowLen; i++ {
		bar += "[green]█[white]"
	}

	return fmt.Sprintf("%s (%d total)", bar, total)
}

// GetView returns the underlying tview primitive
func (p *StatsPanel) GetView() tview.Primitive {
	return p.view
}

// SetTitle sets the panel title
func (p *StatsPanel) SetTitle(title string) {
	p.title = title
	// TextView doesn't have SetBorderTitle, so we'll just update the text
	// The border title is set during construction
	_ = title // suppress unused warning
}

// SetBorderColor sets the border color
func (p *StatsPanel) SetBorderColor(color tcell.Color) {
	p.view.SetBorderColor(color)
}

// Focus sets focus to this panel
func (p *StatsPanel) Focus(delegate func(p tview.Primitive)) {
	delegate(p.view)
}

// HasFocus returns true if this panel has focus
func (p *StatsPanel) HasFocus() bool {
	return p.view.HasFocus()
}

// SetVisible sets the panel visibility
func (p *StatsPanel) SetVisible(v bool) {
	p.visible = v
	if v {
		p.view.SetBorder(true)
	} else {
		p.view.SetBorder(false)
	}
}

// IsVisible returns whether the panel is visible
func (p *StatsPanel) IsVisible() bool {
	return p.visible
}
