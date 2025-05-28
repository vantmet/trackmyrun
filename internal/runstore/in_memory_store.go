package runstore

import "github.com/google/uuid"

type InMemoryRunnerStore struct {
	runs []Run
}

func (i *InMemoryRunnerStore) GetRunnerRuns() []Run {
	return i.runs
}

func (i *InMemoryRunnerStore) RecordRun(r Run) {
	i.runs = append(i.runs, r)
}

func (i *InMemoryRunnerStore) GetRunnerStravaToken(userid uuid.UUID) (StravaToken, error) {
	st := StravaToken{}
	return st, nil
}

func (i *InMemoryRunnerStore) NewRunnerStravaToken(token StravaToken) (StravaToken, error) {
	st := StravaToken{}
	return st, nil
}

func (i *InMemoryRunnerStore) UpdateRunnerStravaToken(token StravaToken) (StravaToken, error) {
	st := StravaToken{}
	return st, nil
}
