package report

import (
	"fmt"
	"sort"

	"github.com/ItserX/biathlon_competions/internal/events"
)

func GenerateReport(rp *events.RaceProcessor) string {
	completeString := ""
	ids := make([]int, 0)
	for id := range rp.Competitors {
		ids = append(ids, id)
	}

	sort.Slice(ids, func(i, j int) bool {
		compI := rp.Competitors[ids[i]]
		compJ := rp.Competitors[ids[j]]
		return compI.TotalTime.Before(compJ.TotalTime)
	})

	for _, id := range ids {
		val := rp.Competitors[id]

		var fTime string
		if val.Status == "" {
			fTime = events.FormatTime(rp.Competitors[id].TotalTime)
		} else {
			fTime = val.Status
		}

		fLaps := FormatLaps(val.Laps)

		fPenaltyLaps := FormatLaps(val.PenaltyLaps)

		completeString += fmt.Sprintf("[%s] %d %s %s %d/%d\n", fTime, id, fLaps, fPenaltyLaps, val.Hits, val.Shots)
	}
	return completeString

}

func FormatLaps(laps []*events.Lap) string {
	var completeString string
	for i, val := range laps {
		if val.AvgSpeed == 0 {
			completeString += "{,}"
			continue
		}
		if i < len(laps)-1 {
			completeString += fmt.Sprintf("{%s, %.3f}, ", events.FormatTime(val.TotalTime), float64(int(val.AvgSpeed*1000))/1000)
		} else {
			completeString += fmt.Sprintf("{%s, %.3f}", events.FormatTime(val.TotalTime), float64(int(val.AvgSpeed*1000))/1000)
		}
	}

	if len(laps) > 1 {
		return fmt.Sprintf("[%s]", completeString)
	}
	return completeString
}
