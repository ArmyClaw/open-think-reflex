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
	view     *tview.InputField
	theme    *Theme
	callback InputCallback
	prompt   string
}

// NewInputView creates a new input view
func NewInputView(theme *Theme, callback InputCallback) *InputView {
	v := &InputView{
		theme:    theme,
		callback: callback,
		prompt:   ">",
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
					v.callback(text)
					v.view.SetText("")
				}
			}
		})
	
	return v
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
