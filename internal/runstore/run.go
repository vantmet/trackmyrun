package runstore

import (
	"fmt"
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
	Distance float32 //All distances stored in m.
	RunTime  RunTime
}

type PlanRun struct {
	Date     time.Time
	Distance float32 //All distances stored in m.
}

func GetDistanceKm(d float64) string {
	if d <= 0 {
		return "Invalid Distance"
	} else {
		if math.Floor(d) != d {
			return strconv.FormatFloat(d, 'f', 2, 32) + "km"
		} else {
			return strconv.FormatFloat(d, 'f', 0, 32) + "km"
		}
	}
}

func (r Run) GetRunDistanceKm() string {

	return GetDistanceKm(float64(r.Distance) / 1000.0)
}

func (r PlanRun) GetRunDistanceKm() string {
	return GetDistanceKm(float64(r.Distance) / 1000.0)
}

func (r Run) GetRunTime() RunTime {
	return r.RunTime
}

func (r Run) GetRunTimeString() string {
	return fmt.Sprint(r.RunTime.Hours, ":", r.RunTime.Minutes, ":", r.RunTime.Seconds)
}

func (r Run) GetRunPace() string {

	if r.Distance > 0 {
		// Pace is minutes per Km.
		kmDistance := float64(r.Distance) / 1000.0
		// Calculate total time in secs (no division needed)
		timeInSecs := float64(r.RunTime.Seconds)
		timeInSecs += float64(r.RunTime.Minutes * 60)
		timeInSecs += float64(r.RunTime.Hours * 60 * 60)

		pace := timeInSecs / float64(kmDistance) // pace in secs/km
		//divide pace by 60 to give pace in min.
		return strconv.FormatFloat(pace/60, 'f', 2, 64)
	} else {
		return "Invalid Distance"
	}

}

func (r Run) GetRunDateString() string {
	const shortForm = "2006-Jan-02"

	return r.Date.Format(shortForm)
}

func (r PlanRun) GetRunDateString() string {
	const shortForm = "2006-Jan-02"

	return r.Date.Format(shortForm)
}
