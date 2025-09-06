CREATE TABLE IF NOT EXISTS runs(
  date TIMESTAMPTZ NOT NULL,
  distance double precision  NOT NULL,
  runtime integer NOT NULL
);
