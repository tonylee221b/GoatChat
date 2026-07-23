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

-- name: ExistsUserByUsername :one
SELECT EXISTS (
  SELECT 1
  FROM users u
  WHERE u.username = sqlc.arg(username)
);

-- name: UpdateContact :one
UPDATE users
SET
  phone_number = sqlc.arg(phone_number),
  email = sqlc.arg(email),
  updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;
  
