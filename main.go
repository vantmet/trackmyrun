package main

import (
	"log"
	"net/http"
	"time"
)

type InMemoryRunnerStore struct{}

func (i *InMemoryRunnerStore) GetRunnerRuns() []Run {
	const shortForm = "2006-Jan-02"
	date, _ := time.Parse(shortForm, "2013-Feb-03")
	runs := []Run{
		{
			Date:     date,
			Distance: 5.42,
			RunTime:  RunTime{0, 34, 52},
		},
	}
	return runs
}

func main() {
	handler := &RunnerServer{&InMemoryRunnerStore{}}
	log.Fatal(http.ListenAndServe(":5000", handler))
}
