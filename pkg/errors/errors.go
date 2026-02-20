package errors

import (
	"fmt"
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
	Field     string      `json:"field,omitempty"`
	Value     interface{} `json:"value,omitempty"`
	Constraint string    `json:"constraint,omitempty"`
	Provider   string     `json:"provider,omitempty"`
	RetryAfter int        `json:"retry_after,omitempty"`
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

// ValidationError creates a validation error
func ValidationError(field, reason string, value interface{}) *Error {
	return &Error{
		Code:    Code(fmt.Sprintf("PATTERN_ERR_%04d", 1)),
		Message: fmt.Sprintf("validation failed for %s: %s", field, reason),
		Details: Details{
			Field:     field,
			Value:     value,
			Constraint: reason,
		},
		Timestamp: time.Now(),
	}
}

// Common error codes
const (
	// Pattern errors (1000-1099)
	ErrPatternNotFound     = New(Code("STORAGE_ERR_2001"), "Pattern not found")
	ErrInvalidPatternID    = New(Code("PATTERN_ERR_1001"), "Invalid pattern ID format")
	ErrEmptyTrigger        = New(Code("PATTERN_ERR_1002"), "Trigger cannot be empty")
	ErrEmptyResponse       = New(Code("PATTERN_ERR_1003"), "Response cannot be empty")
	ErrInvalidStrength     = New(Code("PATTERN_ERR_1004"), "Strength must be between 0 and 100")

	// Storage errors (2000-2099)
	ErrDatabaseError       = New(Code("STORAGE_ERR_2003"), "Database operation failed")
	ErrSpaceNotFound       = New(Code("STORAGE_ERR_2002"), "Space not found")
	ErrDatabaseLocked      = New(Code("STORAGE_ERR_2009"), "Database locked")

	// AI errors (3000-3099)
	ErrAIProviderUnavailable = New(Code("AI_ERR_3002"), "AI provider unavailable")
	ErrAITimeout            = New(Code("AI_ERR_3006"), "AI request timeout")
	ErrAIRateLimited        = New(Code("AI_ERR_3004"), "Rate limit exceeded")
	ErrAPIKeyMissing        = New(Code("AI_ERR_3001"), "API key not configured")

	// CLI errors (4000-4099)
	ErrUnknownCommand      = New(Code("CLI_ERR_4001"), "Unknown command")
	ErrMissingArgument     = New(Code("CLI_ERR_4002"), "Missing required argument")
	ErrInvalidArgument      = New(Code("CLI_ERR_4003"), "Invalid argument value")

	// Config errors (5000-5099)
	ErrConfigNotFound      = New(Code("CONFIG_ERR_5001"), "Configuration file not found")
	ErrConfigParseError    = New(Code("CONFIG_ERR_5002"), "Configuration parse error")
)

// Is checks if the error matches a given code
func (e *Error) Is(code Code) bool {
	return e.Code == code
}

// IsNotFound checks if the error is a "not found" error
func (e *Error) IsNotFound() bool {
	return e.Code == ErrPatternNotFound.Code || e.Code == ErrSpaceNotFound.Code
}
