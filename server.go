package main

import (
	"net/http"
	"text/template"
)

func RunnerServer(w http.ResponseWriter, r *http.Request) {
	data := struct {
		PageTitle string
	}{
		PageTitle: "My Latest Runs",
	}
	t, _ := template.ParseFiles("html\\GetLatest.html")

	t.Execute(w, data)

}
