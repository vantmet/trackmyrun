package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRefreshToken(t *testing.T) {
	var want StravaToken
	// Mock the ENVVars
	t.Setenv("STRAVA_CLIENT_ID", "137832987")
	t.Setenv("STRAVA_CLIENT_SECRET", "somefancysecret")
	t.Setenv("STRAVA_REFRESH_TOKEN", "atoken")

	want.AccessToken = "12345"

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"token_type":"Bearer","access_token":"12345","expires_at":1745348271,"expires_in":15104,"refresh_token":"arefreshtoken"}`))
		}))
	defer server.Close()
	result, err := refreshToken(server.URL)

	if result.AccessToken != want.AccessToken || err != nil {
		t.Errorf("Unable to refresh token: %q", err)
	}
}

func TestEnv(t *testing.T) {
	want := ""
	url := "https://test.server"
	result, err := refreshToken(url)

	// Mock the ENVVars
	t.Setenv("STRAVA_CLIENT_ID", "")
	t.Setenv("STRAVA_CLIENT_SECRET", "")
	t.Setenv("STRAVA_ACCESS_CODE", "")

	if result.AccessToken != want || err == nil {
		t.Errorf("Invalid Environment failed to error. Want: %q, Error: %q", result, err)
	}
}
