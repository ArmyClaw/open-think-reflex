package utils

import (
	"regexp"
	"strings"
)

// SyntaxHighlighter provides syntax highlighting for code and JSON
type SyntaxHighlighter struct {
	// Color scheme using tview tags
	KeywordColor   string
	StringColor    string
	NumberColor    string
	CommentColor   string
	FunctionColor  string
	TypeColor      string
	Punctuation    string
	Plain          string
}

// DefaultHighlighter returns a highlighter with default colors
func DefaultHighlighter() *SyntaxHighlighter {
	return &SyntaxHighlighter{
		KeywordColor:  "cyan",
		StringColor:   "green",
		NumberColor:   "yellow",
		CommentColor:  "gray",
		FunctionColor: "blue",
		TypeColor:     "magenta",
		Punctuation:   "white",
		Plain:         "white",
	}
}

// color wraps text with tview color tags
func (h *SyntaxHighlighter) color(text, color string) string {
	return "[" + color + "]" + text + "[" + h.Plain + "]"
}

// HighlightJSON highlights JSON content with formatting
func (h *SyntaxHighlighter) HighlightJSON(jsonStr string) string {
	// First format the JSON
	formatted, err := FormatJSON(jsonStr)
	if err != nil {
		return h.highlightGeneric(jsonStr)
	}

	// Then apply highlighting
	return h.highlightJSONFormatted(formatted)
}

// highlightJSONFormatted highlights pre-formatted JSON
func (h *SyntaxHighlighter) highlightJSONFormatted(jsonStr string) string {
	lines := strings.Split(jsonStr, "\n")
	var results []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Process each token
		var tokens []string
		// Use regex to find all JSON tokens
		re := regexp.MustCompile(`("[^"]*")|(\d+\.?\d*)|(\{)|(\})|(\])|(\[)|(:)|(,)|(true)|(false)|(null)`)
		
		matches := re.FindAllStringSubmatchIndex(line, -1)
		if len(matches) == 0 {
			tokens = append(tokens, line)
			continue
		}

		var colored strings.Builder
		lastEnd := 0
		for _, m := range matches {
			// Add text before match
			if m[0] > lastEnd {
				colored.WriteString(line[lastEnd:m[0]])
			}

			// Determine token type and color
			// Groups: 1=string, 2=number, 3={, 4=}, 5=], 6=[, 7=:, 8=,, 9=true, 10=false, 11=null
			token := line[m[0]:m[1]]
			var coloredToken string

			switch {
			case m[2] != -1: // string (key or value)
				// Check if this is a key (followed by :)
				after := strings.TrimSpace(line[m[1]:])
				if strings.HasPrefix(after, ":") {
					coloredToken = h.color(token, h.StringColor) // Key
				} else {
					coloredToken = h.color(token, h.StringColor) // String value
				}
			case m[4] != -1: // number
				coloredToken = h.color(token, h.NumberColor)
			case m[6] != -1: // {
				coloredToken = h.color(token, h.Punctuation)
			case m[7] != -1: // }
				coloredToken = h.color(token, h.Punctuation)
			case m[8] != -1: // ]
				coloredToken = h.color(token, h.Punctuation)
			case m[9] != -1: // [
				coloredToken = h.color(token, h.Punctuation)
			case m[10] != -1: // :
				coloredToken = h.color(token, h.Punctuation)
			case m[11] != -1: // ,
				coloredToken = h.color(token, h.Punctuation)
			case m[12] != -1: // true
				coloredToken = h.color(token, h.KeywordColor)
			case m[13] != -1: // false
				coloredToken = h.color(token, h.KeywordColor)
			case m[14] != -1: // null
				coloredToken = h.color(token, h.KeywordColor)
			default:
				coloredToken = token
			}

			colored.WriteString(coloredToken)
			lastEnd = m[1]
		}

		// Add remaining text
		if lastEnd < len(line) {
			colored.WriteString(line[lastEnd:])
		}

		results = append(results, colored.String())
	}

	return strings.Join(results, "\n")
}

// HighlightCode highlights common programming code
func (h *SyntaxHighlighter) HighlightCode(code, language string) string {
	switch strings.ToLower(language) {
	case "json":
		return h.HighlightJSON(code)
	case "go", "golang":
		return h.highlightGo(code)
	case "python", "py":
		return h.highlightPython(code)
	case "javascript", "js", "typescript", "ts":
		return h.highlightJS(code)
	case "bash", "shell", "sh":
		return h.highlightBash(code)
	default:
		return h.highlightGeneric(code)
	}
}

