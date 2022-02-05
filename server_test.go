package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGETRuns(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/runs", nil)
	response := httptest.NewRecorder()

	RunnerServer(response, request)

	t.Run("Body Contains 'Latest Runs'", func(t *testing.T) {
		got := response.Body.String()
		want := "Latest Runs"

		if !strings.Contains(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("Returns 200", func(t *testing.T) {
		got := response.Code
		want := 200

		if got != want {
			t.Errorf("got %d, want %d", got, want)
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
