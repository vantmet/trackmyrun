-- name: GetStravaToken :one
SELECT * FROm strava_tokens WHERE id=$1;