// highlightGo highlights Go code
func (h *SyntaxHighlighter) highlightGo(code string) string {
	keywords := []string{"package", "import", "func", "return", "var", "const", "type", "struct",
		"interface", "map", "chan", "select", "case", "default", "if", "else", "for", "range",
		"switch", "break", "continue", "goto", "defer", "go", "fallthrough", "nil", "true", "false"}

	types := []string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32",
		"uint64", "float32", "float64", "complex64", "complex128", "byte", "rune", "string", "bool",
		"error", "any"}

	code = h.highlightComments(code, "//", "\n")
	code = h.highlightComments(code, "/*", "*/")
	code = h.highlightStrings(code, `"`, `"`)
	code = h.highlightStrings(code, "`", "`")
	code = h.highlightKeywords(code, keywords)
	code = h.highlightTypes(code, types)

	return code
}

// highlightPython highlights Python code
func (h *SyntaxHighlighter) highlightPython(code string) string {
	keywords := []string{"def", "class", "return", "if", "elif", "else", "for", "while", "try",
		"except", "finally", "with", "as", "import", "from", "pass", "break", "continue", "raise",
		"yield", "lambda", "and", "or", "not", "in", "is", "True", "False", "None", "global",
		"nonlocal", "assert", "del", "async", "await"}

	builtins := []string{"print", "len", "range", "str", "int", "float", "list", "dict", "set",
		"tuple", "bool", "type", "input", "open", "file", "map", "filter", "zip", "enumerate",
		"sorted", "reversed", "sum", "min", "max", "abs", "round", "isinstance", "hasattr",
		"getattr", "setattr", "staticmethod", "classmethod", "property"}

	code = h.highlightComments(code, "#", "\n")
	code = h.highlightStrings(code, `"`, `"`)
	code = h.highlightStrings(code, "'''", "'''")
	code = h.highlightStrings(code, `"""`, `"""`)
	code = h.highlightKeywords(code, keywords)
	code = h.highlightBuiltins(code, builtins)

	return code
}

// highlightJS highlights JavaScript/TypeScript code
func (h *SyntaxHighlighter) highlightJS(code string) string {
	keywords := []string{"function", "return", "var", "let", "const", "if", "else", "for", "while",
		"do", "switch", "case", "default", "break", "continue", "try", "catch", "finally", "throw",
		"new", "delete", "typeof", "instanceof", "in", "of", "class", "extends", "super", "this",
		"import", "export", "from", "as", "async", "await", "yield", "true", "false", "null",
		"undefined", "NaN", "Infinity", "static", "get", "set"}

	types := []string{"string", "number", "boolean", "object", "function", "symbol", "bigint",
		"Array", "Object", "String", "Number", "Boolean", "Function", "Symbol", "Map", "Set",
		"WeakMap", "WeakSet", "Promise", "Date", "RegExp", "Error", "JSON", "Math", "console"}

	code = h.highlightComments(code, "//", "\n")
	code = h.highlightComments(code, "/*", "*/")
	code = h.highlightStrings(code, `"`, `"`)
	code = h.highlightStrings(code, "'", "'")
	code = h.highlightStrings(code, "`", "`")
	code = h.highlightKeywords(code, keywords)
	code = h.highlightTypes(code, types)

	return code
}

// highlightBash highlights Bash/Shell code
func (h *SyntaxHighlighter) highlightBash(code string) string {
	keywords := []string{"if", "then", "else", "elif", "fi", "for", "while", "do", "done", "case",
		"esac", "in", "function", "return", "local", "export", "source", "alias", "echo", "read",
		"exit", "break", "continue", "shift", "set", "unset", "true", "false"}

	code = h.highlightComments(code, "#", "\n")
	code = h.highlightStrings(code, `"`, `"`)
	code = h.highlightStrings(code, "'", "'")
	code = h.highlightKeywords(code, keywords)

	// Highlight variables
	varRe := regexp.MustCompile(`\$\{?[\w]+\}?`)
	code = varRe.ReplaceAllString(code, h.color("$0", h.KeywordColor))

	return code
}

// highlightGeneric highlights generic text with basic patterns
func (h *SyntaxHighlighter) highlightGeneric(code string) string {
	code = h.highlightComments(code, "//", "\n")
	code = h.highlightComments(code, "#", "\n")
	code = h.highlightStrings(code, `"`, `"`)
	code = h.highlightStrings(code, "'", "'")

	// Numbers
	numberRe := regexp.MustCompile(`\b\d+\.?\d*\b`)
	code = numberRe.ReplaceAllString(code, h.color("$0", h.NumberColor))

	return code
}

