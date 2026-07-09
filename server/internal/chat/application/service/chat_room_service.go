package service

import (
	"context"
	"errors"
	"log/slog"

	port "GoatChat/GoatChat/internal/chat/application/port"
	"GoatChat/GoatChat/internal/chat/domain"
)

type ChatService struct {
	repo port.ChatRepository
}

func NewChatService(repo port.ChatRepository) (*ChatService, error) {
	if repo == nil {
		return nil, errors.New("Received repo as nil")
	}

	return &ChatService{repo}, nil
}

func (svc *ChatService) CreateChatroom(ctx context.Context, rt domain.RoomType, rn domain.RoomName, d domain.RoomDescription, oid domain.RoomOwnerId) (*domain.ChatRoom, error) {
	cr, err := domain.NewChatroom(rt, rn, d, oid)
	if err != nil {
		slog.Info("failed to create chat room", "err", err.Error())
		return nil, errors.New("failed to create chat room, error: " + err.Error())
	}
	slog.Debug("chatroom domain: ", slog.Any("Chatroom", *cr))

	err = svc.repo.SaveChatroom(ctx, *cr)
	if err != nil {
		slog.Info("failed to save chat room", "err", err.Error())
		return nil, errors.New("failed to save chat room, error: " + err.Error())
	}

	return cr, nil
}
