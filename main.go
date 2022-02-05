package main

import (
	"log"
	"net/http"
)

type InMemoryRunnerStore struct {
	runs []Run
}

func (i *InMemoryRunnerStore) GetRunnerRuns() []Run {
	return i.runs
}

func (i *InMemoryRunnerStore) RecordRun(r Run) {}

func main() {
	handler := &RunnerServer{&InMemoryRunnerStore{}}
	log.Fatal(http.ListenAndServe(":5000", handler))
}
