package application

import (
	"context"
	"errors"
	"log/slog"

	port "GoatChat/GoatChat/internal/chat/application/port"
	"GoatChat/GoatChat/internal/chat/domain"
	dbtx "GoatChat/GoatChat/internal/shared/db_tx"

	"github.com/google/uuid"
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
	var cr *domain.ChatRoom
	err := svc.tx.WithinTx(ctx, func(ctx context.Context) error {
		newCr, err := domain.NewChatroom(rt, rn, d, oid)
		if err != nil {
			slog.Info("failed to create chat room", "err", err.Error())
			return errors.New("failed to create chat room, error: " + err.Error())
		}
		slog.Debug("chatroom domain: ", slog.Any("Chatroom", *newCr))

		err = svc.repo.Save(ctx, *newCr)
		if err != nil {
			slog.Info("failed to save chat room", "err", err.Error())
			return errors.New("failed to save chat room, error: " + err.Error())
		}
		cr = newCr

		return nil
	})

	return cr, err
}

func (svc *ChatService) DeleteChatroom(ctx context.Context, id uuid.UUID) error {
	return svc.repo.Delete(ctx, id)
}

func (svc *ChatService) UpdateChatroom(ctx context.Context,
	id uuid.UUID,
	rt domain.RoomType,
	rn domain.RoomName,
	rd domain.RoomDescription,
) error {
	return svc.tx.WithinTx(ctx, func(ctx context.Context) error {

		cfd, err := svc.repo.FindByChatroomId(ctx, id)
		if err != nil {
			return errors.New("failed to update chat room, error: " + err.Error())
		}

		cfd.UpdateRoomType(rt)
		cfd.UpdateRoomName(rn)
		cfd.UpdateRoomDescription(rd)

		err = svc.repo.Update(ctx, *cfd)
		if err != nil {
			slog.Info("failed to update chat room", "err", err.Error())
			return errors.New("failed to update chat room, error: " + err.Error())
		}

		return nil
	})
}
