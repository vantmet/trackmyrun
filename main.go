package main

import (
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(RunnerServer)
	log.Fatal(http.ListenAndServe(":5000", handler))
}
