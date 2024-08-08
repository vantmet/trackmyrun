package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGETRuns(t *testing.T) {
	const shortForm = "2006-Jan-02"
	date, _ := time.Parse(shortForm, "2013-Feb-03")
	store := StubRunStore{
		[]Run{
			{
				Date:     date,
				Distance: 5420,
				RunTime:  RunTime{0, 34, 52},
			},
		},
		nil,
	}
	server := &RunnerServer{&store, filepath.FromSlash("../../web/html")}
	response := httptest.NewRecorder()

	server.ServeHTTP(response, newGetRunsRequest())

	t.Run("Returns 200", func(t *testing.T) {
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("Body Contains 'Latest Runs'", func(t *testing.T) {
		assertResponseBody(t, response.Body.String(), "Latest Runs")
	})

	t.Run("Title is 'My Latest Runs'", func(t *testing.T) {
		assertResponseBody(t, response.Body.String(), "<title>My Latest Runs</title>")
	})

	t.Run("Contains run table header", func(t *testing.T) {
		wants := [4]string{"<th>Date</th>",
			"<th>Distance</th>",
			"<th>Time</th>",
			"<th>Pace</th>"}

		for _, want := range wants {
			assertResponseBody(t, response.Body.String(), want)
		}
	})
	t.Run("Contains a run", func(t *testing.T) {
		wants := [4]string{"<td>2013-Feb-03</td>",
			"<td>5.42km</td>",
			"<td>0:34:52</td>",
			"<td>6.43</td>"}

		for _, want := range wants {
			assertResponseBody(t, response.Body.String(), want)
		}
	})
}

func TestStoreRun(t *testing.T) {
	const shortForm = "2006-Jan-02"
	date, _ := time.Parse(shortForm, "2013-Feb-03")
	run := Run{
		Date:     date,
		Distance: 5420,
		RunTime:  RunTime{0, 34, 52},
	}

	store := StubRunStore{}
	server := &RunnerServer{&store, filepath.FromSlash("../../web/html")}

	t.Run("it returns accepted on POST", func(t *testing.T) {
		response := httptest.NewRecorder()

		server.ServeHTTP(response, newPostRunRequest(run))
		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.recordRunCalls) != 1 {
			t.Errorf("got %d calls to RecordRun want %d", len(store.recordRunCalls), 1)
		}

		if len(store.runs) == 0 {
			t.Errorf("Date not stored. Runs list empty.")
		}

		if store.runs[len(store.runs)-1].Date != date {
			t.Errorf("Got %q for run date expected %q.", store.runs[len(store.runs)-1].Date, date)
		}

	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

type StubRunStore struct {
	runs           []Run
	recordRunCalls []string
}

func (s *StubRunStore) GetRunnerRuns() []Run {
	return s.runs
}

func (s *StubRunStore) RecordRun(r Run) {
	s.recordRunCalls = append(s.recordRunCalls, "Added")
	s.runs = append(s.runs, r)
}

func newPostRunRequest(run Run) *http.Request {
	jRun, _ := json.Marshal(run)

	req, _ := http.NewRequest(http.MethodPost, "/runs", bytes.NewBuffer(jRun))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func newGetRunsRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/runs", nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
