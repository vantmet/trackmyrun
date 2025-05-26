CREATE TABLE IF NOT EXISTS strava_tokens (
   id UUID PRIMARY KEY,
   acess_token VARCHAR(255) NOT NULL,
   expires_at integer NOT NULL,
   refresh_token VARCHAR(255) NOT NULL
)
