package main

import (
	biathlontracker "github.com/ItserX/biathlon_competions/internal"
)

func main() {
	ev1, _ := biathlontracker.ParseIncomingEvent("[09:05:59.867] 1 1")
	ev2, _ := biathlontracker.ParseIncomingEvent("[09:15:00.841] 2 1 09:30:00.000")
	ev3, _ := biathlontracker.ParseIncomingEvent("[09:29:45.734] 3 1")
	ev4, _ := biathlontracker.ParseIncomingEvent("[09:30:01.005] 4 1")
	ev5, _ := biathlontracker.ParseIncomingEvent("[09:49:31.659] 5 1 1")

	ms := []*biathlontracker.Event{ev1, ev2, ev3, ev4, ev5}
	cfg, _ := biathlontracker.ParseConfig("/home/lucky7788/Загрузки/yadro/sunny_5_skiers/config.json")
	rp := biathlontracker.RaceProcessor{Config: cfg, Competitors: make(map[int]*biathlontracker.Competitor)}
	for _, val := range ms {
		rp.ProcessEvent(val)
	}
}
