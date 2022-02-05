package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
	"time"
)

func (rs *RunnerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		rs.processRun(w)
	case http.MethodGet:
		rs.showRuns(w, r)
	}

}

func (rs *RunnerServer) processRun(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}

func (rs *RunnerServer) showRuns(w http.ResponseWriter, r *http.Request) {
	runs := rs.store.GetRunnerRuns()
	data := struct {
		PageTitle string
		Runs      []Run
	}{
		PageTitle: "My Latest Runs",
		Runs:      runs,
	}
	f := filepath.Join("html", "GetLatest.html")
	t, err := template.ParseFiles(f)

	if err == nil {
		t.Execute(w, data)
	} else {
		fmt.Printf("Template error: %q", err)
	}
}

func GetRunnerRuns() []Run {
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

type RunnerStore interface {
	GetRunnerRuns() []Run
}

type RunnerServer struct {
	store RunnerStore
}
