package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorWithContext(t *testing.T) {
	err := New(ErrTypeValidation, "invalid input")
	
	if err.Type != ErrTypeValidation {
		t.Errorf("Expected validation type, got %v", err.Type)
	}
	if err.Message != "invalid input" {
		t.Errorf("Expected 'invalid input', got %s", err.Message)
	}
}

func TestErrorWithContext_Error(t *testing.T) {
	original := errors.New("original error")
	err := Wrap(original, ErrTypeDatabase, "operation failed")
	
	if !strings.Contains(err.Error(), "operation failed") {
		t.Error("Expected error message to contain 'operation failed'")
	}
	if !strings.Contains(err.Error(), "original error") {
		t.Error("Expected error message to contain 'original error'")
	}
}

func TestErrorWithContext_Unwrap(t *testing.T) {
	original := errors.New("original")
	err := Wrap(original, ErrTypeNotFound, "not found")
	
	unwrapped := err.Unwrap()
	if unwrapped == nil {
		t.Error("Expected original error")
	}
	if unwrapped.Error() != "original" {
		t.Errorf("Expected 'original', got %s", unwrapped.Error())
	}
}

func TestWithSuggestion(t *testing.T) {
	err := New(ErrTypeValidation, "error").
		WithSuggestion("Check your input")
	
	if len(err.Suggestions) != 1 {
		t.Errorf("Expected 1 suggestion, got %d", len(err.Suggestions))
	}
	if err.Suggestions[0] != "Check your input" {
		t.Errorf("Unexpected suggestion: %s", err.Suggestions[0])
	}
}

func TestWithSuggestions(t *testing.T) {
	err := New(ErrTypeValidation, "error").
		WithSuggestions([]string{"tip1", "tip2"})
	
	if len(err.Suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(err.Suggestions))
	}
}

func TestFormatForDisplay(t *testing.T) {
	err := New(ErrTypeValidation, "Invalid username").
		WithSuggestion("Use 3-20 characters")
	
	display := err.FormatForDisplay()
	
	if !strings.Contains(display, "❌") {
		t.Error("Expected error emoji")
	}
	if !strings.Contains(display, "Invalid username") {
		t.Error("Expected error message")
	}
	if !strings.Contains(display, "💡") {
		t.Error("Expected suggestion emoji")
	}
	if !strings.Contains(display, "Use 3-20 characters") {
		t.Error("Expected suggestion text")
	}
}

func TestValidationError(t *testing.T) {
	err := ValidationError("Email is required", "Enter a valid email")
	
	if err.Type != ErrTypeValidation {
		t.Errorf("Expected validation type, got %v", err.Type)
	}
	if len(err.Suggestions) != 1 {
		t.Errorf("Expected 1 suggestion, got %d", len(err.Suggestions))
	}
}

func TestNotFoundError(t *testing.T) {
	err := NotFoundError("pattern", "abc123")
	
	if err.Type != ErrTypeNotFound {
		t.Errorf("Expected not_found type, got %v", err.Type)
	}
	if !strings.Contains(err.Message, "pattern") {
		t.Error("Expected 'pattern' in message")
	}
	if !strings.Contains(err.Message, "abc123") {
		t.Error("Expected 'abc123' in message")
	}
}

func TestDatabaseError(t *testing.T) {
	original := errors.New("connection refused")
	err := DatabaseError(original)
	
	if err.Type != ErrTypeDatabase {
		t.Errorf("Expected database type, got %v", err.Type)
	}
	if err.Original != original {
		t.Error("Expected original error")
	}
	if len(err.Suggestions) == 0 {
		t.Error("Expected suggestion")
	}
}
