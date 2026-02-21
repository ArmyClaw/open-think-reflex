package ui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestDefaultTheme(t *testing.T) {
	theme := DefaultTheme()
	if theme == nil {
		t.Fatal("DefaultTheme should not return nil")
	}

	// Check essential colors are set (not black/white defaults)
	if theme.Primary == tcell.ColorBlack {
		t.Error("Primary color should not be black")
	}
	if theme.Background == tcell.ColorWhite {
		t.Error("Background color should not be white (for dark theme)")
	}
	if theme.Text == tcell.ColorBlack {
		t.Error("Text color should not be black (for dark theme)")
	}
}

func TestLightTheme(t *testing.T) {
	theme := LightTheme()
	if theme == nil {
		t.Fatal("LightTheme should not return nil")
	}

	if theme.Name != "light" {
		t.Errorf("Expected name 'light', got '%s'", theme.Name)
	}
}

func TestGetTheme(t *testing.T) {
	// Test dark theme
	dark := GetTheme("dark")
	if dark == nil || dark.Name != "dark" {
		t.Error("GetTheme(dark) should return dark theme")
	}

	// Test light theme
	light := GetTheme("light")
	if light == nil || light.Name != "light" {
		t.Error("GetTheme(light) should return light theme")
	}

	// Test unknown theme (should return dark)
	unknown := GetTheme("unknown")
	if unknown == nil || unknown.Name != "dark" {
		t.Error("GetTheme(unknown) should return dark theme as default")
	}
}

func TestAvailableThemes(t *testing.T) {
	themes := AvailableThemes()
	if len(themes) != 2 {
		t.Errorf("Expected 2 themes, got %d", len(themes))
	}

	// Check both themes are present
	foundDark := false
	foundLight := false
	for _, t := range themes {
		if t == "dark" {
			foundDark = true
		}
		if t == "light" {
			foundLight = true
		}
	}

	if !foundDark {
		t.Error("dark theme should be available")
	}
	if !foundLight {
		t.Error("light theme should be available")
	}
}
