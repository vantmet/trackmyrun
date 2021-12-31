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
		want := "<th>Date</th>\r\n                <th>Distance</th>\r\n                <th>Time</th>\r\n                <th>Pace</th>\r\n"

		if !strings.Contains(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("Contains a run", func(t *testing.T) {
		got := response.Body.String()
		want := "<td>2013-02-03 00:00:00 +0000 UTC</td>\r\n                <td>5.42km</td>\r\n                <td>{0 34 52}</td>\r\n                <td>6.43</td>\r\n"

		if !strings.Contains(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	//TODO GetRunTime should return a properly formatted string.

	//TODO GetRunDate should return 2021-10-31
}
