package domain

import "github.com/google/uuid"

type ChatRoom struct {
	id     uuid.UUID
	member RoomMember
}

func NewChatroom(member RoomMember) (*ChatRoom, error) {
	return &ChatRoom{
		id:     uuid.New(),
		member: member,
	}, nil
}
