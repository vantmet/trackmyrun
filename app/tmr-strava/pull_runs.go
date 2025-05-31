package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/vantmet/trackmyrun/internal/runstore"
)

// create a new struct to hold the run data
type StravaActivity struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Distance     float64   `json:"distance"`
	StartDate    time.Time `json:"start_date"`
	ActivityType string    `json:"type"`
	MovingTime   int       `json:"moving_time"`
	ElapsedTime  int       `json:"elapsed_time"`
}

func main() {
	var st runstore.StravaToken
	var store runstore.Store
	var err error
	tokenid, _ := uuid.Parse("891b5b6d-ee44-4dd4-b288-81ee766338c5")
	ctx := context.Background()

	url := "https://www.strava.com/oauth/token"

	// load env vars
	// First try to load from the local environemnet
	if os.Getenv("STRAVA_CLIENT_ID") == "" || os.Getenv("STRAVA_CLIENT_SECRET") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Environment not loaded.")
		}
	}

	if os.Getenv("TMRENV") == "DEV" {
		store = &runstore.InMemoryRunnerStore{}
	} else {
		store, err = runstore.NewSQLRunerStore(ctx)
		if err != nil {
			panic(err)
		}
	}

	st, err = store.GetRunnerStravaToken(tokenid)
	if err != nil {
		log.Println("Unable to open token.json continuing.")
		tok := requestAccess()
		os.Setenv("STRAVA_ACCESS_TOKEN", tok)
		st, err = exchangeToken(url)
		if err != nil {
			log.Fatal(err)
		}
		st.ID = tokenid
		_, err = store.NewRunnerStravaToken(st)
		if err != nil {
			log.Fatal(err)
		}
	}

	convertedTime := time.Unix(int64(st.ExpiresAt), 0)
	log.Printf("Token Expires: %q", convertedTime)
	if time.Now().Unix() > int64(st.ExpiresAt) {
		log.Println("Token Expired, refreshing...")
		st, err = getrefreshedToken(url, st.RefreshToken)
		st.ID = tokenid
		_, err = store.UpdateRunnerStravaToken(st)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Access Token: valid.")

	// get the strava runs
	stravaRuns := getStravaRuns(store, st.AccessToken)
	// log the number of runs collected
	log.Printf("Retrieved %d runs.", len(stravaRuns))
	runs := convertStravaRuns(stravaRuns)
	for _, run := range runs {
		store.RecordRun(run)
	}
}

// create a new function to get the strava runs
func getStravaRuns(store runstore.Store, token string) []StravaActivity {
	var limit string
	lastrun, err := store.GetLastRunnerRun()
	if err == pgx.ErrNoRows {
		limit = ""
	} else {
		limit = fmt.Sprintf("?after=%d", lastrun.Date.Unix())
	}

	// create a new http client
	client := &http.Client{}
	base := "https://www.strava.com/api/v3/athlete/activities"
	lookup := base + limit

	req, err := http.NewRequest("GET", lookup, nil)
	if err != nil {
		log.Println("Error reading request", err)
		return nil
	}

	//add headers
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "TMR-Strava")
	// log.Println(req.Header)

	//add query params
	q := req.URL.Query()
	q.Add("per_page", "10")
	req.URL.RawQuery = q.Encode()

	log.Printf("Getting Runs, %s", req.URL)
	// send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request", err)
		return nil
	}

	log.Println("Response Status:", resp.Status)
	// check if client is authorized
	if resp.StatusCode == 401 {
		log.Println("Unauthorized")
		return nil
	}

	// read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response", err)
		return nil
	}
	// log.Println("Response Body:", string(body))
	// parse the response
	// create variable to hold runs
	var runs []StravaActivity
	err = json.Unmarshal(body, &runs)
	if err != nil {
		log.Println("Error parsing response", err)
		return nil
	}

	return runs
}

func convertStravaRuns(runs []StravaActivity) []runstore.Run {
	var tmrRuns []runstore.Run

	for _, stravaRun := range runs {
		rt := stravaRun.ElapsedTime
		run := runstore.Run{
			Date:     stravaRun.StartDate,
			Distance: stravaRun.Distance,
			Runtime:  int32(rt)}
		tmrRuns = append(tmrRuns, run)
	}

	return tmrRuns
}
