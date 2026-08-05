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

-- name: UpdateChatRoomByID :exec
UPDATE chat_rooms
SET room_type = sqlc.arg(room_type),
    name = sqlc.narg(name),
    description = sqlc.narg(description),
    last_message_id = sqlc.narg(last_message_id),
    updated_at = now()
WHERE id = sqlc.arg(id);

-- name: DeleteChatRoomByID :exec
UPDATE chat_rooms
SET deleted_at = now()
WHERE id = sqlc.arg(id);
