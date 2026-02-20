# Open-Think-Reflex Error Codes

> **Version**: v1.0  
> **Status**: Active  
> **Scope**: Error codes and handling for v1.0 implementation

---

## Table of Contents

1. [Error Code Structure](#1-error-code-structure)
2. [Validation Errors (1xxx)](#2-validation-errors-1xxx)
3. [Storage Errors (2xxx)](#3-storage-errors-2xxx)
4. [AI Errors (3xxx)](#4-ai-errors-3xxx)
5. [CLI Errors (4xxx)](#5-cli-errors-4xxx)
6. [Configuration Errors (5xxx)](#6-configuration-errors-5xxx)
7. [System Errors (9xxx)](#7-system-errors-9xxx)
8. [Error Handling Patterns](#8-error-handling-patterns)
9. [Error Recovery Strategies](#9-error-recovery-strategies)

---

## 1. Error Code Structure

### 1.1 Error Code Format

```
[PREFIX]_[SEVERITY]_[MODULE]_[NUMBER]

Example: PATTERN_ERR_001
- PATTERN: Component/module
- ERR: Severity (ERR/WARN/INFO)
- 001: Unique error number
```

### 1.2 Severity Levels

| Level | Prefix | Description | User Action |
|-------|--------|-------------|-------------|
| **ERROR** | ERR | Operation failed completely | Must fix before retry |
| **WARNING** | WARN | Operation partially succeeded | May need attention |
| **INFO** | INFO | Informational message | No action needed |

### 1.3 Module Codes

| Module | Code | Description |
|--------|------|-------------|
| **Validation** | PATTERN | Pattern validation |
| **Storage** | STORAGE | Data persistence layer |
| **AI** | AI | AI provider interactions |
| **CLI** | CLI | Command-line interface |
| **Config** | CONFIG | Configuration loading |
| **System** | SYSTEM | System-level errors |

---

## 2. Validation Errors (1xxx)

### 2.1 Pattern Validation (1000-1099)

| Error Code | Message | HTTP Status | Cause |
|------------|---------|-------------|-------|
| `PATTERN_ERR_1001` | Trigger cannot be empty | 400 | Pattern trigger is empty |
| `PATTERN_ERR_1002` | Trigger exceeds 500 characters | 400 | Trigger too long |
| `PATTERN_ERR_1003` | Response cannot be empty | 400 | Pattern response is empty |
| `PATTERN_ERR_1004` | Strength must be between 0 and 100 | 400 | Invalid strength value |
| `PATTERN_ERR_1005` | Threshold must be between 0 and 100 | 400 | Invalid threshold value |
| `PATTERN_ERR_1006` | Decay rate must be between 0 and 1 | 400 | Invalid decay rate |
| `PATTERN_ERR_1007` | Invalid pattern ID format | 400 | UUID parsing failed |
| `PATTERN_ERR_1008` | Duplicate trigger in project | 409 | Trigger already exists |
| `PATTERN_ERR_1009` | Invalid tag format | 400 | Tag contains invalid characters |
| `PATTERN_ERR_1010` | Project name too long | 400 | Project name > 100 chars |

### 2.2 Input Validation (1100-1199)

| Error Code | Message | HTTP Status | Cause |
|------------|---------|-------------|-------|
| `INPUT_ERR_1101` | Query string too long | 400 | Query > 1000 chars |
| `INPUT_ERR_1102` | Invalid query format | 400 | Query contains invalid chars |
| `INPUT_ERR_1103` | Missing required field | 400 | Required field is empty |
| `INPUT_ERR_1104` | Invalid JSON format | 400 | JSON parsing failed |
| `INPUT_ERR_1105` | Invalid YAML format | 400 | YAML parsing failed |

### 2.3 Validation Error Response Format

```json
{
  "error": {
    "code": "PATTERN_ERR_1001",
    "message": "Trigger cannot be empty",
    "details": {
      "field": "trigger",
      "value": "",
      "constraint": "must not be empty"
    },
    "timestamp": "2026-02-20T10:00:00Z",
    "request_id": "req_abc123"
  }
}
```

---

## 3. Storage Errors (2xxx)

### 3.1 Database Errors (2000-2099)

| Error Code | Message | HTTP Status | Cause |
|------------|---------|-------------|-------|
| `STORAGE_ERR_2001` | Pattern not found | 404 | Pattern ID doesn't exist |
| `STORAGE_ERR_2002` | Space not found | 404 | Space ID doesn't exist |
| `STORAGE_ERR_2003` | Database connection failed | 503 | Cannot connect to database |
| `STORAGE_ERR_2004` | Database transaction failed | 500 | Transaction rolled back |
| `STORAGE_ERR_2005` | Concurrent modification detected | 409 | Optimistic lock failed |
| `STORAGE_ERR_2006` | Database schema mismatch | 500 | Migration required |
| `STORAGE_ERR_2007` | Disk full | 507 | Insufficient storage |
| `STORAGE_ERR_2008` | Permission denied | 403 | File/dir access denied |
| `STORAGE_ERR_2009` | Database locked | 503 | SQLite lock contention |

### 3.2 Cache Errors (2100-2199)

| Error Code | Message | Cause |
|------------|---------|-------|
| `CACHE_ERR_2101` | Cache miss | Key not found |
| `CACHE_ERR_2102` | Cache corruption | Checksum verification failed |
| `CACHE_ERR_2103` | Cache size exceeded | LRU eviction failed |
| `CACHE_ERR_2104` | Invalid cache key | Key format invalid |

### 3.3 Export/Import Errors (2200-2299)

| Error Code | Message | Cause |
|------------|---------|-------|
| `EXPORT_ERR_2201` | Unsupported format | Export format not supported |
| `EXPORT_ERR_2202` | Export failed | Write operation failed |
| `IMPORT_ERR_2211` | Unsupported format | Import format not supported |
| `IMPORT_ERR_2212` | Import validation failed | Schema mismatch |
| `IMPORT_ERR_2213` | Import duplicate detected | Pattern already exists |
| `IMPORT_ERR_2214` | Import partial success | Some patterns skipped |

---

## 4. AI Errors (3xxx)

### 4.1 Provider Errors (3000-3099)

| Error Code | Message | HTTP Status | Cause |
|------------|---------|-------------|-------|
| `AI_ERR_3001` | Provider not configured | 400 | API key missing |
| `AI_ERR_3002` | Provider unavailable | 503 | Service down |
| `AI_ERR_3003` | Invalid API key | 401 | Authentication failed |
| `AI_ERR_3004` | Rate limit exceeded | 429 | Too many requests |
| `AI_ERR_3005` | Quota exceeded | 402 | Billing limit reached |
| `AI_ERR_3006` | Request timeout | 504 | AI service too slow |
| `AI_ERR_3007` | Invalid model | 400 | Model not available |
| `AI_ERR_3008` | Context length exceeded | 400 | Prompt too long |

### 4.2 Response Errors (3100-3199)

| Error Code | Message | Cause |
|------------|---------|-------|
| `AI_ERR_3101` | Response parsing failed | Invalid JSON from AI |
| `AI_ERR_3102` | Response content filtered | Safety filter triggered |
| `AI_ERR_3103` | Response incomplete | Stream ended early |
| `AI_ERR_3104` | Response too long | Max tokens exceeded |

### 4.3 Anthropic-Specific (3200-3299)

| Error Code | Message | Cause |
|------------|---------|-------|
| `AI_ERR_3201` | Anthropic API error | See API response |
| `AI_ERR_3202` | Claude version deprecated | Model needs update |

### 4.4 AI Error Response Format

```json
{
  "error": {
    "code": "AI_ERR_3004",
    "message": "Rate limit exceeded",
    "details": {
      "provider": "anthropic",
      "retry_after": 60,
      "current_usage": 45000,
      "rate_limit": 50000
    },
    "timestamp": "2026-02-20T10:00:00Z",
    "request_id": "req_abc123"
  }
}
```

---

## 5. CLI Errors (4xxx)

### 5.1 Command Errors (4000-4099)

| Error Code | Message | Cause |
|------------|---------|-------|
| `CLI_ERR_4001` | Unknown command | Command not found |
| `CLI_ERR_4002` | Missing required argument | Required arg missing |
| `CLI_ERR_4003` | Invalid argument value | Arg validation failed |
| `CLI_ERR_4004` | Too many arguments | Extra args provided |
| `CLI_ERR_4005` | Flag not found | Invalid flag |
| `CLI_ERR_4006` | Flag value invalid | Flag parsing failed |

### 5.2 UI Errors (4100-4199)

| Error Code | Message | Cause |
|------------|---------|-------|
| `UI_ERR_4101` | Terminal not supported | TUI requires terminal |
| `UI_ERR_4102` | Terminal resize failed | Window too small |
| `UI_ERR_4103` | Input interrupted | Ctrl+C pressed |
| `UI_ERR_4104` | Copy to clipboard failed | No clipboard support |

---

## 6. Configuration Errors (5xxx)

### 6.1 Config File Errors (5000-5099)

| Error Code | Message | Cause |
|------------|---------|-------|
| `CONFIG_ERR_5001` | Config file not found | Missing config |
| `CONFIG_ERR_5002` | Config parse error | YAML/JSON syntax error |
| `CONFIG_ERR_5003` | Config schema error | Invalid config structure |
| `CONFIG_ERR_5004` | Config version mismatch | Needs migration |
| `CONFIG_ERR_5005` | Config permission denied | Cannot read/write |
| `CONFIG_ERR_5006` | Deprecated config key | Key renamed/removed |

### 6.2 Environment Errors (5100-5199)

| Error Code | Message | Cause |
|------------|---------|-------|
| `CONFIG_ERR_5101` | Environment variable not set | Required env var missing |
| `CONFIG_ERR_5102` | Invalid environment value | Env var format invalid |

---

## 7. System Errors (9xxx)

### 9.1 Internal Errors (9000-9099)

| Error Code | Message | Cause |
|------------|---------|-------|
| `SYSTEM_ERR_9001` | Internal error | Unexpected panic |
| `SYSTEM_ERR_9002` | Memory allocation failed | Out of memory |
| `SYSTEM_ERR_9003` | Stack overflow | Infinite recursion |
| `SYSTEM_ERR_9004` | Deadlock detected | Lock ordering issue |

### 9.2 Signal/Interrupt Errors (9100-9199)

| Error Code | Message | Cause |
|------------|---------|-------|
| `SYSTEM_ERR_9101` | Interrupted by user | Ctrl+C received |
| `SYSTEM_ERR_9102` | Shutdown signal received | SIGTERM/SIGHUP |

---

## 8. Error Handling Patterns

### 8.1 Error Type Definition

```go
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
    Stack     string   `json:"stack,omitempty"`
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

// WithStack adds stack trace to error
func (e *Error) WithStack() *Error {
    e.Stack = stackTrace()
    return e
}
```

### 8.2 Error Creation Helpers

```go
package errors

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
        Code:      Code(fmt.Sprintf("PATTERN_ERR_%04d", 1)),
        Message:   fmt.Sprintf("Validation failed for %s: %s", field, reason),
        Details: Details{
            Field:     field,
            Value:     value,
            Constraint: reason,
        },
        Timestamp: time.Now(),
    }
}
```

### 8.3 Common Error Wrapping

```go
// Pattern not found
var ErrPatternNotFound = New(
    Code("STORAGE_ERR_2001"),
    "Pattern not found",
)

// Wrap database error
func WrapDBError(err error) *Error {
    if err == sql.ErrNoRows {
        return ErrPatternNotFound.WithRequestID(generateRequestID())
    }
    return Wrap(err, Code("STORAGE_ERR_2003"), "Database operation failed")
}

// Wrap AI error
func WrapAIError(err error, provider string) *Error {
    // Extract retry-after if available
    retryAfter := extractRetryAfter(err)
    
    return &Error{
        Code:    Code("AI_ERR_3001"),
        Message: fmt.Sprintf("%s API error: %v", provider, err),
        Details: Details{
            Provider: provider,
            RetryAfter: retryAfter,
        },
        Timestamp: time.Now(),
    }
}
```

### 8.4 Error Handling Middleware

```go
// Handler wraps HTTP handlers with error handling
func Handler(fn func(w http.ResponseWriter, r *http.Request) *Error) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := fn(w, r)
        if err == nil {
            return
        }
        
        status := http.StatusInternalServerError
        switch err.Code[0:1] {
        case "1": // Validation
            status = http.StatusBadRequest
        case "2": // Storage
            if strings.Contains(string(err.Code), "NOT_FOUND") {
                status = http.StatusNotFound
            } else {
                status = http.StatusInternalServerError
            }
        case "3": // AI
            if strings.Contains(string(err.Code), "RATE_LIMIT") {
                status = http.StatusTooManyRequests
            } else {
                status = http.StatusBadGateway
            }
        case "4": // CLI
            status = http.StatusBadRequest
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(status)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "error": err,
        })
    }
}
```

---

## 9. Error Recovery Strategies

### 9.1 Retry Logic

```go
// RetryConfig defines retry behavior
type RetryConfig struct {
    MaxRetries   int
    InitialDelay time.Duration
    MaxDelay     time.Duration
    Multiplier   float64
}

// WithRetry executes a function with retry
func WithRetry(fn func() error, config RetryConfig) error {
    var lastErr error
    
    delay := config.InitialDelay
    for attempt := 0; attempt <= config.MaxRetries; attempt++ {
        if err := fn(); err != nil {
            lastErr = err
            
            // Check if retryable
            if !isRetryable(err) {
                return err
            }
            
            // Check if max delay reached
            if delay > config.MaxDelay {
                delay = config.MaxDelay
            }
            
            // Wait before retry
            time.Sleep(delay)
            delay *= time.Duration(config.Multiplier)
            continue
        }
        return nil
    }
    
    return fmt.Errorf("max retries (%d) exceeded: %v", config.MaxRetries, lastErr)
}

// isRetryable checks if error can be retried
func isRetryable(err error) bool {
    // Retryable errors
    retryableCodes := map[Code]bool{
        "AI_ERR_3002": true, // Provider unavailable
        "AI_ERR_3004": true, // Rate limit
        "AI_ERR_3006": true, // Timeout
        "STORAGE_ERR_2003": true, // Connection failed
        "STORAGE_ERR_2009": true, // Lock contention
    }
    
    if e, ok := err.(*Error); ok {
        return retryableCodes[e.Code]
    }
    
    // Also retry network errors
    return net.Error(err) != nil
}
```

### 9.2 Circuit Breaker

```go
// CircuitBreaker prevents cascade failures
type CircuitBreaker struct {
    name        string
    failures    int
    successes   int
    lastFailure time.Time
    state       State
    
    threshold   int           // Failures before opening
    timeout     time.Duration  // Time to stay open
    interval    time.Duration  // Success interval
}

// State represents circuit breaker state
type State string

const (
    StateClosed   State = "closed"   // Normal operation
    StateOpen     State = "open"     // Failing, reject requests
    StateHalfOpen State = "half_open" // Testing recovery
)

// Execute runs function through circuit breaker
func (cb *CircuitBreaker) Execute(fn func() error) error {
    if cb.state == StateOpen {
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = StateHalfOpen
        } else {
            return fmt.Errorf("circuit breaker is open")
        }
    }
    
    err := fn()
    
    if err != nil {
        cb.failures++
        if cb.failures >= cb.threshold {
            cb.state = StateOpen
            cb.lastFailure = time.Now()
        }
        return err
    }
    
    cb.successes++
    if cb.state == StateHalfOpen {
        cb.failures = 0
        cb.state = StateClosed
    }
    
    return nil
}
```

---

## Error Code Quick Reference

| Range | Category | Example |
|-------|----------|---------|
| 1000-1099 | Pattern Validation | `PATTERN_ERR_1001` |
| 1100-1199 | Input Validation | `INPUT_ERR_1101` |
| 2000-2099 | Database | `STORAGE_ERR_2001` |
| 2100-2199 | Cache | `CACHE_ERR_2101` |
| 2200-2299 | Export/Import | `EXPORT_ERR_2201` |
| 3000-3099 | AI Provider | `AI_ERR_3001` |
| 3100-3199 | AI Response | `AI_ERR_3101` |
| 4000-4099 | CLI Commands | `CLI_ERR_4001` |
| 4100-4199 | UI | `UI_ERR_4101` |
| 5000-5099 | Config File | `CONFIG_ERR_5001` |
| 5100-5199 | Environment | `CONFIG_ERR_5101` |
| 9000-9099 | Internal | `SYSTEM_ERR_9001` |

---

**Document Version**: v1.0  
**Created**: 2026-02-20  
**Project**: open-think-reflex
