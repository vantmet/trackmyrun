package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/vantmet/trackmyrun/internal/runstore"
)

func TestRecordingAndRetrievingRuns(t *testing.T) {
	store := runstore.InMemoryRunnerStore{}
	server := RunnerServer{&store, filepath.FromSlash("../../web/html")}

	const shortForm = "2006-Jan-02"
	date1, _ := time.Parse(shortForm, "2013-Feb-03")
	date2, _ := time.Parse(shortForm, "2013-Feb-04")
	date3, _ := time.Parse(shortForm, "2013-Feb-05")

	run1 := runstore.Run{
		Date:     date1,
		Distance: 5.42,
		Runtime:  2094,
	}
	run2 := runstore.Run{
		Date:     date2,
		Distance: 5.42,
		Runtime:  2094,
	}
	run3 := runstore.Run{
		Date:     date3,
		Distance: 5.42,
		Runtime:  2094,
	}

	server.ServeHTTP(httptest.NewRecorder(), newPostRunRequest(run1))
	server.ServeHTTP(httptest.NewRecorder(), newPostRunRequest(run2))
	server.ServeHTTP(httptest.NewRecorder(), newPostRunRequest(run3))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetRunsRequest())

	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "2013-Feb-03")
	assertResponseBody(t, response.Body.String(), "2013-Feb-04")
	assertResponseBody(t, response.Body.String(), "2013-Feb-05")
}
