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

type ExchangeToken struct {
	StravaClientID     string `json:"client_id"`
	StravaClientSecret string `json:"client_secret"`
	Code               string `json:"code"`
	GrantType          string `json:"grant_type"`
}

// create a new authentication struct
type RefreshToken struct {
	StravaClientID     string `json:"client_id"`
	StravaClientSecret string `json:"client_secret"`
	RefreshToken       string `json:"refresh_token"`
	GrantType          string `json:"grant_type"`
}

// create a struct to hold token refresh respone
type StravaToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	// refresh the token
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
		// TODO
		// 1) Save the token we get after exchange to a file.
		// 2) Check the file to see if we need to request access using master creds
		// 3) Check to see if we need to refresh the token,
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
	if time.Now().Unix() > int64(st.ExpiresAt) {
		st, err = refreshToken(url)
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Access Token: valid.")
	//Write token out to file
	output, _ := json.Marshal(st)

	err = os.WriteFile("token.json", output, 0644)
	// get the strava runs
	runs := getStravaRuns(st.AccessToken)

	// print the runs
	for _, run := range runs {
		log.Println(run)
	}
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
