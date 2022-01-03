package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
)

func RunnerServer(w http.ResponseWriter, r *http.Request) {
	const shortForm = "2006-Jan-02"
	date, _ := time.Parse(shortForm, "2013-Feb-03")
	run := Run{date, 5.42, RunTime{0, 34, 52}}
	data := struct {
		PageTitle string
		Runs      []Run
	}{
		PageTitle: "My Latest Runs",
		Runs: []Run{
			run,
		},
	}
	t, err := template.ParseFiles("html\\GetLatest.html")

	if err == nil {
		t.Execute(w, data)
	} else {
		fmt.Printf("Template error: %q", err)
	}

}
