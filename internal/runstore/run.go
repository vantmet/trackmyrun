package runstore

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type Store interface {
	GetRunnerStravaToken(int) (StravaToken, error)
	GetRunnerRuns() []Run
	RecordRun(Run)
}

type PlanRun struct {
	Date     time.Time
	Distance float32 //All distances stored in m.
}

type StravaToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
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

func (r Run) GetRuntime() int {
	return int(r.Runtime)
}

func (r Run) GetRuntimeString() string {
	return fmt.Sprintf("%s", time.Duration(int64(r.Runtime)*1000000000))
}

func (r Run) GetRunPace() string {

	if r.Distance > 0 {
		// Pace is minutes per Km.
		kmDistance := float64(r.Distance) / 1000.0
		pace := float64(r.Runtime) / float64(kmDistance) // pace in secs/km

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
