package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
	"time"
)

func RunnerServer(w http.ResponseWriter, r *http.Request) {
	run := GetRunnerRuns()
	data := struct {
		PageTitle string
		Runs      []Run
	}{
		PageTitle: "My Latest Runs",
		Runs: []Run{
			run,
		},
	}
	f := filepath.Join("html", "GetLatest.html")
	t, err := template.ParseFiles(f)

	if err == nil {
		t.Execute(w, data)
	} else {
		fmt.Printf("Template error: %q", err)
	}

}

func GetRunnerRuns() Run {
	const shortForm = "2006-Jan-02"
	date, _ := time.Parse(shortForm, "2013-Feb-03")
	run := Run{date, 5.42, RunTime{0, 34, 52}}
	return run
}
