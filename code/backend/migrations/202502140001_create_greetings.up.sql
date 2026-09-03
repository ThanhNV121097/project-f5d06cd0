-- Create greetings table for persisted demo messages.
CREATE TABLE IF NOT EXISTS greetings (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL CHECK (length(trim(name)) > 0 AND length(name) <= 80),
  message TEXT NOT NULL CHECK (length(trim(message)) > 0 AND length(message) <= 500),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS greetings_created_at_idx ON greetings (created_at DESC, id DESC);
