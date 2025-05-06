package events

import (
	"fmt"
	"time"
)

func FormatTime(t time.Time) string {
	hour, min, sec := t.Clock()
	millis := t.Nanosecond() / 1e6
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hour, min, sec, millis)
}

func (rp *RaceProcessor) isTooLaterToStart(currentTime, startTime time.Time) bool {
	return currentTime.Before(startTime.Add(rp.Config.StartDelta.Sub(time.Time{})))
}

func (rp *RaceProcessor) calculateTotalTime(startTime, finishTime time.Time) time.Time {
	totalTime := time.Time{}.Add(finishTime.Sub(startTime))
	return totalTime
}

func (rp *RaceProcessor) calculateAvgSpeed(len int, totalTime time.Time) float64 {
	return float64(len) / (float64(totalTime.Hour()*3600+totalTime.Minute()*60+totalTime.Second()) +
		float64(totalTime.Nanosecond())/1e9)

}
