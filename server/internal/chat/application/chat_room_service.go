package application

import (
	"context"
	"errors"
	"log/slog"

	port "GoatChat/GoatChat/internal/chat/application/port"
	"GoatChat/GoatChat/internal/chat/domain"
	dbtx "GoatChat/GoatChat/internal/shared/db_tx"
)

type ChatService struct {
	tx   dbtx.Tx
	repo port.ChatRepository
}

func NewChatService(tx dbtx.Tx, repo port.ChatRepository) *ChatService {
	return &ChatService{tx, repo}
}

func (svc *ChatService) CreateChatroom(ctx context.Context,
	rt domain.RoomType,
	rn domain.RoomName,
	d domain.RoomDescription,
	oid domain.RoomOwnerId,
) (*domain.ChatRoom, error) {
	cr := domain.NewChatroom(rt, rn, d, oid)
	slog.Debug("chatroom domain: ", slog.Any("Chatroom", *cr))

	err := svc.repo.SaveChatroom(ctx, *cr)
	if err != nil {
		slog.Info("failed to save chat room", "err", err.Error())
		return nil, errors.New("failed to save chat room, error: " + err.Error())
	}

	return cr, nil
}
