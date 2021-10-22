package main

import (
	"strconv"
	"time"
)

type RunTime struct {
	Hours   int
	Minutes int
	Seconds float32
}

type Run struct {
	Date     time.Time
	Distance int //All distances stored in km.
	RunTime  RunTime
}

func GetRunDistanceKm(r Run) string {
	return strconv.Itoa(r.Distance) + "km"
}

func GetRunTime(r Run) RunTime {
	return r.RunTime
}

func GetRunPace(r Run) string {

	if r.Distance > 0 {
		// Pace is minutes per Km.
		timeInSecs := float64(r.RunTime.Seconds)
		timeInSecs += float64(r.RunTime.Minutes * 60)
		timeInSecs += float64(r.RunTime.Hours * 60 * 60)

		pace := timeInSecs / float64(r.Distance)
		return strconv.FormatFloat(pace/60, 'f', 2, 64)
	} else {
		return "Invalid Distance"
	}

}

func GetRunDate(r Run) time.Time {
	return r.Date
}