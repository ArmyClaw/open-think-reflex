package models

import (
	"testing"
)

func TestNoteValidate(t *testing.T) {
	tests := []struct {
		name    string
		note    Note
		wantErr bool
	}{
		{
			name: "valid note",
			note: Note{
				Title:   "Test",
				Content: "Content",
			},
			wantErr: false,
		},
		{
			name: "empty title",
			note: Note{
				Title:   "",
				Content: "Content",
			},
			wantErr: true,
		},
		{
			name: "empty content",
			note: Note{
				Title:   "Title",
				Content: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.note.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNoteCalculateStats(t *testing.T) {
	note := Note{
		Content: "Hello world this is a test",
	}
	note.CalculateStats()

	if note.WordCount != 6 {
		t.Errorf("Expected 6 words, got %d", note.WordCount)
	}
	if note.CharCount != 26 {
		t.Errorf("Expected 26 chars, got %d", note.CharCount)
	}
}

func TestNotePreview(t *testing.T) {
	note := Note{
		Content: "This is a long content that should be truncated",
	}

	// Test with long content
	preview := note.Preview(10)
	if len(preview) != 13 { // 10 + "..."
		t.Errorf("Preview length = %d, want 13", len(preview))
	}

	// Test with short content
	note.Content = "Short"
	preview = note.Preview(10)
	if preview != "Short" {
		t.Errorf("Preview = %s, want Short", preview)
	}
}

func TestNotePatternLinks(t *testing.T) {
	note := Note{}

	// Test AddPattern
	note.AddPattern("pattern-1")
	if len(note.PatternIDs) != 1 {
		t.Errorf("Expected 1 pattern, got %d", len(note.PatternIDs))
	}

	// Test duplicate AddPattern
	note.AddPattern("pattern-1")
	if len(note.PatternIDs) != 1 {
		t.Errorf("Expected 1 pattern after duplicate add, got %d", len(note.PatternIDs))
	}

	// Test HasPattern
	if !note.HasPattern("pattern-1") {
		t.Error("Expected HasPattern to return true")
	}
	if note.HasPattern("pattern-2") {
		t.Error("Expected HasPattern to return false")
	}

	// Test RemovePattern
	note.RemovePattern("pattern-1")
	if len(note.PatternIDs) != 0 {
		t.Errorf("Expected 0 patterns after remove, got %d", len(note.PatternIDs))
	}
}

func TestNoteCategories(t *testing.T) {
	if CategoryThought != "thought" {
		t.Errorf("CategoryThought = %s", CategoryThought)
	}
	if CategoryIdea != "idea" {
		t.Errorf("CategoryIdea = %s", CategoryIdea)
	}
	if CategoryTodo != "todo" {
		t.Errorf("CategoryTodo = %s", CategoryTodo)
	}
	if CategoryMemory != "memory" {
		t.Errorf("CategoryMemory = %s", CategoryMemory)
	}
	if CategoryQuestion != "question" {
		t.Errorf("CategoryQuestion = %s", CategoryQuestion)
	}
	if CategoryNote != "note" {
		t.Errorf("CategoryNote = %s", CategoryNote)
	}
}

func TestNoteNew(t *testing.T) {
	note := NewNote("Test Title", "Test Content")

	if note.Title != "Test Title" {
		t.Errorf("Title = %s, want Test Title", note.Title)
	}
	if note.Content != "Test Content" {
		t.Errorf("Content = %s, want Test Content", note.Content)
	}
	if note.Category != CategoryNote {
		t.Errorf("Category = %s, want note", note.Category)
	}
	if note.ID == "" {
		t.Error("Expected ID to be set")
	}
	if note.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}
