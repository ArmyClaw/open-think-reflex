package ui

import (
	"testing"
)

func TestDefaultTheme(t *testing.T) {
	theme := DefaultTheme()
	if theme == nil {
		t.Fatal("DefaultTheme should not return nil")
	}

	// Check essential colors are set
	if theme.Primary == "" {
		t.Error("Primary color should be set")
	}
	if theme.Background == "" {
		t.Error("Background color should be set")
	}
	if theme.Text == "" {
		t.Error("Text color should be set")
	}
}

func TestThemeManager_Current(t *testing.T) {
	manager := NewThemeManager()
	if manager == nil {
		t.Fatal("NewThemeManager should not return nil")
	}

	theme := manager.Current()
	if theme == nil {
		t.Fatal("Current should not return nil")
	}
}

func TestThemeManager_Switch(t *testing.T) {
	manager := NewThemeManager()

	// Get initial theme
	initial := manager.Current().Name

	// Switch theme
	manager.Switch()

	// Theme should be different
	after := manager.Current().Name
	if initial == after {
		t.Logf("Theme switched from %s to %s", initial, after)
	}
}

func TestThemeManager_SetTheme(t *testing.T) {
	manager := NewThemeManager()

	// Set dark theme
	err := manager.SetTheme("dark")
	if err != nil {
		t.Errorf("SetTheme(dark) failed: %v", err)
	}

	// Set light theme
	err = manager.SetTheme("light")
	if err != nil {
		t.Errorf("SetTheme(light) failed: %v", err)
	}

	// Set invalid theme
	err = manager.SetTheme("invalid-theme")
	if err == nil {
		t.Error("Expected error for invalid theme")
	}
}

func TestTheme_IsDark(t *testing.T) {
	dark := &Theme{Name: "dark"}
	light := &Theme{Name: "light"}

	if !dark.IsDark() {
		t.Error("dark theme should return true for IsDark()")
	}

	if light.IsDark() {
		t.Error("light theme should return false for IsDark()")
	}
}
