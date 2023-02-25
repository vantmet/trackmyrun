package main

type InMemoryRunnerStore struct {
	runs []Run
}

func (i *InMemoryRunnerStore) GetRunnerRuns() []Run {
	return i.runs
}

func (i *InMemoryRunnerStore) RecordRun(r Run) {
	i.runs = append(i.runs, r)
}
