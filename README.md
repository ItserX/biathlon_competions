# System prototype for biathlon competitions
The prototype must be able to work with a configuration file and a set of external events of a certain format.

## Configuration (json)

- **Laps**        - Amount of laps for main distance
- **LapLen**      - Length of each main lap
- **PenaltyLen**  - Length of each penalty lap
- **FiringLines** - Number of firing lines per lap
- **Start**       - Planned start time for the first competitor
- **StartDelta**  - Planned interval between starts

## Events
All events are characterized by time and event identifier. Outgoing events are events created during program operation. Events related to the "incoming" category cannot be generated and are output in the same form as they were submitted in the input file.
```
Incoming events
EventID | extraParams | Comments
1       |             | The competitor registered
2       | startTime   | The start time was set by a draw
3       |             | The competitor is on the start line
4       |             | The competitor has started
5       | firingRange | The competitor is on the firing range
6       | target      | The target has been hit
7       |             | The competitor left the firing range
8       |             | The competitor entered the penalty laps
9       |             | The competitor left the penalty laps
10      |             | The competitor ended the main lap
11      | comment     | The competitor can`t continue
```
An competitor is disqualified if he/she does not start during his/her start interval. This marked as **NotStarted** in final report.
If the competitor can`t continue it should be marked in final report as **NotFinished**

```
Outgoing events
EventID | extraParams | Comments
32      |             | The competitor is disqualified
33      |             | The competitor has finished
```

## Final report
The final report should contain the list of all registered competitors
sorted by ascending time.
- Total time includes the difference between scheduled and actual start time or **NotStarted**/**NotFinished** marks
- Time taken to complete each lap
- Average speed for each lap [m/s]
- Time taken to complete penalty laps
- Average speed over penalty laps [m/s]
- Number of hits/number of shots

## Build and run
### Run directly
```bash
cd biathlon_competitions
go run cmd/main.go -events=<path-to-events-file> -config=<path-to-config-file>
```
### Build and run executable
```bash
cd biathlon_competions
go build -o bin/app ./cmd/main.go
./bin/app -events=<path-to-events-file> -config=<path-to-config-file> 
```

## Tests
```bash 
go test -cover ./...
```
### Coverage
```
        github.com/ItserX/biathlon_competions/cmd               coverage: 0.0% of statements
ok      github.com/ItserX/biathlon_competions/internal/config   0.002s  coverage: 100.0% of statements
?       github.com/ItserX/biathlon_competions/internal/constants        [no test files]
ok      github.com/ItserX/biathlon_competions/internal/events   0.002s  coverage: 89.5% of statements
ok      github.com/ItserX/biathlon_competions/internal/report   0.002s  coverage: 86.2% of statements
```
