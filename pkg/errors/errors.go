package errors

import (
	"fmt"
	"strings"
)

// ErrorType represents the type of error
type ErrorType string

const (
	ErrTypeValidation ErrorType = "validation"
	ErrTypeNotFound   ErrorType = "not_found"
	ErrTypeConflict  ErrorType = "conflict"
	ErrTypeDatabase  ErrorType = "database"
	ErrTypeNetwork   ErrorType = "network"
	ErrTypeAuth     ErrorType = "auth"
	ErrTypeUnknown  ErrorType = "unknown"
)

// ErrorWithContext is an error with additional context (Iter 57)
type ErrorWithContext struct {
	Type        ErrorType
	Code        string
	Message     string
	Original    error
	Suggestions []string
}

// Error implements the error interface
func (e *ErrorWithContext) Error() string {
	var sb strings.Builder
	sb.WriteString(string(e.Type))
	sb.WriteString(": ")
	sb.WriteString(e.Message)
	if e.Original != nil {
		sb.WriteString(" (")
		sb.WriteString(e.Original.Error())
		sb.WriteString(")")
	}
	return sb.String()
}

// Unwrap returns the original error
func (e *ErrorWithContext) Unwrap() error {
	return e.Original
}

// WithSuggestion adds a suggestion to the error
func (e *ErrorWithContext) WithSuggestion(suggestion string) *ErrorWithContext {
	e.Suggestions = append(e.Suggestions, suggestion)
	return e
}

// WithSuggestions adds multiple suggestions
func (e *ErrorWithContext) WithSuggestions(suggestions []string) *ErrorWithContext {
	e.Suggestions = append(e.Suggestions, suggestions...)
	return e
}

// FormatForDisplay returns a user-friendly formatted error
func (e *ErrorWithContext) FormatForDisplay() string {
	var sb strings.Builder
	
	sb.WriteString("❌ ")
	sb.WriteString(e.Message)
	sb.WriteString("\n")
	
	if e.Code != "" {
		sb.WriteString("   Code: ")
		sb.WriteString(e.Code)
		sb.WriteString("\n")
	}
	
	if len(e.Suggestions) > 0 {
		sb.WriteString("   💡 Suggestions:\n")
		for _, s := range e.Suggestions {
			sb.WriteString("   - ")
			sb.WriteString(s)
			sb.WriteString("\n")
		}
	}
	
	return sb.String()
}

// New creates a new error with context
func New(errType ErrorType, message string) *ErrorWithContext {
	return &ErrorWithContext{
		Type:    errType,
		Message: message,
	}
}

// Wrap wraps an existing error with context
func Wrap(err error, errType ErrorType, message string) *ErrorWithContext {
	return &ErrorWithContext{
		Type:     errType,
		Message:  message,
		Original: err,
	}
}

// ValidationError creates a validation error
func ValidationError(message string, suggestions ...string) *ErrorWithContext {
	return New(ErrTypeValidation, message).WithSuggestions(suggestions)
}

// NotFoundError creates a not found error
func NotFoundError(resource string, id string) *ErrorWithContext {
	return New(ErrTypeNotFound, fmt.Sprintf("%s not found: %s", resource, id)).
		WithSuggestion(fmt.Sprintf("Check if the %s exists", resource))
}

// DatabaseError creates a database error
func DatabaseError(err error) *ErrorWithContext {
	return Wrap(err, ErrTypeDatabase, "Database operation failed").
		WithSuggestion("Check database connection and try again")
}
