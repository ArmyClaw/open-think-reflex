package models

import (
	"errors"
	"time"
)

// Note represents a thought/note in the system.
// Notes are lightweight thoughts that can be organized and searched.
type Note struct {
	// Identification
	ID        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Content
	Title   string `json:"title" db:"title"`     // Note title
	Content string `json:"content" db:"content"` // Note content (markdown supported)

	// Organization
	SpaceID string   `json:"space_id" db:"space_id"` // Associated Space
	Tags    []string `json:"tags" db:"-"`              // Tags for categorization

	// Status
	IsPinned bool   `json:"is_pinned" db:"is_pinned"` // Pinned to top
	Category string `json:"category" db:"category"`     // Category: thought/idea/todo/memory

	// Metadata
	WordCount int        `json:"word_count" db:"word_count"` // Content word count
	CharCount int        `json:"char_count" db:"char_count"` // Content character count
	LastViewed *time.Time `json:"last_viewed,omitempty" db:"last_viewed"` // Last viewed timestamp
}

// NoteCategory constants
const (
	CategoryThought  = "thought"  // General thoughts
	CategoryIdea     = "idea"     // Ideas and concepts
	CategoryTodo     = "todo"     // Tasks and todos
	CategoryMemory   = "memory"   // Memories to remember
	CategoryQuestion = "question" // Questions to explore
	CategoryNote     = "note"     // General notes
)

// ErrInvalidNote represents an invalid note error
var ErrInvalidNote = errors.New("invalid note")

// Validate checks if the note is valid
func (n *Note) Validate() error {
	if n.Title == "" {
		return ErrInvalidNote
	}
	if n.Content == "" {
		return ErrInvalidNote
	}
	return nil
}

// WordCount calculates the word count of the content
func (n *Note) CalculateStats() {
	n.CharCount = len(n.Content)
	// Simple word count: split by whitespace
	if n.Content == "" {
		n.WordCount = 0
		return
	}
	
	wordCount := 0
	inWord := false
	for _, r := range n.Content {
		if r == ' ' || r == '\n' || r == '\t' || r == '\r' {
			inWord = false
		} else if !inWord {
			inWord = true
			wordCount++
		}
	}
	n.WordCount = wordCount
}

// Preview returns a preview of the note content
func (n *Note) Preview(maxLen int) string {
	if len(n.Content) <= maxLen {
		return n.Content
	}
	return n.Content[:maxLen] + "..."
}
