package ui

import (
	"fmt"

	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/rivo/tview"
)

// ThoughtChainView displays the thought chain tree (Layer 1)
type ThoughtChainView struct {
	view     *tview.TextView
	theme    *Theme
	results  []contracts.MatchResult
	selected int
}

// NewThoughtChainView creates a new thought chain view
func NewThoughtChainView(theme *Theme) *ThoughtChainView {
	v := &ThoughtChainView{
		theme:    theme,
		selected: 0,
	}
	
	v.view = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	v.view.SetBorder(true).SetTitle("💭 Thought Chain")
	
	v.view.SetBackgroundColor(theme.Background)
	v.view.SetBorderColor(theme.Border)
	v.view.SetTitleColor(theme.Primary)
	
	v.render()
	
	return v
}

// SetResults sets the match results to display
func (v *ThoughtChainView) SetResults(results []contracts.MatchResult) {
	v.results = results
	v.selected = 0
	v.render()
}

// Clear clears the thought chain
func (v *ThoughtChainView) Clear() {
	v.results = nil
	v.selected = 0
	v.view.SetText("")
}

// Selected returns the currently selected index
func (v *ThoughtChainView) Selected() int {
	return v.selected
}

// SelectNext selects the next item
func (v *ThoughtChainView) SelectNext() {
	if len(v.results) > 0 {
		v.selected = (v.selected + 1) % len(v.results)
		v.render()
	}
}

// SelectPrev selects the previous item
func (v *ThoughtChainView) SelectPrev() {
	if len(v.results) > 0 {
		v.selected = (v.selected - 1 + len(v.results)) % len(v.results)
		v.render()
	}
}

// GetSelectedResult returns the currently selected result
func (v *ThoughtChainView) GetSelectedResult() *contracts.MatchResult {
	if v.selected >= 0 && v.selected < len(v.results) {
		return &v.results[v.selected]
	}
	return nil
}

func (v *ThoughtChainView) render() {
	if len(v.results) == 0 {
		v.view.SetText(fmt.Sprintf("[%s]No matches yet[white]\n\nType a query to see thought branches",
			v.theme.TextDim))
		return
	}
	
	text := ""
	
	for i, r := range v.results {
		prefix := "  "
		selected := false
		
		if i == v.selected {
			prefix = "▶ "
			selected = true
		}
		
		// Format based on match type
		matchType := r.Branch
		
		if selected {
			text += fmt.Sprintf("[%s]%s[%s] %s[white]\n", 
				v.theme.Selected, prefix, v.theme.Accent, r.Pattern.Trigger)
			text += fmt.Sprintf("     ├ Confidence: [%.0f%%] %s\n", r.Confidence, matchType)
			text += fmt.Sprintf("     ├ Response: %s\n", truncate(r.Pattern.Response, 40))
			text += fmt.Sprintf("     └ Strength: %.1f / %.1f\n", r.Pattern.Strength, r.Pattern.Threshold)
		} else {
			text += fmt.Sprintf("%s[%s] %s[white] (%.0f%%)\n", 
				prefix, v.theme.Secondary, r.Pattern.Trigger, r.Confidence)
		}
	}
	
	v.view.SetText(text)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
