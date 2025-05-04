package biathlon

import (
	"errors"
	"fmt"
	"time"
)

type RaceProcessor struct {
	Config      *Config
	Competitors map[int]*Competitor
}

var (
	ErrNotRegistered     = errors.New("The competitor not registered")
	ErrAlreadyRegistered = errors.New("The competitor already registered")
	ErrStartTimeNotSet   = errors.New("Start time not set")
	ErrNotInStartLine    = errors.New("The competitor not in the start line")
	ErrNotStarted        = errors.New("The competitor not started")
	ErrNotContinue       = errors.New("The competitor cannot continue the race")
)

func (rp *RaceProcessor) ProcessEvent(event *Event) error {
	if err := rp.validateEvent(event); err != nil {
		return err
	}

	switch event.ID {
	case CompetitorRegistered:
		return rp.processRegistration(event)
	case StartTimeSet:
		return rp.processStartTimeSet(event)
	case OnStartLine:
		return rp.processOnStartLine(event)
	case CompetitorStarted:
		return rp.processCompetitorStarted(event)
	case OnFiringRange:
		return rp.processOnFiringRange(event)
	default:
		return fmt.Errorf("unknown event type: %d", event.ID)
	}
}

func (rp *RaceProcessor) validateEvent(event *Event) error {
	competitor := rp.Competitors[event.CompetitorID]

	if competitor == nil && event.ID != CompetitorRegistered {
		return ErrNotRegistered
	}

	if competitor != nil && competitor.Status != "" {
		return ErrNotContinue
	}

	return nil
}

func (rp *RaceProcessor) processRegistration(event *Event) error {
	if rp.Competitors[event.CompetitorID] != nil {
		return ErrAlreadyRegistered
	}

	rp.Competitors[event.CompetitorID] = &Competitor{
		LastEventID: CompetitorRegistered,
	}
	fmt.Printf("The competitor(%d) registered\n", event.CompetitorID)
	return nil
}

func (rp *RaceProcessor) processStartTimeSet(event *Event) error {
	competitor := rp.Competitors[event.CompetitorID]
	competitor.SetStartTime = event.StartTime
	competitor.LastEventID = StartTimeSet
	fmt.Printf("The start time for the competitor(%d) was set by a draw to %s\n",
		event.CompetitorID, competitor.SetStartTime)
	return nil
}

func (rp *RaceProcessor) processOnStartLine(event *Event) error {
	competitor := rp.Competitors[event.CompetitorID]

	if err := rp.checkEventSequence(competitor.LastEventID, event.ID, ErrStartTimeNotSet); err != nil {
		return err
	}

	if rp.isTooEarlyToStart(event.CurrentTime, competitor.SetStartTime) {
		competitor.Status = "NotStarted"
		return fmt.Errorf("Competitor(%d) has not started\n", event.CompetitorID)

	}

	competitor.LastEventID = OnStartLine
	fmt.Printf("The competitor(%d) is on the start line\n", event.CompetitorID)
	return nil
}

func (rp *RaceProcessor) processCompetitorStarted(event *Event) error {
	competitor := rp.Competitors[event.CompetitorID]

	if err := rp.checkEventSequence(competitor.LastEventID, event.ID, ErrNotInStartLine); err != nil {
		return err
	}

	if rp.isTooEarlyToStart(event.CurrentTime, competitor.SetStartTime) {
		competitor.Status = "NotStarted"
		return fmt.Errorf("Competitor(%d) has not started\n", event.CompetitorID)
	}

	competitor.ActualStartTime = event.CurrentTime
	competitor.LastEventID = CompetitorStarted
	fmt.Printf("The competitor(%d) has started\n", event.CompetitorID)
	return nil
}

func (rp *RaceProcessor) processOnFiringRange(event *Event) error {
	competitor := rp.Competitors[event.CompetitorID]

	if err := rp.checkEventSequence(competitor.LastEventID, event.ID, ErrNotInStartLine); err != nil {
		return err
	}
	if event.FiringRange > rp.Config.FiringLines {
		return fmt.Errorf("There is no such firing range\n")
	}
	competitor.LastEventID = OnFiringRange
	fmt.Printf("[%s] The competitor(%d) is on the firing range(%d)", event.CurrentTime, event.CompetitorID, event.FiringRange)
	return nil
}

func (rp *RaceProcessor) checkEventSequence(lastEventID, currentEventID EventType, err error) error {
	if currentEventID-1 != lastEventID {
		return err
	}
	return nil
}

func (rp *RaceProcessor) isTooEarlyToStart(currentTime, startTime *time.Time) bool {
	return currentTime.Before(startTime.Add(rp.Config.StartDelta.Sub(time.Time{})))
}
