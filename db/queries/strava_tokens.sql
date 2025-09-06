-- name: GetStravaToken :one
SELECT * FROM strava_tokens WHERE id=$1;

-- name: StoreStravaToken :one
UPDATE strava_tokens SET 
access_token=$2, expires_at=$3, refresh_token=$4
WHERE id=$1
RETURNING *;

-- name: NewStravaToken :one
INSERT INTO strava_tokens (
	id, access_token, expires_at, refresh_token
	) VALUES (
	$1,$2,$3,$4
	)
RETURNING *;
