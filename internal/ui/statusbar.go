package ui

import (
	"fmt"
	"sync"
	"time"

	"github.com/rivo/tview"
)

// Status represents the current application status
type Status int

const (
	StatusIdle Status = iota
	StatusMatching
	StatusThinking
	StatusProcessing
	StatusError
)

// StatusBar displays application status information
type StatusBar struct {
	view         *tview.TextView
	theme        *Theme
	status       Status
	statusText   string
	patternCount int
	matchCount   int
	mu           sync.RWMutex
}

// NewStatusBar creates a new status bar
func NewStatusBar(theme *Theme) *StatusBar {
	sb := &StatusBar{
		theme:      theme,
		status:     StatusIdle,
		statusText: "Ready",
	}

	sb.view = tview.NewTextView()
	sb.view.SetDynamicColors(true)
	sb.render()

	return sb
}

// SetStatus updates the current status
func (sb *StatusBar) SetStatus(status Status, text string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.status = status
	sb.statusText = text
	sb.render()
}

// SetPatternCount updates the pattern count
func (sb *StatusBar) SetPatternCount(count int) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.patternCount = count
	sb.render()
}

// SetMatchCount updates the match count
func (sb *StatusBar) SetMatchCount(count int) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.matchCount = count
	sb.render()
}

// SetTheme updates the theme
func (sb *StatusBar) SetTheme(theme *Theme) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.theme = theme
	sb.render()
}

// View returns the underlying tview primitive
func (sb *StatusBar) View() *tview.TextView {
	return sb.view
}

func (sb *StatusBar) render() {
	sb.mu.RLock()
	defer sb.mu.RUnlock()

	// Status icon and color
	statusIcon := "●"
	statusColor := sb.theme.Success // green for idle

	switch sb.status {
	case StatusIdle:
		statusIcon = "●"
		statusColor = sb.theme.Success
	case StatusMatching:
		statusIcon = "◐"
		statusColor = sb.theme.Primary // cyan
	case StatusThinking:
		statusIcon = "◑"
		statusColor = sb.theme.Warning // yellow
	case StatusProcessing:
		statusIcon = "◒"
		statusColor = sb.theme.Warning
	case StatusError:
		statusIcon = "✕"
		statusColor = sb.theme.Error // red
	}

	// Get status text
	statusText := sb.statusText

	// Get current time
	currentTime := time.Now().Format("15:04:05")

	// Build status bar content using tview dynamic colors
	content := fmt.Sprintf("[%s]%s[] %s | Patterns: %d | Matches: %d | %s",
		statusColor,
		statusIcon,
		statusText,
		sb.patternCount,
		sb.matchCount,
		currentTime,
	)

	sb.view.SetText(content)
	sb.view.SetTextColor(sb.theme.Secondary)
	sb.view.SetBackgroundColor(sb.theme.Background)
}

// GetStatus returns the current status
func (sb *StatusBar) GetStatus() Status {
	sb.mu.RLock()
	defer sb.mu.RUnlock()
	return sb.status
}

// IsIdle returns true if the status is idle
func (sb *StatusBar) IsIdle() bool {
	return sb.GetStatus() == StatusIdle
}
