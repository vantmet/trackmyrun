CREATE TABLE IF NOT EXISTS users(
   id UUID PRIMARY KEY,
   name VARCHAR (50) UNIQUE NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL,
   strava_token UUID
)
