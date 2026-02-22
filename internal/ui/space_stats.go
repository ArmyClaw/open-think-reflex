package ui

import (
	"fmt"
	"time"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SpaceStatsPanel displays space statistics
type SpaceStatsPanel struct {
	view       *tview.TextView
	spaces     []*models.Space
	title      string
	visible    bool
}

// NewSpaceStatsPanel creates a new space stats panel
func NewSpaceStatsPanel() *SpaceStatsPanel {
	statsView := tview.NewTextView()
	statsView.SetDynamicColors(true)
	statsView.SetScrollable(true)
	statsView.SetBorder(true)
	statsView.SetTitle(" Space Statistics ")

	panel := &SpaceStatsPanel{
		view:    statsView,
		title:   "Space Statistics",
		visible: false,
	}

	return panel
}

// SetSpaces updates the spaces to show stats for
func (p *SpaceStatsPanel) SetSpaces(spaces []*models.Space) {
	p.spaces = spaces
	p.updateStats()
}

// updateStats refreshes the space statistics display
func (p *SpaceStatsPanel) updateStats() {
	if len(p.spaces) == 0 {
		p.view.SetText("[yellow]No spaces available[white]")
		return
	}

	var totalPatterns int
	defaultSpaceCount := 0
	recentSpaces := 0

	oneDayAgo := time.Now().Add(-24 * time.Hour)

	for _, space := range p.spaces {
		totalPatterns += space.PatternCount
		if space.DefaultSpace {
			defaultSpaceCount++
		}
		if space.CreatedAt.After(oneDayAgo) {
			recentSpaces++
		}
	}

	// Find space with most patterns
	var topSpace *models.Space
	for _, space := range p.spaces {
		if topSpace == nil || space.PatternCount > topSpace.PatternCount {
			topSpace = space
		}
	}

	// Build stats display
	statsText := fmt.Sprintf(`[bold]📊 Space Statistics[white]

[cyan]Total Spaces:[white] %d
[cyan]Total Patterns:[white] %d

[cyan]Space Details:[white]
%s

[cyan]Quick Stats:[white]
  ✓ Default Spaces: %d
  🆕 Created (24h): %d
  📈 Most Active: %s (%d patterns)

[dim]Last updated: %s[white]`,
		len(p.spaces),
		totalPatterns,
		p.buildSpaceList(),
		defaultSpaceCount,
		recentSpaces,
		topSpace.Name,
		topSpace.PatternCount,
		time.Now().Format("15:04:05"),
	)

	p.view.SetText(statsText)
}

// buildSpaceList creates a formatted list of spaces
func (p *SpaceStatsPanel) buildSpaceList() string {
	list := ""
	for _, space := range p.spaces {
		icon := "  "
		if space.DefaultSpace {
			icon = "✓ "
		}
		list += fmt.Sprintf(`  %s[bold]%s[white]: %d patterns
    [dim]%s[white]
`,
			icon,
			space.Name,
			space.PatternCount,
			space.Description,
		)
	}
	return list
}

// GetView returns the underlying tview primitive
func (p *SpaceStatsPanel) GetView() tview.Primitive {
	return p.view
}

// SetTitle sets the panel title
func (p *SpaceStatsPanel) SetTitle(title string) {
	p.title = title
}

// SetBorderColor sets the border color
func (p *SpaceStatsPanel) SetBorderColor(color tcell.Color) {
	p.view.SetBorderColor(color)
}

// Focus sets focus to this panel
func (p *SpaceStatsPanel) Focus(delegate func(p tview.Primitive)) {
	delegate(p.view)
}

// HasFocus returns true if this panel has focus
func (p *SpaceStatsPanel) HasFocus() bool {
	return p.view.HasFocus()
}

// SetVisible sets the panel visibility
func (p *SpaceStatsPanel) SetVisible(v bool) {
	p.visible = v
	if v {
		p.view.SetBorder(true)
	} else {
		p.view.SetBorder(false)
	}
}

// IsVisible returns whether the panel is visible
func (p *SpaceStatsPanel) IsVisible() bool {
	return p.visible
}
