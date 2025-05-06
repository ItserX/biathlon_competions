package config

import (
	"os"
	"testing"
	"time"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name         string
		configFile   string
		expected     *Config
		expectStatus bool
	}{
		{
			name:       "valid config",
			configFile: "testdata/valid_config.json",
			expected: &Config{
				Laps:        3,
				LapLen:      4000,
				PenaltyLen:  150,
				FiringLines: 5,
				Start:       time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
				StartDelta:  time.Date(0, 1, 1, 0, 1, 0, 0, time.UTC),
			},
			expectStatus: false,
		},
		{
			name:         "invalid file path",
			configFile:   "testdata/nonexistent.json",
			expectStatus: true,
		},
		{
			name:         "invalid json format",
			configFile:   "testdata/invalid_json.json",
			expectStatus: true,
		},
		{
			name:         "invalid time format",
			configFile:   "testdata/invalid_time.json",
			expectStatus: true,
		},
		{
			name:       "missing optional fields",
			configFile: "testdata/minimal_config.json",
			expected: &Config{
				Laps:        3,
				LapLen:      4000,
				PenaltyLen:  150,
				FiringLines: 5,
				Start:       time.Time{},
				StartDelta:  time.Time{},
			},
			expectStatus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ParseConfig(tt.configFile)

			if tt.expectStatus {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if cfg.Laps != tt.expected.Laps {
				t.Errorf("Expected Laps %d, got %d", tt.expected.Laps, cfg.Laps)
			}
			if cfg.LapLen != tt.expected.LapLen {
				t.Errorf("Expected LapLen %d, got %d", tt.expected.LapLen, cfg.LapLen)
			}
			if cfg.PenaltyLen != tt.expected.PenaltyLen {
				t.Errorf("Expected PenaltyLen %d, got %d", tt.expected.PenaltyLen, cfg.PenaltyLen)
			}
			if cfg.FiringLines != tt.expected.FiringLines {
				t.Errorf("Expected FiringLines %d, got %d", tt.expected.FiringLines, cfg.FiringLines)
			}
			if !cfg.Start.Equal(tt.expected.Start) {
				t.Errorf("Expected Start %v, got %v", tt.expected.Start, cfg.Start)
			}
			if !cfg.StartDelta.Equal(tt.expected.StartDelta) {
				t.Errorf("Expected StartDelta %v, got %v", tt.expected.StartDelta, cfg.StartDelta)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		expected     time.Time
		expectStatus bool
	}{
		{
			name:     "valid time",
			input:    []byte(`"10:00:00"`),
			expected: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
		},
		{
			name:         "invalid format",
			input:        []byte(`"10-00-00"`),
			expectStatus: true,
		},
		{
			name:         "not a string",
			input:        []byte(`12345`),
			expectStatus: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ft flexibleTime
			err := ft.UnmarshalJSON(tt.input)

			if tt.expectStatus {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !ft.Time.Equal(tt.expected) {
				t.Errorf("Expected time %v, got %v", tt.expected, ft.Time)
			}
		})
	}
}

func TestMain(m *testing.M) {
	writeTestFile("testdata/valid_config.json", `{
		"laps": 3,
		"lapLen": 4000,
		"penaltyLen": 150,
		"firingLines": 5,
		"start": "10:00:00",
		"startDelta": "00:01:00"
	}`)

	writeTestFile("testdata/invalid_json.json", `{
		"laps": 3,
		"lapLen": 4000,
		"penaltyLen": 150,
		"firingLines": 5,
		"start": "10:00:00",
		"startDelta": "00:01:00",
	}`)

	writeTestFile("testdata/invalid_time.json", `{
		"laps": 3,
		"lapLen": 4000,
		"penaltyLen": 150,
		"firingLines": 5,
		"start": "10:00:99",
		"startDelta": "00:01:00"
	}`)

	writeTestFile("testdata/minimal_config.json", `{
		"laps": 3,
		"lapLen": 4000,
		"penaltyLen": 150,
		"firingLines": 5
	}`)

	code := m.Run()

	os.RemoveAll("testdata")
	os.Exit(code)
}

func writeTestFile(path, content string) {
	os.MkdirAll("testdata", 0755)
	os.WriteFile(path, []byte(content), 0644)
}
