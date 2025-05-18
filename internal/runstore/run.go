package runstore

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type Store interface {
	GetRunnerRuns() []Run
	RecordRun(Run)
}

type RunTime struct {
	Hours   int
	Minutes int
	Seconds float32
}

type Run struct {
	Date     time.Time
	Distance float32 //All distances stored in m.
	RunTime  int     //All times stored in seconds.
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

func (r Run) GetRunTime() int {
	return r.RunTime
}

func (r Run) GetRunTimeString() string {
	return fmt.Sprintf("%s", time.Duration(r.RunTime*1000000000))
}

func (r Run) GetRunPace() string {

	if r.Distance > 0 {
		// Pace is minutes per Km.
		kmDistance := float64(r.Distance) / 1000.0

		pace := float64(r.RunTime) / float64(kmDistance) // pace in secs/km
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

func SecondsToRunTime(seconds int) RunTime {
	secs := seconds % 60
	mins := (seconds - secs) % 60
	var hours int
	if (seconds - secs - mins) >= 60 {
		hours = (seconds - secs - mins) / 60
	} else {
		hours = 0
	}

	r := RunTime{Hours: hours, Minutes: mins, Seconds: float32(secs)}
	return r
}

/* func SecondsToRunTime(s int) RunTime {
 	elapsed := time.Duration(s * int(time.Second))

	return RunTime{Hours: int(elapsed.Hours()), Minutes: int(elapsed.Minutes()), Seconds: float32(elapsed.Seconds())}

}*/
