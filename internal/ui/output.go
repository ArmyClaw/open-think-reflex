package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

// OutputView displays the AI generated content (Layer 2)
type OutputView struct {
	view  *tview.TextView
	theme *Theme
}

// NewOutputView creates a new output view
func NewOutputView(theme *Theme) *OutputView {
	v := &OutputView{
		theme: theme,
	}
	
	v.view = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	v.view.SetBorder(true).SetTitle("📤 Output")
	
	v.view.SetBackgroundColor(theme.Background)
	v.view.SetBorderColor(theme.Border)
	v.view.SetTitleColor(theme.Primary)
	
	v.view.SetText(fmt.Sprintf("[%s]Welcome to Open-Think-Reflex[white]\n\nYour AI-powered input accelerator.\n\nStart typing in the input area below to see pattern matches and AI-generated responses.\n\nPress [?%s] for help[white]",
		theme.Secondary, theme.Accent))
	
	return v
}

// SetOutput sets the output content
func (v *OutputView) SetOutput(text string) {
	v.view.SetText(text)
}

// SetFormattedOutput sets the output with formatting
func (v *OutputView) SetFormattedOutput(title, content string) {
	v.view.SetTitle(fmt.Sprintf("📤 %s", title))
	v.view.SetText(content)
}

// AppendOutput appends text to the output
func (v *OutputView) AppendOutput(text string) {
	current := v.view.GetText(false)
	v.view.SetText(current + "\n" + text)
}

// Clear clears the output
func (v *OutputView) Clear() {
	v.view.SetText("")
}

// SetTitle sets the panel title
func (v *OutputView) SetTitle(title string) {
	v.view.SetTitle(title)
}

// SetStatus sets a status message
func (v *OutputView) SetStatus(status string, success bool) {
	color := v.theme.Success
	if !success {
		color = v.theme.Error
	}
	v.AppendOutput(fmt.Sprintf("[%s]%s[white]", color, status))
}

// GetView returns the underlying TextView
func (v *OutputView) GetView() *tview.TextView {
	return v.view
}
