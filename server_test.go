package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETRuns(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/runs", nil)
	response := httptest.NewRecorder()

	RunnerServer(response, request)

	got := response.Body.String()
	want := "Latest Runs"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
