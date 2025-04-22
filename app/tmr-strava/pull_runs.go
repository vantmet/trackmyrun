package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

	// load env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Environment not loaded.")
	}
	// TODO
	// 1) Save the token we get after exchange to a file.
	// 2) Check the file to see if we need to request access using master creds
	// 3) Check to see if we need to refresh the token,
	tok := requestAccess()
	url := "https://www.strava.com/oauth/token"
	os.Setenv("STRAVA_ACCESS_TOKEN", tok)
	st, err = exchangeToken(url)
	if err != nil {
		log.Fatal(err)
	}

	//st.AccessToken = os.Getenv("STRAVA_ACCESS_TOKEN")
	if time.Now().Unix() > int64(st.ExpiresAt) {
		st, err = refreshToken(url)
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Access Token: valid.")
	// get the strava runs
	runs := getStravaRuns(st.AccessToken)

	// print the runs
	for _, run := range runs {
		log.Println(run)
	}
}

// Initial Access Request
func requestAccess() (access_tok string) {
	fmt.Println("Please go to the following URL and Authorizxe the app. Once redirected copy the 'code' and paste it here. Then hit return.")
	fmt.Printf("https://www.strava.com/oath/authorize?client_id=%s&response_type=code&redirect_uri=http://localhost/ex_tok&approval_prompt=force&scope=read_all,activity:read_all\n: ", os.Getenv("STRAVA_CLIENT_ID"))

	reader := bufio.NewReader(os.Stdin)

	text, _ := reader.ReadString('\n')

	return strings.TrimSpace(text)
}

// create a new function to refresh the token
func exchangeToken(baseURL string) (st StravaToken, err error) {
	et := ExchangeToken{
		StravaClientID:     os.Getenv("STRAVA_CLIENT_ID"),
		StravaClientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
		Code:               os.Getenv("STRAVA_ACCESS_TOKEN"),
		GrantType:          "authorization_code",
	}
	etString, _ := json.Marshal(et)

	// Debug: Print the userAuth struct
	log.Printf("Exchange Token: %+v", et)

	// create a new http client
	client := &http.Client{}

	// create a new request
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(etString))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Add headers and send the request
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Debug: Print the response status and body
	body := new(strings.Builder)
	io.Copy(body, resp.Body)
	log.Printf("Response Status: %s", resp.Status)
	//	log.Printf("Response Body: %s", body)

	if resp.StatusCode != 200 {
		return st, errors.New("Returned Error: " + body.String())
	}

	// Handle the response...
	var tokenResponse StravaToken
	err = json.Unmarshal([]byte(body.String()), &tokenResponse)
	if err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	return tokenResponse, nil
}

// create a new function to refresh the token
func refreshToken(baseURL string) (st StravaToken, err error) {
	userAuth := RefreshToken{
		StravaClientID:     os.Getenv("STRAVA_CLIENT_ID"),
		StravaClientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
		RefreshToken:       os.Getenv("STRAVA_REFRESH_TOKEN"),
		GrantType:          "refresh_token",
	}
	if userAuth.RefreshToken == "" || userAuth.StravaClientID == "" || userAuth.StravaClientSecret == "" {
		err = fmt.Errorf("Environment not complete. %q", userAuth)
		return st, err
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

	return tokenResponse, nil
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
