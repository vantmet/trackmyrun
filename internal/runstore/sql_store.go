package runstore

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLRunnerStore struct {
	handle *Queries
	ctx    context.Context
}

const targetSchemaVersion = 3

func Config() *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5
	DATABASE_URL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("TMRDBHOST"),
		os.Getenv("TMRDBPORT"),
		os.Getenv("TMRDBUSER"),
		os.Getenv("TMRDBPASSWORD"),
		os.Getenv("TMRDBNAME"),
	)

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database!!")
	}

	return dbConfig
}

func NewSQLRunerStore(ctx context.Context) (*SQLRunnerStore, error) {
	connPool, err := pgxpool.NewWithConfig(ctx, Config())
	if err != nil {
		return &SQLRunnerStore{}, err
	}

	queries := New(connPool)

	schemaVersion, err := queries.GetSchemaVersion(ctx)
	if err != nil {
		return &SQLRunnerStore{}, err
	}

	if (schemaVersion.Version) < int32(targetSchemaVersion) {
		log.Println("Incorrect DB Schema")
		return &SQLRunnerStore{}, err
	}
	// defer connPool.Close()
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
	if reflect.DeepEqual(r, run) {
		log.Printf("Unable to save run: %v isnt: %v", run, r)
	}
}

func (rs *SQLRunnerStore) GetRunnerStravaToken(tokenid uuid.UUID) (StravaToken, error) {
	return rs.handle.GetStravaToken(rs.ctx, tokenid)
}

func (rs *SQLRunnerStore) NewRunnerStravaToken(token StravaToken) (StravaToken, error) {
	return rs.handle.NewStravaToken(rs.ctx, NewStravaTokenParams(token))
}

func (rs *SQLRunnerStore) UpdateRunnerStravaToken(token StravaToken) (StravaToken, error) {
	return rs.handle.StoreStravaToken(rs.ctx, StoreStravaTokenParams(token))
}

func (rs *SQLRunnerStore) GetLastRunnerRun() (Run, error) {
	return rs.handle.GetLastRun(rs.ctx)
}
