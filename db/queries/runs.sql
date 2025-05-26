-- name: GetRuns :many
SELECT * FROM runs;

-- name: CreateRun :one
INSERT into runs (
	date, 
	distance, 
	runtime
) values (
	$1,$2,$3
	)
	returning *;
