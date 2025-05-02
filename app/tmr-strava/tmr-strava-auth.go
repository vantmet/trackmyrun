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
)

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
// Initial Access Request
func requestAccess() (access_tok string) {
	fmt.Println("Please go to the following URL and Authorizxe the app. Once redirected copy the 'code' and paste it here. Then hit return.")
	fmt.Printf("https://www.strava.com/oath/authorize?client_id=%s&response_type=code&redirect_uri=http://localhost/ex_tok&approval_prompt=force&scope=read_all,activity:read_all\n: ", os.Getenv("STRAVA_CLIENT_ID"))

	reader := bufio.NewReader(os.Stdin)

	text, _ := reader.ReadString('\n')

	return strings.TrimSpace(text)
}

// create a new function to exchange the access token from access request
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
