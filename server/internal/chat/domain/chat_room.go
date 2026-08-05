package domain

import (
	"github.com/google/uuid"
)

type ChatRoom struct {
	ID          uuid.UUID
	RoomType    RoomType
	RoomName    RoomName
	Description RoomDescription
	OwnerId     RoomOwnerId
}

func NewChatroom(rt RoomType, rn RoomName, d RoomDescription, oid RoomOwnerId) *ChatRoom {
	return &ChatRoom{
		ID:          uuid.New(),
		RoomType:    rt,
		RoomName:    rn,
		Description: d,
		OwnerId:     oid,
	}
}
