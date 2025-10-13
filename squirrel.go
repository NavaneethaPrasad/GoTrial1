package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type Journalentries struct {
	Events   []string `json:"events"`
	Squirrel bool     `json:"squirrel"`
}

func phi(n11, n10, n01, n00 float64) float64 {
	numerator := n11*n00 - n10*n01
	denominator := math.Sqrt((n11 + n10) * (n01 + n00) * (n11 + n01) * (n10 + n00))
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

func contains(eventlist []string, targetevent string) bool {
	for _, currentevent := range eventlist {
		if currentevent == targetevent {
			return true
		}
	}
	return false
}

func main() {
	journal := []Journalentries{}
	jsonData, _ := os.ReadFile("journal.json")
	err := json.Unmarshal(jsonData, &journal)
	if err != nil {
		fmt.Println("Couldn't unmarshal data", err)
		return
	}

	// collect unique events
	uniqueevents := make(map[string]bool)
	for _, events := range journal {
		for _, event := range events.Events {
			uniqueevents[event] = true
		}
	}

	maxEvent := ""
	minEvent := ""
	maxPhi := math.Inf(-1)
	minPhi := math.Inf(1)
	fmt.Println("Event       Correlation")

	//compute correlations
	for event := range uniqueevents {
		var n11, n10, n01, n00 float64
		for _, entry := range journal {
			has := contains(entry.Events, event)
			squirrel := entry.Squirrel
			if has && squirrel {
				n11++
			} else if has && !squirrel {
				n10++
			} else if !has && squirrel {
				n01++
			} else {
				n00++
			}
		}
		correlation := phi(n11, n10, n01, n00)
		fmt.Printf("%-15s %.4f\n", event, correlation)
		if correlation > maxPhi {
			maxPhi = correlation
			maxEvent = event
		}
		if correlation < minPhi {
			minPhi = correlation
			minEvent = event
		}
	}
	fmt.Println("\nMost positively correlated event:", maxEvent, maxPhi)
	fmt.Println("Most negatively correlated event:", minEvent, minPhi)
}
