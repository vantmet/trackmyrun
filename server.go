package main

import (
	"fmt"
	"net/http"
)

func RunnerServer(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "Latest Runs") }
