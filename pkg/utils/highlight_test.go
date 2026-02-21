package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyntaxHighlighter_HighlightJSON(t *testing.T) {
	h := DefaultHighlighter()
	
	jsonStr := `{"name": "test", "count": 42, "active": true}`
	result := h.HighlightJSON(jsonStr)
	
	// Should contain colored tags
	assert.True(t, strings.Contains(result, "[green]"), "Should contain green for strings")
	assert.True(t, strings.Contains(result, "[yellow]"), "Should contain yellow for numbers")
	assert.True(t, strings.Contains(result, "[cyan]"), "Should contain cyan for keywords")
}

func TestSyntaxHighlighter_HighlightJSONNested(t *testing.T) {
	h := DefaultHighlighter()
	
	jsonStr := `{"user": {"name": "Alice", "age": 30}, "active": false}`
	result := h.HighlightJSON(jsonStr)
	
	assert.True(t, strings.Contains(result, "[green]"), "Should contain green for strings")
	assert.True(t, strings.Contains(result, "[yellow]"), "Should contain yellow for numbers")
}

func TestSyntaxHighlighter_HighlightGo(t *testing.T) {
	h := DefaultHighlighter()
	
	goCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}`
	
	result := h.highlightGo(goCode)
	
	// Should highlight keywords
	assert.True(t, strings.Contains(result, "[cyan]package[white]"))
	assert.True(t, strings.Contains(result, "[cyan]func[white]"))
	assert.True(t, strings.Contains(result, "[green]"))
}

func TestSyntaxHighlighter_HighlightPython(t *testing.T) {
	h := DefaultHighlighter()
	
	pyCode := `def hello():
    print("Hello")
    return True`
	
	result := h.highlightPython(pyCode)
	
	// Should highlight keywords
	assert.True(t, strings.Contains(result, "[cyan]def[white]"))
	assert.True(t, strings.Contains(result, "[cyan]return[white]"))
}

func TestSyntaxHighlighter_HighlightJS(t *testing.T) {
	h := DefaultHighlighter()
	
	jsCode := `const x = function() {
	return true;
}`
	
	result := h.highlightJS(jsCode)
	
	// Should highlight keywords (check for presence, not exact format)
	assert.True(t, strings.Contains(result, "[cyan]const[white]"))
	assert.True(t, strings.Contains(result, "[cyan]function[white]") || strings.Contains(result, "[magenta]function[white]"))
	assert.True(t, strings.Contains(result, "[cyan]return[white]"))
}

func TestSyntaxHighlighter_HighlightBash(t *testing.T) {
	h := DefaultHighlighter()
	
	bashCode := `#!/bin/bash
echo "Hello"
if [ -f file ]; then
    echo "exists"
fi`
	
	result := h.highlightBash(bashCode)
	
	// Should highlight keywords
	assert.True(t, strings.Contains(result, "[cyan]echo[white]"))
	assert.True(t, strings.Contains(result, "[cyan]if[white]"))
	assert.True(t, strings.Contains(result, "[cyan]fi[white]"))
}

func TestSyntaxHighlighter_HighlightGeneric(t *testing.T) {
	h := DefaultHighlighter()
	
	code := `// This is a comment
var x = 42
string s = "hello"`
	
	result := h.highlightGeneric(code)
	
	// Should highlight comments and numbers
	assert.True(t, strings.Contains(result, "[gray]"))
	assert.True(t, strings.Contains(result, "[yellow]"))
}

func TestSyntaxHighlighter_HighlightComments(t *testing.T) {
	h := DefaultHighlighter()
	
	code := `line1
// comment
line2`
	
	result := h.highlightComments(code, "//", "\n")
	
	assert.True(t, strings.Contains(result, "[gray]"))
}

func TestSyntaxHighlighter_HighlightStrings(t *testing.T) {
	h := DefaultHighlighter()
	
	code := `hello "world" test`
	
	result := h.highlightStrings(code, `"`, `"`)
	
	assert.True(t, strings.Contains(result, "[green]\"world\"[white]"))
}

func TestSyntaxHighlighter_HighlightKeywords(t *testing.T) {
	h := DefaultHighlighter()
	
	code := `if else while for`
	keywords := []string{"if", "else", "while", "for"}
	
	result := h.highlightKeywords(code, keywords)
	
	assert.True(t, strings.Contains(result, "[cyan]if[white]"))
	assert.True(t, strings.Contains(result, "[cyan]else[white]"))
	assert.True(t, strings.Contains(result, "[cyan]while[white]"))
	assert.True(t, strings.Contains(result, "[cyan]for[white]"))
}

func TestFormatJSON(t *testing.T) {
	jsonStr := `{"name":"test","count":42}`
	result, err := FormatJSON(jsonStr)
	
	assert.NoError(t, err)
	assert.Contains(t, result, "\n")
	assert.Contains(t, result, "  ")
}

func TestFormatJSON_Nested(t *testing.T) {
	jsonStr := `{"user":{"name":"Alice","age":30}}`
	result, err := FormatJSON(jsonStr)
	
	assert.NoError(t, err)
	assert.Contains(t, result, "\n")
	assert.Contains(t, result, "  ")
}

func TestDetectLanguage(t *testing.T) {
	// Test JSON detection
	assert.Equal(t, "json", DetectLanguage(`{"key": "value"}`))
	
	// Test Go detection
	assert.Equal(t, "go", DetectLanguage(`package main
func main() {}`))
	
	// Test Python detection
	assert.Equal(t, "python", DetectLanguage(`def hello():
    pass`))
	
	// Test JavaScript detection
	assert.Equal(t, "javascript", DetectLanguage(`const x = () => {}`))
	
	// Test Bash detection
	assert.Equal(t, "bash", DetectLanguage(`#!/bin/bash
echo "hello"`))
	
	// Test default
	assert.Equal(t, "plain", DetectLanguage(`some random text`))
}

func TestDefaultHighlighter(t *testing.T) {
	h := DefaultHighlighter()
	
	assert.Equal(t, "cyan", h.KeywordColor)
	assert.Equal(t, "green", h.StringColor)
	assert.Equal(t, "yellow", h.NumberColor)
	assert.Equal(t, "gray", h.CommentColor)
	assert.Equal(t, "blue", h.FunctionColor)
	assert.Equal(t, "magenta", h.TypeColor)
	assert.Equal(t, "white", h.Plain)
}

func TestHighlightCode(t *testing.T) {
	h := DefaultHighlighter()
	
	// Test JSON
	result := h.HighlightCode(`{"a":1}`, "json")
	assert.True(t, strings.Contains(result, "["))
	
	// Test Go
	result = h.HighlightCode(`package main`, "go")
	assert.True(t, strings.Contains(result, "["))
	
	// Test unknown language (should fall back to generic)
	result = h.HighlightCode(`test`, "unknown")
	assert.NotEmpty(t, result)
}
