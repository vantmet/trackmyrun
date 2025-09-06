package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	"github.com/vantmet/trackmyrun/internal/runstore"
)

type RunnerStore interface {
	GetRunnerRuns() []runstore.Run
	RecordRun(r runstore.Run)
}

type RunnerServer struct {
	store    RunnerStore
	htmlRoot string
}

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
	var run runstore.Run
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
		fRuntimeString := r.FormValue("runtime")
		log.Printf("Date: %q, Distance: %q, Runtime: %q", fDateString, fDistString, fRuntimeString)

		//parse Date
		const shortForm = "2006-01-02T15:04"
		fDate, _ := time.Parse(shortForm, fDateString)
		//Parse Distance
		fDist, _ := strconv.ParseFloat(fDistString, 32)
		fDist = fDist * 1000.0 //Convert to meters
		//Parse runtime
		duration, err := time.ParseDuration(fRuntimeString)
		if err != nil {
			log.Printf("Unable to parse time.")
		}
		//Populate the run
		run.Date = fDate
		run.Distance = fDist
		run.Runtime = int32(duration.Seconds())
		log.Printf("Saved Run: %v", run)

	}
	rs.store.RecordRun(run)
	w.WriteHeader(http.StatusAccepted)
	rs.showRuns(w, r, true)
}

func (rs *RunnerServer) showRuns(w http.ResponseWriter, r *http.Request, success bool) {
	runs := rs.store.GetRunnerRuns()

	data := struct {
		PageTitle string
		Runs      []runstore.Run
		Status    bool
		Version   string
	}{
		PageTitle: "My Latest Runs",
		Runs:      runs,
		Status:    success,
		Version:   Version,
	}
	f := filepath.Join(rs.htmlRoot, "GetLatest.html")
	t, err := template.ParseFiles(f)

	if err == nil {
		t.Execute(w, data)
	} else {
		log.Printf("Template error: %q", err)
	}
}
