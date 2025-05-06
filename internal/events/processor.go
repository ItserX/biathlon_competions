package events

import (
	"fmt"
	"time"

	"github.com/ItserX/biathlon_competions/internal/constants"
)

func (rp *RaceProcessor) ProcessEvent(event *Event) {
	switch event.ID {
	case constants.CompetitorRegistered:
		rp.processRegistration(event)
	case constants.StartTimeSet:
		rp.processStartTimeSet(event)
	case constants.OnStartLine:
		rp.processOnStartLine(event)
	case constants.CompetitorStarted:
		rp.processCompetitorStarted(event)
	case constants.OnFiringRange:
		rp.processOnFiringRange(event)
	case constants.TargetHit:
		rp.processTargetHit(event)
	case constants.LeftFiringRange:
		rp.processLeftTheFiringRange(event)
	case constants.EnteredPenaltyLaps:
		rp.processEnteredPenaltyLaps(event)
	case constants.LeftPenaltyLaps:
		rp.processLeftPenaltyLaps(event)
	case constants.EndedMainLap:
		rp.processEndedMainLap(event)
	case constants.CannotContinue:
		rp.processCannotContinue(event)
	}
}

func (rp *RaceProcessor) processRegistration(event *Event) {
	rp.Competitors[event.CompetitorID] = &Competitor{Laps: make([]*Lap, rp.Config.Laps), PenaltyLaps: make([]*Lap, 0)}
	fmt.Printf("[%s] The competitor(%d) registered\n", FormatTime(event.CurrentTime), event.CompetitorID)
}

func (rp *RaceProcessor) processStartTimeSet(event *Event) {
	competitor := rp.Competitors[event.CompetitorID]
	competitor.SetStartTime = event.StartTime

	fmt.Printf("[%s] The start time for the competitor(%d) was set by a draw to %s\n",
		FormatTime(event.CurrentTime), event.CompetitorID, event.StartTime)
}

func (rp *RaceProcessor) processOnStartLine(event *Event) {
	competitor := rp.Competitors[event.CompetitorID]

	if rp.isTooLaterToStart(event.CurrentTime, competitor.SetStartTime) {
		competitor.Status = constants.StatusNotStarted
		fmt.Printf("[%s] The competitor(%d) is disqualified\n", FormatTime(event.CurrentTime), event.CompetitorID)
		return
	}

	fmt.Printf("[%s] The competitor(%d) is on the start line\n", FormatTime(event.CurrentTime), event.CompetitorID)
}

func (rp *RaceProcessor) processCompetitorStarted(event *Event) {
	competitor := rp.Competitors[event.CompetitorID]

	if rp.isTooLaterToStart(event.CurrentTime, competitor.SetStartTime) {
		competitor.Status = constants.StatusNotStarted
		fmt.Printf("[%s] The competitor(%d) is disqualified\n", FormatTime(event.CurrentTime), event.CompetitorID)
		return
	}

	competitor.NumCurrentLap = 1
	competitor.Laps[competitor.NumCurrentLap-1] = &Lap{StartTime: competitor.SetStartTime}

	fmt.Printf("[%s] The competitor(%d) has started\n", FormatTime(event.CurrentTime), event.CompetitorID)
}

func (rp *RaceProcessor) processOnFiringRange(event *Event) {
	fmt.Printf("[%s] The competitor(%d) is on the firing range(%d)\n",
		FormatTime(event.CurrentTime), event.CompetitorID, event.FiringRange)
}

func (rp *RaceProcessor) processTargetHit(event *Event) {
	competitor := rp.Competitors[event.CompetitorID]
	competitor.Hits += 1

	fmt.Printf("[%s] The target(%d) has been hit by competitor(%d)\n", FormatTime(event.CurrentTime), event.Target, event.CompetitorID)
}

func (rp *RaceProcessor) processLeftTheFiringRange(event *Event) {
	competitor := rp.Competitors[event.CompetitorID]

	competitor.Shots += constants.TargetsCount

	fmt.Printf("[%s] The competitor(%d) left the firing range\n", FormatTime(event.CurrentTime), event.CompetitorID)
}

func (rp *RaceProcessor) processEnteredPenaltyLaps(event *Event) {
	competitor := rp.Competitors[event.CompetitorID]

	penaltyLap := Lap{StartTime: event.CurrentTime}
	competitor.PenaltyLaps = append(competitor.PenaltyLaps, &penaltyLap)

	fmt.Printf("[%s] The competitor(%d) entered the penalty laps\n", FormatTime(event.CurrentTime), event.CompetitorID)
}

func (rp *RaceProcessor) processLeftPenaltyLaps(event *Event) {
	competitor := rp.Competitors[event.CompetitorID]

	currentPenaltyLap := competitor.PenaltyLaps[len(competitor.PenaltyLaps)-1]
	currentPenaltyLap.TotalTime = rp.calculateTotalTime(currentPenaltyLap.StartTime, event.CurrentTime)
	currentPenaltyLap.AvgSpeed = rp.calculateAvgSpeed(rp.Config.PenaltyLen, currentPenaltyLap.TotalTime)

	fmt.Printf("[%s] The competitor(%d) left the penalty laps\n", FormatTime(event.CurrentTime), event.CompetitorID)
}

func (rp *RaceProcessor) processEndedMainLap(event *Event) {
	competitor := rp.Competitors[event.CompetitorID]

	currentLap := competitor.Laps[competitor.NumCurrentLap-1]
	currentLap.TotalTime = rp.calculateTotalTime(currentLap.StartTime, event.CurrentTime)
	currentLap.AvgSpeed = rp.calculateAvgSpeed(rp.Config.LapLen, currentLap.TotalTime)
	competitor.TotalTime = competitor.TotalTime.Add(currentLap.TotalTime.Sub(time.Time{}))
	competitor.NumCurrentLap += 1
	competitor.LapsCompleted += 1

	if competitor.LapsCompleted >= 1 && competitor.LapsCompleted < rp.Config.Laps {
		competitor.Laps[competitor.NumCurrentLap-1] = &Lap{StartTime: event.CurrentTime}
	}

	fmt.Printf("[%s] The competitor(%d) ended the main lap\n", FormatTime(event.CurrentTime), event.CompetitorID)

	if competitor.LapsCompleted == rp.Config.Laps {
		fmt.Printf("[%s] The competitor(%d) has finished\n", FormatTime(event.CurrentTime), event.CompetitorID)
	}

}

func (rp *RaceProcessor) processCannotContinue(event *Event) {
	competitor := rp.Competitors[event.CompetitorID]
	competitor.Status = constants.StatusNotFinished
	fmt.Printf("[%s] The competitor(%d) can`t continue: %s\n", FormatTime(event.CurrentTime), event.CompetitorID, event.Comment)
}
