package errors

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	e := New(Code("TEST_ERR_001"), "test error")
	if e.Code != "TEST_ERR_001" {
		t.Errorf("expected code TEST_ERR_001, got %s", e.Code)
	}
	if e.Message != "test error" {
		t.Errorf("expected message 'test error', got %s", e.Message)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected timestamp to be set")
	}
}

func TestNewf(t *testing.T) {
	e := Newf(Code("TEST_ERR_002"), "test error: %s", "formatted")
	if e.Message != "test error: formatted" {
		t.Errorf("expected 'test error: formatted', got %s", e.Message)
	}
}

func TestError_Error(t *testing.T) {
	e := New(Code("TEST_ERR_003"), "test error")
	if e.Error() != "test error" {
		t.Errorf("expected 'test error', got %s", e.Error())
	}
}

func TestError_ErrorWithWrapped(t *testing.T) {
	original := errors.New("original error")
	e := Wrap(original, Code("TEST_ERR_004"), "wrapped error")
	expected := "wrapped error: original error"
	if e.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, e.Error())
	}
}

func TestError_Unwrap(t *testing.T) {
	original := errors.New("original error")
	e := Wrap(original, Code("TEST_ERR_005"), "wrapped error")
	if e.Unwrap() != original {
		t.Error("expected unwrap to return original error")
	}
}

func TestError_Is(t *testing.T) {
	e := New(Code("TEST_ERR_006"), "test error")
	if !e.Is(Code("TEST_ERR_006")) {
		t.Error("expected Is to return true for matching code")
	}
	if e.Is(Code("TEST_ERR_007")) {
		t.Error("expected Is to return false for non-matching code")
	}
}

func TestError_IsNotFound(t *testing.T) {
	e := New(ErrPatternNotFound.Code, "not found")
	if !e.IsNotFound() {
		t.Error("expected IsNotFound to return true for pattern not found")
	}

	e2 := New(Code("TEST_ERR_008"), "some other error")
	if e2.IsNotFound() {
		t.Error("expected IsNotFound to return false for non-not-found error")
	}
}

func TestError_IsValidation(t *testing.T) {
	e := New(Code("PATTERN_ERR_1001"), "validation error")
	if !e.IsValidation() {
		t.Error("expected IsValidation to return true")
	}

	e2 := New(Code("STORAGE_ERR_2001"), "storage error")
	if e2.IsValidation() {
		t.Error("expected IsValidation to return false for storage error")
	}
}

func TestError_IsDatabase(t *testing.T) {
	e := New(Code("STORAGE_ERR_2001"), "database error")
	if !e.IsDatabase() {
		t.Error("expected IsDatabase to return true")
	}

	e2 := New(Code("AI_ERR_3001"), "AI error")
	if e2.IsDatabase() {
		t.Error("expected IsDatabase to return false for AI error")
	}
}

func TestError_IsAI(t *testing.T) {
	e := New(Code("AI_ERR_3001"), "AI error")
	if !e.IsAI() {
		t.Error("expected IsAI to return true")
	}

	e2 := New(Code("STORAGE_ERR_2001"), "storage error")
	if e2.IsAI() {
		t.Error("expected IsAI to return false for storage error")
	}
}

func TestError_IsRetryable(t *testing.T) {
	testCases := []struct {
		code    Code
		retryable bool
	}{
		{ErrDatabaseLocked.Code, true},
		{ErrAIProviderUnavailable.Code, true},
		{ErrAITimeout.Code, true},
		{ErrAIRateLimited.Code, true},
		{ErrPatternNotFound.Code, false},
		{ErrEmptyTrigger.Code, false},
	}

	for _, tc := range testCases {
		e := New(tc.code, "test")
		if e.IsRetryable() != tc.retryable {
			t.Errorf("expected IsRetryable=%v for %s, got %v", tc.retryable, tc.code, e.IsRetryable())
		}
	}
}

func TestAsError(t *testing.T) {
	// Test with *Error
	original := New(Code("TEST_ERR_009"), "test error")
	var e *Error
	if !errors.As(original, &e) {
		t.Error("expected errors.As to find *Error")
	}

	// Test with plain error
	plain := errors.New("plain error")
	_, ok := AsError(plain)
	if ok {
		t.Error("expected AsError to return false for plain error")
	}
}

func TestAsErrorf(t *testing.T) {
	// Test wrapping existing *Error - should preserve original
	original := New(Code("TEST_ERR_010"), "original")
	wrapped := AsErrorf(original, Code("NEW_CODE"), "new message")
	// When input is already *Error, it's returned as-is with original code preserved
	if wrapped.Code != "TEST_ERR_010" {
		t.Errorf("expected to preserve original code, got %s", wrapped.Code)
	}

	// Test with plain error - should wrap with new code
	plain := errors.New("plain error")
	wrapped2 := AsErrorf(plain, Code("NEW_CODE_2"), "new message %s", "formatted")
	if wrapped2.Code != "NEW_CODE_2" {
		t.Errorf("expected new code, got %s", wrapped2.Code)
	}
	if wrapped2.Wrapped == nil {
		t.Error("expected wrapped error to be set")
	}

	// Test with nil
	var nilErr error
	wrapped3 := AsErrorf(nilErr, Code("NEW_CODE_3"), "message")
	if wrapped3 != nil {
		t.Error("expected nil for nil error")
	}
}

func TestWrapIf(t *testing.T) {
	// Test with non-nil error
	err := errors.New("test")
	wrapped := WrapIf(err, Code("TEST_ERR_011"), "wrapped")
	if wrapped == nil {
		t.Error("expected non-nil wrapped error")
	}

	// Test with nil error
	var nilErr error
	wrappedNil := WrapIf(nilErr, Code("TEST_ERR_012"), "wrapped")
	if wrappedNil != nil {
		t.Error("expected nil for nil error")
	}
}

func TestError_WithDetails(t *testing.T) {
	e := New(Code("TEST_ERR_013"), "test error")
	details := Details{
		Field:      "trigger",
		Value:      "",
		Constraint: "cannot be empty",
	}
	e2 := e.WithDetails(details)
	if e2.Details.Field != "trigger" {
		t.Errorf("expected Field 'trigger', got %s", e2.Details.Field)
	}
}

func TestError_WithField(t *testing.T) {
	e := New(Code("TEST_ERR_014"), "test error")
	e2 := e.WithField("response", "test value")
	if e2.Details.Field != "response" {
		t.Errorf("expected Field 'response', got %s", e2.Details.Field)
	}
	if e2.Details.Value != "test value" {
		t.Errorf("expected Value 'test value', got %v", e2.Details.Value)
	}
}

func TestError_WithConstraint(t *testing.T) {
	e := New(Code("TEST_ERR_015"), "test error")
	e2 := e.WithConstraint("must be positive")
	if e2.Details.Constraint != "must be positive" {
		t.Errorf("expected Constraint 'must be positive', got %s", e2.Details.Constraint)
	}
}

func TestError_WithRequestID(t *testing.T) {
	e := New(Code("TEST_ERR_016"), "test error")
	e2 := e.WithRequestID("req-123")
	if e2.RequestID != "req-123" {
		t.Errorf("expected RequestID 'req-123', got %s", e2.RequestID)
	}
}
