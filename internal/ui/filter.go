package ui

import (
	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FilterPanel provides search and filter functionality for patterns
type FilterPanel struct {
	view          *tview.Frame
	searchInput   *tview.InputField
	filterSelect  *tview.DropDown
	resultsView   *tview.TextView
	headerFlex    *tview.Flex
	contentFlex   *tview.Flex
	visible       bool
	app           *tview.Application
	
	// Callbacks
	onFilter func(query string, filterType FilterType) []*models.Pattern
	
	// State
	currentQuery    string
	currentFilter   FilterType
	filteredResults []*models.Pattern
}

// FilterType represents the type of filter to apply
type FilterType int

const (
	FilterAll FilterType = iota
	FilterByTrigger
	FilterByResponse
	FilterByStrength
	FilterRecent
	FilterFavorites
)

// NewFilterPanel creates a new filter panel
func NewFilterPanel() *FilterPanel {
	f := &FilterPanel{
		currentFilter: FilterAll,
		visible:       false,
	}
	f.buildView()
	return f
}

func (f *FilterPanel) buildView() {
	// Search input
	f.searchInput = tview.NewInputField()
	f.searchInput.SetLabel("🔍 搜索: ")
	f.searchInput.SetPlaceholder("输入 trigger 或 response 关键词...")
	f.searchInput.SetFieldBackgroundColor(tcell.ColorBlack)
	f.searchInput.SetFieldTextColor(tcell.ColorWhite)
	f.searchInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			f.applyFilter()
		}
	})
	
	// Filter type dropdown
	f.filterSelect = tview.NewDropDown()
	f.filterSelect.SetLabel("类型: ")
	f.filterSelect.SetOptions([]string{
		"全部",
		"按 Trigger",
		"按 Response",
		"按 Strength",
		"最近使用",
		"收藏",
	}, nil)
	f.filterSelect.SetCurrentOption(0)
	f.filterSelect.SetSelectedFunc(func(text string, index int) {
		f.currentFilter = FilterType(index)
		f.applyFilter()
	})
	
	// Results view
	f.resultsView = tview.NewTextView()
	f.resultsView.SetDynamicColors(true)
	f.resultsView.SetBorder(true)
	f.resultsView.SetTitle("过滤结果")
	
	// Header layout
	f.headerFlex = tview.NewFlex()
	f.headerFlex.SetDirection(tview.FlexColumn)
	f.headerFlex.AddItem(f.searchInput, 0, 1, true)
	f.headerFlex.AddItem(f.filterSelect, 180, 0, false)
	
	// Content layout
	f.contentFlex = tview.NewFlex().SetDirection(tview.FlexRow)
	f.contentFlex.AddItem(f.headerFlex, 3, 0, true)
	f.contentFlex.AddItem(f.resultsView, 0, 1, false)
	
	// Set background colors
	f.headerFlex.SetBackgroundColor(tcell.ColorBlack)
	f.contentFlex.SetBackgroundColor(tcell.ColorBlack)
	f.searchInput.SetBackgroundColor(tcell.ColorBlack)
	f.filterSelect.SetBackgroundColor(tcell.ColorBlack)
	f.resultsView.SetBackgroundColor(tcell.ColorBlack)
	
	// Create frame with border
	f.view = tview.NewFrame(f.contentFlex)
	f.view.SetBorder(true)
	f.view.SetTitle("🔍 搜索过滤")
	f.view.SetBackgroundColor(tcell.ColorBlack)
}

// SetFilterCallback sets the callback for filtering patterns
func (f *FilterPanel) SetFilterCallback(callback func(query string, filterType FilterType) []*models.Pattern) {
	f.onFilter = callback
}

func (f *FilterPanel) applyFilter() {
	f.currentQuery = f.searchInput.GetText()
	
	if f.onFilter != nil {
		f.filteredResults = f.onFilter(f.currentQuery, f.currentFilter)
		f.updateResultsView()
	}
}

func (f *FilterPanel) updateResultsView() {
	f.resultsView.Clear()
	
	if len(f.filteredResults) == 0 {
		f.resultsView.SetText("[grey]没有找到匹配的模式[/grey]")
		return
	}
	
	for i, p := range f.filteredResults {
		// Color by strength
		strengthColor := f.getStrengthColor(p.Strength)
		f.resultsView.Write([]byte(
			strengthColor + p.Trigger + "[white] → [grey]" + truncate(p.Response, 50) + "[white]\n",
		))
		
		// Add separator (except for last item)
		if i < len(f.filteredResults)-1 {
			f.resultsView.Write([]byte("[darkgrey]────────────────────────────────────[/darkgrey]\n"))
		}
	}
}

func (f *FilterPanel) getStrengthColor(strength float64) string {
	switch {
	case strength >= 80:
		return "[green]"
	case strength >= 50:
		return "[yellow]"
	case strength >= 20:
		return "[orange]"
	default:
		return "[red]"
	}
}

// GetView returns the underlying tview primitive
func (f *FilterPanel) GetView() tview.Primitive {
	return f.view
}

// GetFilteredResults returns the current filtered results
func (f *FilterPanel) GetFilteredResults() []*models.Pattern {
	return f.filteredResults
}

// IsVisible returns whether the panel is visible
func (f *FilterPanel) IsVisible() bool {
	return f.visible
}

// SetVisible sets the visibility of the panel
func (f *FilterPanel) SetVisible(visible bool) {
	f.visible = visible
}

// Clear resets the filter
func (f *FilterPanel) Clear() {
	f.searchInput.SetText("")
	f.filterSelect.SetCurrentOption(0)
	f.currentQuery = ""
	f.currentFilter = FilterAll
	f.filteredResults = nil
	f.resultsView.Clear()
}

// Focus focuses the search input
func (f *FilterPanel) Focus() {
	if f.app != nil {
		f.app.SetFocus(f.searchInput)
	}
}
