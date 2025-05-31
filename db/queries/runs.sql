-- name: GetRuns :many
SELECT * FROM runs LIMIT 10;

-- name: GetLastRun :one
SELECT * FROM runs ORDER BY date DESC LIMIT 1;

-- name: CreateRun :one
INSERT into runs (
	date, 
	distance, 
	runtime
) values (
	$1,$2,$3
	)
	returning *;
