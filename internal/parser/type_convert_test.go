package parser

import (
	"testing"
	"time"
)

// TestConvertToString tests the ConvertToString function.
func TestConvertToString(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			name: "Standard Time",
			time: time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC),
			want: "2025-03-10T00:00:00Z",
		},
		{
			name: "Leap Year",
			time: time.Date(2020, 2, 29, 12, 30, 45, 0, time.UTC),
			want: "2020-02-29T12:30:45Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToString(tt.time); got != tt.want {
				t.Errorf("ConvertToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConvertToTime tests the ConvertToTime function.
func TestConvertToTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:  "Valid Time",
			input: "2025-03-10T00:00:00Z",
			want:  time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "Invalid Time Format",
			input:   "2025-03-10T00:00:00",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ConvertToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
