package ui

import (
	"github.com/gdamore/tcell/v2"
)

// Theme defines the color scheme for the TUI
type Theme struct {
	Name        string
	Background  tcell.Color
	Foreground  tcell.Color
	Primary     tcell.Color
	Secondary   tcell.Color
	Accent      tcell.Color
	Border      tcell.Color
	Selected    tcell.Color
	Text        tcell.Color
	TextDim     tcell.Color
	Success     tcell.Color
	Warning     tcell.Color
	Error       tcell.Color
}

// DefaultTheme returns the default dark theme
func DefaultTheme() *Theme {
	return &Theme{
		Name:      "dark",
		Background: tcell.ColorBlack,
		Foreground: tcell.ColorWhite,
		Primary:    tcell.ColorBlue,
		Secondary:  tcell.ColorGray,
		Accent:     tcell.ColorGreen,
		Border:     tcell.ColorDarkGray,
		Selected:   tcell.ColorDarkBlue,
		Text:       tcell.ColorWhite,
		TextDim:    tcell.ColorGray,
		Success:    tcell.ColorGreen,
		Warning:    tcell.ColorYellow,
		Error:      tcell.ColorRed,
	}
}

// LightTheme returns a light color scheme
func LightTheme() *Theme {
	return &Theme{
		Name:      "light",
		Background: tcell.ColorWhite,
		Foreground: tcell.ColorBlack,
		Primary:    tcell.ColorBlue,
		Secondary:  tcell.ColorGray,
		Accent:     tcell.ColorGreen,
		Border:     tcell.ColorLightGray,
		Selected:   tcell.ColorLightBlue,
		Text:       tcell.ColorBlack,
		TextDim:    tcell.ColorGray,
		Success:    tcell.ColorGreen,
		Warning:    tcell.ColorYellow,
		Error:      tcell.ColorRed,
	}
}

// GetTheme returns a theme by name
func GetTheme(name string) *Theme {
	switch name {
	case "light":
		return LightTheme()
	case "dark":
		fallthrough
	default:
		return DefaultTheme()
	}
}

// AvailableThemes returns all available theme names
func AvailableThemes() []string {
	return []string{"dark", "light"}
}

// ApplyTheme applies the theme to a primitive
func (t *Theme) ApplyTheme(p interface {
	SetBackgroundColor(tcell.Color)
	SetBorderColor(tcell.Color)
}) {
	p.SetBackgroundColor(t.Background)
	p.SetBorderColor(t.Border)
}

// ThemeManager manages theme switching
type ThemeManager struct {
	current  *Theme
	themes   map[string]*Theme
}

// NewThemeManager creates a new theme manager
func NewThemeManager() *ThemeManager {
	return &ThemeManager{
		current: DefaultTheme(),
		themes: map[string]*Theme{
			"dark":  DefaultTheme(),
			"light": LightTheme(),
		},
	}
}

// Current returns the current theme
func (tm *ThemeManager) Current() *Theme {
	return tm.current
}

// SetTheme switches to a theme by name
func (tm *ThemeManager) SetTheme(name string) bool {
	if theme, ok := tm.themes[name]; ok {
		tm.current = theme
		return true
	}
	return false
}

// Toggle switches between light and dark themes
func (tm *ThemeManager) Toggle() {
	if tm.current.Name == "dark" {
		tm.current = tm.themes["light"]
	} else {
		tm.current = tm.themes["dark"]
	}
}

// Names returns all available theme names
func (tm *ThemeManager) Names() []string {
	names := make([]string, 0, len(tm.themes))
	for name := range tm.themes {
		names = append(names, name)
	}
	return names
}
