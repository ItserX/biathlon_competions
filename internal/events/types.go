package events

import (
	"time"

	"github.com/ItserX/biathlon_competions/internal/config"
	"github.com/ItserX/biathlon_competions/internal/constants"
)

type Lap struct {
	StartTime time.Time
	TotalTime time.Time
	AvgSpeed  float64
}

type Competitor struct {
	Status        string
	SetStartTime  time.Time
	TotalTime     time.Time
	NumCurrentLap int
	LapsCompleted int
	Laps          []*Lap
	PenaltyLaps   []*Lap
	Hits          int
	Shots         int
}

type RaceProcessor struct {
	Config      *config.Config
	Competitors map[int]*Competitor
}

type Event struct {
	CurrentTime  time.Time
	ID           constants.EventType
	CompetitorID int
	StartTime    time.Time
	FiringRange  int
	Target       int
	Comment      string
}
