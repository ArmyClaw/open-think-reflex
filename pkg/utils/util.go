// Package utils provides common utility functions.
package utils

import "time"

// MinInt returns the minimum of two integers
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxInt returns the maximum of two integers
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// UnixToTime converts Unix timestamp to time.Time
func UnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// Ptr returns a pointer to the given value
func Ptr[T any](v T) *T {
	return &v
}
