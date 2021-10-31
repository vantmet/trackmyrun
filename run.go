package main

import (
	"math"
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
	Distance float32 //All distances stored in km.
	RunTime  RunTime
}

func (r Run) GetRunDistanceKm() string {
	if math.Floor(float64(r.Distance)) != float64(r.Distance) {
		return strconv.FormatFloat(float64(r.Distance), 'f', 2, 32) + "km"
	} else {
		return strconv.FormatFloat(float64(r.Distance), 'f', 0, 32) + "km"
	}
}

func (r Run) GetRunTime() RunTime {
	return r.RunTime
}

func (r Run) GetRunPace() string {

	if r.Distance > 0 {
		// Pace is minutes per Km.
		// Calculate total time in secs (no division needed)
		timeInSecs := float64(r.RunTime.Seconds)
		timeInSecs += float64(r.RunTime.Minutes * 60)
		timeInSecs += float64(r.RunTime.Hours * 60 * 60)

		pace := timeInSecs / float64(r.Distance) // pace in secs/km
		//divide pace by 60 to give pace in min.
		return strconv.FormatFloat(pace/60, 'f', 2, 64)
	} else {
		return "Invalid Distance"
	}

}

func (r Run) GetRunDate() time.Time {
	return r.Date
}
