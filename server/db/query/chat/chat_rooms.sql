-- name: CreateChatRoom :one
INSERT INTO chat_rooms (
  id,
  room_type,
  name,
  description,
  owner_id
) VALUES (
  sqlc.arg(id),
  sqlc.arg(room_type),
  sqlc.narg(name),
  sqlc.narg(description),
  sqlc.arg(owner_id)
)
RETURNING *;

-- name: GetChatRoomByID :one
SELECT *
FROM chat_rooms
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;
