package errors

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	e := New(ErrTypeValidation, "test error")
	if e.Type != ErrTypeValidation {
		t.Errorf("expected type validation, got %s", e.Type)
	}
	if e.Message != "test error" {
		t.Errorf("expected message 'test error', got %s", e.Message)
	}
}

func TestError_Error(t *testing.T) {
	e := New(ErrTypeValidation, "test error")
	if e.Error() != "validation: test error" {
		t.Errorf("expected 'validation: test error', got '%s'", e.Error())
	}
}

func TestError_WithWrapped(t *testing.T) {
	original := errors.New("original error")
	e := Wrap(original, ErrTypeNotFound, "wrapped error")
	expected := "not_found: wrapped error (original error)"
	if e.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, e.Error())
	}
}

func TestError_Unwrap(t *testing.T) {
	original := errors.New("original error")
	e := Wrap(original, ErrTypeValidation, "wrapped error")
	if e.Unwrap() != original {
		t.Error("expected unwrap to return original error")
	}
}

func TestValidationError(t *testing.T) {
	e := ValidationError("invalid input", "check the format")
	if e.Type != ErrTypeValidation {
		t.Error("expected validation error type")
	}
	if len(e.Suggestions) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(e.Suggestions))
	}
}

func TestNotFoundError(t *testing.T) {
	e := NotFoundError("pattern", "123")
	if e.Type != ErrTypeNotFound {
		t.Error("expected not found error type")
	}
	if len(e.Suggestions) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(e.Suggestions))
	}
}

func TestDatabaseError(t *testing.T) {
	original := errors.New("connection refused")
	e := DatabaseError(original)
	if e.Type != ErrTypeDatabase {
		t.Error("expected database error type")
	}
	if e.Original != original {
		t.Error("expected original error to be preserved")
	}
}

func TestWithSuggestion(t *testing.T) {
	e := New(ErrTypeValidation, "test error")
	e2 := e.WithSuggestion("try again")
	if len(e2.Suggestions) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(e2.Suggestions))
	}
}

func TestWithSuggestions(t *testing.T) {
	e := New(ErrTypeValidation, "test error")
	e2 := e.WithSuggestions([]string{"suggestion 1", "suggestion 2"})
	if len(e2.Suggestions) != 2 {
		t.Errorf("expected 2 suggestions, got %d", len(e2.Suggestions))
	}
}

func TestFormatForDisplay(t *testing.T) {
	e := New(ErrTypeValidation, "test error")
	e = e.WithSuggestion("try again")
	output := e.FormatForDisplay()
	
	if output == "" {
		t.Error("expected non-empty output")
	}
	
	// Check for expected content
	if len(output) < len("❌ test error") {
		t.Error("output seems too short")
	}
}

func TestWrapIf_NonNil(t *testing.T) {
	err := errors.New("test")
	wrapped := WrapIf(err, ErrTypeValidation, "wrapped")
	if wrapped == nil {
		t.Error("expected non-nil wrapped error")
	}
	if wrapped.Message != "wrapped" {
		t.Errorf("expected message 'wrapped', got '%s'", wrapped.Message)
	}
}

func TestWrapIf_Nil(t *testing.T) {
	var nilErr error
	wrapped := WrapIf(nilErr, ErrTypeValidation, "wrapped")
	if wrapped != nil {
		t.Error("expected nil for nil error")
	}
}
