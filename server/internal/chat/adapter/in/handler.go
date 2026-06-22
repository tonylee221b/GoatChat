package in

import (
	"errors"

	"GoatChat/GoatChat/internal/chat/application/service"
)

type ChatHandler struct {
	crSvc service.ChatRoomCreateService
}

type ChatroomResponse struct {
	name string `json:"name"`
	// ...
}

func NewChatHandler(crSvc service.ChatRoomCreateService) (*ChatHandler, error) {
	return &ChatHandler{crSvc}, nil
}

func (h *ChatHandler) CreateChatroom(name string) (ChatroomResponse, error) {
	cr, err := h.crSvc.Create(name)
	if err != nil {
		return ChatroomResponse{}, errors.New("Sival, something so bad")
	}

	return ChatroomResponse{
		name: cr.name,
	}, nil
}
