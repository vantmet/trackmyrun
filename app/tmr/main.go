package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	chiprometheus "github.com/jamscloud/chi-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	store := InMemoryRunnerStore{}
	server := RunnerServer{&store, "/opt/tmr/html"}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	m := chiprometheus.NewMiddleware("serviceName")
	r.Use(m)
	r.Use(middleware.Recoverer)

	r.Handle("/metrics", promhttp.Handler())
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/runs", server.ServeHTTP)
	r.Post("/runs", server.ServeHTTP)
	http.ListenAndServe(":5000", r)

}
