package main

import (
	"testing"
)

func TestRun(t *testing.T) {
	cases := []struct {
		Description string
		Run         Run
		DistWant    string
		TimeWant    RunTime
		PaceWant    string
	}{
		{"0Km Run", Run{0, RunTime{0, 4, 1}}, "0km", RunTime{0, 4, 1}, "Invalid Distance"},
		{"-10Km Run", Run{-10, RunTime{0, 4, 1}}, "-10km", RunTime{0, 4, 1}, "Invalid Distance"},
		{"1Km Run", Run{1, RunTime{0, 4, 1}}, "1km", RunTime{0, 4, 1}, "4.02"},
		{"5Km Run", Run{5, RunTime{0, 34, 1}}, "5km", RunTime{0, 34, 1}, "6.80"},
		{"10Km Run", Run{10, RunTime{1, 4, 1}}, "10km", RunTime{1, 4, 1}, "6.40"},
		{"100Km Run", Run{100, RunTime{36, 4, 1}}, "100km", RunTime{36, 4, 1}, "21.64"},
	}

	for _, test := range cases {
		t.Run(test.Description+"Dist", func(t *testing.T) {
			got := GetRunDistanceKm(test.Run)
			if got != test.DistWant {
				t.Errorf("got %q want %q", got, test.DistWant)
			}
		})
		t.Run(test.Description+"Time", func(t *testing.T) {
			got := GetRunTime(test.Run)
			if got != test.TimeWant {
				t.Errorf("got %v want %v", got, test.TimeWant)
			}
		})
		t.Run(test.Description+"Pace", func(t *testing.T) {
			got := GetRunPace(test.Run)
			if got != test.PaceWant {
				t.Errorf("got %q want %q", got, test.PaceWant)
			}
		})
	}
}
