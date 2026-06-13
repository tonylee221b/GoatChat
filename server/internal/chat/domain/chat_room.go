package domain

import "github.com/google/uuid"

type ChatRoom struct {
	id     uuid.UUID
	member RoomMember
}
