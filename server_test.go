package main

import (
	"net/http"
	"net/http/httptest"
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
				Distance: 5.42,
				RunTime:  RunTime{0, 34, 52},
			},
		},
	}
	server := &RunnerServer{&store}
	request, _ := http.NewRequest(http.MethodGet, "/runs", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	t.Run("Returns 200", func(t *testing.T) {
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("Body Contains 'Latest Runs'", func(t *testing.T) {
		got := response.Body.String()
		want := "Latest Runs"

		if !strings.Contains(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("title is 'My Latest Runs'", func(t *testing.T) {
		got := response.Body.String()
		want := "<title>My Latest Runs</title>"

		if !strings.Contains(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Contains run table header", func(t *testing.T) {
		got := response.Body.String()
		wants := [4]string{"<th>Date</th>",
			"<th>Distance</th>",
			"<th>Time</th>",
			"<th>Pace</th>"}

		for _, want := range wants {
			if !strings.Contains(got, want) {
				t.Errorf("got %q, want %q", got, want)
			}
		}
	})
	t.Run("Contains a run", func(t *testing.T) {
		got := response.Body.String()
		wants := [4]string{"<td>2013-Feb-03</td>",
			"<td>5.42km</td>",
			"<td>0:34:52</td>",
			"<td>6.43</td>"}

		for _, want := range wants {
			if !strings.Contains(got, want) {
				t.Errorf("got %q, want %q", got, want)
			}
		}
	})
}

func TestStoreRun(t *testing.T) {

	store := StubRunStore{}
	server := &RunnerServer{&store}

	t.Run("it returns accepted on POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/runs", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)

	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

type StubRunStore struct {
	runs []Run
}

func (s *StubRunStore) GetRunnerRuns() []Run {
	return s.runs
}
