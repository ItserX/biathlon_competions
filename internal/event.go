package biathlon

import (
	"strconv"
	"strings"
	"time"
)

type Event struct {
	CurrentTime  *time.Time
	ID           EventType
	CompetitorID int
	StartTime    *time.Time
	FiringRange  int
	Target       int
	Comment      string
}

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

func ParseIncomingEvent(eventString string) (*Event, error) {
	params := strings.Split(eventString, " ")

	currentTime, err := parseTime(params[0])
	if err != nil {

		return nil, err
	}

	eventID, err := strconv.Atoi(params[1])
	if err != nil {
		return nil, err
	}
	eventType := EventType(eventID)

	competitorID, err := strconv.Atoi(params[2])
	if err != nil {
		return nil, err
	}

	event := Event{ID: eventType, CurrentTime: &currentTime, CompetitorID: competitorID}

	if len(params) > 3 {
		switch eventType {
		case StartTimeSet:
			startTime, err := parseTime(params[3])
			if err != nil {
				return nil, err
			}
			event.StartTime = &startTime

		case OnFiringRange:
			firingRange, err := strconv.Atoi(params[3])
			if err != nil {
				return nil, err
			}
			event.FiringRange = firingRange

		case TargetHit:
			target, err := strconv.Atoi(params[3])
			if err != nil {
				return nil, err
			}
			event.Target = target

		case CannotContinue:
			event.Comment = strings.Join(params[3:], " ")
		}
	}

	return &event, nil
}

func parseTime(timeStr string) (time.Time, error) {
	cleanStr := strings.Trim(timeStr, "[]")
	parsedTime, err := time.Parse("15:04:05.000", cleanStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
