package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"
)

type Journalentries struct {
	Events   []string `json:"events"`
	Squirrel bool     `json:"squirrel"`
}

type counts struct {
	n00 uint
	n10 uint
	n01 uint
	n11 uint
}

func getCounts(entries []Journalentries, event string) counts {
	var n00, n10, n01, n11 uint
	for _, i := range entries {
		if slices.Contains(i.Events, event) {
			//fmt.Println(i.Events, i.Squirrel)
			if i.Squirrel {
				n11++
			} else {
				n10++
			}
		} else {
			if i.Squirrel {
				n01++
			} else {
				n00++
			}
		}
	}
	d := counts{n00: n00, n10: n10, n01: n01, n11: n11}
	return d
}

func phi(count counts) float64 {
	n00 := float64(count.n00)
	n10 := float64(count.n10)
	n01 := float64(count.n01)
	n11 := float64(count.n11)
	n := (n11*n00 - n10*n01)
	d := math.Sqrt((n10 + n11) * (n01 + n00) * (n01 + n11) * (n10 + n00))
	if d == 0 {
		return 0
	}
	return n / d
}

func main() {
	journal := []Journalentries{}
	jsonData, _ := os.ReadFile("journal.json")
	err := json.Unmarshal(jsonData, &journal)
	if err != nil {
		fmt.Println("Couldn't unmarshal data", err)
	}
	map1 := make(map[string]float64)
	for _, events := range journal {
		for _, e := range events.Events {
			c := getCounts(journal, e)
			//fmt.Println(phi(c))
			map1[e] = phi(c)
		}
	}
	fmt.Println("Event\t\tCorrelation")
	fmt.Println("--------------------------------------------")
	for event, correlation := range map1 {
		fmt.Printf("%-15s %.4f\n", event, correlation)
	}

	max := -1.0
	min := 1.0
	var maxEvent, minEvent string
	for key, value := range map1 {
		if value > max {
			max = value
			maxEvent = key
		}
		if value < min {
			min = value
			minEvent = key
		}
	}
	fmt.Printf("Most Positively Correlated Event: %s = %f\n", maxEvent, max)
	fmt.Printf("Most Negatively Correlated Event: %s = %f\n", minEvent, min)
}
