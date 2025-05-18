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

func (rs *SQLRunnerStore) GetRunnerStravaToken(userid int) (StravaToken, error) {
	var token StravaToken
	var id int
	sqlStatement := `SELECT * from strava_tokens WHERE token_id=$1;`

	row := rs.handle.QueryRow(sqlStatement, userid)
	switch err := row.Scan(
		&id,
		&token.AccessToken,
		&token.ExpiresAt,
		&token.ExpiresIn,
		&token.RefreshToken,
	); err {
	case sql.ErrNoRows:
		log.Printf("Select Failed: %q", err)
		return token, nil
	case nil:
		return token, nil
	default:
		panic(err)
	}
}

func setupTables(db *sql.DB) error {
	var err error
	// Runs table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS runs(
		date TIMESTAMPTZ NOT NULL,
		distance float(32) NOT NULL,
		runtime VARCHAR NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("could not create table runs: %w", err)
	}

	// Users Table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(
		user_id INT GENERATED ALWAYS AS IDENTITY,
		strava_token INT
	)`)
	if err != nil {
		return fmt.Errorf("could not create table users: %w", err)
	}

	//Strava Token Table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS strava_tokens(
		token_id INT GENERATED ALWAYS AS IDENTITY,
		access_token VARCHAR NOT NULL,
		expires_at INT NOT NULL,
		expires_in INT NOT NULL,
		refresh_token VARCHAR NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("could not create table strava_tokens: %w", err)
	}

	return nil
}
