package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ItserX/biathlon_competions/internal/constants"
)

type flexibleTime struct {
	time.Time
}

type Config struct {
	Laps        int       `json:"laps"`
	LapLen      int       `json:"lapLen"`
	PenaltyLen  int       `json:"penaltyLen"`
	FiringLines int       `json:"firingLines"`
	Start       time.Time `json:"start"`
	StartDelta  time.Time `json:"startDelta"`
}

func (ft *flexibleTime) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}

	t, err := time.Parse(constants.TimeLayoutConfig, timeStr)
	if err == nil {
		ft.Time = t
		return nil
	}

	return fmt.Errorf("invalid time format: %s, expected HH:MM:SS", timeStr)
}

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var tempConfig struct {
		Laps        int           `json:"laps"`
		LapLen      int           `json:"lapLen"`
		PenaltyLen  int           `json:"penaltyLen"`
		FiringLines int           `json:"firingLines"`
		Start       *flexibleTime `json:"start"`
		StartDelta  *flexibleTime `json:"startDelta"`
	}

	if err := json.NewDecoder(file).Decode(&tempConfig); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	config := &Config{
		Laps:        tempConfig.Laps,
		LapLen:      tempConfig.LapLen,
		PenaltyLen:  tempConfig.PenaltyLen,
		FiringLines: tempConfig.FiringLines,
	}

	if tempConfig.Start != nil {
		config.Start = tempConfig.Start.Time
	}
	if tempConfig.StartDelta != nil {
		config.StartDelta = tempConfig.StartDelta.Time
	}

	return config, nil
}
