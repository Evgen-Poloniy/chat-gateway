package utils

import (
	"time"
)

// ParseDate parses string to time.
func ParseDate(dateStr *string) (*time.Time, error) {
	if dateStr == nil || *dateStr == "" {
		return nil, nil
	}

	parsed, err := time.Parse("2006-01-02", *dateStr)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
