-- name: CreateSession :one
INSERT INTO sessions (
  id,
  user_id,
  refresh_token_hash,
  expires_at,
  revoked_at,
) VALUES (
  sqlc.arg(id),
  sqlc.arg(user_id),
  sqlc.arg(refresh_token_hash),
  sqlc.arg(expires_at),
  sqlc.arg(revoked_at)
)
RETURNING *;

-- name: FindSessionByRefreshTokenHash :one
SELECT *
FROM sessions s
WHERE s

-- name: FindSessionByUserID :many
