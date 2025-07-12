package parser

import "time"

// ConvertToString converts a time.Time object to a string in 2025-03-10T00:00:00Z format.
func ConvertToString(t time.Time) string {
	return t.Format("2006-01-02T15:04:05Z")
}
