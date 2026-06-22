package service

import (
	"errors"

	"GoatChat/GoatChat/internal/chat/domain"
)

type ChatRoomCreateService struct {
	repo ChatRoomRepository
}

func NewChatroomCreateService(repo *ChatRoomCreateService) (*ChatRoomCreateService, error) {
	if repo == nil {
		return nil, errors.New("Received repo as nil")
	}

	return &ChatRoomCreateService{repo}, nil
}

func (svc *ChatRoomCreateService) Create(name string) (*domain.ChatRoom, error) {
	if name == "" {
		return nil, errors.New("Name cannot be blank")
	}

	return domain.NewChatroom(domain.RoomMember{})
}
