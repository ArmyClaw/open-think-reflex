// Package ui provides the terminal user interface (TUI) for Open-Think-Reflex.
// Uses tview for rich terminal UI components.
//
// The UI consists of three main layers:
//   - Layer 1: Thought Chain Tree (pattern matching results)
//   - Layer 2: Output View (AI-generated content)
//   - Layer 3: Input Area (user input)
package ui

import (
	"fmt"

	"github.com/ArmyClaw/open-think-reflex/pkg/contracts"
	"github.com/rivo/tview"
)

// ThoughtChainNode represents a node in the thought chain tree.
// Each node corresponds to a match result, optionally with child branches.
type ThoughtChainNode struct {
	Result    *contracts.MatchResult // The match result for this node
	Children  []*ThoughtChainNode   // Child branches (sub-matches)
	Expanded  bool                  // Whether children are visible
	Level     int                   // Tree depth (0 = root)
}

// ThoughtChainView displays the thought chain tree (Layer 1).
// Shows pattern matching results as an interactive tree structure.
type ThoughtChainView struct {
	view     *tview.TextView // tview text view component
	theme    *Theme          // Active theme for colors
	results  []contracts.MatchResult // Raw match results
	nodes    []*ThoughtChainNode     // Tree structure
	selected int          // Currently selected node index
	focused  bool         // Whether this view has focus
}

// NewThoughtChainView creates a new thought chain view with the given theme.
func NewThoughtChainView(theme *Theme) *ThoughtChainView {
	v := &ThoughtChainView{
		theme:    theme,
		selected: 0,
		focused:  false,
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

// SetResults sets the match results to display and builds tree nodes
func (v *ThoughtChainView) SetResults(results []contracts.MatchResult) {
	v.results = results
	v.selected = 0
	
	// Build tree nodes from flat results
	v.nodes = v.buildTree(results)
	v.render()
}

// buildTree converts flat results into a tree structure
func (v *ThoughtChainView) buildTree(results []contracts.MatchResult) []*ThoughtChainNode {
	nodes := make([]*ThoughtChainNode, len(results))
	
	for i, r := range results {
		node := &ThoughtChainNode{
			Result:   &r,
			Children:  nil,
			Expanded:  i < 3, // Auto-expand first 3
			Level:    0,
		}
		
		// Add simulated children for visual tree effect
		if r.Confidence > 80 {
			node.Children = []*ThoughtChainNode{
				{Result: &r, Expanded: false, Level: 1},
			}
		}
		
		nodes[i] = node
	}
	
	return nodes
}

// SetFocused sets the focus state
func (v *ThoughtChainView) SetFocused(focused bool) {
	v.focused = focused
	if focused {
		v.view.SetBorderColor(v.theme.Accent)
	} else {
		v.view.SetBorderColor(v.theme.Border)
	}
}

// Clear clears the thought chain
func (v *ThoughtChainView) Clear() {
	v.results = nil
	v.nodes = nil
	v.selected = 0
	v.view.SetText("")
}

// Selected returns the currently selected index
func (v *ThoughtChainView) Selected() int {
	return v.selected
}

// SelectNext selects the next item
func (v *ThoughtChainView) SelectNext() {
	max := len(v.results)
	if max == 0 {
		return
	}
	v.selected = (v.selected + 1) % max
	v.render()
}

// SelectPrev selects the previous item
func (v *ThoughtChainView) SelectPrev() {
	max := len(v.results)
	if max == 0 {
		return
	}
	v.selected = (v.selected - 1 + max) % max
	v.render()
}

// Expand expands the selected node
func (v *ThoughtChainView) Expand() {
	if v.selected >= 0 && v.selected < len(v.nodes) {
		if v.nodes[v.selected] != nil {
			v.nodes[v.selected].Expanded = true
			v.render()
		}
	}
}

// Collapse collapses the selected node
func (v *ThoughtChainView) Collapse() {
	if v.selected >= 0 && v.selected < len(v.nodes) {
		if v.nodes[v.selected] != nil {
			v.nodes[v.selected].Expanded = false
			v.render()
		}
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
	
	// Build title with mode indicator
	selectedIndicator := ""
	if v.focused {
		selectedIndicator = fmt.Sprintf(" [[%s]NAV[white]]", v.theme.Accent)
	}
	v.view.SetTitle(fmt.Sprintf("💭 Thought Chain%s", selectedIndicator))
	
	text := ""
	
	for i, r := range v.results {
		// Skip hidden children
		if i > 0 && v.nodes[i-1] != nil && !v.nodes[i-1].Expanded && v.nodes[i].Level > 0 {
			continue
		}
		
		node := v.nodes[i]
		if node == nil {
			node = &ThoughtChainNode{Result: &v.results[i], Level: 0}
		}
		
		prefix := v.getTreePrefix(node.Level, i == v.selected, node.Expanded)
		selected := (i == v.selected && v.focused)
		
		// Format based on match type
		matchType := r.Branch
		icon := "🎯"
		
		switch r.Branch {
		case "exact":
			icon = "💯"
		case "keyword":
			icon = "🔑"
		case "fuzzy":
			icon = "🔍"
		}
		
		if selected {
			text += fmt.Sprintf("[%s]%s[white]\n", v.theme.Selected, prefix)
			text += fmt.Sprintf("   ├ [%s]%s[white]\n", v.theme.Accent, r.Pattern.Trigger)
			text += fmt.Sprintf("   ├ Confidence: [%.0f%%] %s\n", r.Confidence, matchType)
			text += fmt.Sprintf("   ├ Response: %s\n", truncate(r.Pattern.Response, 35))
			text += fmt.Sprintf("   └ Strength: %.1f / %.1f\n", r.Pattern.Strength, r.Pattern.Threshold)
		} else {
			text += fmt.Sprintf("%s[%s]%s[white] %s (%.0f%%)\n", 
				prefix, v.theme.Secondary, icon, r.Pattern.Trigger, r.Confidence)
		}
	}
	
	v.view.SetText(text)
}

// getTreePrefix generates tree visualization prefix
func (v *ThoughtChainView) getTreePrefix(level int, selected, expanded bool) string {
	prefix := ""
	
	// Branch indicators
	if level > 0 {
		indent := ""
		for i := 0; i < level; i++ {
			indent += "   "
		}
		
		expandIcon := "▶"
		if expanded {
			expandIcon = "▼"
		}
		
		prefix = fmt.Sprintf("%s[%s]%s[white] ", indent, v.theme.Secondary, expandIcon)
	} else {
		// Top level
		if selected {
			prefix = "▶ "
		} else {
			prefix = "  "
		}
	}
	
	return prefix
}
