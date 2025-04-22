package main

import (
	"bytes"
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

// create a new authentication struct
type Auth struct {
	StravaClientID     string `json:"client_id"`
	StravaClientSecret string `json:"client_secret"`
	Code               string `json:"code"`
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
	url := "https://www.strava.com/api/v3/oauth/token"

	// load env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Environment not loaded.")
	}
	st.AccessToken = os.Getenv("STRAVA_ACCESS_TOKEN")
	st.ExpiresAt = 1723166586
	if time.Now().Unix() > int64(st.ExpiresAt) {
		st.AccessToken, st.ExpiresAt = refreshToken()
	}
	log.Printf("Access Token: %s", st.AccessToken)
	// get the strava runs
	runs := getStravaRuns(st.AccessToken)

	// print the runs
	for _, run := range runs {
		log.Println(run)
	}
}

// create a new function to refresh the token
func refreshToken(baseURL string) (st StravaToken, err error) {
	userAuth := Auth{
		StravaClientID:     os.Getenv("STRAVA_CLIENT_ID"),
		StravaClientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
		Code:               os.Getenv("STRAVA_ACCESS_CODE"),
		GrantType:          "refresh_token",
	}
	authStr, _ := json.Marshal(userAuth)

	// Debug: Print the userAuth struct
	log.Printf("User Auth: %+v", userAuth)

	// create a new http client
	client := &http.Client{}

	// create a new request
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(authStr))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Debug: Print the request body
	log.Printf("Request Body: %s", authStr)

	// Add headers and send the request
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Debug: Print the response status and body
	body, _ := io.ReadAll(resp.Body)
	log.Printf("Response Status: %s", resp.Status)
	log.Printf("Response Body: %s", body)

	// Handle the response...
	var tokenResponse StravaToken
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	return tokenResponse.AccessToken, tokenResponse.ExpiresAt
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

	log.Printf("Getting Runs, %s", req)
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
	log.Println("Response Body:", string(body))
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
