package events

import (
	"testing"
	"time"

	"github.com/ItserX/biathlon_competions/internal/config"
	"github.com/ItserX/biathlon_competions/internal/constants"
)

func TestRaceProcessor(t *testing.T) {
	testConfig := &config.Config{
		Laps:        2,
		LapLen:      3651,
		PenaltyLen:  50,
		FiringLines: 1,
	}

	parseTime := func(timeStr string) time.Time {
		tm, err := time.Parse("15:04:05.000", timeStr)
		if err != nil {
			t.Fatalf("Failed to parse time: %v", err)
		}
		return tm
	}

	t.Run("Full race process", func(t *testing.T) {
		rp := &RaceProcessor{
			Config:      testConfig,
			Competitors: make(map[int]*Competitor),
		}

		event := &Event{
			CurrentTime:  parseTime("09:05:59.867"),
			ID:           constants.CompetitorRegistered,
			CompetitorID: 1,
		}
		rp.ProcessEvent(event)
		if rp.Competitors[1] == nil {
			t.Error("Expected competitor to be registered")
		}
		if len(rp.Competitors[1].Laps) != 2 {
			t.Errorf("Expected 2 laps, got %d", len(rp.Competitors[1].Laps))
		}

		event = &Event{
			CurrentTime:  parseTime("09:15:00.841"),
			ID:           constants.StartTimeSet,
			CompetitorID: 1,
			StartTime:    parseTime("09:30:00.000"),
		}
		rp.ProcessEvent(event)
		if !rp.Competitors[1].SetStartTime.Equal(parseTime("09:30:00.000")) {
			t.Error("Start time not set correctly")
		}

		event = &Event{
			CurrentTime:  parseTime("09:29:45.734"),
			ID:           constants.OnStartLine,
			CompetitorID: 1,
		}
		rp.ProcessEvent(event)

		event = &Event{
			CurrentTime:  parseTime("09:30:01.005"),
			ID:           constants.CompetitorStarted,
			CompetitorID: 1,
		}
		rp.ProcessEvent(event)
		if rp.Competitors[1].NumCurrentLap != 1 {
			t.Errorf("Expected current lap 1, got %d", rp.Competitors[1].NumCurrentLap)
		}
		if rp.Competitors[1].Laps[0] == nil {
			t.Error("First lap not initialized")
		}

		event = &Event{
			CurrentTime:  parseTime("09:30:01.005"),
			ID:           constants.OnFiringRange,
			CompetitorID: 1,
			FiringRange:  1,
		}
		rp.ProcessEvent(event)

		targets := []int{1, 2, 4, 5}
		for i, target := range targets {
			event = &Event{
				CurrentTime:  parseTime("09:49:33.123").Add(time.Duration(i) * time.Second),
				ID:           constants.TargetHit,
				CompetitorID: 1,
				Target:       target,
			}
			rp.ProcessEvent(event)
			if rp.Competitors[1].Hits != i+1 {
				t.Errorf("Expected %d hits, got %d", i+1, rp.Competitors[1].Hits)
			}
		}

		event = &Event{
			CurrentTime:  parseTime("09:49:38.339"),
			ID:           constants.LeftFiringRange,
			CompetitorID: 1,
		}
		rp.ProcessEvent(event)
		if rp.Competitors[1].Shots != constants.TargetsCount {
			t.Errorf("Expected %d shots, got %d", constants.TargetsCount, rp.Competitors[1].Shots)
		}

		event = &Event{
			CurrentTime:  parseTime("09:49:55.915"),
			ID:           constants.EnteredPenaltyLaps,
			CompetitorID: 1,
		}
		rp.ProcessEvent(event)
		if len(rp.Competitors[1].PenaltyLaps) != 1 {
			t.Errorf("Expected 1 penalty lap, got %d", len(rp.Competitors[1].PenaltyLaps))
		}

		event = &Event{
			CurrentTime:  parseTime("09:51:48.391"),
			ID:           constants.LeftPenaltyLaps,
			CompetitorID: 1,
		}
		rp.ProcessEvent(event)
		if rp.Competitors[1].PenaltyLaps[0].TotalTime.IsZero() {
			t.Error("Penalty lap total time not set")
		}
		if rp.Competitors[1].PenaltyLaps[0].AvgSpeed == 0 {
			t.Error("Penalty lap avg speed not calculated")
		}

		event = &Event{
			CurrentTime:  parseTime("09:59:03.872"),
			ID:           constants.EndedMainLap,
			CompetitorID: 1,
		}
		rp.ProcessEvent(event)
		if rp.Competitors[1].LapsCompleted != 1 {
			t.Errorf("Expected 1 lap completed, got %d", rp.Competitors[1].LapsCompleted)
		}
		if rp.Competitors[1].Laps[1] == nil {
			t.Error("Second lap not initialized")
		}

		event = &Event{
			CurrentTime:  parseTime("09:59:05.321"),
			ID:           constants.CannotContinue,
			CompetitorID: 1,
			Comment:      "Lost in the forest",
		}
		rp.ProcessEvent(event)
		if rp.Competitors[1].Status != constants.StatusNotFinished {
			t.Errorf("Expected status %s, got %s", constants.StatusNotFinished, rp.Competitors[1].Status)
		}
	})
}
