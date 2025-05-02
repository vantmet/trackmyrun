package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// create a new struct to hold the run data
type StravaActivity struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Distance     float64   `json:"distance"`
	StartDate    time.Time `json:"start_date"`
	ActivityType string    `json:"type"`
}

type StravaToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	var st StravaToken
	url := "https://www.strava.com/oauth/token"

	// load env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Environment not loaded.")
	}

	//load the TokenCache if available
	stRaw, err := os.ReadFile("token.json")
	if err != nil {
		log.Println("Unable to open token.json continuing.")
		tok := requestAccess()
		os.Setenv("STRAVA_ACCESS_TOKEN", tok)
		st, err = exchangeToken(url)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = json.Unmarshal(stRaw, &st)
		if err != nil {
			log.Fatal("Unable to unmarshall token")
		}
	}
	//st.AccessToken = os.Getenv("STRAVA_ACCESS_TOKEN")
	convertedTime := time.Unix(int64(st.ExpiresAt), 0)
	log.Printf("Token Expires: %q", convertedTime)
	if time.Now().Unix() > int64(st.ExpiresAt) {
		log.Println("Token Expired, refreshing...")
		st, err = getrefreshedToken(url, st.RefreshToken)
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Access Token: valid.")
	//Write token out to file
	output, _ := json.MarshalIndent(st, "", "  ")

	err = os.WriteFile("token.json", output, 0644)
	// get the strava runs
	runs := getStravaRuns(st.AccessToken)

	// log the number of runs collected
	log.Printf("Retrieved %d runs.", len(runs))
}

// create a new function to get the strava runs
func getStravaRuns(token string) []StravaActivity {

	// create a new http client
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://www.strava.com/api/v3/athlete/activities", nil)
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
