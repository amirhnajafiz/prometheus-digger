package parser

import (
	"testing"
	"time"
)

// TestConvertSliceToTime tests the ConvertSliceToTime function with various inputs.
// It checks if the function correctly adjusts the base time according to the input parameters.
func TestConvertSliceToTime(t *testing.T) {
	baseTime := time.Now()

	tests := []struct {
		slice    string
		expected time.Time
	}{
		{"+1:02:03:04", baseTime.AddDate(0, 0, 1).Add(2*time.Hour + 3*time.Minute + 4*time.Second)},
		{"-1:02:03:04", baseTime.AddDate(0, 0, -1).Add(-2*time.Hour - 3*time.Minute - 4*time.Second)},
		{"+0:01:02:03", baseTime.Add(1*time.Hour + 2*time.Minute + 3*time.Second)},
		{"-0:01:02:03", baseTime.Add(-1*time.Hour - 2*time.Minute - 3*time.Second)},
		{"+0:00:00:00", baseTime},
		{"-0:00:00:00", baseTime},
		{"+2:00:00:00", baseTime.AddDate(0, 0, 2)},
		{"-2:00:00:00", baseTime.AddDate(0, 0, -2)},
		{"+0:00:00:01", baseTime.Add(time.Second)},
		{"-0:00:00:01", baseTime.Add(-time.Second)},
		{"+0:00:01:00", baseTime.Add(time.Minute)},
		{"-0:00:01:00", baseTime.Add(-time.Minute)},
		{"+0:01:00:00", baseTime.Add(1 * time.Hour)},
		{"-0:01:00:00", baseTime.Add(-1 * time.Hour)},
		{"+1:00:00:00", baseTime.AddDate(0, 0, 1)},
		{"-1:00:00:00", baseTime.AddDate(0, 0, -1)},
	}

	for _, test := range tests {
		result := ConvertSliceToTime(baseTime, test.slice)
		if !result.Equal(test.expected) {
			t.Errorf("ConvertSliceToTime(%q) = %v; want %v", test.slice, result, test.expected)
		}
	}
}
