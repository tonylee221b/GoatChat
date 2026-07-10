-- name: CreateUser :one
INSERT INTO users (
  id,
  username,
  email,
  phone_number
) VALUES (
  sqlc.arg(id),
  sqlc.arg(username),
  sqlc.arg(email),
  sqlc.arg(phone_number)
)
RETURNING *;

-- name: GetUserByUsername :one
SELECT *
FROM users u
WHERE u.username = sqlc.arg(username)
  AND deleted_at IS NULL;
