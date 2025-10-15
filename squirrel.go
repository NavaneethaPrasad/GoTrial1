package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"
)

type Journalentry struct {
	Events   []string `json:"events"`
	Squirrel bool     `json:"squirrel"`
}

type Counts struct {
	n00 uint
	n10 uint
	n01 uint
	n11 uint
}

type MaxMin struct {
	max      float64
	min      float64
	maxEvent string
	minEvent string
}

func getCounts(entries []Journalentry, event string) Counts {
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
	d := Counts{n00: n00, n10: n10, n01: n01, n11: n11}
	return d
}

func phi(count Counts) float64 {
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

func getCorrelations(journalEntries []Journalentry) map[string]float64 {
	map1 := make(map[string]float64)
	for _, entry := range journalEntries {
		for _, e := range entry.Events {
			c := getCounts(journalEntries, e)
			map1[e] = phi(c)

		}
	}
	return map1
}

func getMaxMin(corrValues map[string]float64) MaxMin {
	var results MaxMin
	results.max = -1.0
	results.min = 1.0
	for key, value := range corrValues {
		if value > results.max {
			results.max = value
			results.maxEvent = key
		}
		if value < results.min {
			results.min = value
			results.minEvent = key
		}

	}
	return results
}

func preprocess(journalEntries []Journalentry) []Journalentry {
	var journal []Journalentry
	for _, entry := range journalEntries {
		hasPeanuts := slices.Contains(entry.Events, "peanuts")
		notBrushedTeeth := !slices.Contains(entry.Events, "brushed teeth")
		if hasPeanuts && notBrushedTeeth {
			entry.Events = append(entry.Events, "dirty teeth")
			journal = append(journal, entry)
		} else {
			journal = append(journal, entry)
		}
	}
	return journal
}

func main() {
	journal := []Journalentry{}
	jsonData, _ := os.ReadFile("journal.json")
	err := json.Unmarshal(jsonData, &journal)
	if err != nil {
		fmt.Println("Couldn't unmarshal data", err)
	}
	journal = preprocess(journal)
	map1 := getCorrelations(journal)
	results := getMaxMin(map1)
	fmt.Printf("%-15s %s\n", "Event", "Correlation")
	fmt.Println("--------------------------------------------")
	for event, correlation := range map1 {
		fmt.Printf("%-15s %.4f\n", event, correlation)
	}
	fmt.Printf("Most Positively Correlated Event: %s = %f\n", results.maxEvent, results.max)
	fmt.Printf("Most Negatively Correlated Event: %s = %f\n", results.minEvent, results.min)
}
