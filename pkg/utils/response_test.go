package utils

import (
	"net/http/httptest"
	"testing"
)

func TestNewPaginatedResponse(t *testing.T) {
	// Test basic pagination
	data := []string{"a", "b", "c"}
	resp := NewPaginatedResponse(data, 1, 10, 25)

	if resp.Page != 1 {
		t.Errorf("Expected page 1, got %d", resp.Page)
	}
	if resp.PageSize != 10 {
		t.Errorf("Expected pageSize 10, got %d", resp.PageSize)
	}
	if resp.TotalCount != 25 {
		t.Errorf("Expected total 25, got %d", resp.TotalCount)
	}
	if resp.TotalPages != 3 {
		t.Errorf("Expected totalPages 3, got %d", resp.TotalPages)
	}
	if !resp.HasNext {
		t.Error("Expected HasNext to be true")
	}
	if resp.HasPrev {
		t.Error("Expected HasPrev to be false")
	}

	// Test last page
	resp2 := NewPaginatedResponse(data, 3, 10, 25)
	if resp2.HasNext {
		t.Error("Expected HasNext to be false on last page")
	}
	if !resp2.HasPrev {
		t.Error("Expected HasPrev to be true")
	}

	// Test first page
	resp3 := NewPaginatedResponse(data, 1, 10, 10)
	if resp3.HasNext {
		t.Error("Expected HasNext to be false when exactly fits")
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	err := WriteJSON(w, 200, map[string]string{"key": "value"})
	if err != nil {
		t.Errorf("WriteJSON failed: %v", err)
	}
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json")
	}
}
