package logger

import (
	"io"
	"log"
	"os"
	"sync"
)

// Level represents log level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Logger is a structured logger (Iter 53)
type Logger struct {
	mu      sync.Mutex
	level   Level
	output  *log.Logger
	prefix  string
	flags   int
}

// New creates a new logger
func New(output io.Writer, prefix string, flags int) *Logger {
	return &Logger{
		level:  INFO,
		output: log.New(output, prefix, flags),
	}
}

// Default returns the default logger
var Default = New(os.Stdout, "", log.LstdFlags)

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// Debug logs a debug message
func (l *Logger) Debug(v ...interface{}) {
	if l.level <= DEBUG {
		l.output.Print("[DEBUG] ", v)
	}
}

// Info logs an info message
func (l *Logger) Info(v ...interface{}) {
	if l.level <= INFO {
		l.output.Print("[INFO] ", v)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(v ...interface{}) {
	if l.level <= WARN {
		l.output.Print("[WARN] ", v)
	}
}

// Error logs an error message
func (l *Logger) Error(v ...interface{}) {
	if l.level <= ERROR {
		l.output.Print("[ERROR] ", v)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(v ...interface{}) {
	l.output.Print("[FATAL] ", v)
	os.Exit(1)
}

// WithFields returns a logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Entry {
	return &Entry{logger: l, fields: fields}
}

// Entry represents a log entry with fields
type Entry struct {
	logger *Logger
	fields map[string]interface{}
}

// Debug logs with fields
func (e *Entry) Debug(v ...interface{}) {
	e.logger.output.Print("[DEBUG] ", e.formatFields(), " ", v)
}

// Info logs with fields
func (e *Entry) Info(v ...interface{}) {
	e.logger.output.Print("[INFO] ", e.formatFields(), " ", v)
}

// Error logs with fields
func (e *Entry) Error(v ...interface{}) {
	e.logger.output.Print("[ERROR] ", e.formatFields(), " ", v)
}

func (e *Entry) formatFields() string {
	if len(e.fields) == 0 {
		return ""
	}
	return ""
}

// Package-level functions

// Debug logs a debug message
func Debug(v ...interface{}) {
	Default.Debug(v...)
}

// Info logs an info message
func Info(v ...interface{}) {
	Default.Info(v...)
}

// Warn logs a warning message
func Warn(v ...interface{}) {
	Default.Warn(v...)
}

// Error logs an error message
func Error(v ...interface{}) {
	Default.Error(v...)
}

// Fatal logs a fatal message and exits
func Fatal(v ...interface{}) {
	Default.Fatal(v...)
}

// WithFields returns a logger with additional fields
func WithFields(fields map[string]interface{}) *Entry {
	return Default.WithFields(fields)
}
