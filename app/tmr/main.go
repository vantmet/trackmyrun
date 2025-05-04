package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	chiprometheus "github.com/jamscloud/chi-prometheus"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Version string

var appVersion prometheus.Gauge

func main() {
	appVersion := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "app_version_info",
		Help:        "App info at buildtime",
		ConstLabels: prometheus.Labels{"version": Version},
	})
	prometheus.Register(appVersion)
	appVersion.Set(1)

	log.Println("App Version ", Version, ", registered in Prometheus")

	// store := InMemoryRunnerStore{}
	// load env vars
	// First try to load from the local environemnet
	if os.Getenv("STRAVA_CLIENT_ID") == "" || os.Getenv("STRAVA_CLIENT_SECRET") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Environment not loaded.")
		}
	}

	store, err := NewSQLRunerStore()
	if err != nil {
		panic(err)
	}
	server := RunnerServer{store, filepath.FromSlash("html")}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	m := chiprometheus.NewMiddleware("Track My Run")

	r.Use(m)
	r.Use(middleware.Recoverer)

	r.Handle("/metrics", promhttp.Handler())
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/runs", server.ServeHTTP)
	r.Post("/runs", server.ServeHTTP)
	log.Println("Handlers Initiated")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Println("Listening Port configured, starting server on:", port)
	http.ListenAndServe(":5000", r)

}