// highlightComments highlights comments
func (h *SyntaxHighlighter) highlightComments(code, startMarker, endMarker string) string {
	if endMarker == "\n" {
		re := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(startMarker) + `.*?$`)
		return re.ReplaceAllStringFunc(code, func(s string) string {
			return h.color(strings.TrimSuffix(s, "\n"), h.CommentColor)
		})
	}

	re := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(startMarker) + `.*?` + regexp.QuoteMeta(endMarker))
	return re.ReplaceAllStringFunc(code, func(s string) string {
		return h.color(s, h.CommentColor)
	})
}

// highlightStrings highlights string literals
func (h *SyntaxHighlighter) highlightStrings(code, startMarker, endMarker string) string {
	re := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(startMarker) + `.*?` + regexp.QuoteMeta(endMarker))
	return re.ReplaceAllStringFunc(code, func(s string) string {
		return h.color(s, h.StringColor)
	})
}

// highlightKeywords highlights keywords
func (h *SyntaxHighlighter) highlightKeywords(code string, keywords []string) string {
	for _, kw := range keywords {
		re := regexp.MustCompile(`\b` + kw + `\b`)
		code = re.ReplaceAllString(code, h.color(kw, h.KeywordColor))
	}
	return code
}

// highlightTypes highlights types
func (h *SyntaxHighlighter) highlightTypes(code string, types []string) string {
	for _, t := range types {
		re := regexp.MustCompile(`\b` + t + `\b`)
		code = re.ReplaceAllString(code, h.color(t, h.TypeColor))
	}
	return code
}

// highlightBuiltins highlights built-in functions
func (h *SyntaxHighlighter) highlightBuiltins(code string, builtins []string) string {
	for _, b := range builtins {
		re := regexp.MustCompile(`\b` + b + `\b`)
		code = re.ReplaceAllString(code, h.color(b, h.FunctionColor))
	}
	return code
}

// FormatJSON formats JSON string with indentation
func FormatJSON(jsonStr string) (string, error) {
	var result strings.Builder
	depth := 0
	inString := false
	var prevChar rune

	for i, char := range jsonStr {
		if prevChar == '\\' {
			prevChar = char
			result.WriteRune(char)
			continue
		}

		if char == '"' && prevChar != '\\' {
			inString = !inString
		}

		if !inString {
			switch char {
			case '{', '[':
				result.WriteRune(char)
				depth++
				if i+1 < len(jsonStr) && jsonStr[i+1] != '}' && jsonStr[i+1] != ']' {
					result.WriteString("\n")
					result.WriteString(strings.Repeat("  ", depth))
				}
			case '}', ']':
				depth--
				if i > 0 && jsonStr[i-1] != '{' && jsonStr[i-1] != '[' {
					result.WriteString("\n")
					result.WriteString(strings.Repeat("  ", depth))
				}
				result.WriteRune(char)
			case ',':
				result.WriteRune(char)
				if i+1 < len(jsonStr) && jsonStr[i+1] != ' ' {
					result.WriteString("\n")
					result.WriteString(strings.Repeat("  ", depth))
				}
			case ':':
				result.WriteString(": ")
			case ' ', '\t', '\n', '\r':
				// Skip whitespace
			default:
				result.WriteRune(char)
			}
		} else {
			result.WriteRune(char)
		}

		prevChar = char
	}

	return result.String(), nil
}

// DetectLanguage detects the language of code snippet
func DetectLanguage(code string) string {
	code = strings.TrimSpace(code)

	// Check for JSON
	if (strings.HasPrefix(code, "{") && strings.HasSuffix(code, "}")) ||
		(strings.HasPrefix(code, "[") && strings.HasSuffix(code, "]")) {
		if strings.Contains(code, ":") && strings.Contains(code, `"`) {
			return "json"
		}
	}

	// Check for Go
	if strings.Contains(code, "package ") && strings.Contains(code, "func ") {
		return "go"
	}

	// Check for Python
	if strings.HasPrefix(code, "def ") || strings.HasPrefix(code, "class ") ||
		strings.Contains(code, "import ") && !strings.Contains(code, "import {") {
		return "python"
	}

	// Check for JavaScript/TypeScript
	if strings.Contains(code, "function ") || strings.Contains(code, "=>") ||
		strings.Contains(code, "const ") || strings.Contains(code, "let ") {
		return "javascript"
	}

	// Check for Bash
	if strings.HasPrefix(code, "#!") || strings.HasPrefix(code, "# ") ||
		strings.Contains(code, "echo ") || strings.Contains(code, "$") {
		return "bash"
	}

	return "plain"
}
