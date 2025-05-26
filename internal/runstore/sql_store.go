package runstore

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type SQLRunnerStore struct {
	handle *Queries
	ctx    context.Context
}

func NewSQLRunerStore() (*SQLRunnerStore, error) {
	ctx := context.Background()
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("TMRDBHOST"),
		os.Getenv("TMRDBPORT"),
		os.Getenv("TMRDBUSER"),
		os.Getenv("TMRDBPASSWORD"),
		os.Getenv("TMRDBNAME"),
	)

	conn, err := pgx.Connect(ctx, psqlInfo)
	if err != nil {
		return &SQLRunnerStore{}, err
	}
	defer conn.Close(ctx)

	queries := New(conn)

	schemaVersion, err := queries.GetSchemaVersion(ctx)
	if err != nil {
		return &SQLRunnerStore{}, err
	}

	if (schemaVersion.Version) <= int32(4) {
		log.Println("Incorrect DB Schema")
		return &SQLRunnerStore{}, err
	}
	return &SQLRunnerStore{handle: queries, ctx: ctx}, nil
}

func (rs *SQLRunnerStore) GetRunnerRuns() []Run {
	userRuns, err := rs.handle.GetRuns(rs.ctx)
	if err != nil {
		log.Printf("Unable to get runs: %q", err)
	}
	return userRuns
}

func (rs *SQLRunnerStore) RecordRun(r Run) {
	run, err := rs.handle.CreateRun(rs.ctx, CreateRunParams{
		Date:     r.Date,
		Distance: r.Distance,
		Runtime:  r.Runtime})
	if err != nil {
		log.Printf("Unable to save run: %q", err)
	}
	if run != r {
		log.Printf("Unable to save run: %q", err)
	}
}
