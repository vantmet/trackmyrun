package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	logPath := "development.log"

	openLogFile(logPath)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	handler := &RunnerServer{&InMemoryRunnerStore{}}
	log.Fatal(http.ListenAndServe(":5000", handler))
}

func openLogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile:", err)
		}

		log.SetOutput(lf)
	}
}
