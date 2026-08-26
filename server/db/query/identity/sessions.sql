-- name: CreateSession :one
INSERT INTO sessions (
  id,
  user_id,
  refresh_token_hash,
  expires_at,
  revoked_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(user_id),
  sqlc.arg(refresh_token_hash),
  sqlc.arg(expires_at),
  sqlc.arg(revoked_at)
)
RETURNING *;

-- name: RevokeSession :exec
UPDATE sessions
SET revoked_at = COALESCE(revoked_at, sqlc.arg(revoked_at))
WHERE id = sqlc.arg(id);

-- name: FindByRefreshTokenHash :one
SELECT *
FROM sessions s
WHERE s.refresh_token_hash = sqlc.arg(refresh_token_hash);

-- name: FindByRefreshTokenHashForUpdate :one
SELECT *
FROM sessions s
WHERE s.refresh_token_hash = sqlc.arg(refresh_token_hash)
LIMIT 1
FOR UPDATE;

-- name: UpdateRefreshTokenHash :one
UPDATE sessions
SET refresh_token_hash = sqlc.arg(new_refresh_token_hash)
WHERE id = sqlc.arg(id)
  AND revoked_at IS NULL
  AND expires_at > now()
RETURNING *;
