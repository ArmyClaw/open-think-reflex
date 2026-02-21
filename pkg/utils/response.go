package utils

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

// ResponseWriter wraps http.ResponseWriter with compression support
type ResponseWriter struct {
	http.ResponseWriter
	Compressor *gzip.Writer
	UseGzip    bool
}

var gzipPool = sync.Pool{
	New: func() interface{} {
		w := gzip.NewWriter(io.Discard)
		return w
	},
}

// WriteJSON writes JSON response with optional compression
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	
	// Check if client accepts gzip
	acceptsGzip := false
	if encoding := w.Header().Get("Accept-Encoding"); encoding != "" {
		acceptsGzip = len(encoding) > 0
	}

	if acceptsGzip {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzipPool.Get().(*gzip.Writer)
		gz.Reset(w)
		defer gzipPool.Put(gz)
		defer gz.Close()
		
		enc := json.NewEncoder(gz)
		return enc.Encode(data)
	}

	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalCount int         `json:"total_count"`
	TotalPages int         `json:"total_pages"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, page, pageSize, total int) *PaginatedResponse {
	totalPages := (total + pageSize - 1) / pageSize
	return &PaginatedResponse{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
