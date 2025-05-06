package constants

type EventType int

const (
	CompetitorRegistered EventType = iota + 1
	StartTimeSet
	OnStartLine
	CompetitorStarted
	OnFiringRange
	TargetHit
	LeftFiringRange
	EnteredPenaltyLaps
	LeftPenaltyLaps
	EndedMainLap
	CannotContinue
)

const (
	StatusNotFinished string = "NotFinished"
	StatusNotStarted  string = "NotStarted"
	TimeLayoutConfig  string = "15:04:05"
	TimeLayoutEvent   string = "15:04:05.000"
)

const TargetsCount = 5
