package logger

import (
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

func TestPackageLevel(t *testing.T) {
	// Just ensure package level funcs don't panic
	Info("package level info")
	Debug("package level debug")
	Warn("package level warn")
	Error("package level error")
}
