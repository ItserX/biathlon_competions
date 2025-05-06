package report

import (
	"testing"
	"time"

	"github.com/ItserX/biathlon_competions/internal/config"
	"github.com/ItserX/biathlon_competions/internal/events"
)

func TestGenerateReport(t *testing.T) {
	rp := events.RaceProcessor{
		Config: &config.Config{Laps: 2,
			LapLen:      3651,
			PenaltyLen:  50,
			FiringLines: 1},
		Competitors: map[int]*events.Competitor{1: {
			Status:        "NotFinished",
			SetStartTime:  time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC),
			TotalTime:     time.Date(1, 1, 1, 0, 29, 3, 872000000, time.UTC),
			NumCurrentLap: 2,
			LapsCompleted: 1,
			Laps: []*events.Lap{
				{
					StartTime: time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC),
					TotalTime: time.Date(1, 1, 1, 0, 29, 3, 872000000, time.UTC),
					AvgSpeed:  2.0936169627128596,
				},
				{
					StartTime: time.Date(0, 1, 1, 9, 59, 3, 872000000, time.UTC),
					TotalTime: time.Time{},
					AvgSpeed:  0,
				},
			},
			PenaltyLaps: []*events.Lap{
				{
					StartTime: time.Date(0, 1, 1, 9, 49, 55, 915000000, time.UTC),
					TotalTime: time.Date(1, 1, 1, 0, 1, 52, 476000000, time.UTC),
					AvgSpeed:  0.4445392794907358,
				},
			},
			Hits:  4,
			Shots: 5,
		},
		},
	}

	output := GenerateReport(&rp)
	if output != "[NotFinished] 1 [{00:29:03.872, 2.093}, {,}] {00:01:52.476, 0.444} 4/5\n" {
		t.Errorf("Incorrect report")
	}
}
