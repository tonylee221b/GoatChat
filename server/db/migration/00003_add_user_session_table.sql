-- +goose Up
ATLER TABLE users 
  ADD COLUMN password_hash TEXT;

UPDATE users
SET password_hash = 'temp-pw'
WHERE password_hash IS NULL;

ALTER TABLE users
  ALTER COLUMN password_hash SET NOT NULL,
    ADD CONSTRAINT users_password_hash_not_blank
    CHECK (btrim(password_hash) <> '');


CREATE TABLE sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),

  refresh_token_hash BYTEA NOT NULL UNIQUE,

  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL,

  CONSTRAINT sessions_expiration_check
    CHECK (expires_at > created_at),

  CONSTRAINT sessions_refresh_token_hash_length_check
    CHECK (octet_length(refresh_token_hash) = 32)
);

CREATE INDEX idx_sessions_user_id
  ON sessions(user_id);

CREATE INDEX idx_sessions_active_user_id
  ON sessions(user_id)
  WHERE revoked_at IS NULL;

CREATE INDEX idx_sessions_expires_at
  ON sessions(expires_at)
  WHERE revoked_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS sessions;
