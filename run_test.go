package main

import (
	"testing"
)

func TestGetRunDistance(t *testing.T) {
	cases := []struct {
		Description string
		Run         Run
		Want        string
	}{
		{"1Km Run", Run{1}, "1km"},
		{"1Km Run", Run{5}, "5km"},
		{"1Km Run", Run{10}, "10km"},
		{"1Km Run", Run{100}, "100km"},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			got := GetRunDistanceKm(test.Run)
			if got != test.Want {
				t.Errorf("got %q want %q", got, test.Want)
			}
		})
	}
}
