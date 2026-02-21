package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/rivo/tview"
)

// OutputView displays the AI generated content (Layer 2)
type OutputView struct {
	view        *tview.TextView
	theme       *Theme
	title       string
	streaming   bool
	loading     bool
	stopLoading chan struct{}
}

// Spinner frames for loading animation
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// NewOutputView creates a new output view
func NewOutputView(theme *Theme) *OutputView {
	v := &OutputView{
		theme: theme,
		title: "Output",
	}
	
	v.view = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	v.view.SetBorder(true).SetTitle("📤 " + v.title)
	
	v.view.SetBackgroundColor(theme.Background)
	v.view.SetBorderColor(theme.Border)
	v.view.SetTitleColor(theme.Primary)
	
	v.view.SetText(fmt.Sprintf("[%s]Welcome to Open-Think-Reflex[white]\n\nYour AI-powered input accelerator.\n\n[Features]:\n  • Pattern matching with confidence scoring\n  • Three-layer terminal UI\n  • Keyboard navigation (Tab to switch modes)\n  • Vim-style shortcuts (h/j/k/l)\n\n[Quick Start]:\n  1. Type a query in the input area\n  2. Press Enter to search patterns\n  3. Use ↑/↓ to navigate results\n  4. Press Enter to select a response\n\nPress [?%s] for help[white]",
		v.theme.Secondary, v.theme.Accent))
	
	return v
}

// SetOutput sets the output content
func (v *OutputView) SetOutput(text string) {
	v.view.SetText(text)
}

// SetFormattedOutput sets the output with formatting
func (v *OutputView) SetFormattedOutput(title, content string) {
	v.title = title
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
	v.title = title
	v.view.SetTitle(fmt.Sprintf("📤 %s", title))
}

// SetStatus sets a status message
func (v *OutputView) SetStatus(status string, success bool) {
	color := v.theme.Success
	if !success {
		color = v.theme.Error
	}
	v.AppendOutput(fmt.Sprintf("[%s]%s[white]", color, status))
}

// StartStreaming begins streaming animation
func (v *OutputView) StartStreaming() {
	v.streaming = true
}

// StopStreaming stops streaming animation
func (v *OutputView) StopStreaming() {
	v.streaming = false
}

// IsStreaming returns streaming state
func (v *OutputView) IsStreaming() bool {
	return v.streaming
}

// StartLoading begins loading animation with spinner
func (v *OutputView) StartLoading(message string) {
	if v.loading {
		return
	}
	v.loading = true
	v.stopLoading = make(chan struct{})
	
	go func() {
		frameIdx := 0
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		
		for {
			select {
			case <-v.stopLoading:
				return
			case <-ticker.C:
				frame := spinnerFrames[frameIdx%len(spinnerFrames)]
				frameIdx++
				v.view.SetText(fmt.Sprintf("[%s]%s %s[white]\n\nWaiting for AI response...", v.theme.Primary, frame, message))
			}
		}
	}()
}

// StopLoading stops loading animation
func (v *OutputView) StopLoading() {
	if !v.loading {
		return
	}
	v.loading = false
	if v.stopLoading != nil {
		close(v.stopLoading)
		v.stopLoading = nil
	}
}

// IsLoading returns loading state
func (v *OutputView) IsLoading() bool {
	return v.loading
}

// ShowTypingEffect displays text with typing effect (caller would handle timing)
func (v *OutputView) ShowTypingEffect(text string, delayMs int) {
	// For now, just show the full text
	// A more advanced implementation would use goroutines
	v.SetOutput(text)
}

// FormatResponse formats a pattern response for display
func (v *OutputView) FormatResponse(trigger, response string, confidence float64, branch string) string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("[%s]Trigger:[white] %s\n\n", v.theme.Primary, trigger))
	sb.WriteString(fmt.Sprintf("[%s]Response:[white]\n%s\n\n", v.theme.Primary, response))
	sb.WriteString(fmt.Sprintf("[%s]Confidence:[white] %.0f%% (%s)\n", v.theme.Secondary, confidence, branch))
	sb.WriteString(fmt.Sprintf("[%s]%s[white]\n", v.theme.TextDim, strings.Repeat("─", 40)))
	
	return sb.String()
}

// ShowMatchList shows a list of matches
func (v *OutputView) ShowMatchList(results []struct {
	Trigger    string
	Confidence float64
	Branch     string
	Response   string
}) {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("[%s]Found %d match(es):[white]\n\n", v.theme.Primary, len(results)))
	
	for i, r := range results {
		icon := "🎯"
		switch r.Branch {
		case "exact":
			icon = "💯"
		case "keyword":
			icon = "🔑"
		case "fuzzy":
			icon = "🔍"
		}
		
		sb.WriteString(fmt.Sprintf("[%d]. %s %s\n", i+1, icon, r.Trigger))
		sb.WriteString(fmt.Sprintf("   Confidence: %.0f%% (%s)\n", r.Confidence, r.Branch))
		sb.WriteString(fmt.Sprintf("   Response: %s\n\n", truncate(r.Response, 50)))
	}
	
	sb.WriteString(fmt.Sprintf("[%s]Use [↑/↓] to navigate, [Enter] to select[white]", v.theme.TextDim))
	
	v.SetOutput(sb.String())
}

// ShowHelp shows help information
func (v *OutputView) ShowHelp() {
	help := `
[? Help]
════════════════════════════════════════════════════════════

[Input Mode]:
  Type query + Enter    Search patterns
  Tab                   Switch to Navigation mode

[Navigation Mode]:
  ↑/↓ or j/k           Navigate results
  ←/→ or h/l           Collapse/Expand
  Enter                Use selected response
  Tab or Esc           Return to Input mode

[General]:
  ?                    Show this help
  q                    Quit

[Tips]:
  • Patterns below threshold are hidden
  • Higher confidence = better match
  • Use 'otr pattern create' to add patterns
`
	v.SetFormattedOutput("Help", help)
}

// ShowWelcome shows welcome message
func (v *OutputView) ShowWelcome() {
	v.SetFormattedOutput("Welcome", fmt.Sprintf(`[%s]Welcome to Open-Think-Reflex[white]

Your AI-powered input accelerator.

[Quick Start]:
  1. Type a query in the input area
  2. Press Enter to search patterns
  3. Use ↑/↓ to navigate results
  4. Press Enter to select

[Commands]:
  otr pattern create --trigger "x" --response "y"
  otr interactive (or otr tui)

Press [?%s] for help[white]`, v.theme.Secondary, v.theme.Accent))
}

// GetView returns the underlying TextView
func (v *OutputView) GetView() *tview.TextView {
	return v.view
}

// TimeFormatter provides time-based formatting utilities
type TimeFormatter struct{}

func (tf *TimeFormatter) FormatDuration(d time.Duration) string {
	if d < time.Second {
		return "< 1s"
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	return fmt.Sprintf("%.1fm", d.Minutes())
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
