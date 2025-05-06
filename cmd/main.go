package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ItserX/biathlon_competions/internal/config"
	"github.com/ItserX/biathlon_competions/internal/events"
	"github.com/ItserX/biathlon_competions/internal/report"
)

func ParseFlags() (string, string, error) {
	eventsPath := flag.String("events", "", "File with incoming events")
	configPath := flag.String("config", "", "File with configuration")
	flag.Parse()

	if *eventsPath == "" || *configPath == "" {
		flag.Usage()
		return "", "", fmt.Errorf("Flags not specified")
	}

	return *eventsPath, *configPath, nil

}

func ProcessEventsFile(rp *events.RaceProcessor, eventsPath string) error {
	file, err := os.Open(eventsPath)
	if err != nil {
		return fmt.Errorf("Error open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		event, err := events.ParseIncomingEvent(line)
		if err != nil {
			return fmt.Errorf("Error in string %d: %v\n", lineNumber, err)
		}
		rp.ProcessEvent(event)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error read file: %v\n", err)
	}

	return nil
}

func main() {

	eventsPath, configPath, err := ParseFlags()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cfg, _ := config.ParseConfig(configPath)
	rp := events.RaceProcessor{Config: cfg, Competitors: make(map[int]*events.Competitor)}

	err = ProcessEventsFile(&rp, eventsPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	rep := report.GenerateReport(&rp)
	fmt.Print(rep)
}
