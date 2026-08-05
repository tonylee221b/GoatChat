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

func NewChatroom(rt RoomType, rn RoomName, d RoomDescription, oid RoomOwnerId) (*ChatRoom, error) {
	return &ChatRoom{
		ID:          uuid.New(),
		RoomType:    rt,
		RoomName:    rn,
		Description: d,
		OwnerId:     oid,
	}, nil
}

func (c *ChatRoom) UpdateRoomType(rt RoomType) {
	c.RoomType = rt
}

func (c *ChatRoom) UpdateRoomName(rn RoomName) {
	c.RoomName = rn
}

func (c *ChatRoom) UpdateRoomDescription(rd RoomDescription) {
	c.Description = rd
}
