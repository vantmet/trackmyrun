package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func (rs *RunnerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.URL, r.Body)

	switch r.Method {
	case http.MethodPost:
		rs.processRun(w, r)
	case http.MethodGet:
		rs.showRuns(w, r, false)
	}

}

func (rs *RunnerServer) processRun(w http.ResponseWriter, r *http.Request) {
	var run Run
	if r.Header["Content-Type"][0] == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&run)
		log.Printf("%v", run)
		if err != nil {
			log.Printf("%q", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if r.Header["Content-Type"][0] == "application/x-www-form-urlencoded" {
		//Form data is a bunch of strings. Convert them to the right thing.
		//Validation is done in the Javascript for the form. TODO
		fDateString := r.FormValue("date")
		fDistString := r.FormValue("distance")
		fRunTimeString := r.FormValue("runtime")
		log.Printf("Date: %q, Distance: %q, RunTime: %q", fDateString, fDistString, fRunTimeString)

		//parse Date
		const shortForm = "2006-01-02T15:04"
		fDate, _ := time.Parse(shortForm, fDateString)
		//Pase Distance.
		fDist, _ := strconv.ParseFloat(fDistString, 32)
		//Parse runtime
		chunks := strings.Split(fRunTimeString, ":")
		fHours, _ := strconv.ParseInt(chunks[0], 10, 32)
		fMins, _ := strconv.ParseInt(chunks[1], 10, 32)
		fSecs, _ := strconv.ParseFloat(chunks[2], 32)

		//Populate the run
		run.Date = fDate
		run.Distance = float32(fDist)
		run.RunTime.Hours = int(fHours)
		run.RunTime.Minutes = int(fMins)
		run.RunTime.Seconds = float32(fSecs)

	}
	rs.store.RecordRun(run)
	w.WriteHeader(http.StatusAccepted)
	rs.showRuns(w, r, true)
}

func (rs *RunnerServer) showRuns(w http.ResponseWriter, r *http.Request, success bool) {
	runs := rs.store.GetRunnerRuns()

	data := struct {
		PageTitle string
		Runs      []Run
		Status    bool
	}{
		PageTitle: "My Latest Runs",
		Runs:      runs,
		Status:    success,
	}
	f := filepath.Join("html", "GetLatest.html")
	t, err := template.ParseFiles(f)

	if err == nil {
		t.Execute(w, data)
	} else {
		log.Printf("Template error: %q", err)
	}
}

type RunnerStore interface {
	GetRunnerRuns() []Run
	RecordRun(r Run)
}

type RunnerServer struct {
	store RunnerStore
}
