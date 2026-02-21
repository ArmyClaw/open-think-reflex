package errors

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Code represents an error code
type Code string

// Error represents an application error
type Error struct {
	Code      Code     `json:"code"`
	Message   string   `json:"message"`
	Details   Details  `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string   `json:"request_id,omitempty"`
	Wrapped   error    `json:"-"`
}

// Details contains additional error information
type Details struct {
	Field      string      `json:"field,omitempty"`
	Value      interface{} `json:"value,omitempty"`
	Constraint string      `json:"constraint,omitempty"`
	Provider   string      `json:"provider,omitempty"`
	RetryAfter int         `json:"retry_after,omitempty"`
}

// Error returns the error message
func (e *Error) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Wrapped)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *Error) Unwrap() error {
	return e.Wrapped
}

// WithRequestID adds request ID to error
func (e *Error) WithRequestID(id string) *Error {
	e.RequestID = id
	return e
}

// WithDetails adds details to error
func (e *Error) WithDetails(d Details) *Error {
	e.Details = d
	return e
}

// WithField adds field details to error
func (e *Error) WithField(field string, value interface{}) *Error {
	e.Details.Field = field
	e.Details.Value = value
	return e
}

// WithConstraint adds constraint details to error
func (e *Error) WithConstraint(constraint string) *Error {
	e.Details.Constraint = constraint
	return e
}

// New creates a new error with code and message
func New(code Code, message string) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// Newf creates a new formatted error
func Newf(code Code, format string, args ...interface{}) *Error {
	return New(code, fmt.Sprintf(format, args...))
}

// Wrap wraps an existing error
func Wrap(err error, code Code, message string) *Error {
	if err == nil {
		return nil
	}
	e, ok := err.(*Error)
	if !ok {
		e = &Error{
			Code:      code,
			Message:   message,
			Wrapped:   err,
			Timestamp: time.Now(),
		}
	}
	return e
}

// WrapIf wraps an error only if it's not nil
func WrapIf(err error, code Code, message string) *Error {
	if err == nil {
		return nil
	}
	return Wrap(err, code, message)
}

// ValidationError creates a validation error
func ValidationError(field, reason string, value interface{}) *Error {
	return &Error{
		Code:    Code(fmt.Sprintf("PATTERN_ERR_%04d", 1)),
		Message: fmt.Sprintf("validation failed for %s: %s", field, reason),
		Details: Details{
			Field:      field,
			Value:      value,
			Constraint: reason,
		},
		Timestamp: time.Now(),
	}
}

// Common error variables
var (
	// Pattern errors (1000-1099)
	ErrPatternNotFound    = New(Code("STORAGE_ERR_2001"), "Pattern not found")
	ErrInvalidPatternID  = New(Code("PATTERN_ERR_1001"), "Invalid pattern ID format")
	ErrEmptyTrigger      = New(Code("PATTERN_ERR_1002"), "Trigger cannot be empty")
	ErrEmptyResponse     = New(Code("PATTERN_ERR_1003"), "Response cannot be empty")
	ErrInvalidStrength   = New(Code("PATTERN_ERR_1004"), "Strength must be between 0 and 100")

	// Storage errors (2000-2099)
	ErrDatabaseError  = New(Code("STORAGE_ERR_2003"), "Database operation failed")
	ErrSpaceNotFound = New(Code("STORAGE_ERR_2002"), "Space not found")
	ErrDatabaseLocked = New(Code("STORAGE_ERR_2009"), "Database locked")

	// AI errors (3000-3099)
	ErrAIProviderUnavailable = New(Code("AI_ERR_3002"), "AI provider unavailable")
	ErrAITimeout            = New(Code("AI_ERR_3006"), "AI request timeout")
	ErrAIRateLimited        = New(Code("AI_ERR_3004"), "Rate limit exceeded")
	ErrAPIKeyMissing        = New(Code("AI_ERR_3001"), "API key not configured")

	// CLI errors (4000-4099)
	ErrUnknownCommand   = New(Code("CLI_ERR_4001"), "Unknown command")
	ErrMissingArgument = New(Code("CLI_ERR_4002"), "Missing required argument")
	ErrInvalidArgument = New(Code("CLI_ERR_4003"), "Invalid argument value")

	// Config errors (5000-5099)
	ErrConfigNotFound   = New(Code("CONFIG_ERR_5001"), "Configuration file not found")
	ErrConfigParseError = New(Code("CONFIG_ERR_5002"), "Configuration parse error")
)

// Is checks if the error matches a given code
func (e *Error) Is(code Code) bool {
	return e.Code == code
}

// IsNotFound checks if the error is a "not found" error
func (e *Error) IsNotFound() bool {
	return e.Code == ErrPatternNotFound.Code || e.Code == ErrSpaceNotFound.Code
}

// IsValidation checks if the error is a validation error
func (e *Error) IsValidation() bool {
	return strings.HasPrefix(string(e.Code), "PATTERN_ERR_1")
}

// IsDatabase checks if the error is a database error
func (e *Error) IsDatabase() bool {
	return strings.HasPrefix(string(e.Code), "STORAGE_ERR_2")
}

// IsAI checks if the error is an AI-related error
func (e *Error) IsAI() bool {
	return strings.HasPrefix(string(e.Code), "AI_ERR_3")
}

// IsRetryable checks if the error is retryable
func (e *Error) IsRetryable() bool {
	retryableCodes := []Code{
		ErrDatabaseLocked.Code,
		ErrAIProviderUnavailable.Code,
		ErrAITimeout.Code,
		ErrAIRateLimited.Code,
	}
	for _, c := range retryableCodes {
		if e.Code == c {
			return true
		}
	}
	return false
}

// AsError tries to cast error to *Error
func AsError(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// AsErrorf wraps error as *Error if possible, otherwise creates new error
func AsErrorf(err error, code Code, format string, args ...interface{}) *Error {
	if err == nil {
		return nil
	}
	e, ok := AsError(err)
	if ok {
		return e
	}
	return Wrap(err, code, fmt.Sprintf(format, args...))
}
