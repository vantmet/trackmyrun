package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	store := InMemoryRunnerStore{}
	server := RunnerServer{&store}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/runs", server.ServeHTTP)
	r.Post("/runs", server.ServeHTTP)
	http.ListenAndServe(":5000", r)

}
