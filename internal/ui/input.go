package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// InputCallback is called when the user submits input
type InputCallback func(text string)

// InputView handles user input (Layer 3)
type InputView struct {
	view      *tview.InputField
	theme     *Theme
	callback  InputCallback
	prompt    string
	history   []string
	historyPos int
	autocomplete func(string) []string
	autocompletePos int
	showAutocomplete bool
}

// NewInputView creates a new input view
func NewInputView(theme *Theme, callback InputCallback) *InputView {
	v := &InputView{
		theme:     theme,
		callback:  callback,
		prompt:    ">",
		history:   []string{},
		historyPos: -1,
	}
	
	v.view = tview.NewInputField().
		SetLabel(fmt.Sprintf("[%s]%s [white]", v.theme.Accent, v.prompt)).
		SetPlaceholder("Type your query here...").
		SetPlaceholderTextColor(theme.TextDim).
		SetFieldTextColor(theme.Text).
		SetFieldBackgroundColor(theme.Background).
		SetLabelColor(theme.Accent).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				text := v.view.GetText()
				if text != "" {
					v.AddToHistory(text)
					v.callback(text)
					v.view.SetText("")
				}
			}
		})
	
	// Capture key events for history navigation
	v.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyUp {
			// Navigate history up
			v.NavigateHistoryUp()
			return nil
		}
		if event.Key() == tcell.KeyDown {
			// Navigate history down
			v.NavigateHistoryDown()
			return nil
		}
		return event
	})
	
	return v
}

// AddToHistory adds text to the input history
func (v *InputView) AddToHistory(text string) {
	// Don't add duplicates at the end
	if len(v.history) > 0 && v.history[len(v.history)-1] == text {
		return
	}
	v.history = append(v.history, text)
	v.historyPos = len(v.history)
}

// NavigateHistoryUp goes back in input history
func (v *InputView) NavigateHistoryUp() {
	if len(v.history) == 0 {
		return
	}
	
	if v.historyPos > 0 {
		v.historyPos--
		v.view.SetText(v.history[v.historyPos])
	}
}

// NavigateHistoryDown goes forward in input history
func (v *InputView) NavigateHistoryDown() {
	if len(v.history) == 0 {
		return
	}
	
	if v.historyPos < len(v.history)-1 {
		v.historyPos++
		v.view.SetText(v.history[v.historyPos])
	} else {
		// Clear input when at the end
		v.historyPos = len(v.history)
		v.view.SetText("")
	}
}

// SetPrompt changes the input prompt
func (v *InputView) SetPrompt(prompt string) {
	v.prompt = prompt
	v.view.SetLabel(fmt.Sprintf("[%s]%s [white]", v.theme.Accent, v.prompt))
}

// SetPlaceholder sets the placeholder text
func (v *InputView) SetPlaceholder(text string) {
	v.view.SetPlaceholder(text)
}

// GetText returns the current input text
func (v *InputView) GetText() string {
	return v.view.GetText()
}

// SetText sets the input text
func (v *InputView) SetText(text string) {
	v.view.SetText(text)
}

// SetCallback sets the submit callback
func (v *InputView) SetCallback(callback InputCallback) {
	v.callback = callback
}

// SetAutocomplete sets the autocomplete function
func (v *InputView) SetAutocomplete(fn func(string) []string) {
	v.autocomplete = fn
}

// GetHistory returns the input history
func (v *InputView) GetHistory() []string {
	return v.history
}

// ClearHistory clears the input history
func (v *InputView) ClearHistory() {
	v.history = []string{}
	v.historyPos = -1
}

// HistoryLen returns the history length
func (v *InputView) HistoryLen() int {
	return len(v.history)
}

// Focus focuses the input view
func (v *InputView) Focus() {
	v.view.SetBorderColor(v.theme.Accent)
}

// Blur removes focus from the input view
func (v *InputView) Blur() {
	v.view.SetBorderColor(v.theme.Border)
}

// GetView returns the underlying InputField
func (v *InputView) GetView() *tview.InputField {
	return v.view
}
