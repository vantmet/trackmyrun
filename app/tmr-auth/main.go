package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	chiprometheus "github.com/jamscloud/chi-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var Version string

type CognitoClient struct {
	*cip.Client
	AppClientID string
	UserPoolID  string
}

func Init() *CognitoClient {
	// Load the shared AWS config (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	AppID := os.Getenv("COGNITO_APP_CLIENT_ID")
	PoolID := os.Getenv("COGNITO_USER_POOL_ID")

	log.Println("AWS Region ", cfg.Region)
	log.Println("AppID ", AppID)
	log.Println("PoolID ", PoolID)
	log.Println("Configuring Cognito Client ", cfg.Region)
	return &CognitoClient{
		cip.NewFromConfig(cfg),
		AppID,
		PoolID,
	}
}

func main() {
	appVersion := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "app_version_info",
		Help:        "App info at buildtime",
		ConstLabels: prometheus.Labels{"version": Version},
	})
	prometheus.Register(appVersion)
	appVersion.Set(1)

	log.Println("App Version ", Version, ", registered in Prometheus")

	cognitoClient := Init()

	log.Println("Cognito Client Inititated")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger, middleware.WithValue("CognitoClient", cognitoClient))
	r.Use(middleware.Heartbeat("/ping"))
	m := chiprometheus.NewMiddleware("serviceName")
	r.Use(m)
	r.Use(middleware.Recoverer)

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/favicon.ico", faviconHandler)
	r.Get("/", http.HandlerFunc(handleIndex))

	r.Post("/signup", cognitoClient.signUp)

	// r.Post("/signin", signIn)

	// r.Get("/verify", verifyToken)

	log.Println("Handlers Initiated")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	log.Println("Listening Port configured, starting server on:", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	data := struct {
		PageTitle string
		Version   string
	}{
		PageTitle: "Welcome!",
		Version:   Version,
	}
	f := filepath.Join(filepath.FromSlash("html"), "SignUp.html")
	t, err := template.ParseFiles(f)

	if err == nil {
		t.Execute(w, data)
	} else {
		log.Printf("Template error: %q", err)
	}
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/favicon.ico")
}

func (cc *CognitoClient) signUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")
	// phoneNumber := r.Form.Get("phone_number")

	user := &cip.SignUpInput{
		Username: aws.String(username),
		Password: aws.String(password),
		ClientId: aws.String(cc.AppClientID),
	}

	d := time.Now().Add(time.Minute)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	_, err := cc.Client.SignUp(ctx, user)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, fmt.Sprintf("/register?message=%s", err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/otp?username=%s", username), http.StatusFound)
}
