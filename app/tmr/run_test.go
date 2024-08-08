package main

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	const shortForm = "2006-Jan-02"
	dateSet, _ := time.Parse(shortForm, "2013-Feb-03")
	dateWant := "2013-Feb-03"

	cases := []struct {
		Description    string
		Run            Run
		DistWant       string
		TimeWant       RunTime
		TimeStringWant string
		PaceWant       string
		DateWant       string
	}{
		{"0Km Run", Run{dateSet, 0, RunTime{0, 4, 1}}, "Invalid Distance", RunTime{0, 4, 1}, "0:4:1", "Invalid Distance", dateWant},
		{"-10Km Run", Run{dateSet, -10000, RunTime{0, 4, 1}}, "Invalid Distance", RunTime{0, 4, 1}, "0:4:1", "Invalid Distance", dateWant},
		{"1Km Run", Run{dateSet, 1000, RunTime{0, 4, 1}}, "1km", RunTime{0, 4, 1}, "0:4:1", "4.02", dateWant},
		{"5Km Run", Run{dateSet, 5000, RunTime{0, 34, 1}}, "5km", RunTime{0, 34, 1}, "0:34:1", "6.80", dateWant},
		{"5Km Run", Run{dateSet, 5420, RunTime{0, 34, 52}}, "5.42km", RunTime{0, 34, 52}, "0:34:52", "6.43", dateWant},
		{"10Km Run", Run{dateSet, 10000, RunTime{1, 4, 1}}, "10km", RunTime{1, 4, 1}, "1:4:1", "6.40", dateWant},
		{"100Km Run", Run{dateSet, 100000, RunTime{36, 4, 1}}, "100km", RunTime{36, 4, 1}, "36:4:1", "21.64", dateWant},
	}

	for _, test := range cases {
		t.Run(test.Description+"DateString", func(t *testing.T) {
			got := test.Run.GetRunDateString()
			if got != test.DateWant {
				t.Errorf("got %q want %v", got, test.DateWant)
			}
		})
		t.Run(test.Description+"Dist", func(t *testing.T) {
			got := test.Run.GetRunDistanceKm()
			if got != test.DistWant {
				t.Errorf("got %q want %q", got, test.DistWant)
			}
		})
		t.Run(test.Description+"Time", func(t *testing.T) {
			got := test.Run.GetRunTime()
			if got != test.TimeWant {
				t.Errorf("got %v want %v", got, test.TimeWant)
			}
		})
		t.Run(test.Description+"TimeString", func(t *testing.T) {
			got := test.Run.GetRunTimeString()
			if got != test.TimeStringWant {
				t.Errorf("got %q want %q", got, test.TimeStringWant)
			}
		})
		t.Run(test.Description+"Pace", func(t *testing.T) {
			got := test.Run.GetRunPace()
			if got != test.PaceWant {
				t.Errorf("got %q want %q", got, test.PaceWant)
			}
		})
	}
}
func TestPlanRun(t *testing.T) {
	const shortForm = "2006-Jan-02"
	dateSet, _ := time.Parse(shortForm, "2013-Feb-03")
	dateWant := "2013-Feb-03"

	cases := []struct {
		Description string
		PlanRun     PlanRun
		DistWant    string
		DateWant    string
	}{
		{"0Km Run", PlanRun{dateSet, 0}, "Invalid Distance", dateWant},
		{"-10Km Run", PlanRun{dateSet, -10000}, "Invalid Distance", dateWant},
		{"1Km Run", PlanRun{dateSet, 1000}, "1km", dateWant},
		{"5Km Run", PlanRun{dateSet, 5000}, "5km", dateWant},
		{"5Km Run", PlanRun{dateSet, 5420}, "5.42km", dateWant},
		{"10Km Run", PlanRun{dateSet, 10000}, "10km", dateWant},
		{"100Km Run", PlanRun{dateSet, 100000}, "100km", dateWant},
	}

	for _, test := range cases {
		t.Run(test.Description+"DateString", func(t *testing.T) {
			got := test.PlanRun.GetRunDateString()
			if got != test.DateWant {
				t.Errorf("got %q want %v", got, test.DateWant)
			}
		})
		t.Run(test.Description+"Dist", func(t *testing.T) {
			got := test.PlanRun.GetRunDistanceKm()
			if got != test.DistWant {
				t.Errorf("got %q want %q", got, test.DistWant)
			}
		})
	}
}
