package events

import (
	"testing"
)

func TestParseIncomingEvent(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Invalid time format",
			input: "[12:34:56] 1 42",
		},
		{
			name:  "Missing event ID",
			input: "[12:34:56.789]",
		},
		{
			name:  "Invalid event ID",
			input: "[12:34:56.789] abc 42",
		},
		{
			name:  "Missing competitor ID",
			input: "[12:34:56.789] 1",
		},
		{
			name:  "Invalid competitor ID",
			input: "[12:34:56.789] 1 abc",
		},
		{
			name:  "Invalid firing range format",
			input: "[09:15:00.841] 5 1 abc",
		},
		{
			name:  "Invalid target format",
			input: "[09:49:33.123] 6 1 abc",
		},
		{
			name:  "Invalid start time format",
			input: "[12:34:56.789] 2 1 [12:34:56]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseIncomingEvent(tt.input)
			if err == nil {
				t.Errorf("Expected error for input: %q", tt.input)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Empty time string",
			input: "",
		},
		{
			name:  "Malformed time string",
			input: "12:34:56",
		},
		{
			name:  "Invalid time components",
			input: "[25:61:61.999]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseTime(tt.input)
			if err == nil {
				t.Errorf("Expected error for input: %q", tt.input)
			}
		})
	}
}
