package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type SQLRunnerStore struct {
	handle *sql.DB
}

func NewSQLRunerStore() (*SQLRunnerStore, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("TMRDBHOST"),
		os.Getenv("TMRDBPORT"),
		os.Getenv("TMRDBUSER"),
		os.Getenv("TMRDBPASSWORD"),
		os.Getenv("TMRDBNAME"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return &SQLRunnerStore{}, err
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		return &SQLRunnerStore{}, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS runs(
		date TIMESTAMPTZ NOT NULL,
		distance float(32) NOT NULL,
		runtime VARCHAR NOT NULL
	)`)
	if err != nil {
		return &SQLRunnerStore{}, fmt.Errorf("could not create table: %w", err)
	}
	return &SQLRunnerStore{handle: db}, nil
}

func (rs *SQLRunnerStore) GetRunnerRuns() []Run {
	userRuns := []Run{}
	rows, err := rs.handle.Query("SELECT * FROM runs")
	defer rows.Close()
	if err != nil {
		log.Printf("Select Failed: %q", err)
		return []Run{}
	}
	for rows.Next() {
		var run Run
		var tempTime string
		if err := rows.Scan(&run.Date, &run.Distance, &tempTime); err != nil {
			return userRuns
		}
		err = json.Unmarshal([]byte(tempTime), &run.RunTime)
		if err != nil {
			log.Printf("Unable to demarshall runtime: %q", err)
		}
		userRuns = append(userRuns, run)
	}
	return userRuns
}

func (rs *SQLRunnerStore) RecordRun(r Run) {
	time, err := json.Marshal(r.RunTime)
	if err != nil {
		return
	}
	_, err = rs.handle.Exec("INSERT INTO runs VALUES ($1, $2, $3)", r.Date, r.Distance, time)
	if err != nil {
		log.Printf("Error adding run:: %v", err)
	}
}
