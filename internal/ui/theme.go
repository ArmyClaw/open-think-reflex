package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Theme defines the color scheme for the TUI
type Theme struct {
	Background    tcell.Color
	Foreground    tcell.Color
	Primary       tcell.Color
	Secondary     tcell.Color
	Accent        tcell.Color
	Border        tcell.Color
	Selected      tcell.Color
	Text          tcell.Color
	TextDim       tcell.Color
	Success       tcell.Color
	Warning       tcell.Color
	Error         tcell.Color
}

// DefaultTheme returns the default dark theme
func DefaultTheme() *Theme {
	return &Theme{
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

// ApplyTheme applies the theme to a primitive
func (t *Theme) ApplyTheme(p tview.Primitive) {
	if f, ok := p.(interface{ SetBackgroundColor(tcell.Color) }); ok {
		f.SetBackgroundColor(t.Background)
	}
	if f, ok := p.(interface{ SetBorderColor(tcell.Color) }); ok {
		f.SetBorderColor(t.Border)
	}
}
