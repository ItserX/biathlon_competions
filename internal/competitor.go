package biathlon

import "time"

type Lap struct {
	Time     time.Time
	AvgSpeed float64
}

type Competitor struct {
	Status          string
	SetStartTime    *time.Time
	ActualStartTime *time.Time
	TotalTime       *time.Time
	Laps            []*Lap
	PenaltyTime     *time.Time
	Hits            int
	Shots           int
	LastEventID     EventType
}
