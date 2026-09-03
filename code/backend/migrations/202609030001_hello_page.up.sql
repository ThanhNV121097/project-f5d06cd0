-- Create greetings table for persisted demo messages.
CREATE TABLE IF NOT EXISTS greetings (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name TEXT NOT NULL CHECK (length(btrim(name)) > 0 AND length(name) <= 80),
  message TEXT NOT NULL CHECK (length(btrim(message)) > 0 AND length(message) <= 240),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_greetings_created_at_id ON greetings (created_at DESC, id DESC);
