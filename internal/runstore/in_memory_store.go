package runstore

type InMemoryRunnerStore struct {
	runs []Run
}

func (i *InMemoryRunnerStore) GetRunnerRuns() []Run {
	return i.runs
}

func (i *InMemoryRunnerStore) RecordRun(r Run) {
	i.runs = append(i.runs, r)
}

func (i *InMemoryRunnerStore) GetRunnerStravaToken(userid int) (StravaToken, error) {
	st := StravaToken{AccessToken: "", ExpiresAt: 0, ExpiresIn: 0, RefreshToken: ""}
	return st, nil
}
