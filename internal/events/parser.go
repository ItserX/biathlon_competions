package events

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ItserX/biathlon_competions/internal/constants"
)

func ParseIncomingEvent(eventString string) (*Event, error) {
	params := strings.Split(eventString, " ")

	if len(params) < 3 {
		return nil, fmt.Errorf("Incorrect eventString")
	}

	currentTime, err := parseTime(params[0])
	if err != nil {

		return nil, err
	}

	eventID, err := strconv.Atoi(params[1])
	if err != nil {
		return nil, err
	}
	eventType := constants.EventType(eventID)

	competitorID, err := strconv.Atoi(params[2])
	if err != nil {
		return nil, err
	}

	event := Event{ID: eventType, CurrentTime: currentTime, CompetitorID: competitorID}

	if len(params) > 3 {
		switch eventType {
		case constants.StartTimeSet:
			startTime, err := parseTime(params[3])
			if err != nil {
				return nil, err
			}
			event.StartTime = startTime

		case constants.OnFiringRange:
			firingRange, err := strconv.Atoi(params[3])
			if err != nil {
				return nil, err
			}
			event.FiringRange = firingRange

		case constants.TargetHit:
			target, err := strconv.Atoi(params[3])
			if err != nil {
				return nil, err
			}
			event.Target = target

		case constants.CannotContinue:
			event.Comment = strings.Join(params[3:], " ")
		default:
			return nil, fmt.Errorf("Incorrect eventID")
		}
	}

	return &event, nil
}

func parseTime(timeStr string) (time.Time, error) {
	cleanStr := strings.Trim(timeStr, "[]")

	parsedTime, err := time.Parse(constants.TimeLayoutEvent, cleanStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
