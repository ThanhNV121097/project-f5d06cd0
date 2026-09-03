-- Create greetings table for persisted demo messages.
CREATE TABLE IF NOT EXISTS greetings (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name TEXT NOT NULL,
  message TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT ck_greetings_name_not_blank CHECK (length(btrim(name)) > 0),
  CONSTRAINT ck_greetings_name_short CHECK (length(name) <= 80),
  CONSTRAINT ck_greetings_message_not_blank CHECK (length(btrim(message)) > 0),
  CONSTRAINT ck_greetings_message_short CHECK (length(message) <= 240)
);

CREATE INDEX IF NOT EXISTS idx_greetings_created_at_id ON greetings (created_at DESC, id DESC);
