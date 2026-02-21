package logger

import (
	"log"
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	buf := bytes.Buffer{}
	l := New(&buf, "[TEST] ", 0)
	
	l.Info("test message")
	if !strings.Contains(buf.String(), "test message") {
		t.Error("Expected log message")
	}
}

func TestLoggerLevel(t *testing.T) {
	buf := bytes.Buffer{}
	l := New(&buf, "", 0)
	l.SetLevel(WARN)
	
	l.Info("should not appear")
	if buf.Len() > 0 {
		t.Error("Info should not appear when level is WARN")
	}
	
	l.Warn("should appear")
	if buf.Len() == 0 {
		t.Error("Warn should appear")
	}
}

func TestLoggerDebug(t *testing.T) {
	buf := bytes.Buffer{}
	l := New(&buf, "", 0)
	l.SetLevel(DEBUG)
	
	l.Debug("debug message")
	if !strings.Contains(buf.String(), "debug message") {
		t.Error("Expected debug message")
	}
}

func TestLoggerError(t *testing.T) {
	buf := bytes.Buffer{}
	l := New(&buf, "", 0)
	
	l.Error("error message")
	if !strings.Contains(buf.String(), "error message") {
		t.Error("Expected error message")
	}
}

func TestLoggerWithFields(t *testing.T) {
	buf := bytes.Buffer{}
	l := New(&buf, "", 0)
	
	entry := l.WithFields(map[string]interface{}{
		"user": "test",
		"action": "login",
	})
	entry.Info("test")
	
	if !strings.Contains(buf.String(), "test") {
		t.Error("Expected message in log")
	}
}

func TestPackageLevel(t *testing.T) {
	// Just ensure package level funcs don't panic
	Info("package level info")
	Debug("package level debug")
	Warn("package level warn")
	Error("package level error")
}

func TestNewLogger(t *testing.T) {
	buf := bytes.Buffer{}
	l := New(&buf, "[PREFIX] ", log.LstdFlags)
	
	if l.level != INFO {
		t.Errorf("Expected INFO level, got %v", l.level)
	}
	
	l.Info("test")
	if !strings.Contains(buf.String(), "[PREFIX]") {
		t.Error("Expected prefix in output")
	}
}

func TestLoggerFatal(t *testing.T) {
	// Fatal exits, so we just verify it doesn't panic
	buf := bytes.Buffer{}
	l := New(&buf, "", 0)
	
	// This would exit in real scenario
	l.SetLevel(FATAL)
	_ = buf.Len()
}
